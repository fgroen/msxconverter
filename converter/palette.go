package converter

import (
	"encoding/binary"
	"fmt"
)

const (
	// BinaryHeaderSize is the size of the MSX binary file header
	BinaryHeaderSize = 7
	// PaletteSize is the size of a palette (16 colors × 2 bytes)
	PaletteSize = 32
	// MagicByteBinary is the magic byte for MSX binary files
	MagicByteBinary = 0xFE
	// SC5PaletteOffset is the VRAM offset for SCREEN 5 palette (0x7680)
	SC5PaletteOffset = 0x7680
)

// ProcessPaletteFile processes a palette file with format detection
// Returns the palette data (32 bytes) or an error
func ProcessPaletteFile(paletteData []byte) ([]byte, error) {
	if len(paletteData) == 0 {
		return nil, fmt.Errorf("palette file is empty")
	}

	// Case 1: File is exactly 32 bytes - use as-is
	if len(paletteData) == PaletteSize {
		return paletteData, nil
	}

	// Case 2: File starts with 0xFE (binary header)
	if paletteData[0] == MagicByteBinary {
		if len(paletteData) < BinaryHeaderSize {
			return nil, fmt.Errorf("palette file with binary header is too small: got %d bytes, need at least %d", len(paletteData), BinaryHeaderSize)
		}

		// Check if data after header is exactly 32 bytes
		dataAfterHeader := paletteData[BinaryHeaderSize:]
		if len(dataAfterHeader) == PaletteSize {
			// Use the 32 bytes after the header as palette
			return dataAfterHeader, nil
		}

		// Check if it's a SCREEN 5 file and extract palette from offset 0x7680
		if isScreen5File(paletteData) {
			return extractScreen5Palette(paletteData)
		}

		return nil, fmt.Errorf("palette file with binary header: expected 32 bytes after header or SCREEN 5 file, got %d bytes after header", len(dataAfterHeader))
	}

	return nil, fmt.Errorf("palette file format not recognized: file size is %d bytes (expected 32 bytes, or binary file with 0xFE header)", len(paletteData))
}

// isScreen5File checks if the file is a SCREEN 5 file
// SCREEN 5 files have header: [0xFE][begin low][begin high][end low][end high][?][?]
// where end address is typically 0x769F (which includes the palette at 0x7680-0x769F)
func isScreen5File(data []byte) bool {
	if len(data) < BinaryHeaderSize {
		return false
	}

	// Check magic byte
	if data[0] != MagicByteBinary {
		return false
	}

	// Parse end address from header (bytes 3-4, little-endian)
	endAddress := binary.LittleEndian.Uint16(data[3:5])

	// SCREEN 5 palette is at 0x7680-0x769F, so end address should be at least 0x769F
	// This indicates the file includes the palette area
	if endAddress >= SC5PaletteOffset+PaletteSize-1 {
		return true
	}

	return false
}

// extractScreen5Palette extracts the palette from a SCREEN 5 file at offset 0x7680-0x769F
func extractScreen5Palette(data []byte) ([]byte, error) {
	// The palette is at VRAM address 0x7680-0x769F
	// In the file, we need to calculate the offset based on the begin address
	if len(data) < BinaryHeaderSize {
		return nil, fmt.Errorf("file too small for SCREEN 5 header")
	}

	// Parse begin address from header (bytes 1-2, little-endian)
	beginAddress := binary.LittleEndian.Uint16(data[1:3])

	// Calculate file offset for palette
	// Palette VRAM address: 0x7680
	// File offset = (VRAM address - begin address) + header size
	if SC5PaletteOffset < beginAddress {
		return nil, fmt.Errorf("palette address 0x%04X is before begin address 0x%04X", SC5PaletteOffset, beginAddress)
	}

	paletteFileOffset := int(SC5PaletteOffset-beginAddress) + BinaryHeaderSize

	// Check if file is large enough
	if len(data) < paletteFileOffset+PaletteSize {
		return nil, fmt.Errorf("file too small to contain palette: need %d bytes, got %d", paletteFileOffset+PaletteSize, len(data))
	}

	// Extract palette
	palette := make([]byte, PaletteSize)
	copy(palette, data[paletteFileOffset:paletteFileOffset+PaletteSize])

	return palette, nil
}

