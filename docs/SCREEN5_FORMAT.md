# SCREEN 5 (SC5 / GE5) File Format

## Overview

SC5 and GE5 files are SCREEN 5 image formats. They store screen data as a binary dump from VRAM, including pixel data and optionally a palette. The typical resolution is 256×212 pixels.

## File Structure

### Header (7 bytes)

```
Byte 0:   0xFE (binary file)
Byte 1-2: Begin address - normally 0x0000
Byte 3-4: End address - normally 0x769F or 0x7FFF (with palette) or 0x69FF (without palette)
Byte 5-6: Execution address - not used for video files, normally 0x0000
```

Common header values:
- `0xFE 0x00 0x00 0x9F 0x76 0x00 0x00` - Full screen with palette (begin: 0x0000, end: 0x769F)
- `0xFE 0x00 0x00 0xFF 0x7F 0x00 0x00` - Full screen with palette (begin: 0x0000, end: 0x7FFF)
- `0xFE 0x00 0x00 0xFF 0x69 0x00 0x00` - Full screen without palette (begin: 0x0000, end: 0x69FF)

Other values are valid too - the header just defines which VRAM address range to load. Only data for those addresses is included in the file.

### Pixel Data (starts at byte 7)

After the header, pixel data follows.

SCREEN 5 uses nibble mode: each byte holds 2 pixels. The high nibble (bits 7-4) is the left pixel, the low nibble (bits 3-0) is the right pixel. Each nibble is a color index from 0-15.

The screen is 256 pixels wide, which means 128 bytes per scanline. Pixel data is stored row by row, starting from the top-left corner.

### Palette (optional)

If the end address is 0x769F or higher, the file includes palette data at VRAM address 0x7680-0x769F. That's 32 bytes total (16 colors × 2 bytes each).

Each color is stored as 2 bytes:
- Byte 0 (low byte): bits 0-2 = Blue (0-7), bits 4-6 = Red (0-7)
- Byte 1 (high byte): bits 0-2 = Green (0-7)


If the file doesn't include a palette (end address < 0x7680), you can provide one as a separate file, or the default MSX SCREEN 5 palette is used.


