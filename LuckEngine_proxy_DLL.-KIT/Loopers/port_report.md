# LOOPERS proxy DLL port report

## Target

- Game: LOOPERS Steam
- Exe: `C:\Program Files (x86)\Steam\steamapps\common\LOOPERS\LOOPERS.exe`
- Proxy: `version.dll`
- Architecture: x64
- `.rdata` raw -> RVA delta: `0x1400`

## Result

- `patches.py` validates 107 UI string patches.
- The table is based on the HarmoniaHD proxy table, with LOOPERS-specific
  offsets and a few source string variants.
- Additional LOOPERS options covered:
  - `Window Frame Color` -> `Couleur cadre`
  - `Dialogue Only` -> `Dialogue seul`
  - `Low resolution` -> `Basse res.`
  - `Movie Quality` -> `Qualite video`
  - Duplicate Mouse action-list entries for `Disable`, `Hide Window`,
    `System Menu`, and `Message Log`
  - Touch labels `Rewind`, `Jump (Forward)`, and `Jump (Backward)`
  - Touch help text for skip, rewind, and jump/swipe actions

## LOOPERS-specific variants

- `Initial Cursor Position` is `Initial cursor position`.
- `Initial Cursor Position: ...` is `Initial cursor position: ...`.
- The voice tooltip has no space before `(It stops...)`.
- `Wheel Up` / `Wheel Down` are `Wheel up` / `Wheel down`.
- `Jump and Switch Pages` is `Jump and switch pages`.
- `Snap Pointer` is `Snap pointer`.
- The Mouse action list contains a second cluster of `Disable`, `Hide Window`,
  `System Menu`, and `Message Log` strings used by wheel/button choices.
- Touch control labels use separate standalone `Rewind`, `Jump (Forward)`, and
  `Jump (Backward)` strings.
- Touch tooltips for skip, rewind, and outside-window swipe are separate from
  the visible labels.

## Deliberately skipped

These HarmoniaHD entries were not ported because the same exact UI string is
not present as an English slot in LOOPERS:

- `Basic`
- `Text1`
- `Text2`
- `Sound`
- `Mouse`
- `System`
- `Orange`

`Mouse` appears only inside `Mouse Wheel Button`, so it must not be patched as a
tab label.

`Off` was also skipped: in LOOPERS the matching string sits beside title/music
labels (`Tautology`, `Special Spell`), not in the options UI block used by the
proxy.

`Read Text Color` uses the shortened French target `Coul. texte lu` because the
LOOPERS slot budget is 15 bytes.
