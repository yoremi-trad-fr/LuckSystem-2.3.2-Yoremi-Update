# Vietnamese Font Patch GUI Guide

This guide explains how to generate Vietnamese font PAKs for AIR / Planetarian SG using the dedicated GUI workflow.

No command line is required.

## Important

Use the new GUI page:

```text
VIET FONT -> AIR / SG Patch
```

Do not use an old separately compiled `vietnamesefont.exe` / `vietfontpatch.exe` if it produces broken menus. The GUI workflow includes the corrected patch code directly.

## What This GUI Page Does

The GUI automatically:

- reads the original game font PAKs,
- keeps the Vietnamese-compatible characters already present in the game,
- injects only the missing Vietnamese characters,
- preserves the special `FONT__INFO.PAK` layout required by AIR / Planetarian SG,
- lets you test several vertical `Y` alignments,
- writes ready-to-test `FONT__INFO.PAK` and `FONT_GOTHIC1.PAK` files.

## Required Files

You need:

- the game `files` folder,
- a Vietnamese-capable `.ttf` or `.otf` font,
- the full Vietnamese charset file:

```text
examples\AIR_vietnamese_full_134.txt
```

Use the full 134-character charset file with this GUI page.

Do not use the 102-character missing-only charset here. That one is only for manual `Font Edit -> Insert at index` tests.

## Step 1 - Back Up The Original PAKs

Before testing, back up:

```text
files\font_win32_1280\FONT__INFO.PAK
files\font_win32_1280\FONT_GOTHIC1.PAK
```

For example, copy them somewhere safe or rename them:

```text
FONT__INFO.PAK.bak
FONT_GOTHIC1.PAK.bak
```

## Step 2 - Open The GUI Page

Start `LuckSystemGUI.exe`.

In the left menu, select:

```text
VIET FONT -> AIR / SG Patch
```

## Step 3 - Select Inputs

### Game Files Folder

Select the game `files` folder.

Correct examples:

```text
C:\Games\AIR\files
C:\Games\Planetarian SG\files
```

The selected folder must contain:

```text
font_win32_1280
```

Do not select the `font_win32_1280` folder itself.

### Full Vietnamese Charset File

Select:

```text
examples\AIR_vietnamese_full_134.txt
```

This file contains all 134 requested Vietnamese characters. The GUI will automatically keep the 32 characters already present and inject only the missing 102.

### TTF / OTF Font File

Select the font you want to test.

If the GUI reports that the font does not contain a character, use another TTF/OTF.

### Output Folder

Select an empty output folder, for example:

```text
C:\AIR_FONT_TEST\out
```

The GUI creates one subfolder per selected `Y` value.

## Step 4 - Recommended First Test

Use these settings first:

```text
Target slot: English slot
Family: GOTHIC1 quick test
Y alignment: Y+2
```

Then click:

```text
Generate Vietnamese Font PAKs
```

Expected output folder example:

```text
Arial-Unicode-MS_en_GOTHIC1_Y+2
```

Inside it, you should get:

```text
FONT__INFO.PAK
FONT_GOTHIC1.PAK
```

## Step 5 - Test In Game

Copy the generated files into:

```text
files\font_win32_1280
```

Replace:

```text
FONT__INFO.PAK
FONT_GOTHIC1.PAK
```

Start the game and check:

- the menu,
- a translated Vietnamese dialogue line,
- Vietnamese tone marks alignment.

## Step 6 - Test Several Y Values

If the font appears too high or too low, generate several values at once:

```text
Y-2
Y-1
Y+0
Y+1
Y+2
Y+3
```

The GUI creates separate folders for each value. Test them one by one by copying the two PAK files from the selected output folder into:

```text
files\font_win32_1280
```

The previously validated AIR value was:

```text
Y+2
```

Planetarian SG may need the same value or a nearby one.

## Step 7 - Generate More Families

After `GOTHIC1` works, you can generate all English-slot families:

```text
Target slot: English slot
Family: All English families
Y alignment: the best tested value
```

This creates:

```text
FONT__INFO.PAK
FONT_GOTHIC1.PAK
FONT_GOTHIC2.PAK
FONT_GOTHIC3.PAK
FONT_MINCHO.PAK
FONT_MODERN.PAK
```

Copy them into:

```text
files\font_win32_1280
```

For quick dialogue tests, `FONT__INFO.PAK` + `FONT_GOTHIC1.PAK` is usually enough.

## Troubleshooting

### The menu is corrupted

Restore the original backup PAKs first.

This usually means the generated `FONT__INFO.PAK` is not using the required AIR / Planetarian SG info layout. Use the new GUI page, not an old standalone tool.

### Dialogue text disappears

Check that:

- both files were copied together,
- `FONT__INFO.PAK` and `FONT_GOTHIC1.PAK` come from the same generated output folder,
- the selected game folder was the correct `files` folder,
- the full 134-character charset was used.

### Missing Vietnamese characters remain missing

Check that your game script is using the English slot and that the generated PAKs were copied to:

```text
files\font_win32_1280
```

not to:

```text
files\fontzc_win32_1280
```

### The game does not start

Restore:

```text
FONT__INFO.PAK.bak
FONT_GOTHIC1.PAK.bak
```

Then generate again from clean original font PAKs.

### Some marks are clipped

Try another TTF/OTF. Some fonts technically contain Vietnamese characters but have metrics that do not fit well in the game font cells.

## Safe Rule

For first tests, always use:

```text
English slot
GOTHIC1 quick test
Y+2
Full 134-character charset
```

Then adjust only one thing at a time: first `Y`, then TTF, then more font families.
