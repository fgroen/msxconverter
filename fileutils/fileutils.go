package fileutils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ReadInput(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	return buf.Bytes(), nil
}

func WriteOutput(fileName, data string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", fileName, err)
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return fmt.Errorf("failed to write data to %s: %w", fileName, err)
	}
	return nil
}

func WriteOutputBytes(fileName string, data []byte) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", fileName, err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to %s: %w", fileName, err)
	}
	return nil
}

func GenerateOutputFilename(inputFile, extension string) string {
	return strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + extension
}
