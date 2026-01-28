package qast

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"wodge/internal/services"
)

type QastDriver struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type composerRequest struct {
	Query          string `json:"query"`
	UserID         string `json:"user_id"`
	ExpertiseLevel string `json:"expertise_level"`
}

type composerResponse struct {
	Answer  string   `json:"answer"`
	Context []string `json:"context,omitempty"`
	Error   string   `json:"error,omitempty"`
}

func NewQastDriver(baseURL string, apiKey string) *QastDriver {
	if apiKey == "" {
		apiKey = "dev-token-bypass"
	}
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}

	return &QastDriver{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Transport: transport},
	}
}

// Ensure QastDriver implements services.QastService
var _ services.QastService = (*QastDriver)(nil)

func (q *QastDriver) Ask(ctx context.Context, query, userId, expertise string) (string, []string, error) {
	if q == nil || q.httpClient == nil {
		return "", nil, fmt.Errorf("qast driver is nil")
	}

	reqBody := composerRequest{
		Query:          query,
		UserID:         userId,
		ExpertiseLevel: expertise,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/composer/ask", q.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if q.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+q.apiKey)
	}

	resp, err := q.httpClient.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("failed to call qast api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			if errMsg, ok := errResp["error"].(string); ok {
				return "", nil, fmt.Errorf("qast api error (%d): %s", resp.StatusCode, errMsg)
			}
		}
		return "", nil, fmt.Errorf("qast api returned status: %d", resp.StatusCode)
	}

	var respBody composerResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if respBody.Error != "" {
		return "", nil, fmt.Errorf("qast api error: %s", respBody.Error)
	}

	return respBody.Answer, respBody.Context, nil
}

type ingestRequest struct {
	Text         string `json:"text"`
	TemplateName string `json:"template_name"`
	UserID       string `json:"user_id"`
}

type ingestResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func (q *QastDriver) IngestGraph(ctx context.Context, text, userId string) (interface{}, error) {
	if q == nil || q.httpClient == nil {
		return nil, fmt.Errorf("qast driver is nil")
	}

	reqBody := ingestRequest{
		Text:         text,
		UserID:       userId,
		TemplateName: "extract_knowledge_graph", // Hardcoded for now
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/privacy/extract", q.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if q.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+q.apiKey)
	}

	resp, err := q.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call qast api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			if errMsg, ok := errResp["error"].(string); ok {
				return nil, fmt.Errorf("qast api error (%d): %s", resp.StatusCode, errMsg)
			}
		}
		return nil, fmt.Errorf("qast api returned status: %d", resp.StatusCode)
	}

	var respBody ingestResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if respBody.Error != "" {
		return nil, fmt.Errorf("qast api error: %s", respBody.Error)
	}

	return respBody.Result, nil
}

type secureChatResponse struct {
	LLMResponse string            `json:"llm_response"`
	TokenMap    map[string]string `json:"token_map"`
}

type secureChatRequest struct {
	Text      string `json:"text"`
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id,omitempty"`
}

// SecureChat now returns a ReadCloser for the SSE stream
func (q *QastDriver) SecureChat(ctx context.Context, text, userId, sessionId, token string) (io.ReadCloser, error) {
	if q == nil || q.httpClient == nil {
		return nil, fmt.Errorf("qast driver is nil")
	}

	reqBody := secureChatRequest{
		Text:      text,
		UserID:    userId,
		SessionID: sessionId,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/pipeline/chat", q.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	} else if q.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+q.apiKey)
	}

	// Important: we return the body to be streamed
	log.Printf("[QastDriver] Sending request to %s", url)
	resp, err := q.httpClient.Do(req)
	if err != nil {
		log.Printf("[QastDriver] httpClient.Do failed: %v", err)
		return nil, fmt.Errorf("failed to call qast api: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("qast api returned status: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// -- History Methods --

func (q *QastDriver) CreateSession(ctx context.Context, userID, title string) (interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/history/sessions", q.baseURL)
	reqBody := map[string]string{"user_id": userID, "title": title}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if q.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+q.apiKey)
	}

	resp, err := q.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create session: %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (q *QastDriver) GetSessions(ctx context.Context, userID string) (interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/history/sessions?user_id=%s", q.baseURL, userID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	if q.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+q.apiKey)
	}

	resp, err := q.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get sessions: %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (q *QastDriver) GetSession(ctx context.Context, sessionID string) (interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/history/sessions/%s", q.baseURL, sessionID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	if q.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+q.apiKey)
	}

	resp, err := q.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get session: %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (q *QastDriver) DeleteSession(ctx context.Context, sessionID string) error {
	url := fmt.Sprintf("%s/api/v1/history/sessions/%s", q.baseURL, sessionID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}
	if q.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+q.apiKey)
	}

	resp, err := q.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete session: %d", resp.StatusCode)
	}
	return nil
}

func (q *QastDriver) ShareSession(ctx context.Context, sessionID, targetUsername string) (interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/history/sessions/%s/share", q.baseURL, sessionID)
	reqBody := map[string]string{"target_username": targetUsername}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if q.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+q.apiKey)
	}

	resp, err := q.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to share session: %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (q *QastDriver) SearchUsers(ctx context.Context, query string) (interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/users/search?q=%s", q.baseURL, query)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	if q.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+q.apiKey)
	}

	resp, err := q.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to search users: %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (q *QastDriver) SyncUser(ctx context.Context, id, email, username, firstName, lastName string) error {
	url := fmt.Sprintf("%s/api/v1/users", q.baseURL)
	reqBody := map[string]string{
		"id":         id,
		"email":      email,
		"username":   username,
		"first_name": firstName,
		"last_name":  lastName,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if q.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+q.apiKey)
	}

	resp, err := q.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusConflict {
		// Ignore conflict/existing
		return fmt.Errorf("failed to sync user: %d", resp.StatusCode)
	}
	return nil
}

func (q *QastDriver) UpdateContext(ctx context.Context, id, content string) error {
	url := fmt.Sprintf("%s/api/v1/context/%s", q.baseURL, id)
	reqBody := map[string]string{"content": content}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if q.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+q.apiKey)
	}

	resp, err := q.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update context: %d", resp.StatusCode)
	}
	return nil
}

func (q *QastDriver) GetContext(ctx context.Context, id string) (interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/context/%s", q.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	if q.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+q.apiKey)
	}

	resp, err := q.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get context: %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
