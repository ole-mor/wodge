package qast

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	Text   string `json:"text"`
	UserID string `json:"user_id"`
}

func (q *QastDriver) SecureChat(ctx context.Context, text, userId string) (string, map[string]string, error) {
	if q == nil || q.httpClient == nil {
		return "", nil, fmt.Errorf("qast driver is nil")
	}

	reqBody := secureChatRequest{
		Text:   text,
		UserID: userId,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/privacy/chat", q.baseURL)
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

	var respBody secureChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return respBody.LLMResponse, respBody.TokenMap, nil
}
