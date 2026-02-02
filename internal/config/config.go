package config

import (
	"encoding/json"
	"os"
	"runtime"
	"time"
)

// DefaultWeeklyBudget is the default weekly token budget (5 million tokens).
const DefaultWeeklyBudget int64 = 5_000_000

// Source constants for credential sources.
const (
	SourceClaude   = "claude"
	SourceOpenCode = "opencode"
)

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

	// Source is the credential source: "claude" or "opencode".
	// OpenCode is only supported on Linux.
	// If empty, auto-detects based on available credential files.
	Source string `json:"source,omitempty"`
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	return &Config{
		RefreshInterval:        5 * time.Minute,
		RefreshIntervalSeconds: 300,
		WeeklyBudgetTokens:     DefaultWeeklyBudget,
		ClaudeStatsPath:        "",
		ClaudeCredentialsPath:  "",
		Source:                 detectDefaultSource(),
	}
}

// detectDefaultSource determines the default credential source based on available files.
// On Linux, if OpenCode credentials exist but Claude credentials don't, use OpenCode.
// Otherwise, default to Claude.
func detectDefaultSource() string {
	// OpenCode is only supported on Linux
	if runtime.GOOS != "linux" {
		return SourceClaude
	}

	claudeExists := fileExists(GetClaudeCredentialsPath())
	openCodeExists := fileExists(GetOpenCodeCredentialsPath())

	// If only OpenCode exists, use it
	if openCodeExists && !claudeExists {
		return SourceOpenCode
	}

	// Default to Claude (including when both exist)
	return SourceClaude
}

// fileExists checks if a file exists at the given path.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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

	// If source is empty (old config file), auto-detect
	if cfg.Source == "" {
		cfg.Source = detectDefaultSource()
	}

	// OpenCode is only supported on Linux, fallback to Claude on other platforms
	if cfg.Source == SourceOpenCode && runtime.GOOS != "linux" {
		cfg.Source = SourceClaude
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
// If Source is "opencode" and we're on Linux, returns OpenCode path.
func (c *Config) GetCredentialsPath() string {
	if c.ClaudeCredentialsPath != "" {
		return c.ClaudeCredentialsPath
	}
	if c.Source == SourceOpenCode && runtime.GOOS == "linux" {
		return GetOpenCodeCredentialsPath()
	}
	return GetClaudeCredentialsPath()
}

// IsOpenCode returns true if the current source is OpenCode.
func (c *Config) IsOpenCode() bool {
	return c.Source == SourceOpenCode
}

// ToggleSource switches between Claude and OpenCode sources.
// Only works on Linux; no-op on other platforms.
func (c *Config) ToggleSource() {
	if runtime.GOOS != "linux" {
		return
	}
	if c.Source == SourceOpenCode {
		c.Source = SourceClaude
	} else {
		c.Source = SourceOpenCode
	}
}

// GetSourceDisplayName returns a human-readable name for the current source.
func (c *Config) GetSourceDisplayName() string {
	if c.Source == SourceOpenCode {
		return "OpenCode"
	}
	return "Claude Code"
}
