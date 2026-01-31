// generate-icons generates application icons for all platforms
// Run with: go run scripts/generate-icons/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"claude-usage/internal/icon"
)

func main() {
	// Get the project root (assuming we're run from project root)
	assetsDir := "assets"

	// Create directories
	dirs := []string{
		filepath.Join(assetsDir, "linux"),
		filepath.Join(assetsDir, "macos"),
		filepath.Join(assetsDir, "windows"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", dir, err)
			os.Exit(1)
		}
	}

	// Generate Linux icons (PNG)
	linuxSizes := []int{16, 24, 32, 48, 64, 128, 256}
	for _, size := range linuxSizes {
		data, err := icon.RenderAppIconPNG(size)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating Linux icon %dx%d: %v\n", size, size, err)
			os.Exit(1)
		}
		path := filepath.Join(assetsDir, "linux", fmt.Sprintf("claude-usage-%d.png", size))
		if err := os.WriteFile(path, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("Generated: %s\n", path)
	}

	// Generate main icon for Linux desktop file
	data, err := icon.RenderAppIconPNG(256)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating main icon: %v\n", err)
		os.Exit(1)
	}
	mainIconPath := filepath.Join(assetsDir, "linux", "claude-usage.png")
	if err := os.WriteFile(mainIconPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", mainIconPath, err)
		os.Exit(1)
	}
	fmt.Printf("Generated: %s\n", mainIconPath)

	// Generate Windows ICO (multi-resolution)
	// Windows ICO typically includes: 16, 32, 48, 256
	windowsSizes := []int{16, 32, 48, 256}
	icoData, err := icon.RenderMultiResolutionICO(windowsSizes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating Windows ICO: %v\n", err)
		os.Exit(1)
	}
	icoPath := filepath.Join(assetsDir, "windows", "claude-usage.ico")
	if err := os.WriteFile(icoPath, icoData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", icoPath, err)
		os.Exit(1)
	}
	fmt.Printf("Generated: %s\n", icoPath)

	// Generate macOS PNG icons (for iconutil to create .icns)
	// macOS requires specific sizes: 16, 32, 64, 128, 256, 512, 1024
	// with @2x versions for Retina
	macosSizes := []struct {
		size int
		name string
	}{
		{16, "icon_16x16.png"},
		{32, "icon_16x16@2x.png"},
		{32, "icon_32x32.png"},
		{64, "icon_32x32@2x.png"},
		{128, "icon_128x128.png"},
		{256, "icon_128x128@2x.png"},
		{256, "icon_256x256.png"},
		{512, "icon_256x256@2x.png"},
		{512, "icon_512x512.png"},
		{1024, "icon_512x512@2x.png"},
	}

	iconsetDir := filepath.Join(assetsDir, "macos", "AppIcon.iconset")
	if err := os.MkdirAll(iconsetDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating iconset directory: %v\n", err)
		os.Exit(1)
	}

	for _, s := range macosSizes {
		data, err := icon.RenderAppIconPNG(s.size)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating macOS icon %s: %v\n", s.name, err)
			os.Exit(1)
		}
		path := filepath.Join(iconsetDir, s.name)
		if err := os.WriteFile(path, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("Generated: %s\n", path)
	}

	fmt.Println("\nIcon generation complete!")
	fmt.Println("\nTo create macOS .icns file, run on macOS:")
	fmt.Println("  iconutil -c icns assets/macos/AppIcon.iconset -o assets/macos/AppIcon.icns")
}
