// Claude Usage - A system tray app for monitoring Claude Code token usage.
//
// This application displays your weekly Claude Code token usage as a glowing
// neon orb in the system tray. Hover over it to see detailed statistics
// including model breakdown and your subscription plan.
//
// Supported platforms: Linux, Windows, macOS
package main

import (
	"log"
	"os"

	"claude-usage/internal/app"
	"claude-usage/internal/config"
)

// Version is set at build time via -ldflags
var Version = "dev"

func main() {
	// Setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Claude Usage %s starting on %s", Version, config.GetOS())
	log.Printf("Claude data path: %s", config.GetClaudeDir())

	// Check if Claude directory exists
	claudeDir := config.GetClaudeDir()
	if _, err := os.Stat(claudeDir); os.IsNotExist(err) {
		log.Printf("Warning: Claude directory not found at %s", claudeDir)
		log.Println("Make sure Claude Code is installed and has been used at least once.")
	}

	// Create and run the app
	application, err := app.New()
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Run blocks until quit
	application.Run()

	log.Println("Claude Usage exiting")
}
