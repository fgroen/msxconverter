package msxbasic

import (
	"testing"
)

// exponent 0 and any mantissa results in 0!
func TestExponentEquals0(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
	}{
		{input: []byte{0x00, 0x34, 0x56, 0x78}, expected: "0!"},
		{input: []byte{0x00, 0x01, 0x02, 0x03}, expected: "0!"},
		{input: []byte{0x00, 0x88, 0x77, 0x66}, expected: "0!"},
		{input: []byte{0x00, 0xFF, 0xFF, 0xFF}, expected: "0!"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := customBCDToString(tt.input)
			if result != tt.expected {
				t.Errorf("customBCDToString(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNegativeExponent(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
	}{
		{input: []byte{0x01, 0x12, 0x34, 0x56}, expected: "1.23456E-64"},
		{input: []byte{0x02, 0x12, 0x34, 0x56}, expected: "1.23456E-63"},
		{input: []byte{0x03, 0x12, 0x34, 0x56}, expected: "1.23456E-62"},
		{input: []byte{0x04, 0x12, 0x34, 0x56}, expected: "1.23456E-61"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := customBCDToString(tt.input)
			if result != tt.expected {
				t.Errorf("customBCDToString(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
