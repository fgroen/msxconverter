package types

import "bytes"

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

