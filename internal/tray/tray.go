package tray

import (
	"fyne.io/systray"
)

// Tray manages the system tray icon and interactions.
type Tray struct {
	menuItems         *MenuItems
	version           string
	sourceDisplayName string
	onRefresh         func()
	onUpdate          func()
	onSourceToggle    func()
	onQuit            func()
}

// New creates a new Tray manager with the given version string and source display name.
func New(version string, sourceDisplayName string) *Tray {
	return &Tray{
		version:           version,
		sourceDisplayName: sourceDisplayName,
	}
}

// SetOnRefresh sets the callback for the Refresh menu item.
func (t *Tray) SetOnRefresh(fn func()) {
	t.onRefresh = fn
}

// SetOnUpdate sets the callback for the Update menu item.
func (t *Tray) SetOnUpdate(fn func()) {
	t.onUpdate = fn
}

// SetOnSourceToggle sets the callback for the Source toggle menu item (Linux only).
func (t *Tray) SetOnSourceToggle(fn func()) {
	t.onSourceToggle = fn
}

// SetOnQuit sets the callback for the Quit menu item.
func (t *Tray) SetOnQuit(fn func()) {
	t.onQuit = fn
}

// Run starts the system tray. This blocks until Quit is called.
// onReady is called when the tray is initialized and ready.
func (t *Tray) Run(onReady func()) {
	systray.Run(func() {
		// Set initial title
		systray.SetTitle("")
		systray.SetTooltip("Claude Usage - Loading...")

		// Setup menu with version and source
		t.menuItems = SetupMenu(t.version, t.sourceDisplayName)

		// Handle menu events
		HandleMenuEvents(t.menuItems, t.onRefresh, t.onUpdate, t.onSourceToggle, func() {
			if t.onQuit != nil {
				t.onQuit()
			}
			systray.Quit()
		})

		// Call ready callback
		if onReady != nil {
			onReady()
		}
	}, func() {
		// onExit callback - cleanup if needed
	})
}

// UpdateSourceToggle updates the source toggle menu item label.
func (t *Tray) UpdateSourceToggle(sourceDisplayName string) {
	t.sourceDisplayName = sourceDisplayName
	if t.menuItems != nil {
		t.menuItems.UpdateSourceToggle(sourceDisplayName)
	}
}

// SetIcon sets the tray icon from PNG bytes.
func (t *Tray) SetIcon(iconBytes []byte) {
	systray.SetIcon(iconBytes)
}

// SetTooltip sets the tray tooltip text.
func (t *Tray) SetTooltip(text string) {
	systray.SetTooltip(text)
}

// Quit exits the system tray.
func (t *Tray) Quit() {
	systray.Quit()
}

// SetUpdateComplete marks the update as complete and changes the menu item text.
// The menu item is disabled since the user needs to restart.
func (t *Tray) SetUpdateComplete() {
	if t.menuItems != nil && t.menuItems.Update != nil {
		t.menuItems.Update.SetTitle("Restart Required")
		t.menuItems.Update.SetTooltip("Update downloaded. Please restart the application.")
		t.menuItems.Update.Disable()
	}
}
