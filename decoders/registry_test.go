package decoders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()
	assert.NotNil(t, registry)

	expectedFormats := []string{"BAS", "WB2", "SC5", "SC7", "SC8", "S10", "S12", "STP", "CMP"}
	for _, format := range expectedFormats {
		decoder, err := registry.GetDecoder(format)
		assert.NoError(t, err)
		assert.NotNil(t, decoder)
	}
}

func TestRegister(t *testing.T) {
	registry := &Registry{
		decoders: make(map[string]Decoder),
	}

	decoder1 := &msxbasicDecoder{}
	registry.Register("TEST", decoder1)

	decoder, err := registry.GetDecoder("TEST")
	assert.NoError(t, err)
	assert.Equal(t, decoder1, decoder)

	decoder2 := &wbass2Decoder{}
	registry.Register("TEST2", decoder2)

	decoder, err = registry.GetDecoder("TEST2")
	assert.NoError(t, err)
	assert.Equal(t, decoder2, decoder)

	decoder3 := &screen5Decoder{}
	registry.Register("TEST", decoder3)

	decoder, err = registry.GetDecoder("TEST")
	assert.NoError(t, err)
	assert.Equal(t, decoder3, decoder)
	assert.NotEqual(t, decoder1, decoder)
}

func TestGetDecoder(t *testing.T) {
	registry := NewRegistry()

	formats := []string{"BAS", "WB2", "SC5", "SC7", "SC8", "S10", "S12", "STP", "CMP"}
	for _, format := range formats {
		decoder, err := registry.GetDecoder(format)
		assert.NoError(t, err)
		assert.NotNil(t, decoder)
	}

	_, err := registry.GetDecoder("UNKNOWN")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNKNOWN")
	assert.Contains(t, err.Error(), "unknown file format")
}

func TestGetSupportedFormats(t *testing.T) {
	registry := NewRegistry()
	formats := registry.GetSupportedFormats()

	expectedFormats := map[string]bool{
		"BAS": true,
		"WB2": true,
		"SC5": true,
		"SC7": true,
		"SC8": true,
		"S10": true,
		"S12": true,
		"STP": true,
		"CMP": true,
	}

	assert.Equal(t, len(expectedFormats), len(formats))

	for _, format := range formats {
		assert.True(t, expectedFormats[format])
	}

	for format := range expectedFormats {
		found := false
		for _, f := range formats {
			if f == format {
				found = true
				break
			}
		}
		assert.True(t, found)
	}

	customRegistry := &Registry{
		decoders: make(map[string]Decoder),
	}
	customRegistry.Register("CUSTOM1", &msxbasicDecoder{})
	customRegistry.Register("CUSTOM2", &wbass2Decoder{})

	customFormats := customRegistry.GetSupportedFormats()
	assert.Equal(t, 2, len(customFormats))

	customExpected := map[string]bool{"CUSTOM1": true, "CUSTOM2": true}
	for _, format := range customFormats {
		assert.True(t, customExpected[format])
	}
}
