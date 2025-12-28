package cmp

import (
	"fmt"
	"image"
	"msxconverter/decoders/images"
	"msxconverter/decoders/types"
)

// bitStreamReader implements the bit stream reading logic from the Z80 assembly
// It maintains state similar to the alternate register set (EXX) in the assembly
type bitStreamReader struct {
	data        []byte
	offset      int
	bitBuffer   uint8  // B register equivalent - 8-bit buffer
	bitCounter  uint8  // C register equivalent - counts down from 8
	lookupTable []byte // T90D1 equivalent - 8-byte lookup table
	tableIndex  int    // Current position in lookup table
	// In assembly, DE points to lookup table and increments
	// When DE == T90D9 (which is 64 bytes after T90D1), read new table
	// T90D1 is at offset 0, T90D9 is at offset 8 (but T90D9 is a 64-byte buffer)
	// Actually, the check is: if we've exhausted the 8-byte table, read new one
}

// newBitStreamReader creates a new bit stream reader starting at the given offset
func newBitStreamReader(data []byte, offset int) (*bitStreamReader, error) {
	if offset >= len(data) {
		return nil, fmt.Errorf("offset %d exceeds data length %d", offset, len(data))
	}

	reader := &bitStreamReader{
		data:        data,
		offset:      offset,
		lookupTable: make([]byte, 8),
		tableIndex:  0,
	}

	// Initialize by reading the first lookup table
	if err := reader.readLookupTable(); err != nil {
		return nil, err
	}

	return reader, nil
}

// readLookupTable reads an 8-byte lookup table from the stream
// This corresponds to the A9028 function in the assembly
func (r *bitStreamReader) readLookupTable() error {
	if r.offset >= len(r.data) {
		return fmt.Errorf("unexpected end of data while reading lookup table")
	}

	// Read control byte (8 bits)
	controlByte := r.data[r.offset]
	r.offset++

	// For each bit in the control byte, if set, read a byte from stream; else use 0
	for i := 0; i < 8; i++ {
		if (controlByte & (1 << (7 - i))) != 0 {
			// Bit is set - read byte from stream
			if r.offset >= len(r.data) {
				return fmt.Errorf("unexpected end of data while reading lookup table byte %d", i)
			}
			r.lookupTable[i] = r.data[r.offset]
			r.offset++
		} else {
			// Bit is clear - use 0
			r.lookupTable[i] = 0
		}
	}

	// Initialize bit stream state from first byte of lookup table
	r.bitBuffer = r.lookupTable[0]
	r.bitCounter = 8
	r.tableIndex = 1

	return nil
}

// readByte reads the next byte from the bit stream
// This corresponds to the A9016 function in the assembly
func (r *bitStreamReader) readByte() (byte, error) {
	var result byte = 0

	// Shift bit buffer left (SLA B equivalent)
	carry := (r.bitBuffer & 0x80) != 0
	r.bitBuffer <<= 1

	if carry {
		// Carry is set - read byte from stream
		if r.offset >= len(r.data) {
			return 0, fmt.Errorf("unexpected end of data while reading byte from stream (offset=%d, len=%d)", r.offset, len(r.data))
		}
		result = r.data[r.offset]
		r.offset++
	}

	// Decrement counter
	r.bitCounter--

	// If counter reached 0, reload from lookup table
	if r.bitCounter == 0 {
		// Check if we need to read a new lookup table
		// In assembly: if DE (table pointer) == T90D9, read new table
		// Since T90D1 is 8 bytes, when tableIndex >= 8, we've exhausted it
		if r.tableIndex >= len(r.lookupTable) {
			// Table exhausted - try to read new table from stream
			// If we're at the end of the file, use zeros instead
			if r.offset >= len(r.data) {
				// End of file - use zeros for remaining bits
				r.bitBuffer = 0
				r.bitCounter = 8
				r.tableIndex = 0
				// Fill lookup table with zeros
				for i := range r.lookupTable {
					r.lookupTable[i] = 0
				}
			} else {
				// Read new table from stream
				if err := r.readLookupTable(); err != nil {
					return 0, fmt.Errorf("failed to read new lookup table: %w", err)
				}
			}
		} else {
			// Reload from current table
			r.bitBuffer = r.lookupTable[r.tableIndex]
			r.tableIndex++
			r.bitCounter = 8
		}
	}

	return result, nil
}

// DecodeCMP decompresses a CMP file and returns a SCREEN 5 image
// File format:
//   Byte 0: WIDTH (bytes per scanline, where 1 byte = 2 pixels)
//   Byte 1: HEIGHT (number of scanlines)
//   Byte 2: Control byte (number of uncompressed scanlines for Phase 1)
//   Byte 3+: Compressed RLE data stream
func DecodeCMP(data []byte, config types.Config) (types.DecoderResult, error) {
	if len(data) < 3 {
		return types.DecoderResult{}, fmt.Errorf("invalid CMP file: file too small (got %d bytes, need at least 3)", len(data))
	}

	// Read header
	width := int(data[0])   // WIDTH in bytes (1 byte = 2 pixels)
	height := int(data[1])  // HEIGHT in scanlines
	control := int(data[2]) // Number of uncompressed scanlines

	if width <= 0 || height <= 0 {
		return types.DecoderResult{}, fmt.Errorf("invalid CMP file: invalid dimensions (width=%d, height=%d)", width, height)
	}

	if control < 0 || control > height {
		return types.DecoderResult{}, fmt.Errorf("invalid CMP file: invalid control byte %d (height=%d)", control, height)
	}

	// Initialize bit stream reader starting at byte 3
	reader, err := newBitStreamReader(data, 3)
	if err != nil {
		return types.DecoderResult{}, fmt.Errorf("failed to initialize bit stream reader: %w", err)
	}

	// Create output buffer (SCREEN 5: 256 pixels = 128 bytes per scanline)
	// But we use the width from the file
	outputWidth := width * 2 // Convert bytes to pixels
	if outputWidth > images.ScreenWidth {
		outputWidth = images.ScreenWidth
	}

	// Create output buffer for decompressed data
	output := make([]byte, width*height)

	// Phase 1: Direct write phase (uncompressed scanlines)
	for scanline := 0; scanline < control; scanline++ {
		for x := 0; x < width; x++ {
			b, err := reader.readByte()
			if err != nil {
				return types.DecoderResult{}, fmt.Errorf("error reading byte at scanline %d, x %d: %w", scanline, x, err)
			}
			output[scanline*width+x] = b
		}
	}

	// Phase 2: XOR phase (delta compression for remaining scanlines)
	// Each scanline is XORed with the previous scanline's data
	for scanline := control; scanline < height; scanline++ {
		prevScanline := scanline - 1
		for x := 0; x < width; x++ {
			b, err := reader.readByte()
			if err != nil {
				return types.DecoderResult{}, fmt.Errorf("error reading byte at scanline %d, x %d: %w", scanline, x, err)
			}
			// XOR with previous scanline's data (delta compression)
			if prevScanline >= 0 {
				b ^= output[prevScanline*width+x]
			}
			output[scanline*width+x] = b
		}
	}

	// Convert to SCREEN 5 image format (2 pixels per byte, nibble mode)
	// Use palette from config.ExtraData if provided, otherwise use default
	// CMP files don't have embedded palette data, so we pass nil for pixels and 0 for addresses/offset
	palette := images.ResolvePalette(nil, 0, 0, 0, config.ExtraData)
	img := image.NewPaletted(image.Rect(0, 0, outputWidth, height), palette)

	// Convert bytes to pixels (2 pixels per byte)
	for y := 0; y < height; y++ {
		for x := 0; x < width && x*2 < outputWidth; x++ {
			byteVal := output[y*width+x]
			// High nibble (left pixel)
			img.SetColorIndex(x*2, y, byteVal>>4)
			// Low nibble (right pixel)
			if x*2+1 < outputWidth {
				img.SetColorIndex(x*2+1, y, byteVal&0x0F)
			}
		}
	}

	// Finalize image (apply double size if requested, convert to PNG)
	return images.FinalizeImage(img, outputWidth, height, config.DoubleImageSize)
}

