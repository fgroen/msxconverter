package main

import (
	"flag"
	"fmt"
	"log"
	"msxconverter/decoders"
	"msxconverter/decoders/images"
	"msxconverter/decoders/msxbasic"
	"msxconverter/decoders/wbass2"
	"msxconverter/fileutils"
	"msxconverter/format"
	"os"
	"strings"
)

var validTypes = map[string]bool{
	"SC5": true,
	"SC7": true,
	"SC8": true,
	"S10": true,
	"S12": true,
	"WB2": true,
	"BAS": true,
}

func main() {
	typeFlag := flag.String("t", "", "Specify the file type (e.g., BAS, WB2, SC5, SC7, SC8, S10, S12)")
	outputFormatFlag := flag.String("format", "png", "Specify the output format (e.g., png, jpg)")
	doubleSizeFlag := flag.Bool("double", false, "Double the image size")
	verboseFlag := flag.Bool("verbose", false, "Verbose output")

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 || (len(*typeFlag) > 0 && !validTypes[*typeFlag]) {
		fmt.Println("Usage: msxconverter [options] inputfile(s) [outputfile]")
		flag.PrintDefaults()

		if len(*typeFlag) > 0 && !validTypes[*typeFlag] {
			fmt.Println()
			fmt.Println("Error: unsupported type passed:", *typeFlag)
		}
		os.Exit(1)
	}

	if *verboseFlag {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetOutput(os.Stderr)
	}

	var inputs []string = strings.Split(args[0], ",")

	var outputFileName string
	if len(args) > 1 {
		outputFileName = args[1]
	}

	data, err := fileutils.ReadInput(inputs[0])
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
		os.Exit(1)
	}

	var palette []byte
	if len(inputs) > 1 {
		extra, err := fileutils.ReadInput(inputs[1])
		if err != nil {
			log.Fatalf("Error reading palette input: %v", err)
		}
		palette = extra
	}

	config := createDecoderConfig(*outputFormatFlag, *doubleSizeFlag, *verboseFlag, palette)

	// Detect the format of the input file.
	format := format.DetectFormat(data, inputs[0], *typeFlag)
	if format == "" {
		log.Fatalf("Error: could not detect format of input file")
	}

	decoded, err := decodeData(data, format, config)
	if err != nil {
		log.Fatalf("Error decoding data: %v", err)
	}

	err = writeOutput(outputFileName, decoded, inputs[0])
	if err != nil {
		log.Fatalf("Error writing output: %v", err)
	}
}

func createDecoderConfig(outputFormat string, doubleSize, verbose bool, extraData []byte) decoders.Config {
	return decoders.Config{
		OutputFormat:    outputFormat,
		DoubleImageSize: doubleSize,
		VerboseOutput:   verbose,
		ExtraData:       extraData,
	}
}

func decodeData(data []byte, fileType string, config decoders.Config) (decoders.DecoderResult, error) {
	switch fileType {
	case "BAS":
		return msxbasic.DecodeMSXBasic(data)
	case "WB2":
		return wbass2.DecodeWBASS2(data)
	case "SC5":
		return images.DecodeScreen5(data, config)
	case "SC7":
		return images.DecodeScreen7(data, config)
	case "SC8":
		return images.DecodeScreen8(data, config)
	case "S10":
		return images.DecodeScreen10(data, config)
	case "S12":
		return images.DecodeScreen12(data, config)
	default:
		return decoders.DecoderResult{}, fmt.Errorf("unknown file format: %s", fileType)
	}
}

func writeOutput(outputFileName string, decoded decoders.DecoderResult, inputFileName string) error {
	if outputFileName != "" {
		if decoded.IsText {
			return fileutils.WriteOutput(outputFileName, decoded.Text)
		}
		return fileutils.WriteOutputBytes(outputFileName, decoded.Buffer.Bytes())
	}

	// If no output file name is provided, generate one or write to stdout.
	if decoded.IsText {
		fmt.Print(decoded.Text)
	} else {
		outputFileName = fileutils.GenerateOutputFilename(inputFileName, ".png")
		return fileutils.WriteOutputBytes(outputFileName, decoded.Buffer.Bytes())
	}
	return nil
}
