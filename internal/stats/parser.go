package stats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

// ParseOpenCodeCredentials reads OpenCode's auth.json and converts it to the common Credentials format.
func ParseOpenCodeCredentials(path string) (*Credentials, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenCode auth file: %w", err)
	}

	var openCodeCreds OpenCodeCredentials
	if err := json.Unmarshal(data, &openCodeCreds); err != nil {
		return nil, fmt.Errorf("failed to parse OpenCode auth file: %w", err)
	}

	// Convert OpenCode format to common Credentials format
	creds := &Credentials{
		ClaudeAiOauth: OAuthCredentials{
			AccessToken:  openCodeCreds.Anthropic.Access,
			RefreshToken: openCodeCreds.Anthropic.Refresh,
			ExpiresAt:    openCodeCreds.Anthropic.Expires,
			// OpenCode doesn't store subscription info, leave empty
			SubscriptionType: "",
			RateLimitTier:    "",
		},
	}

	return creds, nil
}

// keychainServiceNames are the service names Claude Code uses in macOS Keychain.
// "Claude Code-credentials" is the primary name used by /login.
// "Claude Code" is used by some versions for reading.
var keychainServiceNames = []string{
	"Claude Code-credentials",
	"Claude Code",
}

// ParseKeychainCredentials reads credentials from the macOS Keychain.
// Claude Code stores OAuth credentials in the Keychain under service name
// "Claude Code-credentials". This is used as a fallback when the JSON
// credentials file (~/.claude/.credentials.json) does not exist on macOS.
// Returns an error if not on macOS or if no keychain entry is found.
func ParseKeychainCredentials() (*Credentials, error) {
	if runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("keychain is only available on macOS")
	}

	var output []byte
	var lastErr error

	for _, service := range keychainServiceNames {
		cmd := exec.Command("security", "find-generic-password", "-s", service, "-w")
		out, err := cmd.Output()
		if err == nil {
			output = out
			log.Printf("Found credentials in macOS Keychain (service: %q)", service)
			break
		}
		lastErr = err
	}

	if output == nil {
		return nil, fmt.Errorf("no credentials found in macOS Keychain: %w", lastErr)
	}

	var creds Credentials
	if err := json.Unmarshal(bytes.TrimSpace(output), &creds); err != nil {
		return nil, fmt.Errorf("failed to parse keychain credentials: %w", err)
	}

	return &creds, nil
}

// FileExists checks if a file exists at the given path.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// UpdateKeychainRefreshToken updates the refresh token in the macOS Keychain entry.
// It reads the current keychain entry, updates the refreshToken field, and writes it back.
// The -U flag to security add-generic-password updates an existing entry.
func UpdateKeychainRefreshToken(newRefreshToken string) error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("keychain is only available on macOS")
	}

	// Read current keychain entry (try both service names)
	var output []byte
	var activeService string
	for _, service := range keychainServiceNames {
		cmd := exec.Command("security", "find-generic-password", "-s", service, "-w")
		out, err := cmd.Output()
		if err == nil {
			output = out
			activeService = service
			break
		}
	}
	if output == nil {
		return fmt.Errorf("no existing keychain entry found to update")
	}

	// Parse into a generic map to preserve all fields
	var rawCreds map[string]interface{}
	if err := json.Unmarshal(bytes.TrimSpace(output), &rawCreds); err != nil {
		return fmt.Errorf("failed to parse keychain JSON: %w", err)
	}

	// Get the claudeAiOauth section
	claudeOAuth, ok := rawCreds["claudeAiOauth"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("claudeAiOauth section not found in keychain data")
	}

	// Update only the refresh token
	claudeOAuth["refreshToken"] = newRefreshToken

	// Marshal back to JSON
	updatedJSON, err := json.Marshal(rawCreds)
	if err != nil {
		return fmt.Errorf("failed to marshal updated credentials: %w", err)
	}

	// Get current user for account field
	currentUser := os.Getenv("USER")
	if currentUser == "" {
		currentUser = os.Getenv("LOGNAME")
	}

	// Write back with -U (update existing entry)
	cmd := exec.Command("security", "add-generic-password",
		"-s", activeService,
		"-a", currentUser,
		"-w", string(updatedJSON),
		"-U")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to update keychain: %s: %w", string(out), err)
	}

	return nil
}

// UpdateOpenCodeRefreshToken updates the refresh token in OpenCode's auth.json file.
// It reads the existing file, updates the refresh token, and writes it back atomically.
func UpdateOpenCodeRefreshToken(path string, newRefreshToken string) error {
	// Read the existing file as raw JSON to preserve all fields
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read OpenCode auth file: %w", err)
	}

	// Parse into a generic map to preserve all fields
	var rawCreds map[string]interface{}
	if err := json.Unmarshal(data, &rawCreds); err != nil {
		return fmt.Errorf("failed to parse OpenCode auth file: %w", err)
	}

	// Get the anthropic section
	anthropic, ok := rawCreds["anthropic"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("anthropic section not found or invalid in OpenCode auth file")
	}

	// Update only the refresh token
	anthropic["refresh"] = newRefreshToken

	// Marshal back to JSON with indentation for readability
	updatedData, err := json.MarshalIndent(rawCreds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated OpenCode auth: %w", err)
	}

	// Get original file permissions
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat OpenCode auth file: %w", err)
	}
	perm := fileInfo.Mode().Perm()

	// Write to a temp file first (atomic write pattern)
	dir := filepath.Dir(path)
	tempFile, err := os.CreateTemp(dir, ".auth-*.tmp")
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
		return fmt.Errorf("failed to rename temp file to OpenCode auth file: %w", err)
	}

	return nil
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
