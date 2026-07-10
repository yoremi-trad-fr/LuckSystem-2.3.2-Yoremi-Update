#!/usr/bin/env python3
"""
Luck Engine Steam — string patch table template.
Single source of truth for all source -> target string patches in the game exe.
Generates:
  - patches.h  (C array included by version.c)
  - patches.csv (human-readable review table)

=== HOW TO USE ===

1. Set the configuration variables below (GAME_EXE, RVA_DELTA, etc.).
2. Add your entries to PATCHES. Each entry is a 5-tuple:
       (raw_offset, src_bytes, target_str, context, note)
   - raw_offset  : int  — byte offset of the source string in the exe file
   - src_bytes   : bytes — original bytes to match (used as sentinel + expected)
   - target_str  : str  — translated string (UTF-8 encoded on output)
   - context     : str  — free-form label for the CSV (menu tab, element name…)
   - note        : str  — slot/budget info or any translator note
3. Run:  python3 patches.py
   It validates every offset against the real exe, prints a fit report,
   then writes patches.h and patches.csv.
4. Run:  make
   This compiles version.dll from version.c + patches.h.

=== RULES ===

  - RVA = raw_offset + RVA_DELTA  (computed automatically)
  - target length (UTF-8 bytes) MUST be <= slot_size - 1
    (slot_size is auto-detected from the exe: src bytes + trailing null block)
  - A trailing \\0 is always appended; remaining padding in the slot stays \\0
  - If target length == src length exactly, that is ideal (no wasted bytes)
  - If target length < src length, the extra bytes inside the slot stay \\0 (safe)
  - If target length > budget, the script aborts with exit code 2

=== FINDING RVA_DELTA ===

  For a given exe, open it in a hex editor or use:
      python3 -c "
      import pefile, sys
      pe = pefile.PE(sys.argv[1])
      for s in pe.sections:
          if b'.rdata' in s.Name:
              print(f'.rdata  raw=0x{s.PointerToRawData:X}  va=0x{s.VirtualAddress:X}  delta=0x{s.VirtualAddress - s.PointerToRawData:X}')
      " GameName.exe

  The delta is typically constant for all strings in .rdata.

=== LOGGING ===

  Set the environment variable LUCKPROXY_LOG=1 before launching the game
  (Steam launch options: LUCKPROXY_LOG=1 %command%) to get a luckproxy.log
  next to the exe with timestamped patch events.
"""

# ---------------------------------------------------------------------------
#  Configuration — edit these for your game / translation
# ---------------------------------------------------------------------------

GAME_EXE       = 'GameName.exe'   # exe filename to read offsets from
RVA_DELTA      = 0x000            # raw_offset + RVA_DELTA = RVA  (see above)
PATCH_GAME_NAME = 'GameName'      # written into DllMain log line
PATCH_VERSION  = '0.1'            # written into DllMain log line

# ---------------------------------------------------------------------------
#  Patch table
#  Each entry: (raw_offset, src_bytes, target_str, context, note)
# ---------------------------------------------------------------------------

PATCHES = [
    # Example (remove before publishing your patch):
    # (0x000000, b'Example', 'Translation', 'Menu/Label', 'slot 8, budget 7'),
]

# ---------------------------------------------------------------------------
#  Code below this line does not need editing
# ---------------------------------------------------------------------------

def main():
    import sys

    data = open(GAME_EXE, 'rb').read()

    def slot_size(start):
        i = start
        while data[i] != 0:
            i += 1
        while i < len(data) and data[i] == 0:
            i += 1
        return i - start

    rows = []
    errors = []
    for off, src, target, context, note in PATCHES:
        actual = data[off:off + len(src)]
        if actual != src:
            errors.append(f"0x{off:X}: expected {src!r}, got {actual!r}")
            continue
        target_bytes = target.encode('utf-8')
        slot = slot_size(off)
        budget = slot - 1
        fits = len(target_bytes) <= budget
        rows.append({
            'off':          off,
            'src':          src.decode('utf-8', errors='replace'),
            'target':       target,
            'src_len':      len(src),
            'target_len':   len(target_bytes),
            'slot':         slot,
            'budget':       budget,
            'fits':         fits,
            'context':      context,
            'note':         note,
            'src_bytes':    src,
            'target_bytes': target_bytes,
        })

    if errors:
        print("=== OFFSET MISMATCH (aborting) ===", file=sys.stderr)
        for e in errors:
            print("  " + e, file=sys.stderr)
        sys.exit(1)

    # Print fit report
    print(f"{'off':>8}  {'slot':>4}  {'src':>3}  {'tgt':>3}  {'fit':3}  src -> target")
    print('-' * 100)
    n_ok = n_bad = 0
    for r in rows:
        mark = '✓' if r['fits'] else '✗'
        if r['fits']:
            n_ok += 1
        else:
            n_bad += 1
        print(f"0x{r['off']:06X}  {r['slot']:>4}  {r['src_len']:>3}  {r['target_len']:>3}  {mark}   {r['src']!r} -> {r['target']!r}")
    print(f"\nTotal: {len(rows)}  OK: {n_ok}  FAIL: {n_bad}")
    if n_bad:
        print("\nFailures (target too long):")
        for r in rows:
            if not r['fits']:
                print(f"  0x{r['off']:X}: target {r['target_len']}B > budget {r['budget']}B: {r['target']!r}")
        sys.exit(2)

    # Emit patches.h
    with open('patches.h', 'w', encoding='utf-8') as f:
        f.write('/* Auto-generated from patches.py. Do not edit. */\n')
        f.write('#ifndef LUCKPROXY_PATCHES_H\n#define LUCKPROXY_PATCHES_H\n\n')
        f.write(f'#define PATCH_GAME_NAME "{PATCH_GAME_NAME}"\n')
        f.write(f'#define PATCH_VERSION   "{PATCH_VERSION}"\n\n')
        for i, r in enumerate(rows):
            rva = r['off'] + RVA_DELTA
            write_len = max(len(r['src_bytes']), len(r['target_bytes'])) + 1
            src_padded    = list(r['src_bytes'])    + [0] * (write_len - len(r['src_bytes']))
            target_padded = list(r['target_bytes']) + [0] * (write_len - len(r['target_bytes']))
            assert len(src_padded) == write_len and len(target_padded) == write_len
            src_arr    = ','.join(f'0x{b:02X}' for b in src_padded)
            target_arr = ','.join(f'0x{b:02X}' for b in target_padded)
            f.write(f'static const BYTE s_src_{i:03d}[] = {{ {src_arr} }};\n')
            f.write(f'static const BYTE s_tgt_{i:03d}[] = {{ {target_arr} }};\n')
        f.write('\nstatic const LuckPatch g_patches[] = {\n')
        for i, r in enumerate(rows):
            rva = r['off'] + RVA_DELTA
            write_len = max(len(r['src_bytes']), len(r['target_bytes'])) + 1
            ctx = r['context'] + ': ' + r['src'][:30]
            ctx = (ctx.replace('\\', '\\\\')
                      .replace('"', '\\"')
                      .replace('\n', '\\n')
                      .replace('\r', '\\r')
                      .replace('\t', '\\t'))
            f.write(f'    {{ 0x{rva:06X}, {write_len:>4}, s_src_{i:03d}, s_tgt_{i:03d}, "{ctx}" }},\n')
        f.write('};\n\n#define N_PATCHES (sizeof(g_patches)/sizeof(g_patches[0]))\n')
        f.write('\n#endif\n')
    print(f"\nGenerated patches.h with {len(rows)} entries.")

    # Emit patches.csv
    with open('patches.csv', 'w', encoding='utf-8') as f:
        f.write('raw_offset,rva,slot,budget,src_len,target_len,fits,src,target,context,note\n')
        for r in rows:
            rva = r['off'] + RVA_DELTA
            def esc(s): return '"' + s.replace('"', '""') + '"'
            f.write(
                f'0x{r["off"]:X},0x{rva:X},{r["slot"]},{r["budget"]},'
                f'{r["src_len"]},{r["target_len"]},{r["fits"]},'
                f'{esc(r["src"])},{esc(r["target"])},{esc(r["context"])},{esc(r["note"])}\n'
            )
    print(f"Generated patches.csv with {len(rows)} entries.")


if __name__ == '__main__':
    main()
