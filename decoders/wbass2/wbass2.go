package wbass2

import (
	"bytes"
	"errors"
	"fmt"
	"msxconverter/decoders"
	"strings"
)

var (
	instructions = []string{"LD", "JR", "DJNZ", "CALL", "RET", "JP", "INC", "DEC",
		"PUSH", "POP", "RST", "IN", "OUT", "IM", "EX", "ADD", "ADC", "SUB", "SBC",
		"AND", "XOR", "OR", "CP", "RLC", "RRC", "RL", "RR", "SLA", "SRA", "???",
		"SRL", "BIT", "RES", "SET", "CPD", "CPDR", "CPI", "CPIR", "IND", "INDR",
		"INI", "INIR", "LDD", "LDDR", "LDI", "LDIR", "OUTD", "OTDR", "OUTI",
		"OTIR", "NEG", "RETI", "RETN", "RLD", "RRD", "CCF", "CPL", "DAA", "DI",
		"EI", "EXX", "HALT", "NOP", "RLA", "RLCA", "RRA", "RRCA", "SCF", "ORG",
		"EQU", "END", "DB", "DW", "DS", "DM", "DEFB", "DEFW", "DEFS", "DEFM",
		"GLOBAL", "INCLUDE"}
	registers = []string{"A", "B", "C", "D", "E", "H", "L", "I", "R", "BC", "DE",
		"HL", "SP", "IX", "IY", "AF"}
	condities = []string{"NZ", "Z", "NC", "C", "PO", "PE", "P", "M", "$"}
	logies    = []string{"AND", "XOR", "OR", "MOD"}
	kartab    = []string{",", ")", "(", "+", "-", "*", "/", "^"}
)

func DecodeWBASS2(data []byte) (decoders.DecoderResult, error) {
	decoderResult := decoders.DecoderResult{}

	if len(data) == 0 || data[0] != 0xFD {
		return decoderResult, errors.New("invalid WBASS2 file")
	}

	// see TOKENIZE and DETOK in wbass2-1.asm

	var result bytes.Buffer
	offset := 1

	beglabel := findLabelOffset(data)
	if beglabel == -1 {
		return decoderResult, errors.New("invalid WBASS2 file structure")
	}

	for offset < len(data) {
		length := data[offset]
		offset++

		if length == 0xFF { // end of tokenized content
			break
		}

		if length == 0x00 { // empty line
			result.WriteString("\n")
			continue
		}

		line, offset2 := parseLine(length, data, offset, beglabel)
		offset = offset2
		result.Write(line.Bytes())
	}

	decoderResult.Text = result.String()
	decoderResult.IsText = true

	return decoderResult, nil
}

func parseLine(length byte, data []byte, offset int, beglabel int) (bytes.Buffer, int) {
	var line bytes.Buffer

	if length&128 == 128 { // label?
		length = length & 127
		label := int(data[offset]) | int(data[offset+1])<<8
		offset += 2
		length -= 2

		for i := beglabel + 8*label; (data[i]&127) != 0 && i < beglabel+8*label+6; i++ {
			line.WriteByte(data[i] & 0x7F)
		}
		line.WriteString(":")

		if length == 0 { // only label on this line
			line.WriteString("\n")
			return line, offset
		}
	}

	c := data[offset]
	offset++
	length--

	if c == 1 { // comment
		if line.Len() > 0 {
			line.WriteString(strings.Repeat(" ", 8-line.Len()))
		}
		line.WriteString(";")
		for ; length > 0; length-- {
			line.WriteByte(data[offset])
			offset++
		}
	} else if c > 127 { // instruction
		if line.Len() < 8 {
			line.WriteString(strings.Repeat(" ", 8-line.Len()))
		}

		line.WriteString(instructions[c-128])
		if length > 0 {
			line.WriteString(strings.Repeat(" ", 14-line.Len()))
		}
	}

	endline := offset + int(length&127)
	needspace := false
	for offset < endline {
		c := data[offset]
		offset++
		length--

		switch {
		case c == 1: // comment
			if line.Len() > 14 {
				line.WriteString(" ")
				if line.Len() < 30 {
					line.WriteString(strings.Repeat(" ", 30-line.Len()))
				}
			}

			line.WriteString(";")
			for ; length > 0; length-- {
				line.WriteByte(data[offset])
				offset++
			}

		case c > 1 && c < 14*2: // special character
			line.WriteString(kartab[c/2-1]) //  SRL A, CP 14, kartab-1
			needspace = false

		case c == 34: // quoted string
			if end := bytes.IndexByte(data[offset:endline], 34); end != -1 {
				line.WriteByte('"')
				line.Write(data[offset : offset+end+1])
				offset += end + 1
			}

		case c == 0xC0: // label
			label := int(data[offset]) | int(data[offset+1])<<8
			offset += 2
			length -= 2
			for i := beglabel + 8*label; (data[i]&127) != 0 && i < beglabel+8*label+6; i++ {
				line.WriteByte(data[i] & 0x7F)
			}

		case c == 0xE0 || c == 0xE1 || c == 0xE2: // number formats
			handleSpace(&line, needspace)
			number := int(data[offset]) | int(data[offset+1])<<8
			offset += 2
			length -= 2
			handleNumberFormats(&line, c, number)
			needspace = true

		case c >= 128: // registers, conditions and logies
			handleSpace(&line, needspace)
			var value string
			if c >= 153 {
				value = logies[c-153]
			} else if c >= 144 {
				value = condities[c-144]
			} else {
				value = registers[c-128]
			}

			line.WriteString(value)
			needspace = true
		}
	}

	line.WriteString("\n")

	return line, offset
}

func findLabelOffset(data []byte) int {
	for i := 1; i < len(data); {
		c := data[i]
		if c == 0xFF {
			return i + 1
		}
		i += int(c&127) + 1
	}

	return -1
}

func handleSpace(line *bytes.Buffer, needspace bool) {
	if needspace {
		line.WriteString(" ")
	}
}

func handleNumberFormats(line *bytes.Buffer, c byte, number int) {
	switch c {
	case 0xE0:
		line.WriteString(fmt.Sprintf("%d", number))
	case 0xE1:
		if number > 256 {
			line.WriteString(fmt.Sprintf("&H%04X", number))
		} else {
			line.WriteString(fmt.Sprintf("&H%02X", number))
		}
	case 0xE2:
		if number <= 256 {
			line.WriteString("&B" + fmt.Sprintf("%08b", number))
		} else {
			line.WriteString("&B" + fmt.Sprintf("%016b", number))
		}
	}
}
