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
	Login(ctx context.Context, username, password string) (*AuthResponse, error)
	Register(ctx context.Context, email, username, password, confirmPassword, firstName, lastName, inviteCode string) error
	VerifyToken(ctx context.Context, accessToken string) (*User, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
	Logout(ctx context.Context, accessToken, refreshToken string) error
	ListUsers(ctx context.Context, accessToken string, limit, offset int) (*UserListResponse, error)
	ListRoles(ctx context.Context, accessToken string) (*RoleListResponse, error)
	AddRoleToUser(ctx context.Context, accessToken, userID, roleID string) error
	RemoveRoleFromUser(ctx context.Context, accessToken, userID, roleID string) error
	GenerateInvite(ctx context.Context, accessToken string) (*InviteToken, error)
	ListInvites(ctx context.Context, accessToken string) ([]InviteToken, error)
	ActivateUser(ctx context.Context, accessToken, userID string, active bool) (*User, error)
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
	Username  string `json:"username,omitempty"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Active    bool   `json:"active"`
	Roles     []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"roles,omitempty"`
}

type UserListResponse struct {
	Users      []User `json:"users"`
	TotalCount int64  `json:"total_count"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}

type RoleListResponse struct {
	Roles      []interface{} `json:"roles"`
	TotalCount int64         `json:"total_count"`
}

type InviteToken struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	CreatedBy string    `json:"created_by"`
	UsedBy    *string   `json:"used_by"`
	IsUsed    bool      `json:"is_used"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (d *AstAuthDriver) Login(ctx context.Context, username, password string) (*AuthResponse, error) {
	reqBody := map[string]string{
		"username": username,
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

func (d *AstAuthDriver) Register(ctx context.Context, email, username, password, confirmPassword, firstName, lastName, inviteCode string) error {
	reqBody := map[string]string{
		"email":            email,
		"username":         username,
		"password":         password,
		"confirm_password": confirmPassword,
		"first_name":       firstName,
		"last_name":        lastName,
		"invite_code":      inviteCode,
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

func (d *AstAuthDriver) Logout(ctx context.Context, accessToken, refreshToken string) error {
	reqBody := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", d.BaseURL+"/api/v1/auth/logout", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("logout failed (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (d *AstAuthDriver) ListUsers(ctx context.Context, accessToken string, limit, offset int) (*UserListResponse, error) {
	url := fmt.Sprintf("%s/api/v1/users/all-users?limit=%d&offset=%d", d.BaseURL, limit, offset)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list users failed (status %d): %s", resp.StatusCode, string(body))
	}

	var respData UserListResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, fmt.Errorf("decode users failed: %w", err)
	}
	return &respData, nil
}

func (d *AstAuthDriver) ListRoles(ctx context.Context, accessToken string) (*RoleListResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", d.BaseURL+"/api/v1/users/system-roles", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list roles failed (status %d): %s", resp.StatusCode, string(body))
	}

	var respData RoleListResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, fmt.Errorf("decode roles failed: %w", err)
	}
	return &respData, nil
}

func (d *AstAuthDriver) AddRoleToUser(ctx context.Context, accessToken, userID, roleID string) error {
	reqBody := map[string]string{"user_id": userID, "role_id": roleID}
	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", d.BaseURL+"/api/v1/users/"+userID+"/roles", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("add role failed: %s", string(b))
	}
	return nil
}

func (d *AstAuthDriver) RemoveRoleFromUser(ctx context.Context, accessToken, userID, roleID string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", d.BaseURL+"/api/v1/users/"+userID+"/roles/"+roleID, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := d.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("remove role failed: %s", string(b))
	}
	return nil
}

func (d *AstAuthDriver) GenerateInvite(ctx context.Context, accessToken string) (*InviteToken, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", d.BaseURL+"/api/v1/users/invites", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("generate invite failed: %s", string(body))
	}

	var invite InviteToken
	if err := json.NewDecoder(resp.Body).Decode(&invite); err != nil {
		return nil, err
	}
	return &invite, nil
}

func (d *AstAuthDriver) ListInvites(ctx context.Context, accessToken string) ([]InviteToken, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", d.BaseURL+"/api/v1/users/invites", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list invites failed: %s", string(body))
	}

	var invites []InviteToken
	if err := json.NewDecoder(resp.Body).Decode(&invites); err != nil {
		return nil, err
	}
	return invites, nil
}
func (d *AstAuthDriver) ActivateUser(ctx context.Context, accessToken, userID string, active bool) (*User, error) {
	reqBody := map[string]bool{"active": active}
	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "PUT", d.BaseURL+"/api/v1/users/"+userID+"/activate-user", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("activate user failed: %s", string(b))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUser deletes a user in AstAuth (requires admin token)
func (d *AstAuthDriver) DeleteUser(ctx context.Context, accessToken, userID string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", d.BaseURL+"/api/v1/users/"+userID+"/delete-user", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := d.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete user failed: %s", string(b))
	}
	return nil
}

// AdminResetPassword resets a user's password (requires admin token)
func (d *AstAuthDriver) AdminResetPassword(ctx context.Context, accessToken, userID, newPassword string) error {
	reqBody := map[string]string{"user_id": userID, "new_password": newPassword}
	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "PUT", d.BaseURL+"/api/v1/users/"+userID+"/reset-password", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("admin reset password failed: %s", string(b))
	}
	return nil
}
