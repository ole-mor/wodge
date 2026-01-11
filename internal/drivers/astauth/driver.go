package astauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AstAuthService interface {
	Register(ctx context.Context, email, password, firstName, lastName string) error
	Login(ctx context.Context, email, password string) (*AuthResponse, error)
	// RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
	// VerifyToken(ctx context.Context, token string) (*UserResponse, error)
}

type AstAuthDriver struct {
	BaseURL string
	Client  *http.Client
}

func NewAstAuthDriver(baseURL string) *AstAuthDriver {
	return &AstAuthDriver{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"user"`
}

func (d *AstAuthDriver) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	reqBody := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", d.BaseURL+"/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed (status %d): %s", resp.StatusCode, string(body))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}

func (d *AstAuthDriver) Register(ctx context.Context, email, password, firstName, lastName string) error {
	reqBody := map[string]string{
		"email":      email,
		"password":   password,
		"first_name": firstName,
		"last_name":  lastName,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", d.BaseURL+"/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("register failed (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
