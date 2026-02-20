<<<<<<<< HEAD:Yoremi V3-CHANGELOG-ENG.md
# V3 — Patch 1: CZ1 32-bit Import/Export + CZ0 logging

## Modified files
- `czimage/cz1.go` — Import/Export/Write rewrite
- `czimage/cz.go` — graceful handling of non-CZ files
- `czimage/cz0.go` — added V(0) log in decompress()

## Bugs fixed

### 1. Extended header missing in Write()
The original `Write()` only wrote the 15 bytes of the `CzHeader` struct, ignoring the 13 bytes of extended header (offsets, crop, bounds). The resulting file had the block table at offset 15 instead of 28 → crash on replay.

**Fix**: Save the raw bytes 15→HeaderLength in `ExtendedHeader` in `Load()`, rewrite in `Write()`.

### 2. Import() only handled alpha
The 32-bit Import only compressed channel A (`data[i] = pic.A`), discarding RGB. Result: white/transparent screen in game.

**Fix**: Multi-mode import according to Colorbits (4, 8, 24, 32). The 32-bit mode makes a direct copy of `pic.Pix` (RGBA).

### 3. Colorbits > 32 (8-bit palette)
CZ1 palette files use Colorbits=248 (0xF8), a proprietary Visual Art's marker. LuckSystem did not recognize it → palette ignored → `GetOutputInfo()` read the palette as a block table → crash (`slice bounds out of range`).

**Fix**: Normalization `if Colorbits > 32 → Colorbits = 8` (same approach as lbee-utils).

### 4. Non-CZ files in PAKs
Files without the “CZ” magic number (e.g., トーンカーブ_夕/夜, 768-byte LUTs) caused a `glog.Fatalln(“Unknown Cz image type”)`.

**Fix**: Check the magic number.






========
>>>>>>>> 3c2a0411b72a175df3d2a802bb14d7d721cd6427:Yoremi V2-CHANGELOG-ENG.md
# LuckSystem — Yoremi-Version 2

Fork of [LuckSystem 2.3.2](https://github.com/wetor/LuckSystem) with fixes and additions for visual novel translation support on the ProtoDB/LUCA System engine (AIR, CLANNAD, Kanon, Summer Pockets, Harmonia, etc.).

## Patches

### Patch 1 — Variable-length script import
**File:** `script/script.go`

Importing translated scripts crashed with a panic when the translation had a different length than the original. The code strictly checked `len(paramList) == len(code.Params)`, blocking any longer or shorter translation.

- Removed strict parameter count check
- Added bounds checking in the conversion loop and parameter merge
- Jump offsets (GOTO, IFN, IFY…) are automatically recalculated

### Patch 2 — CZ3 pipeline fixes (PNG export/import)
**Files:** `czimage/cz3.go`, `czimage/imagefix.go`

CZ3 export and import silently corrupted pixel data.

- **Magic byte**: `Write()` overwrote the magic from "CZ3" to "CZ0", making the file unreadable by the game
- **NRGBA format**: Automatic conversion of any PNG format to NRGBA 32-bit before encoding
- **Buffer aliasing**: `DiffLine()` and `LineDiff()` shared slices instead of copying, causing delta data corruption

### Patch 3 — LZW decompressor memory corruption
**File:** `czimage/lzw.go`

The LZW decompressor added dictionary entries that directly referenced the `w` slice instead of making a copy. Old dictionary entries pointed to corrupted data.

- Explicit allocation of `newEntry` with copy of `w` before adding to dictionary

### Patch 4 — Incorrect RawSize in CZ block table
**File:** `czimage/util.go`

Critical bug causing visual CG corruption in-game (color artifacts). `Compress()` and `Compress2()` computed incorrect `RawSize` for each LZW block.

1. **Uncompensated carry-over**: The last LZW element carried to the next block was not deducted from the byte counter.
2. **Go UTF-8 encoding**: `len(string(byte(200)))` returns 2 instead of 1 for bytes > 127, causing ±1 errors on RawSize.

### Patch 5 — CZ4 image format support
**Files:** `czimage/cz4.go` (new), `czimage/imagefix.go`, `czimage/cz.go`

Added CZ4 format decoding and encoding, used in newer games (Little Busters English, LOOPERS, Harmonia, Kanon 2024).

CZ4 differs from CZ3 by storing RGB (w×h×3) and Alpha (w×h) channels separately, each with independent delta line encoding. LZW compression and blockHeight calculation are identical to CZ3.

### Patch 6 — PAK block alignment padding
**File:** `pak/pak.go`

After writing a rebuilt PAK (when replaced files are larger than originals), the file was not aligned to block size, potentially causing read errors.

- Added zero padding at end of file to align to `BlockSize`

### Patch 7 — AIR.py module resolution fix
**File:** `data/AIR.py`

The AIR.py definition script used `from base.air import *` to import functions from `data/base/air.py`. This import consistently failed in `script import` mode because LuckSystem's working directory is not `data/`, preventing Python from resolving the `base/air` path.

Error reproduced with the command documented in usage.md:
```
FileNotFoundError: 'Failed to resolve "base/air"'
panic: runtime error: invalid memory address or nil pointer dereference
```

- Merged all `base/air.py` functions directly into `AIR.py` (IFN, IFY, FARCALL, GOTO, GOSUB, JUMP, etc.)
- Added the missing `ONGOTO` opcode handler
- Removed the `from base.air import *` dependency

## Tested games

- AIR (Steam) — full French translation pipeline, SYSCG.pak 51/51 (CZ3+CZ4), SCRIPT.pak import/export
- Summer Pockets — RawSize fix confirmed (masagrator report)

## Credits

- **wetor** — original LuckSystem
- **masagrator** — RawSize bug identification (CZ3 layers)
- **Yoremi** — patches 1-7, AIR French translation
- **G2-Games** — CZ4 reference (lbee-utils)
