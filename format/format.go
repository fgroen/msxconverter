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
	case 0xFF:
		return "BAS" // MSX Basic
	case 0xFD:
		return "WB2" // WBASS2 file
	case 0xFE:
		if len(data) >= 7 {
			extension := strings.ToUpper(strings.TrimLeft(filepath.Ext(inputFileName), "."))

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
			}

		}
		return "unknown"

	default:
		// just use the extension for filetype, until we have something better..
		return strings.ToUpper(strings.TrimLeft(filepath.Ext(inputFileName), "."))
	}
}
