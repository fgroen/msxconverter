package images

import (
	"encoding/binary"
	"image"
	"image/color"
	"msxconverter/decoders"
)

// decodeScreenNibbles decodes screen data with nibbles to an image.
func decodeScreenNibbles(data []byte, config decoders.Config, width, paletteOffset int) (decoders.DecoderResult, error) {
	beginAddress := binary.LittleEndian.Uint16(data[1:3])
	endAddress := binary.LittleEndian.Uint16(data[3:5])
	height := calculateHeight(endAddress, width/2)

	// Skip the file header and read palette
	pixels := data[7:]
	var palette color.Palette
	if endAddress >= uint16(paletteOffset) {
		palette = getPalette(pixels, paletteOffset-int(beginAddress))
	} else if config.ExtraData != nil {
		palette = getPalette(config.ExtraData, 0)
	} else {
		palette = getPalette(defaultPalette, 0)
	}

	img := image.NewPaletted(image.Rect(0, 0, width, height), palette)
	if uint16(height*(width/2)) < endAddress {
		endAddress = uint16(height * (width / 2))
	}

	for addr := beginAddress; addr < endAddress; addr++ {
		y := int(addr / uint16(width/2))
		x := int(addr % uint16(width/2))
		byteVal := pixels[addr-beginAddress]

		img.SetColorIndex(x*2, y, byteVal>>4)
		img.SetColorIndex(x*2+1, y, byteVal&0x0F)
	}

	var outputImage image.Image
	if config.DoubleImageSize {
		outputImage = doubleSizePaletted(width, height, img)
	} else {
		outputImage = img
	}
	return encodePNG(outputImage)
}

// decodeScreen decodes screen data to an image.
func decodeScreen(data []byte, config decoders.Config, width int) (decoders.DecoderResult, error) {
	beginAddress := binary.LittleEndian.Uint16(data[1:3])
	endAddress := binary.LittleEndian.Uint16(data[3:5])
	height := calculateHeight(endAddress, width)

	// Skip the file header and read pixels
	pixels := data[7:]

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	if uint16(height*width) < endAddress {
		endAddress = uint16(height * width)
	}

	for addr := beginAddress; addr < endAddress; addr++ {
		y := int(addr / uint16(width))
		x := int(addr % uint16(width))
		byteVal := pixels[addr-beginAddress]

		r := color3bitsLookupTable[(byteVal>>2)&0b111]
		g := color3bitsLookupTable[(byteVal>>5)&0b111]
		b := color2bitsLookupTable[byteVal&0b11]

		img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
	}

	var outputImage image.Image
	if config.DoubleImageSize {
		outputImage = doubleSize(width, height, img)
	} else {
		outputImage = img
	}
	return encodePNG(outputImage)
}

// decodeYaeYjk decodes YAE or YJK encoded data to an image.
func decodeYaeYjk(data []byte, config decoders.Config, width, paletteOffset int, isYae bool) (decoders.DecoderResult, error) {
	beginAddress := binary.LittleEndian.Uint16(data[1:3])
	endAddress := binary.LittleEndian.Uint16(data[3:5])
	height := calculateHeight(endAddress, width)

	pixels := data[7:]
	var palette color.Palette
	if paletteOffset > 0 {
		if endAddress >= uint16(paletteOffset) {
			palette = getPalette(pixels, paletteOffset-int(beginAddress))
		} else if config.ExtraData != nil {
			palette = getPalette(config.ExtraData, 0)
		} else {
			palette = getPalette(defaultPalette, 0)
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x += 4 {
			idx := y*width + x
			b1 := pixels[idx+0]
			b2 := pixels[idx+1]
			b3 := pixels[idx+2]
			b4 := pixels[idx+3]

			k := int(b1&7 + (b2&7)<<3)
			j := int(b3&7 + (b4&7)<<3)

			if k > 31 {
				k -= 64
			}
			if j > 31 {
				j -= 64
			}

			for i := 0; i < 4; i++ {
				yVal := int(pixels[idx+i] >> 3)
				if isYae && (yVal&1) == 1 {
					img.Set(x+i, y, palette[yVal>>1])
				} else {
					r := clamp(yVal+j, 0, 31)
					g := clamp(yVal+k, 0, 31)
					b := clamp(5*yVal/4-j/2-k/4, 0, 31)
					img.Set(x+i, y, color.RGBA{color5bitsLookupTable[r], color5bitsLookupTable[g], color5bitsLookupTable[b], 255})
				}
			}
		}
	}

	var outputImage image.Image
	if config.DoubleImageSize {
		outputImage = doubleSize(width, height, img)
	} else {
		outputImage = img
	}
	return encodePNG(outputImage)
}

// DecodeScreen5 decodes screen 5 data.
func DecodeScreen5(data []byte, config decoders.Config) (decoders.DecoderResult, error) {
	return decodeScreenNibbles(data, config, ScreenWidth, PaletteOffset5)
}

// DecodeScreen7 decodes screen 7 data.
func DecodeScreen7(data []byte, config decoders.Config) (decoders.DecoderResult, error) {
	return decodeScreenNibbles(data, config, ScreenWidth7, PaletteOffset)
}

// DecodeScreen8 decodes screen 8 data.
func DecodeScreen8(data []byte, config decoders.Config) (decoders.DecoderResult, error) {
	return decodeScreen(data, config, ScreenWidth)
}

// DecodeScreen10 decodes screen 10 data.
func DecodeScreen10(data []byte, config decoders.Config) (decoders.DecoderResult, error) {
	return decodeYaeYjk(data, config, ScreenWidth, PaletteOffset, true)
}

// DecodeScreen12 decodes screen 12 data.
func DecodeScreen12(data []byte, config decoders.Config) (decoders.DecoderResult, error) {
	return decodeYaeYjk(data, config, ScreenWidth, 0, false)
}
