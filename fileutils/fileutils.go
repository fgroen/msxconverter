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
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, reader)
	return buf.Bytes(), err
}

func WriteOutput(fileName, data string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	return err
}

func WriteOutputBytes(fileName string, data []byte) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func GenerateOutputFilename(inputFile, extension string) string {
	return strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + extension
}
