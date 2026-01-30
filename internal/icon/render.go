package icon

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"runtime"
)

const (
	IconSize = 22
)

// Colors - Neon violet/cyan theme
var (
	chipColor  = color.RGBA{R: 180, G: 0, B: 255, A: 255}   // Neon violet/purple
	cyanAccent = color.RGBA{R: 0, G: 255, B: 255, A: 255}   // Cyan accent
	whiteText  = color.RGBA{R: 255, G: 255, B: 255, A: 255} // White text
)

// Large bold pixel patterns for digits 0-9 (7x9 pixels)
var digitPatterns = map[rune][9][7]int{
	'0': {{0, 1, 1, 1, 1, 1, 0}, {1, 1, 1, 1, 1, 1, 1}, {1, 1, 0, 0, 0, 1, 1}, {1, 1, 0, 0, 0, 1, 1}, {1, 1, 0, 0, 0, 1, 1}, {1, 1, 0, 0, 0, 1, 1}, {1, 1, 0, 0, 0, 1, 1}, {1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 0}},
	'1': {{0, 0, 0, 1, 1, 0, 0}, {0, 0, 1, 1, 1, 0, 0}, {0, 1, 1, 1, 1, 0, 0}, {0, 0, 0, 1, 1, 0, 0}, {0, 0, 0, 1, 1, 0, 0}, {0, 0, 0, 1, 1, 0, 0}, {0, 0, 0, 1, 1, 0, 0}, {0, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 1}},
	'2': {{0, 1, 1, 1, 1, 1, 0}, {1, 1, 1, 1, 1, 1, 1}, {1, 1, 0, 0, 0, 1, 1}, {0, 0, 0, 0, 1, 1, 1}, {0, 0, 1, 1, 1, 1, 0}, {0, 1, 1, 1, 0, 0, 0}, {1, 1, 1, 0, 0, 0, 0}, {1, 1, 1, 1, 1, 1, 1}, {1, 1, 1, 1, 1, 1, 1}},
	'3': {{0, 1, 1, 1, 1, 1, 0}, {1, 1, 1, 1, 1, 1, 1}, {0, 0, 0, 0, 0, 1, 1}, {0, 0, 1, 1, 1, 1, 0}, {0, 0, 1, 1, 1, 1, 0}, {0, 0, 0, 0, 0, 1, 1}, {0, 0, 0, 0, 0, 1, 1}, {1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 0}},
	'4': {{0, 0, 0, 0, 1, 1, 1}, {0, 0, 0, 1, 1, 1, 1}, {0, 0, 1, 1, 0, 1, 1}, {0, 1, 1, 0, 0, 1, 1}, {1, 1, 0, 0, 0, 1, 1}, {1, 1, 1, 1, 1, 1, 1}, {1, 1, 1, 1, 1, 1, 1}, {0, 0, 0, 0, 0, 1, 1}, {0, 0, 0, 0, 0, 1, 1}},
	'5': {{1, 1, 1, 1, 1, 1, 1}, {1, 1, 1, 1, 1, 1, 1}, {1, 1, 0, 0, 0, 0, 0}, {1, 1, 1, 1, 1, 1, 0}, {1, 1, 1, 1, 1, 1, 1}, {0, 0, 0, 0, 0, 1, 1}, {0, 0, 0, 0, 0, 1, 1}, {1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 0}},
	'6': {{0, 1, 1, 1, 1, 1, 0}, {1, 1, 1, 1, 1, 1, 1}, {1, 1, 0, 0, 0, 0, 0}, {1, 1, 1, 1, 1, 1, 0}, {1, 1, 1, 1, 1, 1, 1}, {1, 1, 0, 0, 0, 1, 1}, {1, 1, 0, 0, 0, 1, 1}, {1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 0}},
	'7': {{1, 1, 1, 1, 1, 1, 1}, {1, 1, 1, 1, 1, 1, 1}, {0, 0, 0, 0, 0, 1, 1}, {0, 0, 0, 0, 1, 1, 0}, {0, 0, 0, 1, 1, 0, 0}, {0, 0, 1, 1, 0, 0, 0}, {0, 0, 1, 1, 0, 0, 0}, {0, 0, 1, 1, 0, 0, 0}, {0, 0, 1, 1, 0, 0, 0}},
	'8': {{0, 1, 1, 1, 1, 1, 0}, {1, 1, 1, 1, 1, 1, 1}, {1, 1, 0, 0, 0, 1, 1}, {0, 1, 1, 1, 1, 1, 0}, {0, 1, 1, 1, 1, 1, 0}, {1, 1, 0, 0, 0, 1, 1}, {1, 1, 0, 0, 0, 1, 1}, {1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 0}},
	'9': {{0, 1, 1, 1, 1, 1, 0}, {1, 1, 1, 1, 1, 1, 1}, {1, 1, 0, 0, 0, 1, 1}, {1, 1, 0, 0, 0, 1, 1}, {0, 1, 1, 1, 1, 1, 1}, {0, 0, 0, 0, 0, 1, 1}, {0, 0, 0, 0, 0, 1, 1}, {1, 1, 1, 1, 1, 1, 1}, {0, 1, 1, 1, 1, 1, 0}},
}

// drawChar draws a single character at the given position
func drawChar(img *image.RGBA, char rune, startX, startY int, textColor color.RGBA) {
	pattern, ok := digitPatterns[char]
	if !ok {
		return
	}
	for row := 0; row < 9; row++ {
		for col := 0; col < 7; col++ {
			if pattern[row][col] == 1 {
				x := startX + col
				y := startY + row
				if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
					img.SetRGBA(x, y, textColor)
				}
			}
		}
	}
}

// drawText draws digits centered at the given position
func drawText(img *image.RGBA, text string, centerX, centerY int, textColor color.RGBA) {
	charWidth := 8 // 7 pixels + 1 spacing
	charHeight := 9
	totalWidth := len(text)*charWidth - 1
	startX := centerX - totalWidth/2
	startY := centerY - charHeight/2
	for i, char := range text {
		drawChar(img, char, startX+i*charWidth, startY, textColor)
	}
}

// RenderChipImage creates the chip icon image (without encoding)
func RenderChipImage(c color.RGBA, size int, percentage int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Main body - neon violet chip
	for y := 2; y < size-2; y++ {
		for x := 2; x < size-2; x++ {
			img.SetRGBA(x, y, chipColor)
		}
	}

	// Corner cuts for chip look
	img.SetRGBA(2, 2, color.RGBA{0, 0, 0, 0})
	img.SetRGBA(size-3, 2, color.RGBA{0, 0, 0, 0})
	img.SetRGBA(2, size-3, color.RGBA{0, 0, 0, 0})
	img.SetRGBA(size-3, size-3, color.RGBA{0, 0, 0, 0})

	// Draw cyan pins on all sides
	// Top pins
	for i := 5; i < size-5; i += 3 {
		img.SetRGBA(i, 0, cyanAccent)
		img.SetRGBA(i, 1, cyanAccent)
		img.SetRGBA(i+1, 0, cyanAccent)
		img.SetRGBA(i+1, 1, cyanAccent)
	}
	// Bottom pins
	for i := 5; i < size-5; i += 3 {
		img.SetRGBA(i, size-1, cyanAccent)
		img.SetRGBA(i, size-2, cyanAccent)
		img.SetRGBA(i+1, size-1, cyanAccent)
		img.SetRGBA(i+1, size-2, cyanAccent)
	}
	// Left pins
	for i := 5; i < size-5; i += 3 {
		img.SetRGBA(0, i, cyanAccent)
		img.SetRGBA(1, i, cyanAccent)
		img.SetRGBA(0, i+1, cyanAccent)
		img.SetRGBA(1, i+1, cyanAccent)
	}
	// Right pins
	for i := 5; i < size-5; i += 3 {
		img.SetRGBA(size-1, i, cyanAccent)
		img.SetRGBA(size-2, i, cyanAccent)
		img.SetRGBA(size-1, i+1, cyanAccent)
		img.SetRGBA(size-2, i+1, cyanAccent)
	}

	// Format text - just the number
	var text string
	if percentage >= 100 {
		text = "99"
	} else if percentage < 0 {
		text = "0"
	} else {
		text = fmt.Sprintf("%d", percentage)
	}

	// Draw white text in center (bigger 7x9 font)
	drawText(img, text, size/2, size/2, whiteText)

	return img
}

// RenderChipImageNoText creates the chip icon image without percentage text (for app icon)
func RenderChipImageNoText(size int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Scale factor for larger sizes
	scale := size / 22
	if scale < 1 {
		scale = 1
	}

	borderWidth := 2 * scale
	pinStart := 5 * scale
	pinEnd := size - 5*scale
	pinSpacing := 3 * scale
	pinSize := 2 * scale

	// Main body - neon violet chip
	for y := borderWidth; y < size-borderWidth; y++ {
		for x := borderWidth; x < size-borderWidth; x++ {
			img.SetRGBA(x, y, chipColor)
		}
	}

	// Corner cuts for chip look
	for i := 0; i < scale; i++ {
		for j := 0; j < scale; j++ {
			img.SetRGBA(borderWidth+i, borderWidth+j, color.RGBA{0, 0, 0, 0})
			img.SetRGBA(size-borderWidth-1-i, borderWidth+j, color.RGBA{0, 0, 0, 0})
			img.SetRGBA(borderWidth+i, size-borderWidth-1-j, color.RGBA{0, 0, 0, 0})
			img.SetRGBA(size-borderWidth-1-i, size-borderWidth-1-j, color.RGBA{0, 0, 0, 0})
		}
	}

	// Draw cyan pins on all sides
	// Top and bottom pins
	for i := pinStart; i < pinEnd; i += pinSpacing {
		for px := 0; px < pinSize && i+px < size; px++ {
			for py := 0; py < pinSize; py++ {
				// Top
				if py < size {
					img.SetRGBA(i+px, py, cyanAccent)
				}
				// Bottom
				if size-1-py >= 0 {
					img.SetRGBA(i+px, size-1-py, cyanAccent)
				}
			}
		}
	}
	// Left and right pins
	for i := pinStart; i < pinEnd; i += pinSpacing {
		for py := 0; py < pinSize && i+py < size; py++ {
			for px := 0; px < pinSize; px++ {
				// Left
				if px < size {
					img.SetRGBA(px, i+py, cyanAccent)
				}
				// Right
				if size-1-px >= 0 {
					img.SetRGBA(size-1-px, i+py, cyanAccent)
				}
			}
		}
	}

	// Add a subtle "C" or circuit pattern in the center for the app icon
	// Draw a stylized circuit/connection pattern
	centerX := size / 2
	centerY := size / 2
	circuitColor := cyanAccent

	// Draw small circuit lines in center
	lineWidth := scale
	lineLen := size / 4

	// Horizontal line
	for x := centerX - lineLen/2; x <= centerX+lineLen/2; x++ {
		for w := 0; w < lineWidth; w++ {
			if x >= 0 && x < size && centerY+w >= 0 && centerY+w < size {
				img.SetRGBA(x, centerY+w, circuitColor)
			}
		}
	}

	// Vertical line
	for y := centerY - lineLen/2; y <= centerY+lineLen/2; y++ {
		for w := 0; w < lineWidth; w++ {
			if centerX+w >= 0 && centerX+w < size && y >= 0 && y < size {
				img.SetRGBA(centerX+w, y, circuitColor)
			}
		}
	}

	// Small dots at the ends
	dotSize := scale * 2
	dots := []struct{ x, y int }{
		{centerX - lineLen/2, centerY},
		{centerX + lineLen/2, centerY},
		{centerX, centerY - lineLen/2},
		{centerX, centerY + lineLen/2},
	}
	for _, dot := range dots {
		for dx := -dotSize / 2; dx <= dotSize/2; dx++ {
			for dy := -dotSize / 2; dy <= dotSize/2; dy++ {
				x, y := dot.x+dx, dot.y+dy
				if x >= 0 && x < size && y >= 0 && y < size {
					img.SetRGBA(x, y, circuitColor)
				}
			}
		}
	}

	return img
}

// EncodePNG encodes an image to PNG format
func EncodePNG(img *image.RGBA) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// RenderNeonOrbWithText creates a chip icon with percentage text
// Returns ICO format on Windows, PNG on other platforms
func RenderNeonOrbWithText(c color.RGBA, size int, percentage int) ([]byte, error) {
	img := RenderChipImage(c, size, percentage)

	// On Windows, return ICO format for system tray compatibility
	if runtime.GOOS == "windows" {
		return EncodeICO(img)
	}

	// On other platforms, return PNG
	return EncodePNG(img)
}

// RenderAppIcon creates an application icon (without percentage text)
// Returns the appropriate format for the current platform
func RenderAppIcon(size int) ([]byte, error) {
	img := RenderChipImageNoText(size)

	if runtime.GOOS == "windows" {
		return EncodeICO(img)
	}

	return EncodePNG(img)
}

// RenderAppIconPNG always returns PNG format (for asset generation)
func RenderAppIconPNG(size int) ([]byte, error) {
	img := RenderChipImageNoText(size)
	return EncodePNG(img)
}

// RenderAppIconICO always returns ICO format (for asset generation)
func RenderAppIconICO(size int) ([]byte, error) {
	img := RenderChipImageNoText(size)
	return EncodeICO(img)
}

// RenderMultiResolutionICO creates an ICO with multiple resolutions for Windows
func RenderMultiResolutionICO(sizes []int) ([]byte, error) {
	images := make([]*image.RGBA, len(sizes))
	for i, size := range sizes {
		images[i] = RenderChipImageNoText(size)
	}
	return EncodeMultiResolutionICO(images)
}
