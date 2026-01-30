package icon

import (
	"github.com/user/claude-usage/internal/stats"
)

// Generator creates icons based on usage statistics.
type Generator struct {
	Size int
}

// DefaultGenerator returns a generator with the default icon size.
func DefaultGenerator() *Generator {
	return &Generator{Size: IconSize}
}

// GenerateWithPercentage creates an icon with percentage text overlay.
func (g *Generator) GenerateWithPercentage(weeklyStats *stats.WeeklyStats, percentage int) ([]byte, error) {
	c := ColorGray
	if weeklyStats != nil {
		c = GetColorForTokens(weeklyStats.TotalTokens)
	}

	if percentage < 0 {
		percentage = 0
	}
	if percentage > 99 {
		percentage = 99
	}

	return RenderNeonOrbWithText(c, g.Size, percentage)
}

// GenerateError creates an icon indicating an error state.
func (g *Generator) GenerateError() ([]byte, error) {
	return RenderNeonOrbWithText(ColorNeonPurple, g.Size, 0)
}
