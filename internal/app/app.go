// Package app provides the main application orchestration for claude-usage.
package app

import (
	"log"
	"sync"
	"time"

	"claude-usage/internal/api"
	"claude-usage/internal/config"
	"claude-usage/internal/icon"
	"claude-usage/internal/stats"
	"claude-usage/internal/tray"
	"claude-usage/internal/update"
)

// App is the main application struct that coordinates all components.
type App struct {
	config    *config.Config
	version   string
	tray      *tray.Tray
	iconGen   *icon.Generator
	apiClient *api.Client
	stats     *stats.WeeklyStats
	statsMu   sync.RWMutex
	stopCh    chan struct{}
	refreshCh chan struct{}
}

// New creates a new App instance with the given version string.
func New(version string) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Warning: could not load config, using defaults: %v", err)
		cfg = config.Default()
	}

	return &App{
		config:    cfg,
		version:   version,
		tray:      tray.New(version),
		iconGen:   icon.DefaultGenerator(),
		apiClient: nil, // Will be initialized when we have a token
		stopCh:    make(chan struct{}),
		refreshCh: make(chan struct{}, 1),
	}, nil
}

// Run starts the application. This blocks until the app is quit.
func (a *App) Run() {
	// Set up tray callbacks
	a.tray.SetOnRefresh(func() {
		log.Println("Manual refresh triggered")
		a.triggerRefresh()
	})

	a.tray.SetOnUpdate(func() {
		log.Println("Update triggered")
		a.performUpdate()
	})

	a.tray.SetOnQuit(func() {
		log.Println("Quit triggered")
		close(a.stopCh)
	})

	// Run the tray (this will call onReady when initialized)
	a.tray.Run(a.onReady)
}

// onReady is called when the system tray is initialized and ready.
func (a *App) onReady() {
	log.Println("System tray ready")

	// Initial refresh
	a.refresh()

	// Start the refresh loop
	go a.refreshLoop()
}

// refreshLoop periodically refreshes the stats.
func (a *App) refreshLoop() {
	ticker := time.NewTicker(a.config.RefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-a.stopCh:
			log.Println("Refresh loop stopped")
			return
		case <-ticker.C:
			log.Println("Auto refresh triggered")
			a.refresh()
		case <-a.refreshCh:
			a.refresh()
		}
	}
}

// triggerRefresh requests an immediate refresh.
func (a *App) triggerRefresh() {
	select {
	case a.refreshCh <- struct{}{}:
	default:
		// Channel is full, refresh already pending
	}
}

// refresh reloads stats and updates the tray icon and tooltip.
func (a *App) refresh() {
	log.Println("Refreshing stats...")

	// Parse credentials for plan info and OAuth token (required for API)
	credsPath := a.config.GetCredentialsPath()
	creds, err := stats.ParseCredentials(credsPath)
	if err != nil {
		log.Printf("Error: could not parse credentials: %v", err)
		log.Printf("Credentials path: %s", credsPath)
		a.setError()
		return
	}

	// Verify we have an access token
	if creds.ClaudeAiOauth.AccessToken == "" {
		log.Printf("Error: no access token in credentials file")
		a.setError()
		return
	}

	// Parse stats cache (optional - only used as fallback when API unavailable)
	statsPath := a.config.GetStatsPath()
	cache, err := stats.ParseStatsCache(statsPath)
	if err != nil {
		log.Printf("Note: stats cache not available: %v", err)
		// Continue without stats cache - will use API data
		cache = nil
	}

	// Calculate weekly stats (cache can be nil)
	weeklyStats := stats.CalculateWeeklyStats(cache, creds)

	// Fetch real rate limits from API
	a.fetchAndApplyRateLimits(weeklyStats, creds.ClaudeAiOauth.AccessToken, creds.ClaudeAiOauth.RefreshToken)

	// Store stats
	a.statsMu.Lock()
	a.stats = weeklyStats
	a.statsMu.Unlock()

	// Update tray
	a.updateTray(weeklyStats)

	if weeklyStats.HasAPIData {
		log.Printf("Stats refreshed: %d%% weekly usage (API), %d total tokens", weeklyStats.GetPercentage(), weeklyStats.TotalTokens)
	} else {
		log.Printf("Stats refreshed: %d total tokens this week (estimated)", weeklyStats.TotalTokens)
	}
}

// fetchAndApplyRateLimits fetches rate limits from the API and applies them to weeklyStats.
func (a *App) fetchAndApplyRateLimits(weeklyStats *stats.WeeklyStats, token string, refreshToken string) {
	// Initialize or update API client
	if a.apiClient == nil {
		a.apiClient = api.NewClient(token)

		// Set up callback to persist new refresh tokens when the server rotates them
		credsPath := a.config.GetCredentialsPath()
		a.apiClient.SetRefreshTokenCallback(func(newRefreshToken string) {
			log.Printf("Persisting new refresh token to credentials file...")
			if err := stats.UpdateRefreshToken(credsPath, newRefreshToken); err != nil {
				log.Printf("ERROR: Failed to update credentials file with new refresh token: %v", err)
				log.Printf("The new refresh token is in memory but NOT saved. You may need to re-authenticate on restart.")
			} else {
				log.Printf("Successfully updated credentials file with new refresh token")
			}
		})
	} else {
		a.apiClient.SetToken(token)
	}

	// Always set the refresh token so the client can auto-refresh on 401
	a.apiClient.SetRefreshToken(refreshToken)

	// Fetch rate limits
	rateLimits, err := a.apiClient.FetchRateLimits()
	if err != nil {
		log.Printf("Warning: could not fetch rate limits from API: %v", err)
		return
	}

	// Apply to weekly stats
	weeklyStats.HasAPIData = true
	weeklyStats.FiveHourUtilization = rateLimits.FiveHourUtilization
	weeklyStats.WeeklyUtilization = rateLimits.WeeklyUtilization
	weeklyStats.FiveHourReset = rateLimits.FiveHourReset
	weeklyStats.WeeklyReset = rateLimits.WeeklyReset
	weeklyStats.RateLimitStatus = rateLimits.Status
	weeklyStats.RepresentativeClaim = rateLimits.RepresentativeClaim
	weeklyStats.OpusUtilization = rateLimits.OpusUtilization
	weeklyStats.SonnetUtilization = rateLimits.SonnetUtilization
	weeklyStats.OpusReset = rateLimits.OpusReset
	weeklyStats.SonnetReset = rateLimits.SonnetReset

	log.Printf("API rate limits: 5h=%.1f%%, weekly=%.1f%%, status=%s",
		rateLimits.FiveHourUtilization*100,
		rateLimits.WeeklyUtilization*100,
		rateLimits.Status)
}

// updateTray updates the tray icon and tooltip with current stats.
func (a *App) updateTray(weeklyStats *stats.WeeklyStats) {
	// Get usage percentage
	percentage := weeklyStats.GetPercentage()

	// Generate icon with percentage text overlay
	iconBytes, err := a.iconGen.GenerateWithPercentage(weeklyStats, percentage)
	if err != nil {
		log.Printf("Error generating icon: %v", err)
		return
	}

	// Update icon
	a.tray.SetIcon(iconBytes)

	// Update tooltip with platform-appropriate format (Windows gets compact version)
	tooltip := tray.FormatTooltipForPlatform(weeklyStats)
	a.tray.SetTooltip(tooltip)

	log.Printf("Icon updated: %d%% usage", percentage)
}

// setError sets the tray to an error state.
func (a *App) setError() {
	iconBytes, err := a.iconGen.GenerateError()
	if err != nil {
		log.Printf("Error generating error icon: %v", err)
		return
	}

	a.tray.SetIcon(iconBytes)
	a.tray.SetTooltip("Claude Usage\n━━━━━━━━━━━━━━━━━━\nError loading credentials\nMake sure Claude Code is installed\nand you are logged in")
}

// GetStats returns the current weekly stats (thread-safe).
func (a *App) GetStats() *stats.WeeklyStats {
	a.statsMu.RLock()
	defer a.statsMu.RUnlock()
	return a.stats
}

// performUpdate downloads and installs the latest version.
func (a *App) performUpdate() {
	log.Printf("Starting update from %s", update.GetDownloadURL())

	// Show progress in tooltip
	a.tray.SetTooltip("Downloading update...")

	// Perform the update
	result, err := update.Update()
	if err != nil {
		log.Printf("Update failed: %v", err)
		a.tray.SetTooltip("Update failed: " + err.Error())
		// Restore normal tooltip after a delay
		go func() {
			time.Sleep(5 * time.Second)
			a.statsMu.RLock()
			stats := a.stats
			a.statsMu.RUnlock()
			if stats != nil {
				a.updateTray(stats)
			}
		}()
		return
	}

	log.Printf("Update result: %s", result.Message)
	a.tray.SetTooltip("Update installed. Please restart the application.")

	// Mark update as complete - changes menu item to "Restart Required"
	a.tray.SetUpdateComplete()
}
