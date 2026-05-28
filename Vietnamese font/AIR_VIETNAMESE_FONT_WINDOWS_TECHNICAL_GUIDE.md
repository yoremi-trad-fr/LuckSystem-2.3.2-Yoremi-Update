# AIR Vietnamese Font Technical Guide - Windows

This guide explains how to build the LuckSystem-Yoremi command-line tools on Windows and how to test another TTF for AIR Vietnamese fonts with `fontdiag` and `vietfontpatch`.

It is intended for technical users who are comfortable using PowerShell.

## What The Tools Do

### `fontdiag`

`fontdiag` is a diagnostic roundtrip tool.

It reads an AIR font family PAK, rewrites it without changing the charset, and outputs a rebuilt PAK. Use it to check whether the PAK/CZ2 rebuild path is safe before testing a new TTF.

It does not inject Vietnamese glyphs.

### `vietfontpatch`

`vietfontpatch` is the automated AIR Vietnamese font patcher.

It reads AIR font PAKs, reads a Vietnamese charset file, detects which characters already exist in AIR, injects only the missing characters, applies a vertical `Y` offset, and writes rebuilt PAK files.

Use this tool when testing:

- another TTF,
- another `Y+` value,
- only `GOTHIC1` for quick checks,
- or all English-slot font families for a fuller test.

## Requirements

Install these first:

- Windows 10 or Windows 11.
- PowerShell.
- Go installed and available in `PATH`.
- Optional: Node.js + npm + Wails CLI, only if you also want to rebuild the GUI.

The core CLI module declares Go `1.16`, but using a recent Go version is fine. If you also rebuild the GUI, use Go `1.23` or newer because the GUI module declares Go `1.23`.

Check Go:

```powershell
go version
```

## Folder Example

The commands below use this example layout:

```text
C:\AIR_FONT_TEST\
  LuckSystem-Yoremi\          source code
  ttf\                        test TTF files
  out\                        generated test outputs
```

The AIR font root must be the game `files` folder, not the Steam root folder.

Correct:

```text
C:\...\AIR\files
```

It must contain:

```text
font_win32_1280\FONT__INFO.PAK
font_win32_1280\FONT_GOTHIC1.PAK
```

For the Chinese `ZC` slot, it would also contain:

```text
fontzc_win32_1280\FONTZC__INFO.PAK
fontzc_win32_1280\FONTZC_GOTHIC1.PAK
```

The validated workflow uses the English slot only.

## Build The CLI And Tools

Open PowerShell in the repository root:

```powershell
cd "C:\AIR_FONT_TEST\LuckSystem-Yoremi"
```

Download dependencies:

```powershell
go mod download
```

Run a quick compile/test check:

```powershell
go test ./...
```

Create a build folder:

```powershell
New-Item -ItemType Directory -Force ".\build"
```

Build the normal LuckSystem CLI:

```powershell
go build -o ".\build\lucksystem.exe" .
```

Build the diagnostic tool:

```powershell
go build -o ".\build\fontdiag.exe" ./tools/fontdiag
```

Build the Vietnamese patch tool:

```powershell
go build -o ".\build\vietfontpatch.exe" ./tools/vietfontpatch
```

You should now have:

```text
build\lucksystem.exe
build\fontdiag.exe
build\vietfontpatch.exe
```

## Optional - Build The GUI

This is only needed if you want to rebuild `LuckSystemGUI.exe`.

Install Wails:

```powershell
go install github.com/wailsapp/wails/v2/cmd/wails@v2.11.0
```

Build the GUI:

```powershell
cd "C:\AIR_FONT_TEST\LuckSystem-Yoremi\SourcesGUI-wails"
go mod download
cd ".\frontend"
npm install
cd ".."
wails build
```

The GUI executable is created here:

```text
SourcesGUI-wails\build\bin\LuckSystemGUI.exe
```

For normal GUI use, place `LuckSystemGUI.exe` next to `lucksystem.exe`.

## Charset Files

For `vietfontpatch`, use the full Vietnamese charset. The tool will automatically remove characters already present in AIR and inject only the missing ones.

The repository includes:

```text
examples\AIR_vietnamese_full_134.txt
```

This file contains 134 requested Vietnamese characters:

- 32 already present in AIR,
- 102 missing and injected by the tool.

The GUI-only guide uses:

```text
examples\AIR_vietnamese_missing_102.txt
```

That file is mainly for manual GUI `Insert at index` mode. For `vietfontpatch`, prefer the full 134-character file.

## Step 1 - Back Up Original AIR Font PAKs

Set the AIR `files` folder path:

```powershell
$AirFiles = "C:\Games\Steam\steamapps\common\AIR\files"
```

Back up the English-slot files used for a quick `GOTHIC1` test:

```powershell
Copy-Item "$AirFiles\font_win32_1280\FONT__INFO.PAK" "$AirFiles\font_win32_1280\FONT__INFO.PAK.bak" -Force
Copy-Item "$AirFiles\font_win32_1280\FONT_GOTHIC1.PAK" "$AirFiles\font_win32_1280\FONT_GOTHIC1.PAK.bak" -Force
```

Restore them later with:

```powershell
Copy-Item "$AirFiles\font_win32_1280\FONT__INFO.PAK.bak" "$AirFiles\font_win32_1280\FONT__INFO.PAK" -Force
Copy-Item "$AirFiles\font_win32_1280\FONT_GOTHIC1.PAK.bak" "$AirFiles\font_win32_1280\FONT_GOTHIC1.PAK" -Force
```

## Step 2 - Run A Roundtrip Diagnostic

Before testing another TTF, confirm that a simple rebuild of `FONT_GOTHIC1.PAK` works.

From the repository root:

```powershell
cd "C:\AIR_FONT_TEST\LuckSystem-Yoremi"
```

Run:

```powershell
.\build\fontdiag.exe "$AirFiles" "C:\AIR_FONT_TEST\out\roundtrip" "FONT_GOTHIC1.PAK"
```

Expected output:

```text
C:\AIR_FONT_TEST\out\roundtrip\ROUNDTRIP_FONT_GOTHIC1.PAK
```

To test this in game, copy it over the original family PAK:

```powershell
Copy-Item "C:\AIR_FONT_TEST\out\roundtrip\ROUNDTRIP_FONT_GOTHIC1.PAK" "$AirFiles\font_win32_1280\FONT_GOTHIC1.PAK" -Force
```

Start AIR. If the game starts and menus are not corrupted, the rebuild path is OK.

Restore the original after the roundtrip test:

```powershell
Copy-Item "$AirFiles\font_win32_1280\FONT_GOTHIC1.PAK.bak" "$AirFiles\font_win32_1280\FONT_GOTHIC1.PAK" -Force
```

## Step 3 - Patch `GOTHIC1` With Another TTF

Set paths:

```powershell
$Repo = "C:\AIR_FONT_TEST\LuckSystem-Yoremi"
$AirFiles = "C:\Games\Steam\steamapps\common\AIR\files"
$Charset = "$Repo\examples\AIR_vietnamese_full_134.txt"
$Ttf = "C:\AIR_FONT_TEST\ttf\YourFont.ttf"
$Out = "C:\AIR_FONT_TEST\out\YourFont_GOTHIC1_y2"
```

Generate a quick English-slot `GOTHIC1` test with `Y+2`:

```powershell
cd "$Repo"
.\build\vietfontpatch.exe -slot en -family GOTHIC1 -yoffset 2 "$AirFiles" "$Charset" "$Ttf" "$Out"
```

Expected output files:

```text
C:\AIR_FONT_TEST\out\YourFont_GOTHIC1_y2\FONT__INFO.PAK
C:\AIR_FONT_TEST\out\YourFont_GOTHIC1_y2\FONT_GOTHIC1.PAK
```

Copy the generated files into AIR:

```powershell
Copy-Item "$Out\FONT__INFO.PAK" "$AirFiles\font_win32_1280\FONT__INFO.PAK" -Force
Copy-Item "$Out\FONT_GOTHIC1.PAK" "$AirFiles\font_win32_1280\FONT_GOTHIC1.PAK" -Force
```

Start AIR and check Vietnamese text in the English slot.

## How To Change The `Y+` Value

The vertical adjustment is controlled by:

```text
-yoffset N
```

Examples:

```powershell
.\build\vietfontpatch.exe -slot en -family GOTHIC1 -yoffset 0 "$AirFiles" "$Charset" "$Ttf" "C:\AIR_FONT_TEST\out\YourFont_y0"
.\build\vietfontpatch.exe -slot en -family GOTHIC1 -yoffset 1 "$AirFiles" "$Charset" "$Ttf" "C:\AIR_FONT_TEST\out\YourFont_y1"
.\build\vietfontpatch.exe -slot en -family GOTHIC1 -yoffset 2 "$AirFiles" "$Charset" "$Ttf" "C:\AIR_FONT_TEST\out\YourFont_y2"
.\build\vietfontpatch.exe -slot en -family GOTHIC1 -yoffset 3 "$AirFiles" "$Charset" "$Ttf" "C:\AIR_FONT_TEST\out\YourFont_y3"
```

For the previously validated AIR test, `Y+2` was the best result:

```text
-yoffset 2
```

If the Vietnamese marks are still too high or too low, try:

```text
-yoffset -2
-yoffset -1
-yoffset 0
-yoffset 1
-yoffset 2
-yoffset 3
```

Use a different output folder for each value.

## Generate Several `Y+` Tests At Once

This PowerShell loop generates multiple `GOTHIC1` test folders:

```powershell
$Repo = "C:\AIR_FONT_TEST\LuckSystem-Yoremi"
$AirFiles = "C:\Games\Steam\steamapps\common\AIR\files"
$Charset = "$Repo\examples\AIR_vietnamese_full_134.txt"
$Ttf = "C:\AIR_FONT_TEST\ttf\YourFont.ttf"

cd "$Repo"

foreach ($Y in -2,-1,0,1,2,3) {
    $Out = "C:\AIR_FONT_TEST\out\YourFont_GOTHIC1_y$Y"
    .\build\vietfontpatch.exe -slot en -family GOTHIC1 -yoffset $Y "$AirFiles" "$Charset" "$Ttf" "$Out"
}
```

Then test one output folder at a time by copying its two files:

```powershell
$Out = "C:\AIR_FONT_TEST\out\YourFont_GOTHIC1_y2"
Copy-Item "$Out\FONT__INFO.PAK" "$AirFiles\font_win32_1280\FONT__INFO.PAK" -Force
Copy-Item "$Out\FONT_GOTHIC1.PAK" "$AirFiles\font_win32_1280\FONT_GOTHIC1.PAK" -Force
```

## Generate All English-Slot Font Families

Once `GOTHIC1` looks good, generate all English-slot families with the chosen TTF and `Y` value:

```powershell
$Out = "C:\AIR_FONT_TEST\out\YourFont_EN_all_y2"
.\build\vietfontpatch.exe -slot en -family all -yoffset 2 "$AirFiles" "$Charset" "$Ttf" "$Out"
```

Expected output:

```text
FONT__INFO.PAK
FONT_GOTHIC1.PAK
FONT_GOTHIC2.PAK
FONT_GOTHIC3.PAK
FONT_MINCHO.PAK
FONT_MODERN.PAK
```

Copy all generated English-slot font PAKs:

```powershell
Copy-Item "$Out\FONT*.PAK" "$AirFiles\font_win32_1280" -Force
```

## Available Slot And Family Options

Slots:

```text
-slot en    English/Japanese font slot, uses font_win32_1280
-slot zc    Chinese ZC font slot, uses fontzc_win32_1280
-slot all   Both slots
```

For the current AIR Vietnamese workflow, use:

```text
-slot en
```

Families:

```text
-family GOTHIC1
-family GOTHIC2
-family GOTHIC3
-family MINCHO
-family MODERN
-family all
```

For fast tests, use:

```text
-family GOTHIC1
```

For final English-slot output, use:

```text
-family all
```

## Testing Another TTF

For each new TTF:

1. Put the TTF in a simple path, for example `C:\AIR_FONT_TEST\ttf\FontName.ttf`.
2. Run `fontdiag` once if you have not already confirmed the rebuild path.
3. Run `vietfontpatch` with `-slot en -family GOTHIC1`.
4. Test several `-yoffset` values.
5. Copy only `FONT__INFO.PAK` and `FONT_GOTHIC1.PAK` for quick tests.
6. When satisfied, run `-family all`.

Example:

```powershell
$Ttf = "C:\AIR_FONT_TEST\ttf\AnotherFont.ttf"
$Out = "C:\AIR_FONT_TEST\out\AnotherFont_GOTHIC1_y2"

.\build\vietfontpatch.exe -slot en -family GOTHIC1 -yoffset 2 "$AirFiles" "$Charset" "$Ttf" "$Out"
Copy-Item "$Out\FONT__INFO.PAK" "$AirFiles\font_win32_1280\FONT__INFO.PAK" -Force
Copy-Item "$Out\FONT_GOTHIC1.PAK" "$AirFiles\font_win32_1280\FONT_GOTHIC1.PAK" -Force
```

## Troubleshooting

### The game does not start

Restore the original backups:

```powershell
Copy-Item "$AirFiles\font_win32_1280\FONT__INFO.PAK.bak" "$AirFiles\font_win32_1280\FONT__INFO.PAK" -Force
Copy-Item "$AirFiles\font_win32_1280\FONT_GOTHIC1.PAK.bak" "$AirFiles\font_win32_1280\FONT_GOTHIC1.PAK" -Force
```

Then regenerate with a clean output folder.

### `cannot find the path specified`

Check that `$AirFiles` points to the `files` folder:

```text
...\AIR\files
```

not:

```text
...\AIR
```

### `font file does not contain character`

The TTF does not support one of the Vietnamese characters. Try another TTF.

### `file count mismatch`

The info PAK and family PAK do not belong to the same AIR font folder, or one of them is not original. Restore original PAKs and try again.

### Vietnamese characters appear but are vertically wrong

Regenerate with a different `-yoffset` value.

Recommended quick range:

```text
-2, -1, 0, 1, 2, 3
```

### Existing accents look worse

Make sure you used the full charset with `vietfontpatch`, not manual GUI `Redraw all`. The tool keeps already-present characters mapped to their original glyphs and injects only missing characters.

## Quick Command Summary

Build:

```powershell
cd "C:\AIR_FONT_TEST\LuckSystem-Yoremi"
go mod download
go test ./...
New-Item -ItemType Directory -Force ".\build"
go build -o ".\build\lucksystem.exe" .
go build -o ".\build\fontdiag.exe" ./tools/fontdiag
go build -o ".\build\vietfontpatch.exe" ./tools/vietfontpatch
```

Roundtrip diagnostic:

```powershell
.\build\fontdiag.exe "$AirFiles" "C:\AIR_FONT_TEST\out\roundtrip" "FONT_GOTHIC1.PAK"
```

Patch one TTF with `Y+2`:

```powershell
.\build\vietfontpatch.exe -slot en -family GOTHIC1 -yoffset 2 "$AirFiles" "$Charset" "$Ttf" "$Out"
```

Copy quick-test output:

```powershell
Copy-Item "$Out\FONT__INFO.PAK" "$AirFiles\font_win32_1280\FONT__INFO.PAK" -Force
Copy-Item "$Out\FONT_GOTHIC1.PAK" "$AirFiles\font_win32_1280\FONT_GOTHIC1.PAK" -Force
```
