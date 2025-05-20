package images

import (
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"msxconverter/decoders"
)

// DecodeSTP decodes Dynamic Publisher STP stamp files.
// STP files start with width and height as little endian words
// followed by packed 2-bit pixel values (4 pixels per byte).
// The palette is monochrome.
func DecodeSTP(data []byte, config decoders.Config) (decoders.DecoderResult, error) {
	if len(data) < 4 {
		return decoders.DecoderResult{}, errors.New("invalid STP file")
	}

	width := int(binary.LittleEndian.Uint16(data[0:2]))
	height := int(binary.LittleEndian.Uint16(data[2:4]))
	pixels := data[4:]

	expected := (width*height + 3) / 4
	if len(pixels) < expected {
		return decoders.DecoderResult{}, errors.New("incomplete STP pixel data")
	}

	palette := color.Palette{
		color.RGBA{0, 0, 0, 255},
		color.RGBA{255, 255, 255, 255},
	}

	img := image.NewPaletted(image.Rect(0, 0, width, height), palette)

	idx := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			byteIdx := idx / 4
			shift := 6 - 2*(idx%4)
			val := (pixels[byteIdx] >> shift) & 0x01
			img.SetColorIndex(x, y, 1-val)
			idx++
		}
	}

	// STP pixels are intended for a 512x212 screen where the horizontal
	// resolution is twice the vertical one.  To preserve the correct
	// aspect ratio, each line needs to be duplicated.
	doubledHeightImg := doubleHeightPaletted(width, height, img)

	var output image.Image
	if config.DoubleImageSize {
		// Double both dimensions after the aspect ratio fix.
		output = doubleSizePaletted(width, height*2, doubledHeightImg)
	} else {
		output = doubledHeightImg
	}

	return encodePNG(output)
}
