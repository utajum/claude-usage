package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	// usageEndpoint is the OAuth usage endpoint
	usageEndpoint = "https://api.anthropic.com/api/oauth/usage"

	// tokenEndpoint is the OAuth token refresh endpoint
	tokenEndpoint = "https://platform.claude.com/v1/oauth/token"

	// clientID is the official Claude Code CLI OAuth client ID.
	// This is the same client ID used by Anthropic's official Claude Code binary
	// (verified in claude-code v2.1.27). This is a public identifier and is not
	// considered sensitive. The official CLI also supports overriding this via
	// the CLAUDE_CODE_OAUTH_CLIENT_ID environment variable if needed.
	clientID = "9d1c250a-e61b-44d9-88ed-5944d1962f5e"

	// anthropicBeta is the required beta header for OAuth endpoints
	anthropicBeta = "oauth-2025-04-20"

	// userAgent mimics the Claude CLI
	userAgent = "claude-code/2.1.27"

	// maxRetries is the maximum number of retry attempts for token refresh
	maxRetries = 5
)

// Client is a client for fetching rate limit information from the Anthropic API.
type Client struct {
	httpClient   *http.Client
	token        string
	refreshToken string
}

// NewClient creates a new API client with the given OAuth token.
func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		token: token,
	}
}

// SetToken updates the OAuth token.
func (c *Client) SetToken(token string) {
	c.token = token
}

// SetRefreshToken updates the OAuth refresh token.
func (c *Client) SetRefreshToken(refreshToken string) {
	c.refreshToken = refreshToken
}

// usageResponse represents the response from /api/oauth/usage
type usageResponse struct {
	FiveHour struct {
		Utilization float64 `json:"utilization"`
		ResetsAt    string  `json:"resets_at"`
	} `json:"five_hour"`
	SevenDay struct {
		Utilization float64 `json:"utilization"`
		ResetsAt    string  `json:"resets_at"`
	} `json:"seven_day"`
	SevenDayOauthApps *usageBucket `json:"seven_day_oauth_apps"`
	SevenDayOpus      *usageBucket `json:"seven_day_opus"`
	SevenDaySonnet    *usageBucket `json:"seven_day_sonnet"`
	SevenDayCowork    *usageBucket `json:"seven_day_cowork"`
	ExtraUsage        struct {
		IsEnabled    bool     `json:"is_enabled"`
		MonthlyLimit *float64 `json:"monthly_limit"`
		UsedCredits  *float64 `json:"used_credits"`
		Utilization  *float64 `json:"utilization"`
	} `json:"extra_usage"`
}

type usageBucket struct {
	Utilization float64 `json:"utilization"`
	ResetsAt    string  `json:"resets_at"`
}

// tokenRefreshRequest represents the request body for token refresh
type tokenRefreshRequest struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
	Scope        string `json:"scope"`
}

// tokenRefreshResponse represents the response from token refresh
type tokenRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// FetchRateLimits fetches usage data from the OAuth usage endpoint.
// This is a free endpoint that doesn't consume any tokens.
func (c *Client) FetchRateLimits() (*RateLimitData, error) {
	return c.fetchRateLimitsWithRetry(0)
}

// fetchRateLimitsWithRetry implements retry logic with automatic token refresh on 401
func (c *Client) fetchRateLimitsWithRetry(attempt int) (*RateLimitData, error) {
	if c.token == "" {
		return nil, fmt.Errorf("no OAuth token configured")
	}

	req, err := http.NewRequest("GET", usageEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to match Claude CLI
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("anthropic-beta", anthropicBeta)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Handle 401 Unauthorized with token refresh
	if resp.StatusCode == http.StatusUnauthorized {
		if attempt >= maxRetries {
			return nil, fmt.Errorf("max retries (%d) exceeded after token refresh attempts", maxRetries)
		}

		if c.refreshToken == "" {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("token expired and no refresh token available. Status %d: %s", resp.StatusCode, string(body))
		}

		log.Printf("Token expired (attempt %d/%d), refreshing...", attempt+1, maxRetries)

		// Attempt to refresh the token
		newToken, err := c.RefreshAccessToken()
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}

		// Token already updated in RefreshAccessToken, but being explicit
		c.token = newToken
		log.Printf("Token refreshed successfully, retrying request")

		// Retry the request with the new token
		return c.fetchRateLimitsWithRetry(attempt + 1)
	}

	// Check for other errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var usage usageResponse
	if err := json.NewDecoder(resp.Body).Decode(&usage); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to RateLimitData
	return parseUsageResponse(&usage), nil
}

// RefreshAccessToken uses the refresh token to obtain a new access token.
// Returns the new access token on success.
func (c *Client) RefreshAccessToken() (string, error) {
	if c.refreshToken == "" {
		return "", fmt.Errorf("no refresh token available")
	}

	// Prepare the refresh request
	reqBody := tokenRefreshRequest{
		GrantType:    "refresh_token",
		RefreshToken: c.refreshToken,
		ClientID:     clientID,
		Scope:        "user:inference user:profile user:sessions:claude_code user:mcp_servers",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal refresh request: %w", err)
	}

	req, err := http.NewRequest("POST", tokenEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create refresh request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)

	// Make the refresh request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make refresh request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token refresh failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var refreshResp tokenRefreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&refreshResp); err != nil {
		return "", fmt.Errorf("failed to parse refresh response: %w", err)
	}

	// Update the access token
	c.token = refreshResp.AccessToken

	// Check if we got a new refresh token (some OAuth servers rotate them)
	if refreshResp.RefreshToken != "" && refreshResp.RefreshToken != c.refreshToken {
		c.refreshToken = refreshResp.RefreshToken
		log.Printf("WARNING: Received a new refresh token from the server")

		// Write a warning file next to the binary
		if err := c.writeRefreshTokenWarning(refreshResp.RefreshToken); err != nil {
			log.Printf("Failed to write refresh token warning file: %v", err)
		}
	}

	return refreshResp.AccessToken, nil
}

// writeRefreshTokenWarning creates a warning file next to the binary when a new refresh token is received
func (c *Client) writeRefreshTokenWarning(newRefreshToken string) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	warningPath := filepath.Join(filepath.Dir(exe), "NEW_REFRESH_TOKEN_WARNING.txt")

	content := fmt.Sprintf(`WARNING: New Refresh Token Received
========================================

A new refresh token was received at: %s

The ~/.claude/.credentials.json file was NOT updated.
The new refresh token is stored IN MEMORY ONLY for this session.

New refresh token: %s

If you restart claude-usage, you may need to re-authenticate if the
old refresh token has been invalidated by the server.

This file is for informational purposes only.
`,
		time.Now().Format(time.RFC3339),
		newRefreshToken,
	)

	if err := os.WriteFile(warningPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write warning file: %w", err)
	}

	log.Printf("New refresh token warning written to: %s", warningPath)
	return nil
}

// parseUsageResponse converts the API response to RateLimitData
func parseUsageResponse(usage *usageResponse) *RateLimitData {
	data := &RateLimitData{
		FetchedAt: time.Now(),
		Status:    "allowed", // Default to allowed
	}

	// Parse utilization (API returns percentage 0-100, we store as 0.0-1.0)
	data.FiveHourUtilization = usage.FiveHour.Utilization / 100.0
	data.WeeklyUtilization = usage.SevenDay.Utilization / 100.0

	// Parse reset times (ISO 8601 format)
	if t, err := time.Parse(time.RFC3339, usage.FiveHour.ResetsAt); err == nil {
		data.FiveHourReset = t
	}
	if t, err := time.Parse(time.RFC3339, usage.SevenDay.ResetsAt); err == nil {
		data.WeeklyReset = t
	}

	// Parse model-specific utilization
	if usage.SevenDayOpus != nil {
		data.OpusUtilization = usage.SevenDayOpus.Utilization / 100.0
		if t, err := time.Parse(time.RFC3339, usage.SevenDayOpus.ResetsAt); err == nil {
			data.OpusReset = t
		}
	}
	if usage.SevenDaySonnet != nil {
		data.SonnetUtilization = usage.SevenDaySonnet.Utilization / 100.0
		if t, err := time.Parse(time.RFC3339, usage.SevenDaySonnet.ResetsAt); err == nil {
			data.SonnetReset = t
		}
	}

	// Determine which window is limiting (whichever is higher)
	if data.FiveHourUtilization > data.WeeklyUtilization {
		data.RepresentativeClaim = "five_hour"
	} else {
		data.RepresentativeClaim = "seven_day"
	}

	// Check if throttled (utilization >= 100%)
	if data.FiveHourUtilization >= 1.0 || data.WeeklyUtilization >= 1.0 {
		data.Status = "throttled"
	}

	return data
}
