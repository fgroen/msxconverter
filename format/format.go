package format

import (
	"path/filepath"
	"strings"
)

func DetectFormat(data []byte, inputFileName string, fileType string) string {
	if len(fileType) > 0 {
		// don't try to detect when file type is passed
		return fileType
	}

	if len(data) == 0 {
		return "unknown"
	}

	switch data[0] {
	case MagicByteBAS:
		return "BAS"
	case MagicByteWB2:
		return "WB2"
	case MagicByteBinary:
		if len(data) >= 7 {
			extension := strings.ToUpper(strings.TrimLeft(filepath.Ext(inputFileName), "."))

			// First, try to detect by extension (most reliable)
			if extension == "GE5" || extension == "SC5" || extension == "SR5" {
				return "SC5"
			} else if extension == "SC7" || extension == "SR7" {
				return "SC7"
			} else if extension == "SC8" || extension == "PIC" || extension == "SR8" {
				return "SC8"
			} else if extension == "S10" || extension == "SCA" {
				return "S10"
			} else if extension == "S12" || extension == "SCC" || extension == "SRS" {
				return "S12"
			} else if extension == "CMP" {
				return "CMP"
			}

			// Extension not recognized - check for known binary headers
			if len(data) >= len(HeaderSC5) {
				if matchesHeader(data, HeaderSC5) {
					return "SC5"
				}
				// Add more header checks here as we discover them
			}
		}
		return "unknown"

	default:
		extension := strings.ToUpper(strings.TrimLeft(filepath.Ext(inputFileName), "."))
		if extension != "" {
			return extension
		}
		return "unknown"
	}
}
