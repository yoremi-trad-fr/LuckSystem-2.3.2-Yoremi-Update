# LuckSystem 2.3.2 — Yoremi Fork

Fork de [LuckSystem](https://github.com/wetor/LuckSystem) avec corrections de bugs et ajouts pour la traduction de visual novels Visual Art's/Key.

Fork of [LuckSystem](https://github.com/wetor/LuckSystem) with bug fixes and additions for Visual Art's/Key visual novel translation.

---

## Supported engines / Moteurs supportés

ProtoDB / LUCA System — AIR, CLANNAD, Kanon, Little Busters, Summer Pockets, Harmonia, LOOPERS, LUNARiA, Planetarian, etc.

## Features / Fonctionnalités

| Feature | Format | Status |
|---------|--------|--------|
| Script decompile/compile | SCRIPT.PAK | ✅ |
| CZ0 image export | CZ0 | ✅ |
| CZ1 image export/import (32-bit + 8-bit palette) | CZ1 | ✅ |
| CZ3 image export/import | CZ3 | ✅ |
| CZ4 image export/import | CZ4 | ✅ (new) |
| Font extract/edit | FONT.PAK (CZ2) | ✅ |
| PAK extract/replace | *.PAK | ✅ |

## Patches

### Version 2 (7 patches)

1. **Variable-length script import** — `script/script.go`
2. **CZ3 pipeline fixes** (magic byte, NRGBA, buffer aliasing) — `czimage/cz3.go`, `imagefix.go`
3. **LZW decompressor memory corruption** — `czimage/lzw.go`
4. **RawSize carry-over + UTF-8 length** — `czimage/util.go`
5. **CZ4 format support** (new) — `czimage/cz4.go`
6. **PAK block alignment padding** — `pak/pak.go`
7. **AIR.py module resolution** — `data/AIR.py`

### Version 3 — Patch 1

8. **CZ1 32-bit Import/Export rewrite** — `czimage/cz1.go`
9. **CZ1 8-bit palette support** (Colorbits > 32 normalization) — `czimage/cz1.go`
10. **Non-CZ files graceful handling** — `czimage/cz.go`
11. **CZ0 logging visibility** — `czimage/cz0.go`

### Merged upstream (PR #35)

12. **CZ2 font decompressor crash fix** — `czimage/lzw.go` (boundary check in `decompressLZW2`)

## Documentation

| Document | Description |
|----------|-------------|

| [Yoremi V3-CHANGELOG-FR.md](Yoremi%20V3-CHANGELOG-FR.md) | Changelog V3 Patch 1 (français) |
| [Yoremi V3-CHANGELOG-ENG.md](Yoremi%20V3-CHANGELOG-ENG.md) | V3 Patch 1 changelog (English) |

## GUI

A graphical interface is available in a separate repository:
**[LuckSystem GUI](https://github.com/yoremi-trad-fr/LuckSystemGUI)** — Wails + Svelte wrapper, calls `lucksystem.exe` via subprocess.

## Tested games / Jeux testés

- **AIR** (Steam) — French translation complete (scripts + CG + UI)
- **Summer Pockets** — RawSize fix confirmed
- **Kanon** — CZ2 font fix confirmed

## Usage

See [Usage.md](Usage.md) for CLI commands.

```bash
# Decompile scripts
lucksystem script decompile -s SCRIPT.PAK -c UTF-8 -O data/AIR.txt -p data/AIR.py -o Export

# Import translated scripts
lucksystem script import -s SCRIPT.PAK -c UTF-8 -O data/AIR.txt -p data/AIR.py -i Export -o SCRIPT_FR.PAK

# Export CZ image to PNG
lucksystem image export -i image.cz3 -o image.png

# Extract FONT.PAK
lucksystem pak extract -i FONT.PAK -o list.txt --all ./fonts/

# Edit font with TTF
lucksystem font edit -s 明朝32 -S info32 -f Arial.ttf -o 明朝32_out -O info32_out -r
```

## Credits

- **[wetor](https://github.com/wetor)** — LuckSystem original
- **masagrator** — RawSize bug identification (CZ3 layers)
- **[G2-Games](https://github.com/G2-Games)** — CZ4 reference ([lbee-utils](https://github.com/G2-Games/lbee-utils))
- **Yoremi** — patches 1-12, AIR French translation, GUI
