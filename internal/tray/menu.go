package tray

import (
	"fyne.io/systray"
)

// MenuItems holds references to menu items for updating.
type MenuItems struct {
	Version *systray.MenuItem
	Refresh *systray.MenuItem
	Update  *systray.MenuItem
	Quit    *systray.MenuItem
}

// SetupMenu creates the tray menu with Version display, Refresh, Update, and Quit options.
// The version parameter is displayed as a non-clickable menu item.
// Returns the menu items for event handling.
func SetupMenu(version string) *MenuItems {
	items := &MenuItems{}

	// Version display (disabled/grayed out - not clickable)
	items.Version = systray.AddMenuItem("Version: "+version, "Current application version")
	items.Version.Disable()

	// Separator
	systray.AddSeparator()

	// Refresh option
	items.Refresh = systray.AddMenuItem("Refresh", "Refresh usage statistics")

	// Update option
	items.Update = systray.AddMenuItem("Update", "Download and install the latest version")

	// Separator
	systray.AddSeparator()

	// Quit option
	items.Quit = systray.AddMenuItem("Quit", "Exit Claude Usage")

	return items
}

// HandleMenuEvents starts goroutines to handle menu item clicks.
// onRefresh is called when Refresh is clicked.
// onUpdate is called when Update is clicked.
// onQuit is called when Quit is clicked.
func HandleMenuEvents(items *MenuItems, onRefresh, onUpdate, onQuit func()) {
	go func() {
		for {
			select {
			case <-items.Refresh.ClickedCh:
				if onRefresh != nil {
					onRefresh()
				}
			case <-items.Update.ClickedCh:
				if onUpdate != nil {
					onUpdate()
				}
			case <-items.Quit.ClickedCh:
				if onQuit != nil {
					onQuit()
				}
				return
			}
		}
	}()
}
