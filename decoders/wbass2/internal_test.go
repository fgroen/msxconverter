package wbass2

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// tests the parsing of a label definition
// the line ends after the label definition
func TestParseLine_Label(t *testing.T) {
	data := []byte{
		0xFD,
		0x82, 0x00, 0x00, 0xFF,
		0x54, 0x45, 0x53, 0x54, 0x00, 0x00, 0x00, 0x00, 0xFF} // Label definition
	beglabel := 5

	line, _ := parseLine(data[1], data, 2, beglabel)

	fmt.Printf("line: %v\n", line.String())

	expected := "TEST:\n"
	assert.Equal(t, expected, line.String(), "Parsed line mismatch")
}

// tests the parsing of a label definition and an instruction
// the line ends after the instruction, spacing should be correct
func TestParseLine_LabelAndInstruction(t *testing.T) {
	data := []byte{
		0xFD,
		0x88, 0x00, 0x00, 0x80, 0x80, 0x02, 0xE0, 0x01, 0x00, 0xFF,
		0x54, 0x45, 0x53, 0x54, 0x00, 0x00, 0x00, 0x00, 0xFF} // Label definition
	beglabel := 11

	line, _ := parseLine(data[1], data, 2, beglabel)

	fmt.Printf("line: %v\n", line.String())

	expected := "TEST:   LD    A,1\n"
	assert.Equal(t, expected, line.String(), "Parsed line mismatch")
}

func TestParseLine_Comment(t *testing.T) {
	data := []byte{0x05, 0x01, 'T', 'e', 's', 't', 0xFF} // Comment
	beglabel := 0

	line, _ := parseLine(data[0], data, 1, beglabel)

	fmt.Printf("line: %v\n", line.String())

	expected := ";Test\n"
	assert.Equal(t, expected, line.String(), "Parsed line mismatch")
}

func TestParseLine_Instruction(t *testing.T) {
	data := []byte{0x06, 0x80, 0x80, 0x02, 0xE0, 0x01, 0x00} // Instruction LD A,1
	beglabel := 0

	line, _ := parseLine(data[0], data, 1, beglabel)

	fmt.Printf("line: %v\n", line.String())

	expected := "        LD    A,1\n"
	assert.Equal(t, expected, line.String(), "Parsed line mismatch")
}

func TestHandleNumberFormats(t *testing.T) {
	var buffer bytes.Buffer

	handleNumberFormats(&buffer, 0xE0, 123)
	assert.Equal(t, "123", buffer.String(), "Decimal number format mismatch")

	buffer.Reset()
	handleNumberFormats(&buffer, 0xE1, 0x1A)
	assert.Equal(t, "&H1A", buffer.String(), "Hexadecimal number format mismatch")

	buffer.Reset()
	handleNumberFormats(&buffer, 0xE1, 0xBEEF)
	assert.Equal(t, "&HBEEF", buffer.String(), "Hexadecimal number format mismatch")

	buffer.Reset()
	handleNumberFormats(&buffer, 0xE2, 0b1010)
	assert.Equal(t, "&B00001010", buffer.String(), "Binary number format mismatch")

	buffer.Reset()
	handleNumberFormats(&buffer, 0xE2, 0xAAAA)
	assert.Equal(t, "&B1010101010101010", buffer.String(), "Binary number format mismatch")

}
