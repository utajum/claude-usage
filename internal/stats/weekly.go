package stats

import (
	"math/rand"
	"time"
)

// Known weekly token limits by plan type (rough estimates)
// These are approximations since Anthropic doesn't publish exact limits
var planLimits = map[string]int64{
	// Rate limit tiers (from .credentials.json rateLimitTier field)
	"default_claude_pro":     45_000_000,  // ~45M tokens/week for Pro
	"default_claude_max_5x":  225_000_000, // ~225M tokens/week for Max 5x
	"default_claude_max_20x": 900_000_000, // ~900M tokens/week for Max 20x

	// Subscription types (fallback)
	"pro":         45_000_000,  // ~45M tokens/week for Pro
	"max_5x":      225_000_000, // ~225M tokens/week for Max 5x
	"team":        225_000_000, // Team (assume Max 5x level)
	"team_max_5x": 225_000_000, // Team Max 5x
	"enterprise":  500_000_000, // Enterprise (varies)
	"free":        8_000_000,   // Free tier
}

// GetWeeklyLimit returns the estimated weekly token limit based on plan type.
// Returns 0 if plan is unknown.
func GetWeeklyLimit(subscriptionType, rateLimitTier string) int64 {
	// Check rate limit tier first (more specific)
	if rateLimitTier != "" {
		if limit, ok := planLimits[rateLimitTier]; ok {
			return limit
		}
	}

	// Fall back to subscription type
	if subscriptionType != "" {
		if limit, ok := planLimits[subscriptionType]; ok {
			return limit
		}
	}

	// Default to Pro limits
	return planLimits["pro"]
}

// GetWeekBounds returns the start and end of the current ISO week (Monday-Sunday).
// Times are in UTC.
func GetWeekBounds() (start, end time.Time) {
	now := time.Now().UTC()

	// Find the Monday of the current week
	weekday := now.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	daysFromMonday := int(weekday) - 1

	// Start of week (Monday 00:00:00)
	start = time.Date(
		now.Year(), now.Month(), now.Day()-daysFromMonday,
		0, 0, 0, 0, time.UTC,
	)

	// End of week (Sunday 23:59:59)
	end = start.AddDate(0, 0, 6)
	end = time.Date(
		end.Year(), end.Month(), end.Day(),
		23, 59, 59, 999999999, time.UTC,
	)

	return start, end
}

// CalculateWeeklyStats computes token usage for the current week.
func CalculateWeeklyStats(cache *StatsCache, creds *Credentials) *WeeklyStats {
	weekStart, weekEnd := GetWeekBounds()

	stats := &WeeklyStats{
		WeekStart:     weekStart,
		WeekEnd:       weekEnd,
		TokensByModel: make(map[string]int64),
		TotalTokens:   0,
	}

	// Set subscription info from credentials
	if creds != nil {
		stats.SubscriptionType = creds.ClaudeAiOauth.SubscriptionType
		stats.RateLimitTier = creds.ClaudeAiOauth.RateLimitTier
	}

	if cache == nil {
		return stats
	}

	// Parse and sum tokens for each day in the week
	for _, daily := range cache.DailyModelTokens {
		date, err := time.Parse("2006-01-02", daily.Date)
		if err != nil {
			continue
		}

		// Check if this date falls within the current week
		if date.Before(weekStart) || date.After(weekEnd) {
			continue
		}

		// Sum tokens by model
		for model, tokens := range daily.TokensByModel {
			stats.TokensByModel[model] += tokens
			stats.TotalTokens += tokens
		}
	}

	return stats
}

// GetDaysRemainingInWeek returns the number of days left in the current week.
func GetDaysRemainingInWeek() int {
	now := time.Now().UTC()
	weekday := now.Weekday()
	if weekday == time.Sunday {
		return 0
	}
	return 7 - int(weekday)
}

// GetWeekProgress returns a value from 0.0 to 1.0 representing progress through the week.
func GetWeekProgress() float64 {
	now := time.Now().UTC()
	start, end := GetWeekBounds()

	total := end.Sub(start).Seconds()
	elapsed := now.Sub(start).Seconds()

	if elapsed < 0 {
		return 0.0
	}
	if elapsed > total {
		return 1.0
	}
	return elapsed / total
}

// GetPercentage returns the usage percentage (0-99).
// Prefers real API data (WeeklyUtilization) over estimates based on token counts.
func (w *WeeklyStats) GetPercentage() int {
	if w == nil {
		return 0
	}

	// Prefer real API data if available
	if w.HasAPIData && w.WeeklyUtilization > 0 {
		percentage := int(w.WeeklyUtilization * 100)
		if percentage > 99 {
			percentage = 99
		}
		if percentage < 0 {
			percentage = 0
		}
		return percentage
	}

	// Fallback to estimated calculation based on token counts
	if w.TotalTokens == 0 {
		return 0
	}

	// Get the limit for this plan
	limit := GetWeeklyLimit(w.SubscriptionType, w.RateLimitTier)

	if limit == 0 {
		// Can't calculate without a limit, return placeholder
		// Use a random value that changes daily (seeded by day)
		seed := time.Now().UTC().YearDay()
		r := rand.New(rand.NewSource(int64(seed)))
		return r.Intn(100)
	}

	// Calculate percentage
	percentage := int((w.TotalTokens * 100) / limit)

	// Clamp to 0-99 (we show 99 for >=100%)
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 99 {
		percentage = 99
	}

	return percentage
}

// GetPercentageFloat returns the exact usage percentage as a float.
// Prefers real API data over estimates.
func (w *WeeklyStats) GetPercentageFloat() float64 {
	if w == nil {
		return 0.0
	}

	// Prefer real API data if available
	if w.HasAPIData && w.WeeklyUtilization > 0 {
		return w.WeeklyUtilization * 100.0
	}

	// Fallback to estimated calculation
	if w.TotalTokens == 0 {
		return 0.0
	}

	limit := GetWeeklyLimit(w.SubscriptionType, w.RateLimitTier)
	if limit == 0 {
		return 0.0
	}

	return float64(w.TotalTokens) / float64(limit) * 100.0
}

// GetFiveHourPercentage returns the 5-hour window usage percentage (0-99).
func (w *WeeklyStats) GetFiveHourPercentage() int {
	if w == nil || !w.HasAPIData {
		return 0
	}
	percentage := int(w.FiveHourUtilization * 100)
	if percentage > 99 {
		percentage = 99
	}
	if percentage < 0 {
		percentage = 0
	}
	return percentage
}

// IsThrottled returns true if currently rate limited.
func (w *WeeklyStats) IsThrottled() bool {
	if w == nil {
		return false
	}
	return w.RateLimitStatus == "throttled"
}

// IsLimitedByFiveHour returns true if the 5-hour window is the limiting factor.
func (w *WeeklyStats) IsLimitedByFiveHour() bool {
	if w == nil {
		return false
	}
	return w.RepresentativeClaim == "five_hour"
}
