#!/bin/bash
# Claude Usage - Linux Installation Script
# Usage: curl -sL https://raw.githubusercontent.com/utajum/claude-usage/master/scripts/install-linux.sh | bash
#
# Options:
#   --no-autostart    Don't enable autostart
#   --no-desktop      Don't install desktop entry
#   --uninstall       Remove claude-usage

set -e

APP_NAME="claude-usage"
REPO="utajum/claude-usage"
INSTALL_DIR="$HOME/.local/bin"
DESKTOP_DIR="$HOME/.local/share/applications"
AUTOSTART_DIR="$HOME/.config/autostart"
ICONS_DIR="$HOME/.local/share/icons/hicolor"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Parse arguments
NO_AUTOSTART=false
NO_DESKTOP=false
UNINSTALL=false

for arg in "$@"; do
    case $arg in
        --no-autostart)
            NO_AUTOSTART=true
            ;;
        --no-desktop)
            NO_DESKTOP=true
            ;;
        --uninstall)
            UNINSTALL=true
            ;;
        --help|-h)
            echo "Claude Usage - Linux Installer"
            echo ""
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --no-autostart    Don't enable autostart on login"
            echo "  --no-desktop      Don't install desktop entry (app menu)"
            echo "  --uninstall       Remove claude-usage completely"
            echo "  --help, -h        Show this help"
            exit 0
            ;;
    esac
done

print_status() {
    echo -e "${CYAN}[*]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[+]${NC} $1"
}

print_error() {
    echo -e "${RED}[-]${NC} $1"
}

detect_arch() {
    ARCH=$(uname -m)
    case $ARCH in
        x86_64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        *)
            print_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
}

get_latest_version() {
    curl -sI "https://github.com/$REPO/releases/latest" | grep -i "^location:" | sed 's/.*tag\///' | tr -d '\r\n'
}

uninstall() {
    print_status "Uninstalling Claude Usage..."
    
    # Stop running instance
    pkill -f "$APP_NAME" 2>/dev/null || true
    
    # Remove binary
    if [ -f "$INSTALL_DIR/$APP_NAME" ]; then
        rm -f "$INSTALL_DIR/$APP_NAME"
        print_success "Removed binary"
    fi
    
    # Remove desktop entry
    if [ -f "$DESKTOP_DIR/$APP_NAME.desktop" ]; then
        rm -f "$DESKTOP_DIR/$APP_NAME.desktop"
        print_success "Removed desktop entry"
    fi
    
    # Remove autostart entry
    if [ -f "$AUTOSTART_DIR/$APP_NAME.desktop" ]; then
        rm -f "$AUTOSTART_DIR/$APP_NAME.desktop"
        print_success "Removed autostart entry"
    fi
    
    # Remove icons
    for size in 16 24 32 48 64 128 256; do
        rm -f "$ICONS_DIR/${size}x${size}/apps/$APP_NAME.png" 2>/dev/null || true
    done
    print_success "Removed icons"
    
    # Update icon cache
    gtk-update-icon-cache "$ICONS_DIR" 2>/dev/null || true
    update-desktop-database "$DESKTOP_DIR" 2>/dev/null || true
    
    print_success "Uninstallation complete!"
    exit 0
}

install() {
    echo ""
    echo "  ░█▀▀░█░░░█▀█░█░█░█▀▄░█▀▀░░░█░█░█▀▀░█▀█░█▀▀░█▀▀"
    echo "  ░█░░░█░░░█▀█░█░█░█░█░█▀▀░░░█░█░▀▀█░█▀█░█░█░█▀▀"
    echo "  ░▀▀▀░▀▀▀░▀░▀░▀▀▀░▀▀░░▀▀▀░░░▀▀▀░▀▀▀░▀░▀░▀▀▀░▀▀▀"
    echo ""
    
    # Detect architecture
    ARCH=$(detect_arch)
    print_status "Detected architecture: $ARCH"
    
    # Get latest version
    print_status "Fetching latest version..."
    VERSION=$(get_latest_version)
    if [ -z "$VERSION" ]; then
        print_error "Could not determine latest version"
        exit 1
    fi
    print_status "Latest version: $VERSION"
    
    # Create directories
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$DESKTOP_DIR"
    mkdir -p "$AUTOSTART_DIR"
    
    # Download binary
    BINARY_URL="https://github.com/$REPO/releases/download/$VERSION/$APP_NAME-linux-$ARCH"
    print_status "Downloading $APP_NAME-linux-$ARCH..."
    
    if ! curl -sL "$BINARY_URL" -o "$INSTALL_DIR/$APP_NAME"; then
        print_error "Failed to download binary"
        exit 1
    fi
    chmod +x "$INSTALL_DIR/$APP_NAME"
    print_success "Installed binary to $INSTALL_DIR/$APP_NAME"
    
    # Download and install icons
    if [ "$NO_DESKTOP" = false ]; then
        print_status "Downloading icons..."
        
        # We'll generate simple icons using the binary itself or download from repo
        # For now, create icon directories
        for size in 16 24 32 48 64 128 256; do
            mkdir -p "$ICONS_DIR/${size}x${size}/apps"
        done
        
        # Download icon from repo (we'll use the 256px one and let the system scale)
        ICON_URL="https://raw.githubusercontent.com/$REPO/master/assets/linux/claude-usage.png"
        
        # Try to download pre-generated icons, fall back to generating them
        TMP_DIR=$(mktemp -d)
        
        # Download the generate-icons script and run it, or just use a placeholder
        for size in 16 24 32 48 64 128 256; do
            ICON_FILE="$ICONS_DIR/${size}x${size}/apps/$APP_NAME.png"
            ICON_DL_URL="https://raw.githubusercontent.com/$REPO/master/assets/linux/claude-usage-${size}.png"
            
            # Try to download, ignore if fails (icons are generated in CI)
            curl -sL "$ICON_DL_URL" -o "$ICON_FILE" 2>/dev/null || true
        done
        
        # If we didn't get icons, generate a simple placeholder
        if [ ! -s "$ICONS_DIR/48x48/apps/$APP_NAME.png" ]; then
            print_status "Generating icons locally..."
            # Run the binary with a special flag to generate icons, or use ImageMagick if available
            if command -v convert &> /dev/null; then
                # Create a simple violet square icon with ImageMagick
                convert -size 256x256 xc:'#B400FF' -fill '#00FFFF' \
                    -draw "rectangle 0,0 10,256 rectangle 246,0 256,256 rectangle 0,0 256,10 rectangle 0,246 256,256" \
                    "$TMP_DIR/icon-256.png" 2>/dev/null || true
                
                for size in 16 24 32 48 64 128 256; do
                    convert "$TMP_DIR/icon-256.png" -resize ${size}x${size} \
                        "$ICONS_DIR/${size}x${size}/apps/$APP_NAME.png" 2>/dev/null || true
                done
            fi
        fi
        
        rm -rf "$TMP_DIR"
        
        # Update icon cache
        gtk-update-icon-cache "$ICONS_DIR" 2>/dev/null || true
        print_success "Installed icons"
    fi
    
    # Create desktop entry
    if [ "$NO_DESKTOP" = false ]; then
        print_status "Creating desktop entry..."
        cat > "$DESKTOP_DIR/$APP_NAME.desktop" << EOF
[Desktop Entry]
Type=Application
Name=Claude Usage
GenericName=Token Usage Monitor
Comment=Monitor Claude Code API token usage in system tray
Exec=$INSTALL_DIR/$APP_NAME
Icon=$APP_NAME
Terminal=false
Categories=Utility;Monitor;System;
Keywords=claude;usage;tokens;api;anthropic;monitor;tray;
StartupNotify=false
StartupWMClass=claude-usage
EOF
        update-desktop-database "$DESKTOP_DIR" 2>/dev/null || true
        print_success "Created desktop entry (find 'Claude Usage' in your apps menu)"
    fi
    
    # Create autostart entry
    if [ "$NO_AUTOSTART" = false ]; then
        print_status "Enabling autostart..."
        cat > "$AUTOSTART_DIR/$APP_NAME.desktop" << EOF
[Desktop Entry]
Type=Application
Name=Claude Usage
Comment=Monitor Claude Code API token usage in system tray
Exec=$INSTALL_DIR/$APP_NAME
Icon=$APP_NAME
Terminal=false
Categories=Utility;
StartupNotify=false
X-GNOME-Autostart-enabled=true
EOF
        print_success "Enabled autostart on login"
    fi
    
    # Check if ~/.local/bin is in PATH
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        echo ""
        print_status "NOTE: $INSTALL_DIR is not in your PATH"
        echo "      Add this to your ~/.bashrc or ~/.zshrc:"
        echo ""
        echo "      export PATH=\"\$HOME/.local/bin:\$PATH\""
        echo ""
    fi
    
    echo ""
    print_success "Installation complete!"
    echo ""
    echo "  To start now:     $APP_NAME"
    echo "  To uninstall:     $0 --uninstall"
    echo ""
    
    # Offer to start the app
    read -p "Start Claude Usage now? [Y/n] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]] || [[ -z $REPLY ]]; then
        print_status "Starting Claude Usage..."
        nohup "$INSTALL_DIR/$APP_NAME" > /dev/null 2>&1 &
        print_success "Claude Usage is running! Look for it in your system tray."
    fi
}

# Main
if [ "$UNINSTALL" = true ]; then
    uninstall
else
    install
fi
