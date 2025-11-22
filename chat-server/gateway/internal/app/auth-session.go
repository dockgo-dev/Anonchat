package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	lib "mew-gateway/internal/libs"
	"net/http"
)

type (
	// Request model for token endpoints
	tokenRequest struct {
		Token string `json:"token"`
	}

	// Response model for validate endpoint
	userData struct {
		UserID int64  `json:"user_id"`
		Login  string `json:"login"`
		Email  string `json:"email"`
	}

	// Response model with user data
	validateResponse struct {
		Status  string   `json:"status"`
		Message string   `json:"message"`
		Data    userData `json:"data"`
	}
)

// Refresh refreshes access and refresh tokens via the authorization service REST API
// Returns new access token, new refresh token, and error
func Refresh(config *lib.Config, refreshToken string) (string, string, error) {
	if config == nil {
		return "", "", fmt.Errorf("config is nil")
	}

	if config.AuthService.Mode != "on" {
		return "", "", fmt.Errorf("authorization service is disabled")
	}

	if refreshToken == "" {
		return "", "", fmt.Errorf("refresh token is required")
	}

	// Build request URL
	url := fmt.Sprintf("http://%s/v1/refresh", config.AuthService.Addr)

	// Create request body
	reqBody := tokenRequest{
		Token: refreshToken,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", fmt.Errorf("marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		var errorResp authServiceResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return "", "", fmt.Errorf("%s: %s", errorResp.Status, errorResp.Message)
		}
		return "", "", fmt.Errorf("authorization service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response authServiceResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", "", fmt.Errorf("unmarshal response: %w", err)
	}

	// Check response status
	if response.Status != "success" {
		return "", "", fmt.Errorf("%s: %s", response.Status, response.Message)
	}

	return response.Data.AcessToken, response.Data.RefreshToken, nil
}

// Validate validates an access token via the authorization service REST API
// Returns error if token is invalid
func Validate(config *lib.Config, accessToken string) (*userData, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	if config.AuthService.Mode != "on" {
		return nil, fmt.Errorf("authorization service is disabled")
	}

	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	// Build request URL
	url := fmt.Sprintf("http://%s/v1/validate", config.AuthService.Addr)

	// Create request body
	reqBody := tokenRequest{
		Token: accessToken,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		var errorResp validateResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return nil, fmt.Errorf("%s: %s", errorResp.Status, errorResp.Message)
		}
		return nil, fmt.Errorf("authorization service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response validateResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	// Check response status
	if response.Status != "success" {
		return nil, fmt.Errorf("%s: %s", response.Status, response.Message)
	}

	return &response.Data, nil
}

// Logout logs out a user by removing the session via the authorization service REST API
// Returns error if logout fails
func Logout(config *lib.Config, refreshToken string) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}

	if config.AuthService.Mode != "on" {
		return fmt.Errorf("authorization service is disabled")
	}

	if refreshToken == "" {
		return fmt.Errorf("refresh token is required")
	}

	// Build request URL
	url := fmt.Sprintf("http://%s/v1/logout", config.AuthService.Addr)

	// Create request body
	reqBody := tokenRequest{
		Token: refreshToken,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		var errorResp authServiceResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return fmt.Errorf("%s: %s", errorResp.Status, errorResp.Message)
		}
		return fmt.Errorf("authorization service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response authServiceResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	// Check response status
	if response.Status != "success" {
		return fmt.Errorf("%s: %s", response.Status, response.Message)
	}

	return nil
}
