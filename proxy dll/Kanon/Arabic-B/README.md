# Kanon Arabic-B proxy DLL

Experimental mixed Arabic/English system-menu proxy for Kanon Steam.

This variant uses Arabic-B text where the menu font can display it. Places that render Arabic-B as blank use short English fallback text instead.

## Files

- `version.dll` - built proxy DLL to place next to `Kanon.exe`.
- `patches.py` - Arabic-B translation source.
- `patches.csv` - review table with slot budgets and final patched strings.
- `patches.h` - generated C patch table.

## Notes

- Many exe slots are very small, so some labels are abbreviated.
- This build patches the Japanese language slot only.
- English menu strings are intentionally left untouched, so the English slot stays English.
- Some option values, help texts, and confirmation messages use short English fallbacks because those menu font sizes can render Arabic-B glyphs as blank.
- Some top tab titles may still be texture/resource driven rather than normal exe strings.
- Some technical strings are intentionally left unchanged.
- This does not modify `Kanon.exe` on disk; strings are patched in memory after SteamStub finishes.
- If text appears as boxes or blank glyphs, the menu is using a font size/family that does not contain the injected Arabic Presentation Forms.

## Install

Back up any existing `version.dll`, then copy this folder's `version.dll` next to `Kanon.exe` and launch through Steam.

For logs, set this Steam launch option:

```text
LUCKPROXY_LOG=1 %command%
```
