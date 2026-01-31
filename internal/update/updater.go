// Package update provides self-update functionality for the application.
package update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"claude-usage/internal/config"
)

// Result represents the outcome of an update operation.
type Result struct {
	Success      bool
	Message      string
	NeedsRestart bool
}

// GetPlatformBinaryName returns the appropriate binary name for the current platform.
func GetPlatformBinaryName() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	switch goos {
	case "windows":
		return "claude-usage-windows-amd64.zip"
	case "darwin":
		// macOS uses universal binary names
		return fmt.Sprintf("claude-usage-darwin-%s", goarch)
	case "linux":
		return fmt.Sprintf("claude-usage-linux-%s", goarch)
	default:
		return fmt.Sprintf("claude-usage-%s-%s", goos, goarch)
	}
}

// GetDownloadURL returns the full URL for downloading the latest binary.
func GetDownloadURL() string {
	baseURL := config.GetUpdateURL()
	binaryName := GetPlatformBinaryName()
	return fmt.Sprintf("%s/%s", baseURL, binaryName)
}

// Update downloads the latest version and replaces the current binary.
// Returns a Result indicating success/failure and whether restart is needed.
func Update() (*Result, error) {
	// Get current executable path
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	// Resolve symlinks to get the real path
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve executable path: %w", err)
	}

	// Windows uses ZIP files which need different handling
	if runtime.GOOS == "windows" {
		return nil, fmt.Errorf("Windows auto-update not yet supported. Please download manually from GitHub releases")
	}

	// Download the new binary
	downloadURL := GetDownloadURL()
	tmpFile, err := downloadBinary(downloadURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download update: %w", err)
	}
	defer os.Remove(tmpFile) // Clean up temp file on error

	// Verify the download is a valid executable (basic check: has content)
	info, err := os.Stat(tmpFile)
	if err != nil || info.Size() < 1000 {
		return nil, fmt.Errorf("downloaded file appears invalid (size: %d)", info.Size())
	}

	// Backup the current binary
	backupPath := exePath + ".backup"
	if err := os.Rename(exePath, backupPath); err != nil {
		return nil, fmt.Errorf("failed to backup current binary: %w", err)
	}

	// Move new binary into place
	if err := os.Rename(tmpFile, exePath); err != nil {
		// Try to restore backup on failure
		os.Rename(backupPath, exePath)
		return nil, fmt.Errorf("failed to install update: %w", err)
	}

	// Make executable
	if err := os.Chmod(exePath, 0755); err != nil {
		// Non-fatal, try to continue
		fmt.Printf("Warning: could not set executable permission: %v\n", err)
	}

	// Remove backup (optional, ignore errors)
	os.Remove(backupPath)

	return &Result{
		Success:      true,
		Message:      "Update installed successfully. Please restart the application.",
		NeedsRestart: true,
	}, nil
}

// downloadBinary downloads the binary from the given URL to a temporary file.
// Returns the path to the temporary file.
func downloadBinary(url string) (string, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Minute,
		// Follow redirects (GitHub releases use redirects)
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	// Make request
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Create temporary file in the same directory as the executable
	// This ensures we can rename it later (same filesystem)
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)

	tmpFile, err := os.CreateTemp(exeDir, "claude-usage-update-*")
	if err != nil {
		// Fallback to system temp dir
		tmpFile, err = os.CreateTemp("", "claude-usage-update-*")
		if err != nil {
			return "", fmt.Errorf("failed to create temp file: %w", err)
		}
	}
	tmpPath := tmpFile.Name()

	// Download to temp file
	_, err = io.Copy(tmpFile, resp.Body)
	tmpFile.Close()
	if err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to download file: %w", err)
	}

	return tmpPath, nil
}
