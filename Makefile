# Claude Usage - Makefile
# Cross-platform build targets for Linux, Windows, and macOS

BINARY_NAME=claude-usage
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_FLAGS=-ldflags "-X main.Version=$(VERSION)"

# Output directories
DIST_DIR=dist
INSTALL_DIR=$(HOME)/.local/bin

# Linux paths
DESKTOP_FILE=$(HOME)/.config/autostart/claude-usage.desktop
APPLICATIONS_DIR=$(HOME)/.local/share/applications
ICONS_DIR=$(HOME)/.local/share/icons/hicolor

.PHONY: all build clean install uninstall run test lint \
        build-linux build-windows build-macos build-all \
        autostart autostart-remove desktop-install desktop-remove \
        install-linux install-macos generate-icons \
        build-macos-app help

# Default target
all: build

# Build for current platform
build:
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) ./cmd/claude-usage

# Run the application
run: build
	./$(BINARY_NAME)

# Install to ~/.local/bin
install: build
	@mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)/
	@echo "Installed to $(INSTALL_DIR)/$(BINARY_NAME)"
	@echo "Make sure $(INSTALL_DIR) is in your PATH"

# Uninstall
uninstall: autostart-remove desktop-remove
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Removed $(INSTALL_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf $(DIST_DIR)

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# ============================================
# Icon Generation
# ============================================

# Generate all icon assets
generate-icons:
	go run scripts/generate-icons/main.go

# ============================================
# Cross-platform builds
# ============================================

# Create dist directory
$(DIST_DIR):
	mkdir -p $(DIST_DIR)

# Build for Linux (amd64)
build-linux-amd64: $(DIST_DIR)
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/claude-usage

# Build for Linux (arm64)
build-linux-arm64: $(DIST_DIR)
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/claude-usage

# Build for Linux (all architectures)
build-linux: build-linux-amd64 build-linux-arm64

# Build for Windows (amd64)
build-windows-amd64: $(DIST_DIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build $(BUILD_FLAGS) -ldflags "-X main.Version=$(VERSION) -H=windowsgui" -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/claude-usage

# Build for Windows
build-windows: build-windows-amd64

# Build for macOS (Intel)
build-macos-amd64: $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/claude-usage

# Build for macOS (Apple Silicon)
build-macos-arm64: $(DIST_DIR)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/claude-usage

# Build for macOS (all architectures)
build-macos: build-macos-amd64 build-macos-arm64

# Build macOS .app bundle (run on macOS only)
build-macos-app: build-macos
	./scripts/build-macos-app.sh $(VERSION)

# Build for all platforms
build-all: build-linux build-windows build-macos
	@echo "Built binaries for all platforms in $(DIST_DIR)/"
	@ls -la $(DIST_DIR)/

# ============================================
# Linux Desktop Integration
# ============================================

# Install desktop entry (applications menu)
desktop-install: install generate-icons
	@mkdir -p $(APPLICATIONS_DIR)
	@mkdir -p $(ICONS_DIR)/16x16/apps
	@mkdir -p $(ICONS_DIR)/24x24/apps
	@mkdir -p $(ICONS_DIR)/32x32/apps
	@mkdir -p $(ICONS_DIR)/48x48/apps
	@mkdir -p $(ICONS_DIR)/64x64/apps
	@mkdir -p $(ICONS_DIR)/128x128/apps
	@mkdir -p $(ICONS_DIR)/256x256/apps
	@cp assets/linux/claude-usage-16.png $(ICONS_DIR)/16x16/apps/claude-usage.png
	@cp assets/linux/claude-usage-24.png $(ICONS_DIR)/24x24/apps/claude-usage.png
	@cp assets/linux/claude-usage-32.png $(ICONS_DIR)/32x32/apps/claude-usage.png
	@cp assets/linux/claude-usage-48.png $(ICONS_DIR)/48x48/apps/claude-usage.png
	@cp assets/linux/claude-usage-64.png $(ICONS_DIR)/64x64/apps/claude-usage.png
	@cp assets/linux/claude-usage-128.png $(ICONS_DIR)/128x128/apps/claude-usage.png
	@cp assets/linux/claude-usage-256.png $(ICONS_DIR)/256x256/apps/claude-usage.png
	@sed "s|Exec=claude-usage|Exec=$(INSTALL_DIR)/$(BINARY_NAME)|" assets/linux/claude-usage.desktop > $(APPLICATIONS_DIR)/claude-usage.desktop
	@gtk-update-icon-cache $(ICONS_DIR) 2>/dev/null || true
	@update-desktop-database $(APPLICATIONS_DIR) 2>/dev/null || true
	@echo "Desktop entry installed to $(APPLICATIONS_DIR)/claude-usage.desktop"
	@echo "Icons installed to $(ICONS_DIR)/*/apps/claude-usage.png"

# Remove desktop entry
desktop-remove:
	@rm -f $(APPLICATIONS_DIR)/claude-usage.desktop
	@rm -f $(ICONS_DIR)/16x16/apps/claude-usage.png
	@rm -f $(ICONS_DIR)/24x24/apps/claude-usage.png
	@rm -f $(ICONS_DIR)/32x32/apps/claude-usage.png
	@rm -f $(ICONS_DIR)/48x48/apps/claude-usage.png
	@rm -f $(ICONS_DIR)/64x64/apps/claude-usage.png
	@rm -f $(ICONS_DIR)/128x128/apps/claude-usage.png
	@rm -f $(ICONS_DIR)/256x256/apps/claude-usage.png
	@gtk-update-icon-cache $(ICONS_DIR) 2>/dev/null || true
	@update-desktop-database $(APPLICATIONS_DIR) 2>/dev/null || true
	@echo "Desktop entry removed"

# Install autostart entry (Linux)
autostart: install generate-icons
	@mkdir -p $(HOME)/.config/autostart
	@sed "s|Exec=claude-usage|Exec=$(INSTALL_DIR)/$(BINARY_NAME)|" assets/linux/claude-usage.desktop > $(DESKTOP_FILE)
	@echo "X-GNOME-Autostart-enabled=true" >> $(DESKTOP_FILE)
	@echo "Autostart entry created at $(DESKTOP_FILE)"

# Remove autostart entry
autostart-remove:
	@rm -f $(DESKTOP_FILE)
	@echo "Removed autostart entry"

# Full Linux install (binary + desktop entry + autostart)
install-linux: desktop-install autostart
	@echo ""
	@echo "Linux installation complete!"
	@echo "  - Binary: $(INSTALL_DIR)/$(BINARY_NAME)"
	@echo "  - Desktop entry: $(APPLICATIONS_DIR)/claude-usage.desktop"
	@echo "  - Autostart: $(DESKTOP_FILE)"
	@echo ""
	@echo "You can now find 'Claude Usage' in your applications menu."

# ============================================
# macOS Installation (run on macOS only)
# ============================================

LAUNCHAGENT_DIR=$(HOME)/Library/LaunchAgents
LAUNCHAGENT_PLIST=com.github.utajum.claude-usage.plist

# Install macOS app bundle and launchagent
install-macos: build-macos-app
	@echo "Installing Claude Usage.app to /Applications..."
	@rm -rf "/Applications/Claude Usage.app"
	@cp -r "$(DIST_DIR)/Claude Usage.app" /Applications/
	@mkdir -p $(LAUNCHAGENT_DIR)
	@cp assets/macos/$(LAUNCHAGENT_PLIST) $(LAUNCHAGENT_DIR)/
	@launchctl unload $(LAUNCHAGENT_DIR)/$(LAUNCHAGENT_PLIST) 2>/dev/null || true
	@launchctl load $(LAUNCHAGENT_DIR)/$(LAUNCHAGENT_PLIST)
	@echo ""
	@echo "macOS installation complete!"
	@echo "  - App: /Applications/Claude Usage.app"
	@echo "  - Autostart: $(LAUNCHAGENT_DIR)/$(LAUNCHAGENT_PLIST)"
	@echo ""
	@echo "Claude Usage will start automatically on login."

# Uninstall macOS app
uninstall-macos:
	@launchctl unload $(LAUNCHAGENT_DIR)/$(LAUNCHAGENT_PLIST) 2>/dev/null || true
	@rm -f $(LAUNCHAGENT_DIR)/$(LAUNCHAGENT_PLIST)
	@rm -rf "/Applications/Claude Usage.app"
	@echo "macOS uninstallation complete"

# ============================================
# Development helpers
# ============================================

# Show version
version:
	@echo $(VERSION)

# Format code
fmt:
	go fmt ./...

# Update dependencies
deps:
	go mod tidy
	go mod download

# Show help
help:
	@echo "Claude Usage - Build Targets"
	@echo ""
	@echo "Development:"
	@echo "  make build          - Build for current platform"
	@echo "  make run            - Build and run"
	@echo "  make test           - Run tests"
	@echo "  make lint           - Run linter"
	@echo "  make fmt            - Format code"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make generate-icons - Generate icon assets"
	@echo ""
	@echo "Installation (Linux):"
	@echo "  make install        - Install binary to ~/.local/bin"
	@echo "  make install-linux  - Full install (binary + desktop + autostart)"
	@echo "  make desktop-install- Install desktop entry (app menu)"
	@echo "  make autostart      - Enable autostart on login"
	@echo "  make uninstall      - Remove everything"
	@echo ""
	@echo "Installation (macOS):"
	@echo "  make install-macos  - Install .app bundle + autostart"
	@echo "  make uninstall-macos- Remove app and autostart"
	@echo ""
	@echo "Installation (Windows):"
	@echo "  Run: powershell -File scripts/install-windows.ps1"
	@echo "  Uninstall: powershell -File scripts/install-windows.ps1 -Uninstall"
	@echo ""
	@echo "Cross-compilation:"
	@echo "  make build-linux    - Build for Linux (amd64, arm64)"
	@echo "  make build-windows  - Build for Windows (amd64)"
	@echo "  make build-macos    - Build for macOS (Intel, Apple Silicon)"
	@echo "  make build-macos-app- Build macOS .app bundle"
	@echo "  make build-all      - Build for all platforms"
