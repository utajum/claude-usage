```
╔══════════════════════════════════════════════════════════════════════════════╗
║  > RELATED_SYSTEMS                                                           ║
║  ├─ 📧 gmail-notifier :: KDE plasma notification daemon for Gmail            ║
║  └─ [https://github.com/utajum/gmail-notifier](https://github.com/utajum/gmail-notifier) ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

---

<div align="center">

![Screenshot](screenshot.png)

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

| PLATFORM | ARCH | STATUS | BINARY |
|----------|------|--------|--------|
| 🐧 `LINUX` | x64 | `[SUPPORTED]` | `claude-usage-linux-amd64` |
| 🐧 `LINUX` | ARM64 | `[SUPPORTED]` | `claude-usage-linux-arm64` |
| 🪟 `WINDOWS` | x64 | `[SUPPORTED]` | `claude-usage-windows-amd64.exe` |
| 🍎 `MACOS` | Intel | `[SUPPORTED]` | `claude-usage-darwin-amd64` |
| 🍎 `MACOS` | Apple Silicon | `[SUPPORTED]` | `claude-usage-darwin-arm64` |

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

### 🐧 `> TARGET: LINUX_x64`

```bash
# Download binary
$ curl -sL https://github.com/utajum/claude-usage/releases/latest/download/claude-usage-linux-amd64 -o claude-usage

# Initialize and execute
$ chmod +x claude-usage && ./claude-usage

# [OK] PROCESS SPAWNED :: MONITORING ACTIVE
```

### 🐧 `> TARGET: LINUX_ARM64`

```bash
$ curl -sL https://github.com/utajum/claude-usage/releases/latest/download/claude-usage-linux-arm64 -o claude-usage
$ chmod +x claude-usage && ./claude-usage
```

### 🪟 `> TARGET: WINDOWS_x64`

```
1. Download claude-usage-windows-amd64.exe from releases
```
   [https://github.com/utajum/claude-usage/releases](https://github.com/utajum/claude-usage/releases)

```
2. Execute binary

3. Locate icon in system tray
   > [OK] UPLINK ESTABLISHED
```

### 🍎 `> TARGET: MACOS_ARM64 (Apple Silicon)`

```bash
$ curl -sL https://github.com/utajum/claude-usage/releases/latest/download/claude-usage-darwin-arm64 -o claude-usage
$ chmod +x claude-usage && ./claude-usage
```

### 🍎 `> TARGET: MACOS_x64 (Intel)`

```bash
$ curl -sL https://github.com/utajum/claude-usage/releases/latest/download/claude-usage-darwin-amd64 -o claude-usage
$ chmod +x claude-usage && ./claude-usage
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
> SCANNING LOCAL FILESYSTEM FOR CLAUDE TELEMETRY...
```

| OS | DATA_PATH | STATUS |
|----|-----------|--------|
| 🐧 `LINUX` | `~/.claude/stats-cache.json` | `[ACTIVE]` |
| 🍎 `MACOS` | `~/.claude/stats-cache.json` | `[ACTIVE]` |
| 🪟 `WINDOWS` | `%USERPROFILE%\.claude\stats-cache.json` | `[ACTIVE]` |

```
> PREREQUISITES:
  ├─ [REQUIRED] Claude Code CLI installed
  └─ [REQUIRED] At least one Claude session executed (generates stats)
```

---

## `░▒▓█ 0x05 :: COMPILE FROM SOURCE █▓▒░`

```
> INITIATING BUILD SEQUENCE...
```

```bash
# Clone repository
$ git clone https://github.com/utajum/claude-usage
$ cd claude-usage

# Build for current platform
$ make build

# Or cross-compile for all targets
$ make build-all

# Deploy to ~/.local/bin
$ make install

# Enable persistence (Linux)
$ make autostart

# [OK] BUILD COMPLETE
```

### `> BUILD REQUIREMENTS`

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  DEPENDENCY          VERSION        STATUS                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│  Go                  1.24+          [REQUIRED]                              │
│  External libs       none           [PURE GO BUILD]                         │
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
| 🐧 `LINUX` | `~/.config/claude-usage/config.json` |
| 🍎 `MACOS` | `~/Library/Application Support/claude-usage/config.json` |
| 🪟 `WINDOWS` | `%APPDATA%\claude-usage\config.json` |

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

### 🐧 `> LINUX :: AUTOSTART`

```bash
$ make autostart          # Enable persistence
$ make autostart-remove   # Disable persistence
```

### 🪟 `> WINDOWS :: AUTOSTART`

```
> ADD SHORTCUT TO:
> %APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup
```

### 🍎 `> MACOS :: AUTOSTART`

```
> System Preferences → Users & Groups → Login Items
> OR create LaunchAgent plist
```

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

## `░▒▓█ 0x0A :: LICENSE █▓▒░`

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
