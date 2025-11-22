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
	// Response models from authorization service
	authServiceResponse struct {
		Status  string    `json:"status"`
		Message string    `json:"message"`
		Data    tokenData `json:"data"`
	}

	tokenData struct {
		AcessToken   string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	// Request models for authorization service
	registerRequest struct {
		Login    string `json:"login"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)

// RegisterUser registers a new user via the authorization service REST API
// Returns access token, refresh token, and error
func RegisterUser(config *lib.Config, login string, email string, password string) (string, string, error) {
	if config == nil {
		return "", "", fmt.Errorf("config is nil")
	}

	if config.AuthService.Mode != "on" {
		return "", "", fmt.Errorf("authorization service is disabled")
	}

	// Build request URL
	url := fmt.Sprintf("http://%s/v1/register", config.AuthService.Addr)

	// Create request body
	reqBody := registerRequest{
		Login:    login,
		Email:    email,
		Password: password,
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
	if resp.StatusCode != http.StatusCreated {
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

// LoginUser authenticates a user via the authorization service REST API
// Returns access token, refresh token, and error
func LoginUser(config *lib.Config, email string, password string) (string, string, error) {
	if config == nil {
		return "", "", fmt.Errorf("config is nil")
	}

	if config.AuthService.Mode != "on" {
		return "", "", fmt.Errorf("authorization service is disabled")
	}

	// Build request URL
	url := fmt.Sprintf("http://%s/v1/login", config.AuthService.Addr)

	// Create request body
	reqBody := loginRequest{
		Email:    email,
		Password: password,
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
