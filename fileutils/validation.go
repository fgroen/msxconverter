package fileutils

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	MaxFileSize = 64 * 1024
)

func ValidateInputFile(filename string) error {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("input file does not exist: %s", filename)
		}
		return fmt.Errorf("cannot access input file %s: %w", filename, err)
	}

	if info.IsDir() {
		return fmt.Errorf("input path is a directory, not a file: %s", filename)
	}

	if info.Size() > MaxFileSize {
		return fmt.Errorf("input file too large: %d bytes (max: %d bytes)", info.Size(), MaxFileSize)
	}

	if info.Size() == 0 {
		return fmt.Errorf("input file is empty: %s", filename)
	}

	return nil
}

func ValidateOutputPath(outputPath string) error {
	if outputPath == "" {
		// Empty path means stdout, which is always valid
		return nil
	}

	dir := filepath.Dir(outputPath)
	if dir == "." || dir == "" {
		// Output in current directory, assume it's writable
		return nil
	}

	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("output directory does not exist: %s", dir)
		}
		return fmt.Errorf("cannot access output directory %s: %w", dir, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("output path parent is not a directory: %s", dir)
	}

	// Check if directory is writable by trying to create a test file
	testFile := filepath.Join(dir, ".msxconverter_write_test")
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("output directory is not writable: %s", dir)
	}
	file.Close()
	os.Remove(testFile)

	return nil
}
