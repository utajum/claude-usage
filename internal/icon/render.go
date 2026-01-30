package icon

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
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

// RenderNeonOrbWithText creates a chip icon with percentage text
func RenderNeonOrbWithText(c color.RGBA, size int, percentage int) ([]byte, error) {
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

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
