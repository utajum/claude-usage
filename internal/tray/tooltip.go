// Package tray provides system tray integration for claude-usage.
package tray

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"claude-usage/internal/stats"
	"claude-usage/pkg/format"
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

// menuMaxLines is the maximum number of lines FormatMenuLines can emit:
// header + plan + optional "STATUS: THROTTLED" + 5-hour bar + weekly bar.
const menuMaxLines = 5

// FormatMenuLines returns the usage summary as individual lines for display in
// the tray menu (macOS only). It emits at most menuMaxLines lines — header,
// plan, optional throttled status, and the 5-hour and weekly block-bar windows
// (with the ◀ limiting marker) — using the same helpers and formatting as
// FormatTooltip. The per-model Opus/Sonnet windows are intentionally omitted so
// the menu stays short enough to avoid the macOS menu scroll indicator; those
// remain visible in the Linux hover tooltip via FormatTooltip.
func FormatMenuLines(weeklyStats *stats.WeeklyStats) []string {
	if weeklyStats == nil {
		return []string{"CLAUDE USAGE", "No data available"}
	}

	lines := make([]string, 0, menuMaxLines)

	// Header
	lines = append(lines, "CLAUDE USAGE")

	// Plan info
	if weeklyStats.SubscriptionType != "" {
		planName := format.FormatPlanName(weeklyStats.SubscriptionType, weeklyStats.RateLimitTier)
		lines = append(lines, fmt.Sprintf("Plan: %s", planName))
	}

	// Status (throttled warning)
	if weeklyStats.IsThrottled() {
		lines = append(lines, "STATUS: THROTTLED")
	}

	if weeklyStats.HasAPIData {
		// 5-hour window
		fiveHourPct := weeklyStats.GetFiveHourPercentage()
		fiveHourBar := makeProgressBar(fiveHourPct, 10)
		fiveHourReset := formatShortDuration(time.Until(weeklyStats.FiveHourReset))
		marker := ""
		if weeklyStats.IsLimitedByFiveHour() {
			marker = " ◀"
		}
		lines = append(lines, fmt.Sprintf("%s %3d%% %s%s", fiveHourBar, fiveHourPct, fiveHourReset, marker))

		// Weekly window
		weeklyPct := weeklyStats.GetPercentage()
		weeklyBar := makeProgressBar(weeklyPct, 10)
		weeklyReset := formatShortDuration(time.Until(weeklyStats.WeeklyReset))
		marker = ""
		if !weeklyStats.IsLimitedByFiveHour() {
			marker = " ◀"
		}
		lines = append(lines, fmt.Sprintf("%s %3d%% %s%s", weeklyBar, weeklyPct, weeklyReset, marker))
	} else {
		// Estimated usage based on token counts
		weeklyPct := weeklyStats.GetPercentage()
		weeklyBar := makeProgressBar(weeklyPct, 10)
		daysRemaining := stats.GetDaysRemainingInWeek()
		lines = append(lines, fmt.Sprintf("%s ~%3d%% %dd", weeklyBar, weeklyPct, daysRemaining))
	}

	return lines
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
