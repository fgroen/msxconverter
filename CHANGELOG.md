# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.1] - 2024-06-14

### Added

- Initial setup of the project.
- Added rudimentary MSX BASIC file conversion.
- Added WBASS2 file conversion.
- Added screen 5, screen 7, and screen 8 image file conversion.
- Added screen 10 and screen 12 image file conversion.

## [0.0.2] - 2025-05-20

### Added

- Support for Dynamic Publisher STP stamp images.

## [0.0.3] - 2025-12-28

### Added

- Support for CMP file format (RLE-compressed SCREEN 5 images from Dot Designer's Club by T&E Soft, 1989).
- Enhanced palette file loading with format detection: supports raw 32-byte palettes, binary header files, and automatic extraction from SCREEN 5 files.
- MSX BASIC decoder now supports octal (&O) and hexadecimal (&H) number parsing. Contributed by [@pliniofm](https://github.com/pliniofm). See [#2](https://github.com/fgroen/msxconverter/issues/2).
- Documentation for file formats: added format specifications for SCREEN 5 (SC5/GE5) and CMP decompression algorithm in the `docs/` folder.