// Package app provides the main application orchestration for claude-usage.
package app

import (
	"log"
	"sync"
	"time"

	"github.com/user/claude-usage/internal/api"
	"github.com/user/claude-usage/internal/config"
	"github.com/user/claude-usage/internal/icon"
	"github.com/user/claude-usage/internal/stats"
	"github.com/user/claude-usage/internal/tray"
)

// App is the main application struct that coordinates all components.
type App struct {
	config    *config.Config
	tray      *tray.Tray
	iconGen   *icon.Generator
	apiClient *api.Client
	stats     *stats.WeeklyStats
	statsMu   sync.RWMutex
	stopCh    chan struct{}
	refreshCh chan struct{}
}

// New creates a new App instance.
func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Warning: could not load config, using defaults: %v", err)
		cfg = config.Default()
	}

	return &App{
		config:    cfg,
		tray:      tray.New(),
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

	// Parse stats cache
	statsPath := a.config.GetStatsPath()
	cache, err := stats.ParseStatsCache(statsPath)
	if err != nil {
		log.Printf("Error parsing stats: %v", err)
		a.setError()
		return
	}

	// Parse credentials for plan info and OAuth token
	credsPath := a.config.GetCredentialsPath()
	creds, err := stats.ParseCredentials(credsPath)
	if err != nil {
		log.Printf("Warning: could not parse credentials: %v", err)
		// Continue without credentials - just won't show plan info
	}

	// Calculate weekly stats from local cache
	weeklyStats := stats.CalculateWeeklyStats(cache, creds)

	// Fetch real rate limits from API (if we have a token)
	if creds != nil && creds.ClaudeAiOauth.AccessToken != "" {
		a.fetchAndApplyRateLimits(weeklyStats, creds.ClaudeAiOauth.AccessToken)
	}

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
func (a *App) fetchAndApplyRateLimits(weeklyStats *stats.WeeklyStats, token string) {
	// Initialize or update API client
	if a.apiClient == nil {
		a.apiClient = api.NewClient(token)
	} else {
		a.apiClient.SetToken(token)
	}

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

	// Update tooltip with full stats including plan info
	tooltip := tray.FormatTooltip(weeklyStats)
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
	a.tray.SetTooltip("Claude Usage\n━━━━━━━━━━━━━━━━━━\nError loading stats\nCheck that Claude Code is installed")
}

// GetStats returns the current weekly stats (thread-safe).
func (a *App) GetStats() *stats.WeeklyStats {
	a.statsMu.RLock()
	defer a.statsMu.RUnlock()
	return a.stats
}
