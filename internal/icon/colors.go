// Package icon provides dynamic icon generation for the system tray.
package icon

import "image/color"

// Color palette
var (
	ColorNeonGreen  = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	ColorNeonYellow = color.RGBA{R: 204, G: 255, B: 0, A: 255}
	ColorNeonOrange = color.RGBA{R: 255, G: 153, B: 0, A: 255}
	ColorNeonRed    = color.RGBA{R: 255, G: 51, B: 102, A: 255}
	ColorNeonPurple = color.RGBA{R: 191, G: 0, B: 255, A: 255}
	ColorGray       = color.RGBA{R: 128, G: 128, B: 128, A: 255}
)

// Token thresholds for color changes
const (
	ThresholdLow     = 500_000
	ThresholdMedium  = 2_000_000
	ThresholdHigh    = 5_000_000
	ThresholdExtreme = 10_000_000
)

// GetColorForTokens returns the appropriate color based on token count.
func GetColorForTokens(tokens int64) color.RGBA {
	switch {
	case tokens < ThresholdLow:
		return ColorNeonGreen
	case tokens < ThresholdMedium:
		return ColorNeonYellow
	case tokens < ThresholdHigh:
		return ColorNeonOrange
	case tokens < ThresholdExtreme:
		return ColorNeonRed
	default:
		return ColorNeonPurple
	}
}
