package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// usageEndpoint is the OAuth usage endpoint
	usageEndpoint = "https://api.anthropic.com/api/oauth/usage"

	// anthropicBeta is the required beta header for OAuth endpoints
	anthropicBeta = "oauth-2025-04-20"

	// userAgent mimics the Claude CLI
	userAgent = "claude-code/2.1.27"
)

// Client is a client for fetching rate limit information from the Anthropic API.
type Client struct {
	httpClient *http.Client
	token      string
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

// FetchRateLimits fetches usage data from the OAuth usage endpoint.
// This is a free endpoint that doesn't consume any tokens.
func (c *Client) FetchRateLimits() (*RateLimitData, error) {
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

	// Check for errors
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
