// Package update provides self-update functionality for the application.
package update

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"claude-usage/internal/config"
)

// Result represents the outcome of an update operation.
type Result struct {
	Success      bool
	Message      string
	NeedsRestart bool
}

// CleanupOldBinary removes any leftover .old file from a previous Windows update.
// This should be called on application startup.
func CleanupOldBinary() {
	if runtime.GOOS != "windows" {
		return
	}

	exePath, err := os.Executable()
	if err != nil {
		return
	}

	oldPath := exePath + ".old"
	if _, err := os.Stat(oldPath); err == nil {
		// .old file exists, try to remove it
		if err := os.Remove(oldPath); err == nil {
			fmt.Printf("Cleaned up old binary: %s\n", oldPath)
		}
	}
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

	// Download the update (ZIP for Windows, binary for others)
	downloadURL := GetDownloadURL()
	tmpFile, err := downloadBinary(downloadURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download update: %w", err)
	}
	defer os.Remove(tmpFile) // Clean up temp file on error

	// For Windows, extract the exe from the ZIP
	var newBinaryPath string
	if runtime.GOOS == "windows" {
		extractedPath, err := extractWindowsZip(tmpFile)
		if err != nil {
			return nil, fmt.Errorf("failed to extract update: %w", err)
		}
		defer os.Remove(extractedPath) // Clean up extracted file on error
		newBinaryPath = extractedPath
	} else {
		newBinaryPath = tmpFile
	}

	// Verify the download is a valid executable (basic check: has content)
	info, err := os.Stat(newBinaryPath)
	if err != nil || info.Size() < 1000 {
		return nil, fmt.Errorf("downloaded file appears invalid (size: %d)", info.Size())
	}

	// On Windows, we can't replace a running executable directly.
	// We rename the current exe to .old, copy new one in, then delete .old on next run.
	if runtime.GOOS == "windows" {
		return updateWindows(exePath, newBinaryPath)
	}

	// Linux/macOS: Backup and replace
	backupPath := exePath + ".backup"
	if err := os.Rename(exePath, backupPath); err != nil {
		return nil, fmt.Errorf("failed to backup current binary: %w", err)
	}

	// Move new binary into place
	if err := os.Rename(newBinaryPath, exePath); err != nil {
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

// extractWindowsZip extracts the .exe file from a ZIP archive.
// Returns the path to the extracted executable.
func extractWindowsZip(zipPath string) (string, error) {
	// Open the ZIP file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", fmt.Errorf("failed to open zip: %w", err)
	}
	defer reader.Close()

	// Find the .exe file in the ZIP
	var exeFile *zip.File
	for _, f := range reader.File {
		if strings.HasSuffix(strings.ToLower(f.Name), ".exe") {
			exeFile = f
			break
		}
	}

	if exeFile == nil {
		return "", fmt.Errorf("no .exe file found in zip")
	}

	// Open the file inside the ZIP
	rc, err := exeFile.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open exe in zip: %w", err)
	}
	defer rc.Close()

	// Create a temp file for the extracted exe
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	tmpFile, err := os.CreateTemp(exeDir, "claude-usage-extracted-*.exe")
	if err != nil {
		// Fallback to system temp
		tmpFile, err = os.CreateTemp("", "claude-usage-extracted-*.exe")
		if err != nil {
			return "", fmt.Errorf("failed to create temp file: %w", err)
		}
	}
	tmpPath := tmpFile.Name()

	// Extract the exe
	_, err = io.Copy(tmpFile, rc)
	tmpFile.Close()
	if err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to extract exe: %w", err)
	}

	return tmpPath, nil
}

// updateWindows handles the Windows-specific update process.
// Windows can't replace a running executable, so we:
// 1. Rename current exe to .old
// 2. Copy new exe to original location
// 3. The .old file will be cleaned up on next restart
func updateWindows(exePath, newBinaryPath string) (*Result, error) {
	oldPath := exePath + ".old"

	// Remove any leftover .old file from previous update
	os.Remove(oldPath)

	// Rename current running exe to .old
	// On Windows, you CAN rename a running executable, just can't delete it
	if err := os.Rename(exePath, oldPath); err != nil {
		return nil, fmt.Errorf("failed to rename current executable: %w", err)
	}

	// Copy new binary to the original location
	// We use copy instead of rename because the temp file might be on a different drive
	if err := copyFile(newBinaryPath, exePath); err != nil {
		// Try to restore on failure
		os.Rename(oldPath, exePath)
		return nil, fmt.Errorf("failed to install update: %w", err)
	}

	return &Result{
		Success:      true,
		Message:      "Update installed successfully. Please restart the application.",
		NeedsRestart: true,
	}, nil
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
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
