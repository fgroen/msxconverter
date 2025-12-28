package decoders

import (
	"fmt"
	"msxconverter/decoders/cmp"
	"msxconverter/decoders/images"
	"msxconverter/decoders/msxbasic"
	"msxconverter/decoders/types"
	"msxconverter/decoders/wbass2"
)

type Registry struct {
	decoders map[string]Decoder
}

func NewRegistry() *Registry {
	r := &Registry{
		decoders: make(map[string]Decoder),
	}

	r.Register("BAS", &msxbasicDecoder{})
	r.Register("WB2", &wbass2Decoder{})
	r.Register("SC5", &screen5Decoder{})
	r.Register("SC7", &screen7Decoder{})
	r.Register("SC8", &screen8Decoder{})
	r.Register("S10", &screen10Decoder{})
	r.Register("S12", &screen12Decoder{})
	r.Register("STP", &stpDecoder{})
	r.Register("CMP", &cmpDecoder{})

	return r
}

func (r *Registry) Register(format string, decoder Decoder) {
	r.decoders[format] = decoder
}

// GetDecoder returns a decoder for the given format, or an error if not found
func (r *Registry) GetDecoder(format string) (Decoder, error) {
	decoder, exists := r.decoders[format]
	if !exists {
		return nil, fmt.Errorf("unknown file format: %s", format)
	}
	return decoder, nil
}

func (r *Registry) GetSupportedFormats() []string {
	formats := make([]string, 0, len(r.decoders))
	for format := range r.decoders {
		formats = append(formats, format)
	}
	return formats
}

type msxbasicDecoder struct{}

func (d *msxbasicDecoder) Decode(data []byte, config types.Config) (types.DecoderResult, error) {
	return msxbasic.DecodeMSXBasic(data)
}

type wbass2Decoder struct{}

func (d *wbass2Decoder) Decode(data []byte, config types.Config) (types.DecoderResult, error) {
	return wbass2.DecodeWBASS2(data)
}

type screen5Decoder struct{}

func (d *screen5Decoder) Decode(data []byte, config types.Config) (types.DecoderResult, error) {
	return images.DecodeScreen5(data, config)
}

type screen7Decoder struct{}

func (d *screen7Decoder) Decode(data []byte, config types.Config) (types.DecoderResult, error) {
	return images.DecodeScreen7(data, config)
}

type screen8Decoder struct{}

func (d *screen8Decoder) Decode(data []byte, config types.Config) (types.DecoderResult, error) {
	return images.DecodeScreen8(data, config)
}

type screen10Decoder struct{}

func (d *screen10Decoder) Decode(data []byte, config types.Config) (types.DecoderResult, error) {
	return images.DecodeScreen10(data, config)
}

type screen12Decoder struct{}

func (d *screen12Decoder) Decode(data []byte, config types.Config) (types.DecoderResult, error) {
	return images.DecodeScreen12(data, config)
}

type stpDecoder struct{}

func (d *stpDecoder) Decode(data []byte, config types.Config) (types.DecoderResult, error) {
	return images.DecodeSTP(data, config)
}

type cmpDecoder struct{}

func (d *cmpDecoder) Decode(data []byte, config types.Config) (types.DecoderResult, error) {
	return cmp.DecodeCMP(data, config)
}
