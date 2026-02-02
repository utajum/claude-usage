// Package config provides configuration management for claude-usage.
package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetHomeDir returns the user's home directory.
// Works on Linux, macOS, and Windows.
func GetHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback for different OSes
		if runtime.GOOS == "windows" {
			return os.Getenv("USERPROFILE")
		}
		return os.Getenv("HOME")
	}
	return home
}

// GetClaudeDir returns the Claude configuration directory.
// - Linux/macOS: ~/.claude
// - Windows: %USERPROFILE%\.claude
func GetClaudeDir() string {
	return filepath.Join(GetHomeDir(), ".claude")
}

// GetClaudeStatsPath returns the path to Claude's stats-cache.json.
// - Linux/macOS: ~/.claude/stats-cache.json
// - Windows: %USERPROFILE%\.claude\stats-cache.json
func GetClaudeStatsPath() string {
	return filepath.Join(GetClaudeDir(), "stats-cache.json")
}

// GetClaudeCredentialsPath returns the path to Claude's credentials file.
// - Linux/macOS: ~/.claude/.credentials.json
// - Windows: %USERPROFILE%\.claude\.credentials.json
func GetClaudeCredentialsPath() string {
	return filepath.Join(GetClaudeDir(), ".credentials.json")
}

// GetOpenCodeCredentialsPath returns the path to OpenCode's auth file.
// - Linux: ~/.local/share/opencode/auth.json
// Note: OpenCode is only supported on Linux.
func GetOpenCodeCredentialsPath() string {
	return filepath.Join(GetHomeDir(), ".local", "share", "opencode", "auth.json")
}

// GetConfigDir returns the app's config directory.
// - Linux: ~/.config/claude-usage (respects XDG_CONFIG_HOME)
// - macOS: ~/Library/Application Support/claude-usage
// - Windows: %APPDATA%\claude-usage
func GetConfigDir() string {
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "claude-usage")
		}
		return filepath.Join(GetHomeDir(), "AppData", "Roaming", "claude-usage")

	case "darwin":
		return filepath.Join(GetHomeDir(), "Library", "Application Support", "claude-usage")

	default: // Linux and other Unix-like systems
		configHome := os.Getenv("XDG_CONFIG_HOME")
		if configHome == "" {
			configHome = filepath.Join(GetHomeDir(), ".config")
		}
		return filepath.Join(configHome, "claude-usage")
	}
}

// GetConfigPath returns the path to the app's config file.
func GetConfigPath() string {
	return filepath.Join(GetConfigDir(), "config.json")
}

// ExpandPath expands ~ to the user's home directory.
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(GetHomeDir(), path[2:])
	}
	if path == "~" {
		return GetHomeDir()
	}
	return path
}

// EnsureConfigDir creates the config directory if it doesn't exist.
func EnsureConfigDir() error {
	return os.MkdirAll(GetConfigDir(), 0755)
}

// GetOS returns a human-readable OS name.
func GetOS() string {
	switch runtime.GOOS {
	case "windows":
		return "Windows"
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	default:
		return runtime.GOOS
	}
}
