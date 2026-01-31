package config

import (
	"encoding/json"
	"os"
	"time"
)

// DefaultWeeklyBudget is the default weekly token budget (5 million tokens).
const DefaultWeeklyBudget int64 = 5_000_000

// Config holds the application configuration.
type Config struct {
	// RefreshInterval is how often to refresh stats.
	RefreshInterval time.Duration `json:"-"`

	// RefreshIntervalSeconds is the JSON-serializable version.
	RefreshIntervalSeconds int `json:"refresh_interval_seconds"`

	// WeeklyBudgetTokens is the user's weekly token budget for percentage calculation.
	// The icon will show percentage = (used / budget) * 100.
	// Default is 5 million tokens.
	WeeklyBudgetTokens int64 `json:"weekly_budget_tokens"`

	// ClaudeStatsPath is the path to Claude's stats-cache.json.
	// If empty, uses the default path.
	ClaudeStatsPath string `json:"claude_stats_path,omitempty"`

	// ClaudeCredentialsPath is the path to Claude's credentials file.
	// If empty, uses the default path.
	ClaudeCredentialsPath string `json:"claude_credentials_path,omitempty"`
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	return &Config{
		RefreshInterval:        5 * time.Minute,
		RefreshIntervalSeconds: 300,
		WeeklyBudgetTokens:     DefaultWeeklyBudget,
		ClaudeStatsPath:        "",
		ClaudeCredentialsPath:  "",
	}
}

// Load reads configuration from the config file.
// If the file doesn't exist, returns defaults.
func Load() (*Config, error) {
	cfg := Default()

	configPath := GetConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// Convert seconds to duration
	cfg.RefreshInterval = time.Duration(cfg.RefreshIntervalSeconds) * time.Second

	// Expand paths
	if cfg.ClaudeStatsPath != "" {
		cfg.ClaudeStatsPath = ExpandPath(cfg.ClaudeStatsPath)
	}
	if cfg.ClaudeCredentialsPath != "" {
		cfg.ClaudeCredentialsPath = ExpandPath(cfg.ClaudeCredentialsPath)
	}

	return cfg, nil
}

// Save writes the configuration to the config file.
func (c *Config) Save() error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	// Update seconds from duration
	c.RefreshIntervalSeconds = int(c.RefreshInterval.Seconds())

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(GetConfigPath(), data, 0644)
}

// GetStatsPath returns the effective stats path (config or default).
func (c *Config) GetStatsPath() string {
	if c.ClaudeStatsPath != "" {
		return c.ClaudeStatsPath
	}
	return GetClaudeStatsPath()
}

// GetCredentialsPath returns the effective credentials path (config or default).
func (c *Config) GetCredentialsPath() string {
	if c.ClaudeCredentialsPath != "" {
		return c.ClaudeCredentialsPath
	}
	return GetClaudeCredentialsPath()
}
