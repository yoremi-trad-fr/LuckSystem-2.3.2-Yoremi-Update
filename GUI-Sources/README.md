# LuckSystem GUI - Yoremi fork v3

A graphical interface for [LuckSystem](https://github.com/wetor/LuckSystem), the Visual Art's/Key visual novel translation toolkit.

## Architecture

The GUI is a **standalone wrapper** around `lucksystem.exe`. It does NOT embed LuckSystem source code — it calls the CLI executable via subprocess, exactly like you would from a terminal.

```
LuckSystemGUI.exe  ←→  lucksystem.exe (subprocess)
   (Wails/Go)              (CLI tool)
```

## Setup

1. Build or download `lucksystem.exe` from the [LuckSystem releases](https://github.com/wetor/LuckSystem/releases)
2. Place `lucksystem.exe` next to `LuckSystemGUI.exe`
3. Run `LuckSystemGUI.exe`

The GUI auto-detects `lucksystem.exe` in:
- Same directory as the GUI
- Current working directory
- System PATH

You can also manually locate it by clicking the path indicator in the title bar.

## Features

| Operation | Description |
|-----------|-------------|
| **Script Decompile** | Extract scripts from SCRIPT.PAK to text files |
| **Script Compile** | Repack translated scripts into a new SCRIPT.PAK |
| **PAK Extract** | Extract all files from any .PAK archive |
| **PAK Replace** | Replace files inside a .PAK archive |
| **Font Extract** | Export CZ font atlas to PNG + charset list |
| **Font Edit** | Redraw/append characters using a TTF font |
| **Image Export** | Convert CZ images to PNG (single or batch) |
| **Image Import** | Convert PNG back to CZ format (single or batch) |

## Development

Requires: Go 1.23+, Node.js, [Wails CLI](https://wails.io/)

## Important notes

### Script Compile — folder selection
When selecting the translated scripts folder, point to the **parent folder** (e.g. `TRAD`), not to `TRAD\SCRIPT.PAK`. The tool automatically appends `\SCRIPT.PAK\` to the path internally.

Lors de la sélection du dossier de scripts traduits, pointez vers le **dossier parent** (ex: `TRAD`), pas vers `TRAD\SCRIPT.PAK`. L'outil ajoute automatiquement `\SCRIPT.PAK\` au chemin.

```bash
cd frontend && npm install && cd ..
go mod tidy
wails dev          # Development with hot-reload
wails build        # Build to build/bin/LuckSystemGUI.exe
```

**Important**: Do NOT run `npm audit fix --force` — it upgrades Svelte/Vite to incompatible versions.

## Credits

- **LuckSystem** by [wetor](https://github.com/wetor/LuckSystem) — the core CLI tool
- **GUI** by Yoremi — Wails + Svelte wrapper
