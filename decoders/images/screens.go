package images

import (
	"image"
	"image/color"
	"msxconverter/decoders/types"
)

func FinalizeImage(img image.Image, width, height int, double bool) (types.DecoderResult, error) {
	outputImage := img
	if double {
		switch v := img.(type) {
		case *image.Paletted:
			outputImage = doubleSizePaletted(width, height, v)
		case *image.RGBA:
			outputImage = doubleSize(width, height, v)
		}
	}
	return encodePNG(outputImage)
}

func decodeScreenNibbles(data []byte, config types.Config, width, paletteOffset int) (types.DecoderResult, error) {
	header, err := parseScreenHeader(data)
	if err != nil {
		return types.DecoderResult{}, err
	}

	height := calculateHeight(header.endAddress, width/2)
	palette := ResolvePalette(header.pixels, header.beginAddress, header.endAddress, paletteOffset, config.ExtraData)

	img := image.NewPaletted(image.Rect(0, 0, width, height), palette)
	if uint16(height*(width/2)) < header.endAddress {
		header.endAddress = uint16(height * (width / 2))
	}

	for addr := header.beginAddress; addr < header.endAddress; addr++ {
		y := int(addr / uint16(width/2))
		x := int(addr % uint16(width/2))
		byteVal := header.pixels[addr-header.beginAddress]

		img.SetColorIndex(x*2, y, byteVal>>4)
		img.SetColorIndex(x*2+1, y, byteVal&0x0F)
	}

	return FinalizeImage(img, width, height, config.DoubleImageSize)
}

func decodeScreen(data []byte, config types.Config, width int) (types.DecoderResult, error) {
	header, err := parseScreenHeader(data)
	if err != nil {
		return types.DecoderResult{}, err
	}

	height := calculateHeight(header.endAddress, width)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	if uint16(height*width) < header.endAddress {
		header.endAddress = uint16(height * width)
	}

	for addr := header.beginAddress; addr < header.endAddress; addr++ {
		y := int(addr / uint16(width))
		x := int(addr % uint16(width))
		byteVal := header.pixels[addr-header.beginAddress]

		r := color3bitsLookupTable[(byteVal>>2)&0b111]
		g := color3bitsLookupTable[(byteVal>>5)&0b111]
		b := color2bitsLookupTable[byteVal&0b11]

		img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
	}

	return FinalizeImage(img, width, height, config.DoubleImageSize)
}

func decodeYaeYjk(data []byte, config types.Config, width, paletteOffset int, isYae bool) (types.DecoderResult, error) {
	header, err := parseScreenHeader(data)
	if err != nil {
		return types.DecoderResult{}, err
	}

	height := calculateHeight(header.endAddress, width)
	var palette color.Palette
	if paletteOffset > 0 {
		palette = ResolvePalette(header.pixels, header.beginAddress, header.endAddress, paletteOffset, config.ExtraData)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x += 4 {
			idx := y*width + x
			b1 := header.pixels[idx+0]
			b2 := header.pixels[idx+1]
			b3 := header.pixels[idx+2]
			b4 := header.pixels[idx+3]

			k := int(b1&7 + (b2&7)<<3)
			j := int(b3&7 + (b4&7)<<3)

			if k > 31 {
				k -= 64
			}
			if j > 31 {
				j -= 64
			}

			for i := 0; i < 4; i++ {
				yVal := int(header.pixels[idx+i] >> 3)
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

	return FinalizeImage(img, width, height, config.DoubleImageSize)
}

func DecodeScreen5(data []byte, config types.Config) (types.DecoderResult, error) {
	return decodeScreenNibbles(data, config, ScreenWidth, PaletteOffset5)
}

func DecodeScreen7(data []byte, config types.Config) (types.DecoderResult, error) {
	return decodeScreenNibbles(data, config, ScreenWidth7, PaletteOffset)
}

func DecodeScreen8(data []byte, config types.Config) (types.DecoderResult, error) {
	return decodeScreen(data, config, ScreenWidth)
}

func DecodeScreen10(data []byte, config types.Config) (types.DecoderResult, error) {
	return decodeYaeYjk(data, config, ScreenWidth, PaletteOffset, true)
}

func DecodeScreen12(data []byte, config types.Config) (types.DecoderResult, error) {
	return decodeYaeYjk(data, config, ScreenWidth, 0, false)
}
