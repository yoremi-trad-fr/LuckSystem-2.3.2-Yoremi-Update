LuckSystem 2.3.2 — Yoremi Fork

Fork de LuckSystem avec corrections de bugs et ajouts pour la traduction de visual novels Visual Art's/Key.

Fork of LuckSystem with bug fixes and additions for Visual Art's/Key visual novel translation.
Supported engines / Moteurs supportés

ProtoDB / LUCA System — AIR, CLANNAD, Kanon, Little Busters, Summer Pockets, Harmonia, LOOPERS, LUNARiA, Planetarian, etc.
Features / Fonctionnalités
Feature 	Format 	Status
Script decompile/compile 	SCRIPT.PAK 	✅
CZ0 image export 	CZ0 	✅
CZ1 image export/import (32-bit + 8-bit palette) 	CZ1 	✅
CZ3 image export/import 	CZ3 	✅
CZ4 image export/import 	CZ4 	✅ (new)
Font extract/edit 	FONT.PAK (CZ2) 	✅
Font edit append/insert (CZ2 resize fix) 	FONT.PAK (CZ2) 	✅ (v3p3)
PAK extract/replace 	*.PAK 	✅
Patches
Version 2 (7 patches)

    Variable-length script import — script/script.go
    CZ3 pipeline fixes (magic byte, NRGBA, buffer aliasing) — czimage/cz3.go, imagefix.go
    LZW decompressor memory corruption — czimage/lzw.go
    RawSize carry-over + UTF-8 length — czimage/util.go
    CZ4 format support (new) — czimage/cz4.go
    PAK block alignment padding — pak/pak.go
    AIR.py module resolution — data/AIR.py

Version 3 — Patch 1

    CZ1 32-bit Import/Export rewrite — czimage/cz1.go
    CZ1 8-bit palette support (Colorbits > 32 normalization) — czimage/cz1.go
    Non-CZ files graceful handling — czimage/cz.go
    CZ0 logging visibility — czimage/cz0.go

Version 3 — Patch 3

    CZ2 font import resize fix — czimage/cz2.go, font/font.go
        Import(): update CzHeader dimensions instead of silent nil return when image is resized
        Added SetDimensions() method on Cz2Image
        Write(): sync header before Import() call

Merged upstream (PR #35)

    CZ2 font decompressor crash fix — czimage/lzw.go (boundary check in decompressLZW2)

Documentation
Document 	Description
CHANGELOG.md 	Full changelog — all versions (EN + FR)
TECHNICAL.md 	Technical analysis — all patches
GUI

A graphical interface is available in a separate repository: **LuckSystem GUI — Wails + Svelte wrapper, calls lucksystem.exe via subprocess.
Tested games / Jeux testés

    AIR (Steam) — French translation complete (scripts + CG + UI)
    Summer Pockets — RawSize fix confirmed
    Kanon — CZ2 font fix confirmed
    Little Busters English — CZ4 confirmed

Usage

See Usage.md for CLI commands.

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

Credits

    wetor — LuckSystem original
    masagrator — RawSize bug identification (CZ3 layers)
    G2-Games — CZ4 reference (lbee-utils)
    Yoremi — patches 1-12, AIR French translation, GUI
