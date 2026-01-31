package stats

import (
	"encoding/json"
	"os"
	"testing"
)

func TestUpdateRefreshToken(t *testing.T) {
	// Create a temp credentials file with the same structure as real one
	tempFile, err := os.CreateTemp("", "credentials-test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Initial credentials (mimics real structure)
	initialCreds := map[string]interface{}{
		"claudeAiOauth": map[string]interface{}{
			"accessToken":      "old-access-token",
			"refreshToken":     "old-refresh-token",
			"expiresAt":        1769876156710,
			"scopes":           []string{"user:inference", "user:profile"},
			"subscriptionType": "team",
			"rateLimitTier":    "default_claude_max_5x",
		},
		"mcpOAuth": map[string]interface{}{
			"plugin:atlassian": map[string]interface{}{
				"serverName": "plugin:atlassian:atlassian",
				"serverUrl":  "https://mcp.atlassian.com/v1/mcp",
				"clientId":   "testClientId",
			},
		},
	}

	data, _ := json.Marshal(initialCreds)
	if _, err := tempFile.Write(data); err != nil {
		t.Fatalf("Failed to write initial data: %v", err)
	}
	tempFile.Close()

	// Update the refresh token
	newRefreshToken := "NEW-REFRESH-TOKEN-12345"
	if err := UpdateRefreshToken(tempFile.Name(), newRefreshToken); err != nil {
		t.Fatalf("UpdateRefreshToken failed: %v", err)
	}

	// Read and verify the update
	updatedData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read updated file: %v", err)
	}

	var updatedCreds map[string]interface{}
	if err := json.Unmarshal(updatedData, &updatedCreds); err != nil {
		t.Fatalf("Failed to parse updated data: %v", err)
	}

	claudeOAuth := updatedCreds["claudeAiOauth"].(map[string]interface{})

	// Verify refresh token updated
	if claudeOAuth["refreshToken"] != newRefreshToken {
		t.Errorf("refreshToken not updated! Got: %s, expected: %s",
			claudeOAuth["refreshToken"], newRefreshToken)
	}

	// Verify access token preserved
	if claudeOAuth["accessToken"] != "old-access-token" {
		t.Error("accessToken was changed when it should be preserved")
	}

	// Verify expiresAt preserved
	if claudeOAuth["expiresAt"].(float64) != 1769876156710 {
		t.Error("expiresAt was changed when it should be preserved")
	}

	// Verify subscriptionType preserved
	if claudeOAuth["subscriptionType"] != "team" {
		t.Error("subscriptionType was changed when it should be preserved")
	}

	// Verify mcpOAuth preserved
	if _, ok := updatedCreds["mcpOAuth"]; !ok {
		t.Error("mcpOAuth section was lost")
	}

	t.Log("All checks passed:")
	t.Log("- refreshToken updated correctly")
	t.Log("- accessToken preserved")
	t.Log("- expiresAt preserved")
	t.Log("- subscriptionType preserved")
	t.Log("- mcpOAuth section preserved")
}

func TestUpdateRefreshToken_NonExistentFile(t *testing.T) {
	err := UpdateRefreshToken("/nonexistent/path/credentials.json", "new-token")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestUpdateRefreshToken_InvalidJSON(t *testing.T) {
	tempFile, err := os.CreateTemp("", "credentials-test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write invalid JSON
	tempFile.WriteString("not valid json{")
	tempFile.Close()

	err = UpdateRefreshToken(tempFile.Name(), "new-token")
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestUpdateRefreshToken_MissingClaudeAiOauth(t *testing.T) {
	tempFile, err := os.CreateTemp("", "credentials-test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write JSON without claudeAiOauth section
	data, _ := json.Marshal(map[string]interface{}{
		"someOtherField": "value",
	})
	tempFile.Write(data)
	tempFile.Close()

	err = UpdateRefreshToken(tempFile.Name(), "new-token")
	if err == nil {
		t.Error("Expected error for missing claudeAiOauth section, got nil")
	}
}
