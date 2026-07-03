# LOOPERS proxy DLL port report

## Target

- Game: LOOPERS Steam
- Exe: `C:\Program Files (x86)\Steam\steamapps\common\LOOPERS\LOOPERS.exe`
- Proxy: `version.dll`
- Architecture: x64
- `.rdata` raw -> RVA delta: `0x1400`

## Result

- `patches.py` validates 97 UI string patches.
- The table is based on the HarmoniaHD proxy table, with LOOPERS-specific
  offsets and a few source string variants.
- Additional LOOPERS options covered:
  - `Window Frame Color` -> `Couleur cadre`
  - `Dialogue Only` -> `Dialogue seul`
  - `Low resolution` -> `Basse res.`
  - `Movie Quality` -> `Qualite video`

## LOOPERS-specific variants

- `Initial Cursor Position` is `Initial cursor position`.
- `Initial Cursor Position: ...` is `Initial cursor position: ...`.
- The voice tooltip has no space before `(It stops...)`.
- `Wheel Up` / `Wheel Down` are `Wheel up` / `Wheel down`.
- `Jump and Switch Pages` is `Jump and switch pages`.
- `Snap Pointer` is `Snap pointer`.

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
