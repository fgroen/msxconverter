package msxbasic

import (
	"bytes"
	"errors"
	"fmt"
	"msxconverter/decoders"
	"strings"
)

// basic tokens 0x81 ..
var tokenMap = []string{
	"END", "FOR", "NEXT", "DATA", "INPUT", "DIM", "READ", "LET", "GOTO", "RUN",
	"IF", "RESTORE", "GOSUB", "RETURN", "REM", "STOP", "PRINT", "CLEAR", "LIST",
	"NEW", "ON", "WAIT", "DEF", "POKE", "CONT", "CSAVE", "CLOAD", "OUT", "LPRINT",
	"LLIST", "CLS", "WIDTH", "ELSE", "TRON", "TROFF", "SWAP", "ERASE", "ERROR",
	"RESUME", "DELETE", "AUTO", "RENUM", "DEFSTR", "DEFINT", "DEFSNG", "DEFDBL",
	"LINE", "OPEN", "FIELD", "GET", "PUT", "CLOSE", "LOAD", "MERGE", "FILES",
	"LSET", "RSET", "SAVE", "LFILES", "CIRCLE", "COLOR", "DRAW", "PAINT", "BEEP",
	"PLAY", "PSET", "PRESET", "SOUND", "SCREEN", "VPOKE", "SPRITE", "VDP", "BASE",
	"CALL", "TIME", "KEY", "MAX", "MOTOR", "BLOAD", "BSAVE", "DSKO$",
	"SET", "NAME", "KILL", "IPL", "COPY", "CMD", "LOCATE",
	"TO", "THEN", "TAB(", "STEP", "USR", "FN", "SPC(", "NOT", "ERL", "ERR",
	"STRING$", "USING", "INSTR", "'", "VARPTR", "CSRLIN", "ATTR$", "DSKI$", "OFF",
	"INKEY$", "POINT", ">", "=", "<", "+", "-", "*", "/", "^", "AND", "OR", "XOR",
	"EQV", "IMP", "MOD", "\\",
}

var tokenMapFF = []string{
	"LEFT$", "RIGHT$", "MID$", "SGN", "INT", "ABS", "SQR", "RND", "SIN", "LOG",
	"EXP", "COS", "TAN", "ATN", "FRE", "INP", "POS", "LEN", "STR$", "VAL", "ASC",
	"CHR$", "PEEK", "VPEEK", "SPACES$", "OCT$", "HEX$", "LPOS", "BIN$", "CINT",
	"CSNG", "CDBL", "FIX", "STICK", "STRIG", "PDL", "PAD", "DSKF", "FPOS", "CVI",
	"CVS", "CVD", "EOF", "LOC", "LOF", "MKI$", "MK$", "MKD$",
}

// tokens from https://www.msx.org/wiki/Internal_Structure_Of_BASIC_listing

func DecodeMSXBasic(data []byte) (decoders.DecoderResult, error) {
	decoderResult := decoders.DecoderResult{}

	if len(data) == 0 || data[0] != 0xFF {
		return decoderResult, errors.New("invalid MSX Basic file")
	}

	var result bytes.Buffer
	offset := 1

	for {
		if offset+4 > len(data) {
			break
		}

		// Skip address of next line (2 bytes)
		offset += 2

		// Read line number (2 bytes)
		lineNumber := int(data[offset]) + int(data[offset+1])*256
		offset += 2
		result.WriteString(fmt.Sprintf("%d ", lineNumber))

		// Read tokens until 0x00 (end of line)
		for offset < len(data) && data[offset] != 0x00 {
			token := data[offset]
			if token == 0x0E || token == 0x1C {
				line := int(data[offset+1]) | int(data[offset+2])<<8
				offset += 2
				result.WriteString(fmt.Sprintf("%d", line))
			} else if token == 0x0F {
				value := int(data[offset+1])
				offset++
				result.WriteString(fmt.Sprintf("%d", value))
			} else if token == 0x1D {
				//1D 3F 50 00 00 = .05
				result.WriteString(customBCDToString(data[offset+1 : offset+5]))
				offset += 4
			} else if token == 0x3A {
				if offset+1 < len(data) {
					nextToken := data[offset+1]
					if nextToken == 0x8F { // :REM'
						// skip :REM
						offset++
					} else if nextToken == 0xA1 { // :ELSE
						// skip :
					} else {
						result.WriteByte(token)
					}
				} else {
					result.WriteByte(token)
				}
			} else if token == 255 {
				// Skip this byte and use the token map FF for the next byte
				offset++
				if offset < len(data) {
					nextToken := data[offset]
					result.WriteString(tokenMapFF[nextToken-0x81])
				}
			} else if token >= 128 {
				if int(token-0x81) < len(tokenMap) {
					result.WriteString(tokenMap[token-0x81])
				} else {
					result.WriteString(fmt.Sprintf("-%d-", token))
				}
			} else if token == 34 { // quoted string
				result.WriteByte(token)
				offset++
				for {
					token = data[offset]
					result.WriteByte(token)
					if token == 34 {
						break
					}
					if data[offset+1] == 0 {
						break
					}
					offset++
				}
			} else if token >= 32 {
				result.WriteByte(token)
			} else if token >= 17+0 && token <= 17+9 {
				result.WriteString(fmt.Sprintf("%d", token-17))
			}
			offset++
		}

		// End of line
		if offset < len(data) && data[offset] == 0x00 {
			result.WriteString("\n")
			offset++
		}

		// Check if next line address is 0 (end of file)
		if offset+2 <= len(data) && data[offset] == 0x00 && data[offset+1] == 0x00 {
			break
		}
	}

	decoderResult.Text = result.String()
	decoderResult.IsText = true

	return decoderResult, nil
}

func customBCDToString(b []byte) string {
	if len(b) != 4 {
		return ""
	}

	sign := ""
	if b[0]&0x80 != 0 {
		sign = "-"
	}

	exponent := int(b[0]&0x7F) - 64
	mantissa := fmt.Sprintf("%02X%02X%02X", b[1], b[2], b[3])

	mantissaString := insertDecimalPoint(mantissa, 1)
	mantissaString = RemoveTrailingZeros(mantissaString)

	switch {
	case exponent == -64:
		return "0!"
	case exponent >= -63 && exponent < -1:
		return fmt.Sprintf("%s%sE%03d", sign, mantissaString, exponent-1)
	case exponent == -1:
		mantissaString = ".0" + mantissa
		mantissaString = RemoveTrailingZeros(mantissaString)
		return fmt.Sprintf("%s%s", sign, mantissaString)
	case exponent >= 0 && exponent < 15:
		mantissaString = shiftPointRight(mantissa, exponent)
		mantissaString = RemoveTrailingZeros(mantissaString)
		return fmt.Sprintf("%s%s", sign, mantissaString)
	case exponent >= 15 && exponent <= 63:
		return fmt.Sprintf("%s%sE+%02d", sign, mantissaString, exponent-1)
	default:
		return "XXX"
	}
}

func insertDecimalPoint(mantissa string, pos int) string {
	runes := []rune(mantissa)
	runes = append(runes[:pos], append([]rune{'.'}, runes[pos:]...)...)
	return string(runes)
}

func shiftPointRight(mantissa string, shift int) string {
	runes := []rune(mantissa)
	if len(runes) <= shift {
		return mantissa + strings.Repeat("0", shift-len(runes)) + "!"
	}
	return string(runes[:shift]) + "." + string(runes[shift:])
}

func RemoveTrailingZeros(numStr string) string {
	numStr = strings.TrimRight(numStr, "0")
	numStr = strings.TrimSuffix(numStr, ".")
	return numStr
}
