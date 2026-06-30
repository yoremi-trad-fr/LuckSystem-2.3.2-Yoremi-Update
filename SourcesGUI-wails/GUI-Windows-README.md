# LuckSystem GUI (Windows) — Yoremi fork v3.22

Graphical interface for [LuckSystem](https://github.com/wetor/LuckSystem), the Visual Art's/Key visual novel translation toolkit.

Interface graphique pour [LuckSystem](https://github.com/wetor/LuckSystem), l'outil de traduction de visual novels Visual Art's/Key.

![LuckSystem GUI](screenshot.png)

## Architecture

The GUI is a **standalone wrapper** — it does NOT embed LuckSystem source code. It calls `lucksystem.exe` via subprocess, exactly like you would from a terminal.

```
LuckSystemGUI.exe  ←→  lucksystem.exe (subprocess)
   (Wails/Go)              (CLI tool)
```

This design follows [wetor's recommendation](https://github.com/wetor/LuckSystem) to keep the GUI separated from the core tool for cross-platform compatibility and maintainability.

## Setup

1. Download `lucksystem.exe` from [LuckSystem releases](https://github.com/wetor/LuckSystem/releases) (or build from the [Yoremi fork](https://github.com/yoremi-trad-fr/LuckSystem-2.3.2-Yoremi-Update))
2. Place `lucksystem.exe` next to `LuckSystemGUI.exe`
3. Run `LuckSystemGUI.exe`

The GUI auto-detects `lucksystem.exe` in the same directory, current working directory, or system PATH. You can also manually locate it by clicking the path indicator in the title bar.

## Features

| Operation | Description |
|-----------|-------------|
| **Script Decompile** | Extract scripts from SCRIPT.PAK to text files |
| **Script Compile** | Repack translated scripts into a new SCRIPT.PAK |
| **Siglus -> Luca** | Import translated Siglus script text into Luca scripts and export Luca-only/review TSV files |
| **PAK Extract** | Extract all files from any .PAK archive |
| **PAK Replace** | Replace files inside a .PAK archive |
| **BGMOVIE Extract** | Extract Luca Engine BGMOVIE.PAK videos to WebM |
| **Font Extract** | Export CZ font atlas to PNG + charset list |
| **Font Edit** | Redraw/append characters using a TTF font |
| **Image Export** | Convert CZ images to PNG (single or batch) |
| **Image Import** | Convert PNG back to CZ format (single or batch) |
| **Dialogue Extract** | Extract translatable dialogue from decompiled scripts to TSV (single file or batch) |
| **Dialogue Import** | Reimport translated dialogue from TSV back into scripts (single file or batch) |

### BGMOVIE Extract

`BGMOVIE.PAK` is the Luca Engine video archive found across Visual Art's/Key games. The GUI extracts the raw `MVT` entries first, then writes the embedded VP9/WebM payloads to a `webm` subfolder.

### Siglus -> Luca

This workflow keeps decompiled Luca scripts as the master structure and imports matching translated lines from Siglus `.ss.txt` exports. Select the Luca scripts folder, the Siglus `Full` folder, an output folder, and the target language column. The GUI writes patched scripts plus `hd_candidates.tsv` and `review.tsv`.

### Dialogue Extract / Import

The Dialogue Extract and Import functions provide a streamlined translation workflow based on TSV files, replacing manual script editing.

**Extract** scans decompiled script files (`.txt`) for translatable lines (`MESSAGE`, `LOG_BEGIN`, and `SELECT` entries) and exports them to tab-separated `.tsv` files. The language columns are numbered (Lang 1, Lang 2, Lang 3, Lang 4) rather than named, since the order of languages varies between games. You select which columns to extract via checkboxes.

**Import** reads a translated `.tsv` file and reinjects the text back into the corresponding decompiled script. You select which column number contains the target language. Matching is done by sequential ID for robustness.

Both operations support single-file and batch modes. The format auto-detection scans the script to determine the number of available language columns.

Script decompile/compile also auto-selects the sibling Python plugin when the selected OPCODE file follows the standard `data/GAME.txt` + `data/GAME.py` layout. This avoids repacking translated script text with the generic fallback parser by mistake.

**TSV format example:**
```
ID	TAG	Lang 2
1	MESSAGE	`Rin@❝Stop bullying the weak!❞
2	MESSAGE	`Riki@❝Masato... where are you going?❞
3	LOG_BEGIN	Chapter 1
```

## Supported games

All games using the ProtoDB / LUCA System engine:
AIR, CLANNAD, Kanon, Little Busters, Summer Pockets, Harmonia, LOOPERS, LUNARiA, Planetarian, etc.

## Build from source

Requires: [Go 1.23+](https://go.dev/), [Node.js](https://nodejs.org/), [Wails CLI](https://wails.io/)

```bash
cd frontend && npm install && cd ..
go mod tidy
wails dev          # Development with hot-reload
wails build        # Build to build/bin/LuckSystemGUI.exe
```

> ⚠️ **Do NOT run `npm audit fix --force`** — it upgrades Svelte/Vite to incompatible major versions.

## Credits

- **[wetor](https://github.com/wetor)** — [LuckSystem](https://github.com/wetor/LuckSystem) core CLI tool
- **Yoremi** — GUI development, [Yoremi fork](https://github.com/yoremi-trad-fr/LuckSystem-2.3.2-Yoremi-Update) patches

## License

MIT
