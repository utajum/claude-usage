// Package api provides a client for fetching Claude API rate limit information.
package api

import "time"

// RateLimitData represents rate limit information from Anthropic API.
type RateLimitData struct {
	// FiveHourUtilization is the 5-hour window usage (0.0-1.0)
	FiveHourUtilization float64

	// WeeklyUtilization is the 7-day window usage (0.0-1.0)
	WeeklyUtilization float64

	// FiveHourReset is when the 5-hour window resets
	FiveHourReset time.Time

	// WeeklyReset is when the 7-day window resets
	WeeklyReset time.Time

	// RepresentativeClaim indicates which window is the limiting factor ("five_hour" or "seven_day")
	RepresentativeClaim string

	// Status is "allowed" or "throttled"
	Status string

	// OverageStatus is "allowed" or "rejected"
	OverageStatus string

	// Model-specific weekly utilization (for plans with per-model limits)
	OpusUtilization   float64
	SonnetUtilization float64
	OpusReset         time.Time
	SonnetReset       time.Time

	// FetchedAt is when this data was fetched
	FetchedAt time.Time
}

// GetWeeklyPercentage returns the weekly utilization as a percentage (0-100).
func (r *RateLimitData) GetWeeklyPercentage() int {
	if r == nil {
		return 0
	}
	pct := int(r.WeeklyUtilization * 100)
	if pct > 99 {
		pct = 99
	}
	if pct < 0 {
		pct = 0
	}
	return pct
}

// GetFiveHourPercentage returns the 5-hour utilization as a percentage (0-100).
func (r *RateLimitData) GetFiveHourPercentage() int {
	if r == nil {
		return 0
	}
	pct := int(r.FiveHourUtilization * 100)
	if pct > 99 {
		pct = 99
	}
	if pct < 0 {
		pct = 0
	}
	return pct
}

// TimeUntilFiveHourReset returns the duration until the 5-hour window resets.
func (r *RateLimitData) TimeUntilFiveHourReset() time.Duration {
	if r == nil {
		return 0
	}
	return time.Until(r.FiveHourReset)
}

// TimeUntilWeeklyReset returns the duration until the 7-day window resets.
func (r *RateLimitData) TimeUntilWeeklyReset() time.Duration {
	if r == nil {
		return 0
	}
	return time.Until(r.WeeklyReset)
}

// IsThrottled returns true if the user is currently rate limited.
func (r *RateLimitData) IsThrottled() bool {
	if r == nil {
		return false
	}
	return r.Status == "throttled"
}

// IsLimitedByFiveHour returns true if the 5-hour window is the limiting factor.
func (r *RateLimitData) IsLimitedByFiveHour() bool {
	if r == nil {
		return false
	}
	return r.RepresentativeClaim == "five_hour"
}

// IsLimitedByWeekly returns true if the weekly window is the limiting factor.
func (r *RateLimitData) IsLimitedByWeekly() bool {
	if r == nil {
		return false
	}
	return r.RepresentativeClaim == "seven_day"
}
