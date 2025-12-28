package images

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"msxconverter/decoders/types"
)

// Lookup tables for color conversion
var (
	color3bitsLookupTable = [8]uint8{0x00, 0x24, 0x49, 0x6d, 0x92, 0xb6, 0xdb, 0xff}
	color2bitsLookupTable = [4]uint8{0x00, 0x55, 0xaa, 0xff}
	color5bitsLookupTable = [32]uint8{
		0, 8, 16, 24, 33, 41, 49, 57,
		66, 74, 82, 90, 99, 107, 115, 123,
		132, 140, 148, 156, 165, 173, 181, 189,
		198, 206, 214, 222, 231, 239, 247, 255,
	}
	defaultPalette = []byte{
		0x00, 0x00, 0x00, 0x00, 0x11, 0x06, 0x33, 0x07,
		0x17, 0x01, 0x27, 0x03, 0x51, 0x01, 0x27, 0x06,
		0x71, 0x01, 0x73, 0x03, 0x61, 0x06, 0x64, 0x06,
		0x11, 0x04, 0x65, 0x02, 0x55, 0x05, 0x77, 0x07,
	}
)

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func encodePNG(img image.Image) (types.DecoderResult, error) {
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		return types.DecoderResult{}, err
	}
	return types.DecoderResult{
		Buffer: bytes.NewBuffer(buffer.Bytes()),
		IsText: false,
	}, nil
}

func doubleSize(width, height int, img *image.RGBA) image.Image {
	doubleImg := image.NewRGBA(image.Rect(0, 0, width*2, height*2))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			doubleImg.Set(x*2, y*2, c)
			doubleImg.Set(x*2+1, y*2, c)
			doubleImg.Set(x*2, y*2+1, c)
			doubleImg.Set(x*2+1, y*2+1, c)
		}
	}
	return doubleImg
}

func doubleSizePaletted(width, height int, img *image.Paletted) image.Image {
	doubleImg := image.NewPaletted(image.Rect(0, 0, width*2, height*2), img.Palette)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			doubleImg.Set(x*2, y*2, c)
			doubleImg.Set(x*2+1, y*2, c)
			doubleImg.Set(x*2, y*2+1, c)
			doubleImg.Set(x*2+1, y*2+1, c)
		}
	}
	return doubleImg
}

// doubleHeightPaletted duplicates each row of a paletted image to correct the
// aspect ratio for formats that use rectangular pixels.
func doubleHeightPaletted(width, height int, img *image.Paletted) *image.Paletted {
	doubleImg := image.NewPaletted(image.Rect(0, 0, width, height*2), img.Palette)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.ColorIndexAt(x, y)
			doubleImg.SetColorIndex(x, y*2, c)
			doubleImg.SetColorIndex(x, y*2+1, c)
		}
	}
	return doubleImg
}

func getPalette(data []byte, paletteOffset int) color.Palette {
	paletteData := data[paletteOffset : paletteOffset+32]
	palette := make(color.Palette, 16)
	for i := 0; i < 16; i++ {
		raw := binary.LittleEndian.Uint16(paletteData[i*2 : i*2+2])
		r := color3bitsLookupTable[(raw>>4)&0b111]
		g := color3bitsLookupTable[(raw>>8)&0b111]
		b := color3bitsLookupTable[(raw)&0b111]
		palette[i] = color.RGBA{R: r, G: g, B: b, A: 255}
	}
	return palette
}

func calculateHeight(endAddress uint16, width int) int {
	endLine := endAddress / uint16(width)
	if endLine <= ScreenHeightMSX1 {
		return ScreenHeightMSX1
	}
	return ScreenHeightMSX2
}

type screenHeader struct {
	beginAddress uint16
	endAddress   uint16
	pixels       []byte
}

func parseScreenHeader(data []byte) (screenHeader, error) {
	if len(data) < 7 {
		return screenHeader{}, fmt.Errorf("screen data too small: got %d bytes, need at least 7", len(data))
	}
	return screenHeader{
		beginAddress: binary.LittleEndian.Uint16(data[1:3]),
		endAddress:   binary.LittleEndian.Uint16(data[3:5]),
		pixels:       data[7:],
	}, nil
}

// ResolvePalette resolves the palette to use, checking embedded palette, extraData, or default
func ResolvePalette(pixels []byte, beginAddress, endAddress uint16, paletteOffset int, extraData []byte) color.Palette {
	if paletteOffset > 0 {
		if endAddress >= uint16(paletteOffset) {
			return getPalette(pixels, paletteOffset-int(beginAddress))
		}
		if extraData != nil {
			return getPalette(extraData, 0)
		}
		return getPalette(defaultPalette, 0)
	}
	if extraData != nil {
		return getPalette(extraData, 0)
	}
	return getPalette(defaultPalette, 0)
}

// GetDefaultPalette returns the default MSX SCREEN 5 palette
func GetDefaultPalette() color.Palette {
	return getPalette(defaultPalette, 0)
}

