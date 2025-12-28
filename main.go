package main

import (
	"flag"
	"fmt"
	"log"
	"msxconverter/converter"
	"os"
	"strings"
)

var validTypes = map[string]bool{
	"SC5": true,
	"SC7": true,
	"SC8": true,
	"S10": true,
	"S12": true,
	"STP": true,
	"WB2": true,
	"BAS": true,
	"CMP": true,
}

func main() {
	typeFlag := flag.String("t", "", "Specify the file type (e.g., BAS, WB2, SC5, SC7, SC8, S10, S12, STP, CMP)")
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

	conv := converter.NewConverter()
	opts := converter.ConvertOptions{
		InputFiles:    inputs,
		OutputFile:    outputFileName,
		TypeHint:      *typeFlag,
		OutputFormat:  *outputFormatFlag,
		DoubleSize:    *doubleSizeFlag,
		VerboseOutput: *verboseFlag,
	}

	if err := conv.Convert(opts); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
