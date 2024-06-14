# MSX Converter

MSX Converter is a tool for converting various MSX file formats to modern image and text formats. This tool supports multiple MSX file types, including screen formats and BASIC files.

This is my first project using Go, so I am open to feedback and improvements.

## Features

- Convert MSX screen formats (SC5, SC7, SC8, S10, S12) to PNG images.
- Convert MSX BASIC files (BAS) to text.
- Convert WBASS2 files (WB2) to text.
- Supports additional palette data for accurate color rendering.
- Option to double the size of the output image.
- Verbose output for detailed logging.

## Installation

Clone the repository and build the project using Go:

```sh
git clone https://github.com/fgroen/msxconverter.git
cd msxconverter
go build
```

## Usage

The basic usage of the MSX Converter is as follows:

```sh
msxconverter [options] inputfile(s) [outputfile]
```

### Options

- `-t`: Specify the file type (e.g., BAS, WB2, SC5, SC7, SC8, S10, S12).
- `-double`: Double the image output size.

### Examples

#### Convert an SC5 file to PNG

```sh
msxconverter -t SC5 input.sc5 output.png
```

#### Convert a BAS file to text

```sh
msxconverter -t BAS input.bas output.txt
```

#### Convert an SC7 file with a separate palette to PNG with doubled image size 

```sh
msxconverter -t SC7 -double input.sc7,input.pl5 output.png
```

#### Convert a WB2 file to text

```sh
msxconverter -t WB2 input.wb2 output.txt
```

### Supported Input and Output Formats

#### Input File Types

- **BAS**: MSX BASIC files (can be autodetected).
- **WB2**: WBASS2 files (can be autodetected).
- **SC5**: MSX Screen 5 files.
- **SC7**: MSX Screen 7 files.
- **SC8**: MSX Screen 8 files.
- **S10**: MSX Screen 10 files.
- **S12**: MSX Screen 12 files.

#### Output Formats

- **png**: PNG image format (default for screen files).
- **txt**: Plain text format (default for BASIC and WBASS2 files).

## TODO

- [/] MSX BASIC file conversion to text (number formats need more work).

## References

The following resources were used in the creation of this project:

- [MSX Assembly Page](http://map.grauw.nl)
- [V9938 Programmer's Guide](http://rs.gr8bit.ru/Documentation/V9938-programmers-guide.pdf) - Information on decoding Screen 5, Screen 7, and Screen 8 files.
- [Yamaha V9958 PDF](https://map.grauw.nl/resources/video/yamaha_v9958.pdf) - Information for MSX2+ YJK and YAE to convert Screen 10 and Screen 12 data.
- [WBASS2 Source](https://wbsoft.home.xs4all.nl/msx/) - Source of WBASS2.
