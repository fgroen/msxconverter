package wbass2_test

import (
	"msxconverter/decoders"
	"msxconverter/decoders/wbass2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeWBASS2_EmptyInput(t *testing.T) {
	data := []byte{}
	_, err := wbass2.DecodeWBASS2(data)
	assert.Error(t, err, "Expected error for empty input")
}

func TestDecodeWBASS2_InvalidStartByte(t *testing.T) {
	data := []byte{0x00}
	_, err := wbass2.DecodeWBASS2(data)
	assert.Error(t, err, "Expected error for invalid start byte")
}

func TestDecodeWBASS2_EmptyFile(t *testing.T) {
	data := []byte{0xFD, 0xFF, 0xFF}
	expected := decoders.DecoderResult{Text: "", IsText: true}

	result, err := wbass2.DecodeWBASS2(data)
	assert.NoError(t, err, "Expected no error for valid input")
	assert.Equal(t, expected, result, "Decoded result mismatch")
}

func TestDecodeWBASS2_ValidInput(t *testing.T) {
	data := []byte{0xFD,
		0x06, 0x80, 0x80, 0x02, 0xE0, 0x01, 0x00, 0xFF, 0xFF} // Sample valid input
	expected := decoders.DecoderResult{Text: "        LD    A,1\n", IsText: true}

	result, err := wbass2.DecodeWBASS2(data)
	assert.NoError(t, err, "Expected no error for valid input")
	assert.Equal(t, expected, result, "Decoded result mismatch")
}
