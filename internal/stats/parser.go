package stats

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ParseStatsCache reads and parses Claude's stats-cache.json file.
func ParseStatsCache(path string) (*StatsCache, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read stats file: %w", err)
	}

	var cache StatsCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("failed to parse stats file: %w", err)
	}

	return &cache, nil
}

// ParseCredentials reads and parses Claude's credentials file.
func ParseCredentials(path string) (*Credentials, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse credentials file: %w", err)
	}

	return &creds, nil
}

// FileExists checks if a file exists at the given path.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// UpdateRefreshToken updates only the refresh token in the credentials file.
// It reads the existing file, updates the refresh token, and writes it back atomically.
// This preserves all other fields including mcpOAuth which we don't parse in our structs.
func UpdateRefreshToken(path string, newRefreshToken string) error {
	// Read the existing file as raw JSON to preserve all fields
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read credentials file: %w", err)
	}

	// Parse into a generic map to preserve all fields (including mcpOAuth)
	var rawCreds map[string]interface{}
	if err := json.Unmarshal(data, &rawCreds); err != nil {
		return fmt.Errorf("failed to parse credentials file: %w", err)
	}

	// Get the claudeAiOauth section
	claudeOAuth, ok := rawCreds["claudeAiOauth"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("claudeAiOauth section not found or invalid")
	}

	// Update only the refresh token
	claudeOAuth["refreshToken"] = newRefreshToken

	// Marshal back to JSON with indentation for readability
	updatedData, err := json.Marshal(rawCreds)
	if err != nil {
		return fmt.Errorf("failed to marshal updated credentials: %w", err)
	}

	// Get original file permissions
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat credentials file: %w", err)
	}
	perm := fileInfo.Mode().Perm()

	// Write to a temp file first (atomic write pattern)
	dir := filepath.Dir(path)
	tempFile, err := os.CreateTemp(dir, ".credentials-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()

	// Write the data
	if _, err := tempFile.Write(updatedData); err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Close before rename
	if err := tempFile.Close(); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Set correct permissions on temp file
	if err := os.Chmod(tempPath, perm); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to set permissions on temp file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempPath, path); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to rename temp file to credentials file: %w", err)
	}

	return nil
}
