package tray

import (
	"fyne.io/systray"
)

// Tray manages the system tray icon and interactions.
type Tray struct {
	menuItems *MenuItems
	version   string
	onRefresh func()
	onUpdate  func()
	onQuit    func()
}

// New creates a new Tray manager with the given version string.
func New(version string) *Tray {
	return &Tray{
		version: version,
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

		// Setup menu with version
		t.menuItems = SetupMenu(t.version)

		// Handle menu events
		HandleMenuEvents(t.menuItems, t.onRefresh, t.onUpdate, func() {
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
