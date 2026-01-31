// Package tray provides system tray integration for claude-usage.
package tray

import (
	"fmt"
	"runtime"
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

// FormatTooltipCompact creates a condensed tooltip for Windows (127 char limit).
// Shows only 5-hour and weekly bars, skips Opus/Sonnet to fit within Windows tooltip limit.
func FormatTooltipCompact(weeklyStats *stats.WeeklyStats) string {
	if weeklyStats == nil {
		return "Claude Usage\nNo data"
	}

	var sb strings.Builder

	// Header with plan inline to save space
	sb.WriteString("CLAUDE USAGE")
	if weeklyStats.SubscriptionType != "" {
		planName := format.FormatPlanName(weeklyStats.SubscriptionType, weeklyStats.RateLimitTier)
		sb.WriteString(fmt.Sprintf(" %s", planName))
	}
	sb.WriteString("\n")

	// Status (throttled warning)
	if weeklyStats.IsThrottled() {
		sb.WriteString("THROTTLED\n")
	}

	// Rate Limit Section - only 5-hour and weekly (skip Opus/Sonnet)
	if weeklyStats.HasAPIData {
		// 5-hour window - shorter bar (6 chars) and shorter time format
		fiveHourPct := weeklyStats.GetFiveHourPercentage()
		fiveHourBar := makeProgressBar(fiveHourPct, 6)
		fiveHourReset := formatVeryShortDuration(time.Until(weeklyStats.FiveHourReset))
		marker := ""
		if weeklyStats.IsLimitedByFiveHour() {
			marker = " ◀"
		}
		sb.WriteString(fmt.Sprintf("%s %3d%% %s%s\n", fiveHourBar, fiveHourPct, fiveHourReset, marker))

		// Weekly window - shorter bar (6 chars) and shorter time format
		weeklyPct := weeklyStats.GetPercentage()
		weeklyBar := makeProgressBar(weeklyPct, 6)
		weeklyReset := formatVeryShortDuration(time.Until(weeklyStats.WeeklyReset))
		marker = ""
		if !weeklyStats.IsLimitedByFiveHour() {
			marker = " ◀"
		}
		sb.WriteString(fmt.Sprintf("%s %3d%% %s%s", weeklyBar, weeklyPct, weeklyReset, marker))
	} else {
		// Show estimated usage based on token counts
		weeklyPct := weeklyStats.GetPercentage()
		weeklyBar := makeProgressBar(weeklyPct, 6)
		daysRemaining := stats.GetDaysRemainingInWeek()
		resetStr := fmt.Sprintf("%dd", daysRemaining)
		sb.WriteString(fmt.Sprintf("%s ~%3d%% %s", weeklyBar, weeklyPct, resetStr))
	}

	return sb.String()
}

// formatVeryShortDuration formats duration in ultra-compact format for Windows.
// Examples: "2h", "5d", "1d2h"
func formatVeryShortDuration(d time.Duration) string {
	if d < 0 {
		return "0h"
	}

	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24

	if days > 0 {
		if hours > 0 {
			return fmt.Sprintf("%dd%dh", days, hours)
		}
		return fmt.Sprintf("%dd", days)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%dm", minutes)
}

// FormatTooltipForPlatform returns the appropriate tooltip format based on OS.
// Windows gets a compact version (127 char limit), other platforms get full version.
func FormatTooltipForPlatform(weeklyStats *stats.WeeklyStats) string {
	if runtime.GOOS == "windows" {
		return FormatTooltipCompact(weeklyStats)
	}
	return FormatTooltip(weeklyStats)
}
