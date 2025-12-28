# CMP Decompression Algorithm

## Overview

CMP is an RLE-based compression format for MSX SCREEN 5 images. The compression uses a bit-stream with lookup tables to reduce file size.

This format was used in **Dot Designer's Club** from T&E Soft, released in 1989.

## File Format Structure

```
Byte 0:   WIDTH  (bytes per scanline, where 1 byte = 2 pixels in SCREEN 5)
Byte 1:   HEIGHT (number of scanlines)
Byte 2:   CONTROL (number of uncompressed scanlines for Phase 1)
Byte 3+:  Compressed RLE data stream
```

## Decompression Algorithm

Decompression happens in two phases:

### Phase 1: Direct Write

The first `CONTROL` scanlines are read directly from the bit stream and written to the output buffer without any decompression.

### Phase 2: XOR Phase

The remaining `(HEIGHT - CONTROL)` scanlines are read from the bit stream, then each byte is XORed with the corresponding byte from the previous scanline. This delta compression only stores differences between scanlines.

## Bit Stream Reader

The bit stream uses 8-byte lookup tables:

- **Lookup tables**: A control byte determines which bytes to read from the stream (bit set = read byte, bit clear = use zero)
- **Bit buffer**: An 8-bit buffer tracks when to read bytes from the stream
- **Reloading**: When a lookup table runs out, a new one is read from the stream
- **End-of-file**: If the stream ends while reading a new lookup table, remaining bits are set to zero

## Output Format

The decompressed data is converted to SCREEN 5 format:
- 2 pixels per byte (nibble mode: high nibble = left pixel, low nibble = right pixel)
- 16-color palette (default MSX SCREEN 5 palette)

