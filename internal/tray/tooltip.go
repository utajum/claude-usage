// Package tray provides system tray integration for claude-usage.
package tray

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/claude-usage/internal/stats"
	"github.com/user/claude-usage/pkg/format"
)

// FormatTooltip creates a formatted tooltip string from weekly statistics.
func FormatTooltip(weeklyStats *stats.WeeklyStats) string {
	if weeklyStats == nil {
		return "Claude Usage\nNo data available"
	}

	var sb strings.Builder

	// Header
	sb.WriteString("CLAUDE USAGE\n")

	// Plan info
	if weeklyStats.SubscriptionType != "" {
		planName := format.FormatPlanName(weeklyStats.SubscriptionType, weeklyStats.RateLimitTier)
		sb.WriteString(fmt.Sprintf("Plan: %s\n", planName))
	}

	// Status (throttled warning)
	if weeklyStats.IsThrottled() {
		sb.WriteString("STATUS: THROTTLED\n")
	}

	// Rate Limit Section
	if weeklyStats.HasAPIData {
		// 5-hour window
		fiveHourPct := weeklyStats.GetFiveHourPercentage()
		fiveHourBar := makeProgressBar(fiveHourPct, 10)
		fiveHourReset := formatShortDuration(time.Until(weeklyStats.FiveHourReset))
		marker := ""
		if weeklyStats.IsLimitedByFiveHour() {
			marker = " ◀"
		}
		sb.WriteString(fmt.Sprintf("%s %3d%% %s%s\n", fiveHourBar, fiveHourPct, fiveHourReset, marker))

		// Weekly window
		weeklyPct := weeklyStats.GetPercentage()
		weeklyBar := makeProgressBar(weeklyPct, 10)
		weeklyReset := formatShortDuration(time.Until(weeklyStats.WeeklyReset))
		marker = ""
		if !weeklyStats.IsLimitedByFiveHour() {
			marker = " ◀"
		}
		sb.WriteString(fmt.Sprintf("%s %3d%% %s%s\n", weeklyBar, weeklyPct, weeklyReset, marker))

		// Show model-specific limits if available
		if weeklyStats.OpusUtilization > 0 {
			opusPct := int(weeklyStats.OpusUtilization * 100)
			opusBar := makeProgressBar(opusPct, 10)
			opusReset := formatShortDuration(time.Until(weeklyStats.OpusReset))
			sb.WriteString(fmt.Sprintf("%s %3d%% %s\n", opusBar, opusPct, opusReset))
		}
		if weeklyStats.SonnetUtilization > 0 {
			sonnetPct := int(weeklyStats.SonnetUtilization * 100)
			sonnetBar := makeProgressBar(sonnetPct, 10)
			sonnetReset := formatShortDuration(time.Until(weeklyStats.SonnetReset))
			sb.WriteString(fmt.Sprintf("%s %3d%% %s\n", sonnetBar, sonnetPct, sonnetReset))
		}
	} else {
		// Show estimated usage based on token counts
		weeklyPct := weeklyStats.GetPercentage()
		weeklyBar := makeProgressBar(weeklyPct, 10)
		daysRemaining := stats.GetDaysRemainingInWeek()
		resetStr := fmt.Sprintf("%dd", daysRemaining)
		sb.WriteString(fmt.Sprintf("%s ~%3d%% %s\n", weeklyBar, weeklyPct, resetStr))
	}

	return strings.TrimRight(sb.String(), "\n")
}

// formatShortDuration formats a duration as compact "Xh Ym" or "Xd Yh" format.
func formatShortDuration(d time.Duration) string {
	if d < 0 {
		return "0h 0m"
	}

	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

// makeProgressBar creates a text-based progress bar using Unicode block characters.
func makeProgressBar(percentage int, width int) string {
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}

	filled := (percentage * width) / 100
	empty := width - filled

	return "▕" + strings.Repeat("█", filled) + strings.Repeat("░", empty) + "▏"
}
