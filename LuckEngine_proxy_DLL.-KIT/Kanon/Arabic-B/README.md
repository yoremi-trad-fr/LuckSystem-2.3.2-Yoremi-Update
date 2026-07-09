# Kanon Arabic-B proxy DLL

Experimental Arabic system-menu proxy for Kanon Steam.

This variant uses Arabic-B text: logical Arabic is converted to Arabic Presentation Forms and reversed per line so Luck Engine's left-to-right renderer can display it.

## Files

- `version.dll` - built proxy DLL to place next to `Kanon.exe`.
- `patches.py` - Arabic-B translation source.
- `patches.csv` - review table with slot budgets and final patched strings.
- `patches.h` - generated C patch table.

## Notes

- Many exe slots are very small, so some labels are abbreviated.
- Very long tooltips and technical strings are intentionally left in English for this first test.
- This does not modify `Kanon.exe` on disk; strings are patched in memory after SteamStub finishes.
- If text appears as boxes or blank glyphs, the menu is using a font size/family that does not contain the injected Arabic Presentation Forms.

## Install

Back up any existing `version.dll`, then copy this folder's `version.dll` next to `Kanon.exe` and launch through Steam.

For logs, set this Steam launch option:

```text
LUCKPROXY_LOG=1 %command%
```
