<div align="center">

![Screenshot](screenshot.png)

**Linux**

![Windows Screenshot](screenshot-win.png)

**Windows**

</div>

---

<div align="center">

```
                                                   
  ░█▀▀░█░░░█▀█░█░█░█▀▄░█▀▀░░░█░█░█▀▀░█▀█░█▀▀░█▀▀
  ░█░░░█░░░█▀█░█░█░█░█░█▀▀░░░█░█░▀▀█░█▀█░█░█░█▀▀
  ░▀▀▀░▀▀▀░▀░▀░▀▀▀░▀▀░░▀▀▀░░░▀▀▀░▀▀▀░▀░▀░▀▀▀░▀▀▀
          ╔═══════════════════════════════╗
          ║  SYSTEM TRAY USAGE MONITOR    ║
          ╚═══════════════════════════════╝
```

**`> ESTABLISHING NEURAL UPLINK TO ANTHROPIC SERVERS...`**

**`> CONNECTION SECURED :: MONITORING ACTIVE`**

🔥 *Watch your tokens burn in real-time. Your wallet weeps silently.* 🔥

</div>

---

## `░▒▓█ 0x00 :: SYSTEM OVERVIEW █▓▒░`

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│   Claude Usage is a system tray application that monitors your             │
│   Claude Code API consumption. Tracks token burn rate across models.       │
│                                                                             │
│   > No browser required                                                     │
│   > Real-time progress bars                                                 │
│   > Multi-model breakdown (Opus, Sonnet)                                    │
│   > Reset countdown timers                                                  │
│   > Works on Linux, Windows, macOS                                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## `░▒▓█ 0x01 :: PLATFORM COMPATIBILITY MATRIX █▓▒░`

```
> SCANNING TARGET ARCHITECTURES...
> [OK] COMPATIBILITY CHECK COMPLETE
```

| PLATFORM | ARCH | STATUS | PACKAGE |
|----------|------|--------|---------|
| **Windows** | x64 | `[SUPPORTED]` | `claude-usage-windows-amd64.zip` |
| **macOS** | Universal | `[SUPPORTED]` | `Claude-Usage-*.dmg` |
| **macOS** | Intel | `[SUPPORTED]` | `claude-usage-darwin-amd64` |
| **macOS** | Apple Silicon | `[SUPPORTED]` | `claude-usage-darwin-arm64` |
| **Linux** | x64 | `[SUPPORTED]` | `claude-usage-linux-amd64` |
| **Linux** | ARM64 | `[SUPPORTED]` | `claude-usage-linux-arm64` |

```
> DESKTOP ENVIRONMENT SCAN:
  ├─ [OK] KDE Plasma
  ├─ [OK] GNOME (requires AppIndicator extension)
  ├─ [OK] XFCE
  ├─ [OK] Cinnamon
  ├─ [OK] MATE
  ├─ [OK] Budgie
  └─ [OK] Any DE supporting StatusNotifierItem
```

---

## `░▒▓█ 0x02 :: DEPLOYMENT PROTOCOLS █▓▒░`

```
╔══════════════════════════════════════════════════════════════════════════════╗
║  > BINARY PACKAGES AVAILABLE AT:                                             ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

[https://github.com/utajum/claude-usage/releases](https://github.com/utajum/claude-usage/releases)

---

### **Windows** `> TARGET: WINDOWS_x64`

```
1. Download claude-usage-windows-amd64.zip from releases
2. Extract the ZIP file
3. Right-click install-windows.ps1 → "Run with PowerShell"
4. Find "Claude Usage" in Start Menu
```

**Manual Installation:**
```powershell
# Or run from PowerShell directly:
powershell -ExecutionPolicy Bypass -File install-windows.ps1

# To uninstall:
powershell -ExecutionPolicy Bypass -File install-windows.ps1 -Uninstall
```

**What gets installed:**
- `%LOCALAPPDATA%\Programs\claude-usage\claude-usage.exe`
- Start Menu shortcut
- Startup shortcut (autostart on login)

---

### **macOS** `> TARGET: MACOS (Universal)`

```bash
# Option 1: Download DMG (recommended)
1. Download Claude-Usage-*.dmg from releases
2. Open the DMG file
3. Drag "Claude Usage" to Applications folder
4. Right-click → Open (first time only, to bypass Gatekeeper)
```

```bash
# Option 2: Command line install
curl -sL https://github.com/utajum/claude-usage/releases/latest/download/claude-usage-darwin-arm64 -o claude-usage
chmod +x claude-usage
./claude-usage
```

**Enable autostart:**
```bash
# Copy the LaunchAgent plist (from source repo)
cp assets/macos/com.github.utajum.claude-usage.plist ~/Library/LaunchAgents/
launchctl load ~/Library/LaunchAgents/com.github.utajum.claude-usage.plist
```

---

### **Linux** `> TARGET: LINUX (x64/ARM64)`

**One-liner install (recommended):**
```bash
curl -sL https://raw.githubusercontent.com/utajum/claude-usage/master/scripts/install-linux.sh | bash
```

This automatically:
- Detects your architecture (x64 or ARM64)
- Downloads the latest release
- Installs to `~/.local/bin/`
- Creates desktop entry (shows in applications menu)
- Enables autostart on login

**Install options:**
```bash
# Install without autostart
curl -sL https://raw.githubusercontent.com/utajum/claude-usage/master/scripts/install-linux.sh | bash -s -- --no-autostart

# Install without desktop entry
curl -sL https://raw.githubusercontent.com/utajum/claude-usage/master/scripts/install-linux.sh | bash -s -- --no-desktop

# Uninstall
curl -sL https://raw.githubusercontent.com/utajum/claude-usage/master/scripts/install-linux.sh | bash -s -- --uninstall
```

**Manual download:**
```bash
# x64
curl -sL https://github.com/utajum/claude-usage/releases/latest/download/claude-usage-linux-amd64 -o claude-usage

# ARM64
curl -sL https://github.com/utajum/claude-usage/releases/latest/download/claude-usage-linux-arm64 -o claude-usage

chmod +x claude-usage
./claude-usage
```

**Build from source:**
```bash
git clone https://github.com/utajum/claude-usage
cd claude-usage
make install-linux    # Full install with desktop integration
# or
make install          # Binary only
```

---

## `░▒▓█ 0x03 :: NEURAL INTERFACE █▓▒░`

```
> HOVER OVER TRAY ICON TO ACCESS TELEMETRY FEED
```

Tooltip displays real-time consumption data:

```
┌────────────────────────────────────┐
│ CLAUDE USAGE                       │
│ Plan: Pro (5x)                     │
│ STATUS: THROTTLED                  │
│ ▕████████░░▏  80% 2h 15m ◀         │
│ ▕██████░░░░▏  60% 3d 5h            │
└────────────────────────────────────┘

LEGEND:
├─ First bar  :: 5-hour rolling window
├─ Second bar :: Weekly allocation
├─ ◀ marker   :: Active rate limiter
└─ Time       :: Reset countdown
```

---

## `░▒▓█ 0x04 :: DATA STREAM SOURCES █▓▒░`

```
> SCANNING LOCAL FILESYSTEM FOR CLAUDE CREDENTIALS...
```

The app reads your Claude Code credentials to fetch real-time rate limit data from the API:

| OS | CREDENTIALS FILE |
|----|------------------|
| **Linux** | `~/.claude/.credentials.json` |
| **macOS** | `~/.claude/.credentials.json` |
| **Windows** | `C:\Users\<username>\.claude\.credentials.json` |

```
> DATA FLOW:
  ├─ [1] Read OAuth token from .credentials.json
  ├─ [2] Call Anthropic API for real-time rate limits
  └─ [3] Display usage % and reset timers

> PREREQUISITES:
  ├─ [REQUIRED] Claude Code CLI installed and logged in
  └─ [REQUIRED] Valid credentials file with OAuth token
```

**Note:** The `stats-cache.json` file is only used as a fallback when API data is unavailable.

---

## `░▒▓█ 0x05 :: COMPILE FROM SOURCE █▓▒░`

```
> INITIATING BUILD SEQUENCE...
```

```bash
# Clone repository
git clone https://github.com/utajum/claude-usage
cd claude-usage

# Build for current platform
make build

# Or cross-compile for all targets
make build-all

# Generate icon assets
make generate-icons

# Platform-specific packaging
make build-macos-app  # Creates .app bundle (macOS only)
```

### `> BUILD REQUIREMENTS`

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  DEPENDENCY          VERSION        STATUS                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│  Go                  1.22+          [REQUIRED]                              │
│  CGO                 enabled        [REQUIRED for macOS/Windows]            │
│  External libs       none           [PURE GO BUILD on Linux]                │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## `░▒▓█ 0x06 :: SYSTEM CONFIGURATION █▓▒░`

```
> LOADING CONFIG PARAMETERS...
```

Config file locations:

| OS | CONFIG_PATH |
|----|-------------|
| **Linux** | `~/.config/claude-usage/config.json` |
| **macOS** | `~/Library/Application Support/claude-usage/config.json` |
| **Windows** | `%APPDATA%\claude-usage\config.json` |

```json
{
  "refresh_interval_seconds": 300
}
```

```
> DEFAULT REFRESH RATE: 300 seconds (5 minutes)
```

---

## `░▒▓█ 0x07 :: PERSISTENCE PROTOCOLS █▓▒░`

### **Linux** `> AUTOSTART`

```bash
make autostart          # Enable persistence
make autostart-remove   # Disable persistence
```

### **Windows** `> AUTOSTART`

The installer automatically creates a Startup shortcut. To manage manually:

```
Location: %APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup\Claude Usage.lnk
```

### **macOS** `> AUTOSTART`

```bash
# Enable
launchctl load ~/Library/LaunchAgents/com.github.utajum.claude-usage.plist

# Disable  
launchctl unload ~/Library/LaunchAgents/com.github.utajum.claude-usage.plist
```

Or via System Settings → General → Login Items

---

## `░▒▓█ 0x08 :: CORE LOGIC FLOW █▓▒░`

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│   [1] READ    ──► Parse ~/.claude/stats-cache.json                          │
│                                                                             │
│   [2] PROCESS ──► Calculate token consumption per model                     │
│                                                                             │
│   [3] COMPUTE ──► Aggregate weekly totals (Monday-Sunday cycle)             │
│                                                                             │
│   [4] RENDER  ──► Generate dynamic tray icon based on usage %               │
│                                                                             │
│   [5] DISPLAY ──► System tray with hover tooltip                            │
│                                                                             │
│   [6] LOOP    ──► Auto-refresh every 5 minutes                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## `░▒▓█ 0x09 :: DIAGNOSTICS █▓▒░`

### `> ERROR: ICON_NOT_VISIBLE`

```
LINUX:
├─ Verify DE supports StatusNotifierItem
├─ GNOME: Install AppIndicator extension
└─ Check system tray is enabled

WINDOWS:
├─ Check system tray overflow area (click ^ arrow)
├─ Right-click taskbar → Taskbar settings → System tray icons
└─ Enable "Claude Usage" to always show

MACOS:
├─ Check menu bar (may be hidden by notch on newer Macs)
└─ Try Bartender or similar app to manage menu bar items

ALL PLATFORMS:
└─ Confirm Claude Code is installed and operational
```

### `> ERROR: NO_DATA_AVAILABLE`

```
├─ Execute 'claude' command at least once
├─ Verify ~/.claude/stats-cache.json exists
└─ Check file permissions
```

---

## `░▒▓█ 0x0A :: MAKE TARGETS █▓▒░`

```bash
make help  # Show all available targets
```

```
Development:
  make build          - Build for current platform
  make run            - Build and run
  make test           - Run tests
  make generate-icons - Generate icon assets

Installation (Linux):
  make install        - Install binary to ~/.local/bin
  make install-linux  - Full install (binary + desktop + autostart)
  make desktop-install- Install desktop entry (app menu)
  make autostart      - Enable autostart on login
  make uninstall      - Remove everything

Installation (macOS):
  make install-macos  - Install .app bundle + autostart
  make uninstall-macos- Remove app and autostart

Cross-compilation:
  make build-linux    - Build for Linux (amd64, arm64)
  make build-windows  - Build for Windows (amd64)
  make build-macos    - Build for macOS (Intel, Apple Silicon)
  make build-macos-app- Build macOS .app bundle
  make build-all      - Build for all platforms
```

---

## `░▒▓█ 0x0B :: LICENSE █▓▒░`

```
MIT License

Permission granted to copy, modify, distribute.
No warranty. Use at your own risk.
```

---

<div align="center">

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                                                                              ║
║   ░▒▓█ CONNECTION TERMINATED :: END OF TRANSMISSION █▓▒░                     ║
║                                                                              ║
║   > Stay connected to the grid                                               ║
║   > Monitor your burn rate                                                   ║
║   > Trust no one. Especially your token consumption.                         ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

**`> SESSION_END :: 0x00000000`**

</div>
