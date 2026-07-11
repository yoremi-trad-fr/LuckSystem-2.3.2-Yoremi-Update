# Little Busters! English Edition — mixed string proxy

LBEE differs from the other supported Luck Engine games in three important
ways:

- `LITBUS_WIN32.exe` is a 32-bit x86 executable.
- It imports `winmm.dll`, not `version.dll`.
- Its hardcoded strings mix UTF-8 and UTF-16LE, and one supplied patch is in
  `.rsrc` rather than `.rdata`.

The shared proxy therefore builds as a 32-bit `winmm.dll` for LBEE. It forwards
the timer functions used by the game and dynamically loaded graphics drivers
(`timeBeginPeriod`, `timeEndPeriod`, `timeGetDevCaps`, and `timeGetTime`) to the
real Windows system DLL. Do not install `version.dll` alongside it.

## GUI workflow (v3.29+)

Select **Little Busters! English Edition** in `DLL HOOK -> Luca Menu DLL`, then
select the exact `LITBUS_WIN32.exe` used to build the inventory. The GUI handles
UTF-8 and UTF-16LE budgets separately, computes RVAs from the PE section table,
selects the x86 compiler toolchain, and outputs `winmm.dll` automatically.
Choose **Russe (LBEE)** under **Langue à injecter** to load the validated menu
translations and build the complete bundled community preset automatically.
The GUI runs `mixed_patches.py` against `russian_preset.py`, producing 1,034
effective patches; no Russian strings need to be entered manually.

## Generate the complete Russian preset

The repository includes `russian_preset.py` (the former external `patches2.py`
table under a descriptive name). To run it manually, put it and the matching
`LITBUS_WIN32.exe` in one directory, open a terminal there, then run:

```powershell
python "<LuckSystem repo>\proxy dll\LBEE\mixed_patches.py"
```

The generator detects UTF-8 versus UTF-16LE, computes each RVA from the PE
section table, preserves two-byte UTF-16 terminators, and writes `patches.h`
plus a reviewable `patches.csv`.

It also repairs the original first-match search: duplicate source text is
expanded to every standalone UTF-8/UTF-16LE occurrence in the executable's
string-data sections. This avoids patching substrings such as `Close` inside
`CloseThreadpool` while missing the visible menu label. It corrects 36
non-standalone first-match offsets and removes 10 no-op source=target rows;
The current community table produces 1,034 effective patches.

`CROP_OVERSIZE` is enabled for the supplied proof-of-concept translations.
Rows shortened automatically are marked `cropped=True` in the CSV and should
be rewritten manually before release.

## Build

With 32-bit MinGW-w64:

```bash
make -C "<LuckSystem repo>/proxy dll" PATCH_DIR="<patch directory>" PROXY=winmm ARCH=x86
```

With Visual Studio Build Tools, use an x86 developer environment:

```bat
call "C:\Program Files (x86)\Microsoft Visual Studio\2022\BuildTools\Common7\Tools\VsDevCmd.bat" -arch=x86 -host_arch=x64
cl /nologo /O2 /W3 /LD /D_CRT_SECURE_NO_WARNINGS /DLUCKPROXY_WINMM /I "<patch directory>" /Fe:"<patch directory>\winmm.dll" "<LuckSystem repo>\proxy dll\version.c" /link /DEF:"<LuckSystem repo>\proxy dll\winmm.def" /SUBSYSTEM:WINDOWS /NOLOGO
```

Copy the generated `winmm.dll` next to `LITBUS_WIN32.exe`. For diagnostics,
launch the game with `LUCKPROXY_LOG=1`; `luckproxy.log` will appear beside the
executable.

When the selected executable is already unpacked, all patches are applied
synchronously during DLL attach. This ensures startup menu strings are changed
before LBEE copies them into its runtime UI structures. Packed builds retain
the SteamStub polling worker.
