package decoders

import "msxconverter/decoders/types"

type Decoder interface {
	Decode(data []byte, config types.Config) (types.DecoderResult, error)
}

// Type aliases for backward compatibility with existing code
type Config = types.Config
type DecoderResult = types.DecoderResult
