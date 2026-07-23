module claude-usage

go 1.24.0

toolchain go1.24.12

require fyne.io/systray v1.12.0

// Vendored locally to patch macOS menu display: the upstream implementation
// pops the tray menu up manually with popUpMenuPositioningItem, which makes
// tall menus show a scroll arrow (▲) and shift content on hover. The patch
// assigns the menu to the status item so macOS positions it natively without
// scroll arrows. See third_party/systray/systray_darwin.m.
replace fyne.io/systray => ./third_party/systray

require (
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
)
