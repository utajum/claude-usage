# Claude Usage - Makefile
# Cross-platform build targets for Linux, Windows, and macOS

BINARY_NAME=claude-usage
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_FLAGS=-ldflags "-X main.Version=$(VERSION)"

# Output directories
DIST_DIR=dist
INSTALL_DIR=$(HOME)/.local/bin

.PHONY: all build clean install uninstall run test lint \
        build-linux build-windows build-macos build-all \
        autostart autostart-remove

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
uninstall:
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
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -ldflags "-X main.Version=$(VERSION) -H=windowsgui" -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/claude-usage

# Build for Windows
build-windows: build-windows-amd64

# Build for macOS (Intel)
build-macos-amd64: $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/claude-usage

# Build for macOS (Apple Silicon)
build-macos-arm64: $(DIST_DIR)
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/claude-usage

# Build for macOS (all architectures)
build-macos: build-macos-amd64 build-macos-arm64

# Build for all platforms
build-all: build-linux build-windows build-macos
	@echo "Built binaries for all platforms in $(DIST_DIR)/"
	@ls -la $(DIST_DIR)/

# ============================================
# Autostart (Linux only)
# ============================================

DESKTOP_FILE=$(HOME)/.config/autostart/claude-usage.desktop

# Install autostart entry (Linux)
autostart: install
	@mkdir -p $(HOME)/.config/autostart
	@echo "[Desktop Entry]" > $(DESKTOP_FILE)
	@echo "Type=Application" >> $(DESKTOP_FILE)
	@echo "Name=Claude Usage" >> $(DESKTOP_FILE)
	@echo "Comment=Monitor Claude Code token usage" >> $(DESKTOP_FILE)
	@echo "Exec=$(INSTALL_DIR)/$(BINARY_NAME)" >> $(DESKTOP_FILE)
	@echo "Icon=utilities-system-monitor" >> $(DESKTOP_FILE)
	@echo "Terminal=false" >> $(DESKTOP_FILE)
	@echo "Categories=Utility;" >> $(DESKTOP_FILE)
	@echo "StartupNotify=false" >> $(DESKTOP_FILE)
	@echo "X-GNOME-Autostart-enabled=true" >> $(DESKTOP_FILE)
	@echo "Autostart entry created at $(DESKTOP_FILE)"

# Remove autostart entry
autostart-remove:
	rm -f $(DESKTOP_FILE)
	@echo "Removed autostart entry"

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
	@echo "  make build      - Build for current platform"
	@echo "  make run        - Build and run"
	@echo "  make test       - Run tests"
	@echo "  make lint       - Run linter"
	@echo "  make fmt        - Format code"
	@echo "  make clean      - Remove build artifacts"
	@echo ""
	@echo "Installation:"
	@echo "  make install    - Install to ~/.local/bin"
	@echo "  make uninstall  - Remove from ~/.local/bin"
	@echo "  make autostart  - Enable autostart on login (Linux)"
	@echo "  make autostart-remove - Disable autostart"
	@echo ""
	@echo "Cross-compilation:"
	@echo "  make build-linux    - Build for Linux (amd64, arm64)"
	@echo "  make build-windows  - Build for Windows (amd64)"
	@echo "  make build-macos    - Build for macOS (Intel, Apple Silicon)"
	@echo "  make build-all      - Build for all platforms"
