package tray

import (
	"runtime"

	"fyne.io/systray"
)

// macOSUsageRows is the number of disabled menu rows reserved for the usage
// summary on macOS. It matches menuMaxLines, the maximum number of lines
// FormatMenuLines can emit.
const macOSUsageRows = menuMaxLines

// MenuItems holds references to menu items for updating.
type MenuItems struct {
	Usage        []*systray.MenuItem // Disabled usage-info rows, only populated on macOS
	Version      *systray.MenuItem
	Refresh      *systray.MenuItem
	Update       *systray.MenuItem
	SourceToggle *systray.MenuItem // Only populated on Linux
	Quit         *systray.MenuItem
}

// SetupMenu creates the tray menu with Version display, Refresh, Update, and Quit options.
// The version parameter is displayed as a non-clickable menu item.
// The sourceDisplayName is the current source ("Claude Code" or "OpenCode").
// Returns the menu items for event handling.
func SetupMenu(version string, sourceDisplayName string) *MenuItems {
	items := &MenuItems{}

	// Usage info rows - macOS only. On macOS the tray tooltip (used on Linux to
	// show the usage summary) is a plain OS tooltip that renders the block-char
	// progress bars poorly and is easy to miss, so we surface the same lines as
	// disabled menu rows at the top of the menu instead. They start hidden and
	// are populated by UpdateUsage once stats are available.
	if runtime.GOOS == "darwin" {
		items.Usage = make([]*systray.MenuItem, macOSUsageRows)
		for i := range items.Usage {
			row := systray.AddMenuItem("", "")
			row.Disable()
			row.Hide()
			items.Usage[i] = row
		}
		systray.AddSeparator()
	}

	// Version display (disabled/grayed out - not clickable)
	items.Version = systray.AddMenuItem("Version: "+version, "Current application version")
	items.Version.Disable()

	// Separator
	systray.AddSeparator()

	// Refresh option
	items.Refresh = systray.AddMenuItem("Refresh", "Refresh usage statistics")

	// Update option
	items.Update = systray.AddMenuItem("Update", "Download and install the latest version")

	// Source toggle - Linux only
	if runtime.GOOS == "linux" {
		systray.AddSeparator()
		items.SourceToggle = systray.AddMenuItem("Source: "+sourceDisplayName, "Toggle between Claude Code and OpenCode")
	}

	// Separator
	systray.AddSeparator()

	// Quit option
	items.Quit = systray.AddMenuItem("Quit", "Exit Claude Usage")

	return items
}

// UpdateUsage sets the disabled usage-info rows (macOS only) from the given
// lines. Rows with a corresponding line are shown; surplus rows are hidden.
// It is a no-op when the usage rows were not created (non-macOS), so callers
// don't need their own OS guard.
func (m *MenuItems) UpdateUsage(lines []string) {
	for i, row := range m.Usage {
		if i < len(lines) {
			row.SetTitle(lines[i])
			row.Show()
		} else {
			row.Hide()
		}
	}
}

// UpdateSourceToggle updates the source toggle menu item label.
func (m *MenuItems) UpdateSourceToggle(sourceDisplayName string) {
	if m.SourceToggle != nil {
		m.SourceToggle.SetTitle("Source: " + sourceDisplayName)
	}
}

// HandleMenuEvents starts goroutines to handle menu item clicks.
// onRefresh is called when Refresh is clicked.
// onUpdate is called when Update is clicked.
// onSourceToggle is called when Source toggle is clicked (Linux only).
// onQuit is called when Quit is clicked.
func HandleMenuEvents(items *MenuItems, onRefresh, onUpdate, onSourceToggle, onQuit func()) {
	go func() {
		for {
			// Build select cases dynamically based on available menu items
			if items.SourceToggle != nil {
				select {
				case <-items.Refresh.ClickedCh:
					if onRefresh != nil {
						onRefresh()
					}
				case <-items.Update.ClickedCh:
					if onUpdate != nil {
						onUpdate()
					}
				case <-items.SourceToggle.ClickedCh:
					if onSourceToggle != nil {
						onSourceToggle()
					}
				case <-items.Quit.ClickedCh:
					if onQuit != nil {
						onQuit()
					}
					return
				}
			} else {
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
		}
	}()
}
