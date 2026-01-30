#!/bin/bash
# Build macOS .app bundle for Claude Usage
# Usage: ./scripts/build-macos-app.sh [version]

set -e

VERSION="${1:-dev}"
APP_NAME="Claude Usage"
BUNDLE_NAME="Claude Usage.app"
BUNDLE_ID="com.github.utajum.claude-usage"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
DIST_DIR="$PROJECT_DIR/dist"
ASSETS_DIR="$PROJECT_DIR/assets/macos"

echo "Building $APP_NAME.app (version $VERSION)..."

# Create app bundle structure
APP_DIR="$DIST_DIR/$BUNDLE_NAME"
CONTENTS_DIR="$APP_DIR/Contents"
MACOS_DIR="$CONTENTS_DIR/MacOS"
RESOURCES_DIR="$CONTENTS_DIR/Resources"

rm -rf "$APP_DIR"
mkdir -p "$MACOS_DIR" "$RESOURCES_DIR"

# Build universal binary (if both architectures available)
ARM64_BIN="$DIST_DIR/claude-usage-darwin-arm64"
AMD64_BIN="$DIST_DIR/claude-usage-darwin-amd64"
UNIVERSAL_BIN="$MACOS_DIR/claude-usage"

if [ -f "$ARM64_BIN" ] && [ -f "$AMD64_BIN" ]; then
    echo "Creating universal binary..."
    lipo -create -output "$UNIVERSAL_BIN" "$ARM64_BIN" "$AMD64_BIN"
elif [ -f "$ARM64_BIN" ]; then
    echo "Using ARM64 binary..."
    cp "$ARM64_BIN" "$UNIVERSAL_BIN"
elif [ -f "$AMD64_BIN" ]; then
    echo "Using AMD64 binary..."
    cp "$AMD64_BIN" "$UNIVERSAL_BIN"
else
    echo "Error: No binary found. Run 'make build-macos' first."
    exit 1
fi

chmod +x "$UNIVERSAL_BIN"

# Copy Info.plist and update version
echo "Creating Info.plist..."
sed -e "s/1.0.0/$VERSION/g" "$ASSETS_DIR/Info.plist" > "$CONTENTS_DIR/Info.plist"

# Create .icns file if iconutil is available and iconset exists
ICONSET_DIR="$ASSETS_DIR/AppIcon.iconset"
ICNS_FILE="$RESOURCES_DIR/AppIcon.icns"

if [ -d "$ICONSET_DIR" ]; then
    if command -v iconutil &> /dev/null; then
        echo "Creating AppIcon.icns..."
        iconutil -c icns "$ICONSET_DIR" -o "$ICNS_FILE"
    elif [ -f "$ASSETS_DIR/AppIcon.icns" ]; then
        echo "Copying existing AppIcon.icns..."
        cp "$ASSETS_DIR/AppIcon.icns" "$ICNS_FILE"
    else
        echo "Warning: iconutil not available. App bundle will not have an icon."
        echo "To add icon manually: iconutil -c icns $ICONSET_DIR -o $ICNS_FILE"
    fi
fi

# Create PkgInfo
echo "APPL????" > "$CONTENTS_DIR/PkgInfo"

echo ""
echo "App bundle created: $APP_DIR"
echo ""

# Create DMG if hdiutil is available
if command -v hdiutil &> /dev/null; then
    echo "Creating DMG..."
    DMG_NAME="Claude-Usage-$VERSION.dmg"
    DMG_PATH="$DIST_DIR/$DMG_NAME"
    
    # Create temporary DMG directory
    DMG_TEMP="$DIST_DIR/dmg-temp"
    rm -rf "$DMG_TEMP"
    mkdir -p "$DMG_TEMP"
    cp -r "$APP_DIR" "$DMG_TEMP/"
    
    # Create symlink to Applications
    ln -s /Applications "$DMG_TEMP/Applications"
    
    # Create DMG
    rm -f "$DMG_PATH"
    hdiutil create -volname "$APP_NAME" -srcfolder "$DMG_TEMP" -ov -format UDZO "$DMG_PATH"
    
    # Cleanup
    rm -rf "$DMG_TEMP"
    
    echo "DMG created: $DMG_PATH"
fi

echo ""
echo "Done!"
