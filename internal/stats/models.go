// Package stats provides parsing and calculation of Claude usage statistics.
package stats

import "time"

// StatsCache represents the structure of Claude's stats-cache.json file.
type StatsCache struct {
	Version          int                   `json:"version"`
	LastComputedDate string                `json:"lastComputedDate"`
	DailyActivity    []DailyActivity       `json:"dailyActivity"`
	DailyModelTokens []DailyModelTokens    `json:"dailyModelTokens"`
	ModelUsage       map[string]ModelUsage `json:"modelUsage"`
	TotalSessions    int                   `json:"totalSessions"`
	TotalMessages    int                   `json:"totalMessages"`
	FirstSessionDate string                `json:"firstSessionDate"`
}

// DailyActivity represents activity for a single day.
type DailyActivity struct {
	Date          string `json:"date"`
	MessageCount  int    `json:"messageCount"`
	SessionCount  int    `json:"sessionCount"`
	ToolCallCount int    `json:"toolCallCount"`
}

// DailyModelTokens represents token usage per model for a single day.
type DailyModelTokens struct {
	Date          string           `json:"date"`
	TokensByModel map[string]int64 `json:"tokensByModel"`
}

// ModelUsage represents aggregate usage statistics for a single model.
type ModelUsage struct {
	InputTokens              int64   `json:"inputTokens"`
	OutputTokens             int64   `json:"outputTokens"`
	CacheReadInputTokens     int64   `json:"cacheReadInputTokens"`
	CacheCreationInputTokens int64   `json:"cacheCreationInputTokens"`
	WebSearchRequests        int     `json:"webSearchRequests"`
	CostUSD                  float64 `json:"costUSD"`
}

// Credentials represents the structure of Claude's .credentials.json file.
type Credentials struct {
	ClaudeAiOauth OAuthCredentials `json:"claudeAiOauth"`
}

// OAuthCredentials contains OAuth token information.
type OAuthCredentials struct {
	AccessToken      string   `json:"accessToken"`
	RefreshToken     string   `json:"refreshToken"`
	ExpiresAt        int64    `json:"expiresAt"`
	Scopes           []string `json:"scopes,omitempty"`
	SubscriptionType string   `json:"subscriptionType"`
	RateLimitTier    string   `json:"rateLimitTier"`
}

// OpenCodeCredentials represents OpenCode's auth.json structure.
type OpenCodeCredentials struct {
	Anthropic OpenCodeAnthropicAuth `json:"anthropic"`
}

// OpenCodeAnthropicAuth contains OAuth token information in OpenCode's format.
type OpenCodeAnthropicAuth struct {
	Type    string `json:"type"`    // "oauth"
	Refresh string `json:"refresh"` // refresh token
	Access  string `json:"access"`  // access token
	Expires int64  `json:"expires"` // expiry timestamp in milliseconds
}

// WeeklyStats represents calculated weekly usage statistics.
type WeeklyStats struct {
	// WeekStart is the start of the current week (Monday 00:00 UTC).
	WeekStart time.Time

	// WeekEnd is the end of the current week (Sunday 23:59 UTC).
	WeekEnd time.Time

	// TokensByModel maps model names to their token counts for the week.
	TokensByModel map[string]int64

	// TotalTokens is the sum of all model tokens.
	TotalTokens int64

	// SubscriptionType is the user's subscription (free, pro, team, enterprise).
	SubscriptionType string

	// RateLimitTier provides additional rate limit info (e.g., "default_claude_max_5x").
	RateLimitTier string

	// --- Rate Limit Data from API (real-time) ---

	// FiveHourUtilization is the 5-hour window usage (0.0-1.0) from API
	FiveHourUtilization float64

	// WeeklyUtilization is the 7-day window usage (0.0-1.0) from API
	WeeklyUtilization float64

	// FiveHourReset is when the 5-hour window resets
	FiveHourReset time.Time

	// WeeklyReset is when the 7-day window resets
	WeeklyReset time.Time

	// RateLimitStatus is "allowed" or "throttled"
	RateLimitStatus string

	// RepresentativeClaim indicates which window is limiting ("five_hour" or "seven_day")
	RepresentativeClaim string

	// Model-specific weekly utilization
	OpusUtilization   float64
	SonnetUtilization float64
	OpusReset         time.Time
	SonnetReset       time.Time

	// HasAPIData indicates if we have real API rate limit data
	HasAPIData bool
}

// ModelDisplayName returns a human-friendly name for a model ID.
func ModelDisplayName(modelID string) string {
	displayNames := map[string]string{
		"claude-opus-4-5-20251101":   "Opus 4.5",
		"claude-sonnet-4-5-20250929": "Sonnet 4.5",
		"claude-haiku-3-5-20240307":  "Haiku 3.5",
		"claude-3-opus-20240229":     "Opus 3",
		"claude-3-sonnet-20240229":   "Sonnet 3",
		"claude-3-haiku-20240307":    "Haiku 3",
	}

	if name, ok := displayNames[modelID]; ok {
		return name
	}

	// Fallback: extract model family from ID
	if len(modelID) > 20 {
		return modelID[:20] + "..."
	}
	return modelID
}
