#!/bin/bash
# Binary size comparison tool for claude-usage
# Tests different build configurations and reports sizes

set -e

BINARY="claude-usage"
echo "=== Binary Size Testing ==="
echo ""

# Clean previous test builds
echo "Cleaning previous test builds..."
rm -f test-* 2>/dev/null || true

echo ""
echo "Building different configurations..."
echo ""

# Test 1: Debug build (with symbols)
echo "1. Building with debug symbols..."
go build -ldflags "-X main.Version=test" -o test-debug ./cmd/claude-usage
SIZE_DEBUG=$(stat -c%s test-debug 2>/dev/null || stat -f%z test-debug)
SIZE_DEBUG_H=$(ls -lh test-debug | awk '{print $5}')
echo "   Size: $SIZE_DEBUG_H"

# Test 2: Stripped build
echo ""
echo "2. Building stripped (no debug)..."
go build -ldflags "-s -w -X main.Version=test" -trimpath -o test-stripped ./cmd/claude-usage
SIZE_STRIPPED=$(stat -c%s test-stripped 2>/dev/null || stat -f%z test-stripped)
SIZE_STRIPPED_H=$(ls -lh test-stripped | awk '{print $5}')
echo "   Size: $SIZE_STRIPPED_H"
REDUCTION_1=$(awk "BEGIN {printf \"%.1f\", (1 - $SIZE_STRIPPED / $SIZE_DEBUG) * 100}")
echo "   Reduction: ${REDUCTION_1}%"

# Test 3: Stripped + UPX
if command -v upx >/dev/null 2>&1; then
    echo ""
    echo "3. Building stripped + UPX..."
    cp test-stripped test-upx
    upx --best --lzma test-upx >/dev/null 2>&1
    SIZE_UPX=$(stat -c%s test-upx 2>/dev/null || stat -f%z test-upx)
    SIZE_UPX_H=$(ls -lh test-upx | awk '{print $5}')
    echo "   Size: $SIZE_UPX_H"
    REDUCTION_2=$(awk "BEGIN {printf \"%.1f\", (1 - $SIZE_UPX / $SIZE_DEBUG) * 100}")
    echo "   Reduction: ${REDUCTION_2}%"
else
    echo ""
    echo "UPX not found. Install with:"
    echo "   Linux:   sudo apt install upx-ucl"
    echo "   macOS:   brew install upx"
    echo "   Windows: choco install upx"
fi

# Summary
echo ""
echo "=== Summary ==="
echo ""
printf "%-20s %10s %12s\n" "Build Type" "Size" "Reduction"
printf "%-20s %10s %12s\n" "----------" "----" "---------"
printf "%-20s %10s %12s\n" "Debug (symbols)" "$SIZE_DEBUG_H" "-"
printf "%-20s %10s %12s\n" "Stripped" "$SIZE_STRIPPED_H" "${REDUCTION_1}%"
if command -v upx >/dev/null 2>&1; then
    printf "%-20s %10s %12s\n" "Stripped + UPX" "$SIZE_UPX_H" "${REDUCTION_2}%"
fi

echo ""
echo "Test complete."
echo ""
echo "Recommended builds:"
echo "  - Production: make build       (stripped + UPX)"
echo "  - Quick dev:  make build-fast  (stripped only)"
echo "  - Debugging:  make build-debug (with symbols)"

# Cleanup
echo ""
read -p "Delete test binaries? [Y/n] " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]] || [[ -z $REPLY ]]; then
    rm -f test-debug test-stripped test-upx
    echo "Cleaned up test binaries."
fi
