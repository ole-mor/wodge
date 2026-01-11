package astauth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	User         User   `json:"user"`
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
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

// RefreshToken refreshes the access token using a valid refresh token.
// Corresponds to POST /api/v1/auth/refresh-token
func (d *AstAuthDriver) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	reqBody := map[string]string{
		"refresh_token": refreshToken,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", d.BaseURL+"/api/v1/auth/refresh-token", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	r, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(r.Body).Decode(&errResp)
		msg := errResp["error"]
		if msg == "" {
			msg = "refresh failed"
		}
		return nil, errors.New(msg)
	}

	var resp AuthResponse
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// VerifyToken verifies the validity of an access token by fetching the user profile.
// Corresponds to GET /api/v1/users/me
func (d *AstAuthDriver) VerifyToken(ctx context.Context, accessToken string) (*User, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", d.BaseURL+"/api/v1/users/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	r, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, errors.New("invalid token")
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
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
