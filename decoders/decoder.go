package decoders

import "bytes"

type Decoder interface {
	Decode(data []byte, config Config) (DecoderResult, error)
}

type Config struct {
	OutputFormat    string
	DoubleImageSize bool
	VerboseOutput   bool
	ExtraData       []byte
}

type DecoderResult struct {
	Text   string
	Buffer *bytes.Buffer
	IsText bool
}
