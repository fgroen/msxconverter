package converter

import (
	"fmt"
	"msxconverter/decoders"
	"msxconverter/fileutils"
	"msxconverter/format"
	"os"
)

type Converter struct {
	registry *decoders.Registry
}

func NewConverter() *Converter {
	return &Converter{
		registry: decoders.NewRegistry(),
	}
}

type ConvertOptions struct {
	InputFiles    []string
	OutputFile    string
	TypeHint      string
	OutputFormat  string
	DoubleSize    bool
	VerboseOutput bool
}

func (c *Converter) Convert(opts ConvertOptions) error {
	for _, inputFile := range opts.InputFiles {
		if err := fileutils.ValidateInputFile(inputFile); err != nil {
			return fmt.Errorf("invalid input file %s: %w", inputFile, err)
		}
	}

	if err := fileutils.ValidateOutputPath(opts.OutputFile); err != nil {
		return fmt.Errorf("cannot write to output path: %w", err)
	}

	data, err := fileutils.ReadInput(opts.InputFiles[0])
	if err != nil {
		return fmt.Errorf("cannot read input file %s: %w", opts.InputFiles[0], err)
	}

	var palette []byte
	if len(opts.InputFiles) > 1 {
		paletteData, err := fileutils.ReadInput(opts.InputFiles[1])
		if err != nil {
			return fmt.Errorf("cannot read palette file %s: %w", opts.InputFiles[1], err)
		}
		palette, err = ProcessPaletteFile(paletteData)
		if err != nil {
			return fmt.Errorf("cannot process palette file %s: %w", opts.InputFiles[1], err)
		}
	}

	detectedFormat := format.DetectFormat(data, opts.InputFiles[0], opts.TypeHint)
	if detectedFormat == "" || detectedFormat == "unknown" {
		return fmt.Errorf("could not detect format of input file: %s\nPlease specify the file type using the -t flag (e.g., -t SC5, -t SC7, -t BAS, etc.)", opts.InputFiles[0])
	}

	config := decoders.Config{
		OutputFormat:    opts.OutputFormat,
		DoubleImageSize: opts.DoubleSize,
		VerboseOutput:   opts.VerboseOutput,
		ExtraData:       palette,
	}

	decoder, err := c.registry.GetDecoder(detectedFormat)
	if err != nil {
		return fmt.Errorf("unknown format %s: %w", detectedFormat, err)
	}

	result, err := decoder.Decode(data, config)
	if err != nil {
		return fmt.Errorf("decode error for %s: %w", detectedFormat, err)
	}

	if err := c.writeOutput(opts.OutputFile, result, opts.InputFiles[0]); err != nil {
		return fmt.Errorf("cannot write output: %w", err)
	}

	return nil
}

func (c *Converter) writeOutput(outputFileName string, decoded decoders.DecoderResult, inputFileName string) error {
	if outputFileName != "" {
		if decoded.IsText {
			return fileutils.WriteOutput(outputFileName, decoded.Text)
		}
		return fileutils.WriteOutputBytes(outputFileName, decoded.Buffer.Bytes())
	}

	if decoded.IsText {
		fmt.Fprint(os.Stdout, decoded.Text)
		return nil
	}

	outputFileName = fileutils.GenerateOutputFilename(inputFileName, ".png")
	return fileutils.WriteOutputBytes(outputFileName, decoded.Buffer.Bytes())
}
