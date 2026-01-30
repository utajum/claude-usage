// Package format provides utilities for formatting numbers and text.
package format

import "fmt"

// FormatTokens formats a token count in compact notation (K, M, B).
func FormatTokens(n int64) string {
	switch {
	case n >= 1_000_000_000:
		return fmt.Sprintf("%.1fB", float64(n)/1_000_000_000)
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}

// FormatPlanName formats the subscription type and rate limit tier.
func FormatPlanName(subscriptionType, rateLimitTier string) string {
	if subscriptionType == "" {
		return "Unknown"
	}
	planName := capitalizeFirst(subscriptionType)
	switch rateLimitTier {
	case "default_claude_max_5x":
		planName += " (Max 5x)"
	case "default_claude_max_20x":
		planName += " (Max 20x)"
	}
	return planName
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	if s[0] >= 'a' && s[0] <= 'z' {
		return string(s[0]-32) + s[1:]
	}
	return s
}

// FormatDuration formats a duration in a human-readable way.
func FormatDuration(seconds int64) string {
	if seconds < 0 {
		return "now"
	}
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60

	if days > 0 {
		if hours > 0 {
			return fmt.Sprintf("%dd %dh", days, hours)
		}
		return fmt.Sprintf("%dd", days)
	}
	if hours > 0 {
		if minutes > 0 {
			return fmt.Sprintf("%dh %dm", hours, minutes)
		}
		return fmt.Sprintf("%dh", hours)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	}
	return "<1m"
}
