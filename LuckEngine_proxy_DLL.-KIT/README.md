# LuckEngine Proxy DLL — In-Memory String Patch Toolkit

A `version.dll` proxy for Steam releases of **Visual Art's / Key** games running
on the Luck Engine. Lets you patch hardcoded UI strings in RAM at runtime, with
**zero modifications to the on-disk exe** — SteamStub DRM and file integrity
remain intact.

Tested on: **Kanon** (Steam), **AIR** (Steam). Should work on any Luck Engine
title that ships a `VERSION.dll`-importing exe with SteamStub.

---

## How it works

1. Windows resolves `VERSION.dll` imports from the exe's own directory before
   `System32` (it is not a Known DLL), so our proxy is loaded first.
2. On attach, a worker thread polls a sentinel byte in `.rdata` until SteamStub
   finishes decrypting the section (~200 ms typically, 30 s timeout).
3. Each patch is applied via `VirtualProtect` + `memcpy`, then protection is
   restored. The sentinel is the first entry in `patches.h` — make sure it is
   a reliable, unique string.
4. All `VERSION.dll` exports are forwarded at runtime to the real
   `C:\Windows\System32\version.dll`.

---

## Repository layout

```
patches.py      ← single source of truth: edit this to add/change strings
patches.h       ← auto-generated (do not edit)
patches.csv     ← auto-generated review table
version.c       ← DLL core (proxy + patch engine), shared across all games
version.def     ← export forwarding table
Makefile        ← build recipe
```

One subfolder per game (e.g. `Kanon/`, `AIR/`), each with its own
`patches.py` / `patches.h` / `patches.csv` / pre-built `version.dll`.
The `version.c`, `version.def`, and `Makefile` at the root are shared.

---

## Workflow

### 1. Configure `patches.py`

Open `patches.py` and set the variables at the top:

```python
GAME_EXE        = 'Kanon.exe'   # exe to read offsets from
RVA_DELTA       = 0xC00         # raw_offset + RVA_DELTA = RVA (see below)
PATCH_GAME_NAME = 'Kanon'       # appears in the log
PATCH_VERSION   = '0.1'         # appears in the log
```

### 2. Find `RVA_DELTA` for a new game

```python
python3 -c "
import pefile, sys
pe = pefile.PE(sys.argv[1])
for s in pe.sections:
    if b'.rdata' in s.Name:
        print(f'.rdata  raw=0x{s.PointerToRawData:X}  va=0x{s.VirtualAddress:X}  delta=0x{s.VirtualAddress - s.PointerToRawData:X}')
" GameName.exe
```

### 3. Add patches

Each entry in `PATCHES` is a 5-tuple:

```python
(raw_offset, src_bytes, target_str, context, note)
```

| Field | Type | Description |
|---|---|---|
| `raw_offset` | `int` | Byte offset of the source string in the exe |
| `src_bytes` | `bytes` | Original bytes to match (e.g. `b'Close'`) |
| `target_str` | `str` | Replacement string (UTF-8) |
| `context` | `str` | Label for the CSV (tab name, element) |
| `note` | `str` | Slot/budget info or translator notes |

**Budget rule:** `len(target_str.encode('utf-8')) <= slot_size - 1`

The slot size is auto-detected from the exe (source string + contiguous null
padding). If your translation is longer than the budget, the script exits with
an error.

### 4. Generate and build

```bash
# From the game subfolder (where patches.py and the exe live):
python3 patches.py        # validates offsets, writes patches.h + patches.csv
make -C ..                # compiles version.dll  (requires mingw-w64)
```

Or from the root:
```bash
cd Kanon && python3 patches.py && cd .. && make
```

### 5. Install

1. Back up any existing `version.dll` in the game folder.
2. Copy the newly built `version.dll` into the game folder (next to the `.exe`).
3. Launch via Steam.

---

## Build requirements

- `x86_64-w64-mingw32-gcc` (mingw-w64 cross-compiler)
- Python 3.8+ (for `patches.py`)
- `pefile` Python package if using the RVA delta helper (`pip install pefile`)

On Debian/Ubuntu:
```bash
sudo apt install mingw-w64
pip install pefile
```

---

## Enabling logs

Set the Steam launch option:

```
LUCKPROXY_LOG=1 %command%
```

A `luckproxy.log` file appears next to the exe with timestamped events:

```
[HH:MM:SS.mmm] DLL_PROCESS_ATTACH (Kanon proxy v0.1, 42 patches)
[HH:MM:SS.mmm] Loaded real version.dll from C:\WINDOWS\system32\version.dll
[HH:MM:SS.mmm] Sentinel ready, applying 42 patch(es)
[HH:MM:SS.mmm] Patched RVA 0x... (5 bytes): bottom-right button: Close
...
[HH:MM:SS.mmm] Patch thread done.
```

If the sentinel never matches after 30 s, the log shows the raw bytes found
at that address — useful for diagnosing a wrong `RVA_DELTA` or a game update
that shifted offsets.

---

## Strings that cannot be patched

Some slots are too small for common target-language equivalents. Document them
in `patches.py` as commented-out entries with an explanation:

```python
# (0x499204, b'All', '???', 'Basic/Skip value', 'SKIP: slot 4B / budget 3B — no short equivalent')
```

---

## Iterating

1. Edit `patches.py` (the only file you need to touch for string changes).
2. `python3 patches.py` → regenerates `patches.h` + `patches.csv`.
3. `make` → recompiles `version.dll`.
4. Copy the new `version.dll` to the game folder and relaunch.

---

## Adding a new game

1. Create a subfolder: `mkdir NewGame && cd NewGame`.
2. Copy `patches.py` from an existing game folder (or the root template).
3. Update `GAME_EXE`, `RVA_DELTA`, `PATCH_GAME_NAME`, `PATCH_VERSION`.
4. Clear `PATCHES = []` and start adding entries.
5. Run `python3 patches.py` and `make -C ..`.

---

## Credits

Toolkit developed as part of the [LuckSystem](https://github.com/Yoremi/LuckSystem)
translation toolchain for Visual Art's / Key games.
