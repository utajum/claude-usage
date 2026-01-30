package tray

import (
	"fyne.io/systray"
)

// MenuItems holds references to menu items for updating.
type MenuItems struct {
	Refresh *systray.MenuItem
	Quit    *systray.MenuItem
}

// SetupMenu creates the minimal tray menu with Refresh and Quit options.
// Returns the menu items for event handling.
func SetupMenu() *MenuItems {
	items := &MenuItems{}

	// Refresh option
	items.Refresh = systray.AddMenuItem("Refresh", "Refresh usage statistics")

	// Separator
	systray.AddSeparator()

	// Quit option
	items.Quit = systray.AddMenuItem("Quit", "Exit Claude Usage")

	return items
}

// HandleMenuEvents starts goroutines to handle menu item clicks.
// onRefresh is called when Refresh is clicked.
// onQuit is called when Quit is clicked.
func HandleMenuEvents(items *MenuItems, onRefresh, onQuit func()) {
	go func() {
		for {
			select {
			case <-items.Refresh.ClickedCh:
				if onRefresh != nil {
					onRefresh()
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
