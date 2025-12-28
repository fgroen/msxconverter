package format

const (
	MagicByteBAS     = 0xFF // MSX BASIC file
	MagicByteWB2     = 0xFD // WBASS2 file
	MagicByteBinary  = 0xFE // Binary file
)

// Format: [magic byte][begin address low][begin address high][end address low][end address high][?][?]
var (
	HeaderSC5 = []byte{0xFE, 0x00, 0x00, 0x9F, 0x76, 0x00, 0x00}
)

