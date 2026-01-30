package icon

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
)

// ICO file format structures
// Reference: https://en.wikipedia.org/wiki/ICO_(file_format)

// iconDir is the ICO file header
type iconDir struct {
	Reserved uint16 // Must be 0
	Type     uint16 // 1 for icon, 2 for cursor
	Count    uint16 // Number of images
}

// iconDirEntry describes one image in the ICO file
type iconDirEntry struct {
	Width       uint8  // Width in pixels (0 means 256)
	Height      uint8  // Height in pixels (0 means 256)
	ColorCount  uint8  // Number of colors (0 if >= 8bpp)
	Reserved    uint8  // Must be 0
	Planes      uint16 // Color planes (should be 1)
	BitCount    uint16 // Bits per pixel
	BytesInRes  uint32 // Size of image data
	ImageOffset uint32 // Offset to image data from beginning of file
}

// bitmapInfoHeader is the DIB header for the image data
type bitmapInfoHeader struct {
	Size          uint32
	Width         int32
	Height        int32 // Height * 2 for ICO (includes AND mask)
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}

// EncodeICO encodes an RGBA image to ICO format bytes
func EncodeICO(img *image.RGBA) ([]byte, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Calculate sizes
	// Image data: BGRA pixels (4 bytes per pixel)
	imageDataSize := width * height * 4
	// AND mask: 1 bit per pixel, padded to 4-byte boundary per row
	andMaskRowSize := ((width + 31) / 32) * 4
	andMaskSize := andMaskRowSize * height
	// Total bitmap data size
	bitmapDataSize := 40 + imageDataSize + andMaskSize // 40 = sizeof(bitmapInfoHeader)

	// Create buffer
	buf := new(bytes.Buffer)

	// Write ICONDIR header
	dir := iconDir{
		Reserved: 0,
		Type:     1, // Icon
		Count:    1, // One image
	}
	binary.Write(buf, binary.LittleEndian, dir)

	// Write ICONDIRENTRY
	entry := iconDirEntry{
		Width:       uint8(width),
		Height:      uint8(height),
		ColorCount:  0, // True color
		Reserved:    0,
		Planes:      1,
		BitCount:    32, // 32-bit BGRA
		BytesInRes:  uint32(bitmapDataSize),
		ImageOffset: 6 + 16, // sizeof(iconDir) + sizeof(iconDirEntry)
	}
	// Handle 256x256 case
	if width >= 256 {
		entry.Width = 0
	}
	if height >= 256 {
		entry.Height = 0
	}
	binary.Write(buf, binary.LittleEndian, entry)

	// Write BITMAPINFOHEADER
	bih := bitmapInfoHeader{
		Size:          40,
		Width:         int32(width),
		Height:        int32(height * 2), // Height * 2 for ICO format (image + AND mask)
		Planes:        1,
		BitCount:      32,
		Compression:   0, // BI_RGB
		SizeImage:     uint32(imageDataSize + andMaskSize),
		XPelsPerMeter: 0,
		YPelsPerMeter: 0,
		ClrUsed:       0,
		ClrImportant:  0,
	}
	binary.Write(buf, binary.LittleEndian, bih)

	// Write pixel data (BGRA, bottom-to-top)
	// ICO stores images upside down (bottom row first)
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			c := img.RGBAAt(x, y)
			// Write BGRA (not RGBA)
			buf.WriteByte(c.B)
			buf.WriteByte(c.G)
			buf.WriteByte(c.R)
			buf.WriteByte(c.A)
		}
	}

	// Write AND mask (all zeros = fully opaque, since we use alpha channel)
	// The AND mask is 1-bit per pixel, but we're using 32-bit with alpha,
	// so we just write zeros
	andMask := make([]byte, andMaskSize)
	buf.Write(andMask)

	return buf.Bytes(), nil
}

// EncodeMultiResolutionICO creates an ICO with multiple resolutions
func EncodeMultiResolutionICO(images []*image.RGBA) ([]byte, error) {
	if len(images) == 0 {
		return nil, nil
	}

	// Calculate total size and offsets
	headerSize := 6 + (16 * len(images)) // iconDir + entries

	// Calculate bitmap data for each image
	type imageData struct {
		data   []byte
		width  int
		height int
	}
	imageDatas := make([]imageData, len(images))

	currentOffset := uint32(headerSize)
	for i, img := range images {
		bounds := img.Bounds()
		width := bounds.Dx()
		height := bounds.Dy()

		// Encode single image bitmap data
		bitmapData := encodeBitmapData(img)
		imageDatas[i] = imageData{
			data:   bitmapData,
			width:  width,
			height: height,
		}
	}

	// Create buffer
	buf := new(bytes.Buffer)

	// Write ICONDIR header
	dir := iconDir{
		Reserved: 0,
		Type:     1,
		Count:    uint16(len(images)),
	}
	binary.Write(buf, binary.LittleEndian, dir)

	// Write ICONDIRENTRY for each image
	for i, imgData := range imageDatas {
		entry := iconDirEntry{
			Width:       uint8(imgData.width),
			Height:      uint8(imgData.height),
			ColorCount:  0,
			Reserved:    0,
			Planes:      1,
			BitCount:    32,
			BytesInRes:  uint32(len(imgData.data)),
			ImageOffset: currentOffset,
		}
		if imgData.width >= 256 {
			entry.Width = 0
		}
		if imgData.height >= 256 {
			entry.Height = 0
		}
		binary.Write(buf, binary.LittleEndian, entry)
		currentOffset += uint32(len(imageDatas[i].data))
	}

	// Write bitmap data for each image
	for _, imgData := range imageDatas {
		buf.Write(imgData.data)
	}

	return buf.Bytes(), nil
}

// encodeBitmapData encodes a single image's bitmap data (header + pixels + AND mask)
func encodeBitmapData(img *image.RGBA) []byte {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	imageDataSize := width * height * 4
	andMaskRowSize := ((width + 31) / 32) * 4
	andMaskSize := andMaskRowSize * height

	buf := new(bytes.Buffer)

	// Write BITMAPINFOHEADER
	bih := bitmapInfoHeader{
		Size:          40,
		Width:         int32(width),
		Height:        int32(height * 2),
		Planes:        1,
		BitCount:      32,
		Compression:   0,
		SizeImage:     uint32(imageDataSize + andMaskSize),
		XPelsPerMeter: 0,
		YPelsPerMeter: 0,
		ClrUsed:       0,
		ClrImportant:  0,
	}
	binary.Write(buf, binary.LittleEndian, bih)

	// Write pixel data (BGRA, bottom-to-top)
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			c := img.RGBAAt(x, y)
			buf.WriteByte(c.B)
			buf.WriteByte(c.G)
			buf.WriteByte(c.R)
			buf.WriteByte(c.A)
		}
	}

	// Write AND mask (zeros)
	andMask := make([]byte, andMaskSize)
	buf.Write(andMask)

	return buf.Bytes()
}

// ResizeImage creates a simple nearest-neighbor scaled version of an image
func ResizeImage(src *image.RGBA, newWidth, newHeight int) *image.RGBA {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := x * srcWidth / newWidth
			srcY := y * srcHeight / newHeight
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}

	return dst
}

// CreateImageRGBA creates a new RGBA image of the given size filled with transparent pixels
func CreateImageRGBA(width, height int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, width, height))
}

// SetPixel sets a pixel in an RGBA image
func SetPixel(img *image.RGBA, x, y int, c color.RGBA) {
	if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
		img.SetRGBA(x, y, c)
	}
}
