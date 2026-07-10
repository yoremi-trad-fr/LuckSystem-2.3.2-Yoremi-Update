#!/usr/bin/env python3
"""
AIR Steam (Luck Engine) — string patch table.
Single source of truth for all EN -> target string patches in AIR.exe.
Generates:
  - patches.h  (C array included by version.c)
  - patches.csv (human-readable review table)

Configuration:
  GAME_EXE       — exe to read offsets from
  RVA_DELTA      — raw_offset + RVA_DELTA = RVA  (0x1800 for AIR Steam)
  PATCH_GAME_NAME / PATCH_VERSION — appear in DLL log output

Derived from the Kanon patch table: all strings that appear byte-for-byte
identical in AIR.exe were ported automatically. Kanon-specific strings
(e.g. voice sample labels) are absent.

Rules:
  - Each entry: (raw_offset, src_bytes, target_str, context, note)
  - target length (UTF-8) MUST be <= slot_size - 1  (budget shown in CSV)
  - A trailing \\0 is appended automatically
"""

GAME_EXE        = 'AIR.exe'
RVA_DELTA       = 0x1800
PATCH_GAME_NAME = 'AIR'
PATCH_VERSION   = '0.1'

# (raw_offset, src_bytes, target_str, context, note)

PATCHES = [

    (0x48C534, b'Close',
                                               '',
                                               'bottom-right button', 'slot 8, budget 7'),
    (0x48DC28, b'Defaults',
                                               '',
                                               'bottom-left button', 'budget 15 ok'),
    (0x49E490, b'Basic',
                                               '',
                                               'Basic tab', 'already done via PAK? keep anyway'),
    (0x49E528, b'Shortcut Menu',
                                               '',
                                               'Basic/Shortcut', 'budget 15'),
    (0x49E4CC, b'Hide',
                                               '',
                                               'Basic/Shortcut value', 'budget 11'),
    (0x49E4D8, b'Display',
                                               '',
                                               'Basic/Shortcut value', 'budget 7exact'),
    (0x49E590, b'Skip',
                                               '',
                                               'Basic/Skip', 'budget 7'),
    (0x49E538, b'Previously Read Only',
                                               '',
                                               'Basic/Skip value', 'budget 23'),
    (0x49E5D0, b'Position of Choices',
                                               '',
                                               'Basic/Position', '19B, budget 23 ok'),
    (0x49E5B8, b'Bottom',
                                               '',
                                               'Basic/Position value', 'budget 7'),
    (0x49E5C0, b'Center',
                                               '',
                                               'Basic/Position value', 'budget 7'),
    (0x49E6A8, b'Voice',
                                               '',
                                               'Basic/Voice', 'budget 7'),
    (0x49E658, b'Stop on New Message',
                                               '',
                                               'Basic/Voice value', 'budget 23'),
    (0x49E670, b'      No Stops      ',
                                               '',
                                               'Basic/Voice value', 'keep spaces for alignment, budget 23'),
    (0x49E720, b'Display Date',
                                               '',
                                               'Basic/Date', '13B, budget 15'),
    (0x49E7D8, b'Initial Cursor Position',
                                               '',
                                               'Basic/Cursor', '20B, budget 23'),
    (0x49E730, b'Positioned at \xe2\x9d\x9dYes\xe2\x9d\x9e',
                                               '',
                                               'Basic/Cursor value', 'Keep stylish chevrons if fits; budget 23'),
    (0x49E7A8, b'Positioned at \xe2\x9d\x9dNo\xe2\x9d\x9e',
                                               '',
                                               'Basic/Cursor value', 'budget 23'),
    (0x49E818, b'Controller Rumble Function',
                                               '',
                                               'Basic/Rumble', '17B, budget 27'),
    (0x49E7F8, b'Disable',
                                               '',
                                               'Basic/Rumble value', 'budget 7 (rumble intensity: None/Min/Mid/Max)'),
    (0x49F840, b'Text1',
                                               '',
                                               'Text1 tab', '6B, budget 7'),
    (0x49F848, b'Language',
                                               '',
                                               'Text1/Language', '6B, budget 11'),
    (0x49F908, b'Font',
                                               '',
                                               'Text1/Font', '6B, budget 7'),
    (0x49F8DC, b'Mincho',
                                               '',
                                               'Text1/Font value', 'japanese font family name, keep'),
    (0x49F8F8, b'Modern',
                                               '',
                                               'Text1/Font value', '7B, budget 7 tight!'),
    (0x49F944, b'Solid',
                                               '',
                                               'Text1/Window Transp value', '6B, budget 11'),
    (0x49F96C, b'Clear',
                                               '',
                                               'Text1/Window Transp value', '5B, budget 11'),
    (0x49F988, b'Window Transparency',
                                               '',
                                               'Text1/Window Transp', '12B, budget 23'),
    (0x49F9A0, b'Only Choices',
                                               '',
                                               'Text1/Transp target', '15B, budget 15 exact'),
    (0x49F9C8, b'Previously Read Text',
                                               '',
                                               'Text1/Read target', 'budget 23'),
    (0x49F9F8, b'Color of',
                                               '',
                                               'Text1/Color label', '7B, budget 15'),
    (0x49FA40, b'Green',
                                               '',
                                               'Text1/Color value', '4B, budget 7'),
    (0x49FA48, b'Blue',
                                               '',
                                               'Text1/Color value', '4B, budget 7'),
    (0x49FA60, b'Purple',
                                               '',
                                               'Text1/Color value', '6B, budget 7'),
    (0x49F804, b'Orange',
                                               '',
                                               'Text1/Color value', 'same word'),
    (0x49F830, b'Yellow',
                                               '',
                                               'Text1/Color value', '5B, budget 7'),
    (0x49FA70, b'Read Text Color',
                                               '',
                                               'Text1/Read color', '14B, budget 15'),
    (0x49FE2C, b'Text2',
                                               '',
                                               'Text2 tab', '6B, budget 11'),
    (0x49FEA0, b'Text Speed',
                                               '',
                                               'Text2/Speed', '13B, budget 15'),
    (0x49FE5C, b'Slow',
                                               '',
                                               'Text2/Speed value', '4B, budget 7'),
    (0x49FE64, b'Fast',
                                               '',
                                               'Text2/Speed value', '6B, budget 7'),
    (0x49FEB0, b'0 sec/char',
                                               '',
                                               'Text2/Speed value', '7B, budget 15'),
    (0x49FF00, b'0.1 sec/char',
                                               '',
                                               'Text2/Speed value', '9B, budget 15'),
    (0x49FF20, b'Wait Time Per Character',
                                               '',
                                               'Text2/Wait', 'budget 23'),
    (0x49FF38, b'0 sec',
                                               '',
                                               'Text2/Wait value', '3B, budget 7'),
    (0x49FF48, b'1 sec',
                                               '',
                                               'Text2/Wait value', '3B, budget 7'),
    (0x49FF58, b'2 sec',
                                               '',
                                               'Text2/Wait value', '3B, budget 7'),
    (0x49FF68, b'3 sec',
                                               '',
                                               'Text2/Wait value', '3B, budget 7'),
    (0x49FF78, b'Base Wait Time',
                                               '',
                                               'Text2/Base', '15B, budget 15 exact'),
    (0x49FFCC, b'Sound',
                                               '',
                                               'Sound tab', '3B, budget 11'),
    (0x49FFD8, b'Master Volume',
                                               '',
                                               'Sound/Master', 'budget 15'),
    (0x4A0078, b'System Sounds',
                                               '',
                                               'Sound/System', 'budget 15'),
    (0x4A0110, b'While Pressed',
                                               '',
                                               'Keyboard/mode', '8B, budget 15'),
    (0x4A0120, b'Start/Stop',
                                               '',
                                               'Keyboard/mode', 'budget 15'),
    (0x4A0178, b'C (Skip)',
                                               '',
                                               'Keyboard/key label', '10B, budget 15'),
    (0x4A0188, b'Z (Rewind)',
                                               '',
                                               'Keyboard/key label', '10B, budget 15'),
    (0x4A01E0, b'  Disable  ',
                                               '',
                                               'Keyboard/value', 'retain spaces, budget 15'),
    (0x4A01F0, b'Quick Save',
                                               '',
                                               'Keyboard/key label', '12B, budget 15'),
    (0x4A0240, b'Quick Load',
                                               '',
                                               'Keyboard/key label', '12B, budget 15'),
    (0x4A02A0, b'Switch Language',
                                               '',
                                               'Keyboard/key label', '14B, budget 15'),
    (0x4A0300, b'Up Arrow',
                                               '',
                                               'Keyboard/key label', 'budget 15'),
    (0x4A034C, b'Mouse',
                                               '',
                                               'Mouse tab', '6B, budget 11'),
    (0x4A03A0, b'System Menu',
                                               '',
                                               'Mouse/target', 'budget 15'),
    (0x4A03B0, b'Hide Window',
                                               '',
                                               'Mouse/target', 'budget 15 exact'),
    (0x4A0408, b'Right Click',
                                               '',
                                               'Mouse/binding', '10B, budget 15'),
    (0x4A0428, b'Left+Right Click',
                                               '',
                                               'Mouse/binding', '18B, budget 23'),
    (0x4A0440, b'Mouse Wheel Button',
                                               '',
                                               'Mouse/binding', '14B, budget 23'),
    (0x4A04A0, b'Rewind Once',
                                               '',
                                               'Mouse/target', '13B, budget 15'),
    (0x4A04B0, b'Wheel Up',
                                               '',
                                               'Mouse/binding', '12B, budget 15'),
    (0x4A0510, b'Forward Once',
                                               '',
                                               'Mouse/target', '14B, budget 15'),
    (0x4A0520, b'Wheel Down',
                                               '',
                                               'Mouse/binding', '11B, budget 15'),
    (0x4A05A0, b'Jump and Switch Pages',
                                               '',
                                               'Mouse/target', '14B, budget 23'),
    (0x4A05B8, b'Return/Proceed Button',
                                               '',
                                               'Mouse/target', '16B, budget 23'),
    (0x4A0614, b'Enable',
                                               '',
                                               'Mouse/Gestures value', 'budget 11'),
    (0x4A0630, b'Gestures',
                                               '',
                                               'Mouse/Gestures', '6B, budget 15'),
    (0x4A0670, b'Dialog and Choices',
                                               '',
                                               'Mouse/Snap target', '18B, budget 23'),
    (0x4A06F8, b'Snap Pointer',
                                               '',
                                               'Mouse/Snap', '16B, budget 23'),
    (0x4A0A68, b'Game Ver. ',
                                               '',
                                               'System/info', '12B, budget 15'),
    (0x4A0A8C, b'System',
                                               '',
                                               'System tab', 'budget 11'),
    (0x4A0AD0, b'Window',
                                               '',
                                               'System/Window', '7B, budget 7 exact'),
    (0x4A0AD8, b'Full Screen',
                                               '',
                                               'System/FullScreen', 'budget 15'),
    (0x4A0B00, b'Screen Mode',
                                               '',
                                               'System/ScreenMode', 'budget 15'),
    (0x4A0B40, b'Auto',
                                               '',
                                               'System/value', 'same'),
    (0x4A0B90, b'Window Size',
                                               '',
                                               'System/WindowSize', 'budget 15 exact'),
    (0x49FC40, b'Wait Time Per Character: In Auto Mode, this sets the wait time until the next message is displayed based on the number of characters in text.',
                                               '',
                                               'Texte2 tooltip', 'budget 143'),
    (0x49FCD0, b'Base Wait Time: You can set a Base Wait Time to add to\n\xe2\x9d\x9dWait Time Per Character\xe2\x9d\x9e.',
                                               '',
                                               'Texte2 tooltip', 'budget 95, keep stylized quotes'),
    (0x49EA80, b'Initial Cursor Position: Sets the initial position of the cursor when a Yes/No choice is available.',
                                               '',
                                               'Basic tooltip', 'budget 111'),
    (0x49ECF0, b'Controller Rumble Function: Plug in a controller to use the controller Rumble function.',
                                               '',
                                               'Basic tooltip', 'budget 87'),
    (0x4A0710, b'Gestures: Moving the cursor while holding the left button works the same way as Touch controls.',
                                               '',
                                               'Mouse tooltip', 'budget 95'),
    (0x4A09C0, b'Left+Right Click: Hold left button then right click to switch between languages (English/Simplified Chinese/Japanese).',
                                               '',
                                               'Mouse tooltip', 'budget 119'),
    (0x49E9D0, b'Voice: If you select No Stops, sound will continue to play even if you advance the text during voice playback. (It stops if there is sound on the next message.)',
                                               '',
                                               'Basic tooltip', 'budget 175'),
    (0x48CF40, b'Do you wish to load this save?',
                                               '',
                                               'Save/Load prompt', 'budget 31'),
    (0x48D158, b'Are you sure you wish to overwrite data?',
                                               '',
                                               'Save/Load prompt', 'budget 47'),
    (0x48D110, b'Do you wish to save?',
                                               '',
                                               'Save prompt', 'budget 23, bonus'),
    (0x48D380, b'Are you sure you wish to delete this save data?',
                                               '',
                                               'Save delete prompt', 'budget 47'),
    (0x48D3E0, b'This cannot be deleted.',
                                               '',
                                               'Save delete error', 'budget 23'),
    (0x48DBBC, b'Delete',
                                               '',
                                               'Save menu button', 'budget 11'),
    (0x48DBE4, b'Latest',
                                               '',
                                               'Save menu button', 'budget 11'),
    (0x48DF78, b'Text preview.',
                                               '',
                                               'Text preview label', 'budget 15'),
    (0x48E098, b'Settings such as text speed are reflected.',
                                               '',
                                               'Text preview tooltip', 'budget 47'),
    (0x48C4E8, b'English',
                                               '',
                                               'Language switcher', 'budget 7 hard cap (neighbor is 简体中文). Fill in your ISO language code.'),
    (0x49DA78, b'Off',
                                               '',
                                               'Read text color toggle', 'budget 7'),
    (0x49E6E0, b'On',
                                               '',
                                               'Read text color toggle', 'budget 7'),
    (0x48D488, b'Return to the title screen?',
                                               '',
                                               'Title return prompt', 'budget 31'),
    (0x48D430, b'Return to the menu?',
                                               '',
                                               'Menu return prompt', 'budget 23, bonus'),
    (0x49AE80, b'$A1There is unsaved data.\n$A1Are you sure you wish to quit the game?',
                                               '',
                                               'Quit prompt (unsaved)', 'budget 79, keep $A1 tags'),
    (0x48C508, b'Yes',
                                               '',
                                               'Global dialog button', 'budget 3, exact fit. Used by ALL Yes/No prompts'),
    (0x48C51C, b'No',
                                               '',
                                               'Global dialog button', 'budget 3, exact fit. Used by ALL Yes/No prompts'),
    (0x48D348, b'Save completed.',
                                               '',
                                               'Save confirmation', 'budget 15'),
    (0x49AEC8, b'Are you sure you wish to quit the game?',
                                               '',
                                               'Quit prompt (saved)', 'budget 39 — the short variant, used when no unsaved data'),
    (0x4A0BA0, b'Changing the setting to \xe2\x9d\x9dAuto\xe2\x9d\x9e will open the window at a scale based on the Windows \xe2\x9d\x9dDisplay\xe2\x9d\x9e setting.\nChanging \xe2\x9d\x9d%%\xe2\x9d\x9e will scale the display, with the default resolution being %d\xc3\x97%d pixels.\n (Scale cannot be increased beyond the maximum resolution of your display.)',
                                               '',
                                               'System tab tooltip', 'budget 287, keep %%/%d×%d intact'),

]

def main():
    import sys
    data = open(GAME_EXE, 'rb').read()

    def slot_size(start):
        i = start
        while data[i] != 0: i += 1
        while i < len(data) and data[i] == 0: i += 1
        return i - start

    rows = []
    errors = []
    for off, en, target, context, note in PATCHES:
        actual = data[off:off+len(en)]
        if actual != en:
            errors.append(f"0x{off:X}: expected {en!r}, got {actual!r}")
            continue
        target_bytes = target.encode('utf-8')
        slot = slot_size(off)
        budget = slot - 1
        fits = len(target_bytes) <= budget
        rows.append({
            'off': off,
            'src': en.decode('utf-8', errors='replace'),
            'target': target,
            'src_len': len(en),
            'target_len': len(target_bytes),
            'slot': slot,
            'budget': budget,
            'fits': fits,
            'context': context,
            'note': note,
            'src_bytes': en,
            'target_bytes': target_bytes,
        })

    if errors:
        print("=== OFFSET MISMATCH (aborting) ===", file=sys.stderr)
        for e in errors:
            print("  " + e, file=sys.stderr)
        sys.exit(1)

    # Print status
    print(f"{'off':>8}  {'slot':>4}  {'src':>3}  {'tgt':>3}  {'fit':3}  src -> target")
    print('-' * 100)
    n_ok = n_bad = 0
    for r in rows:
        mark = '✓' if r['fits'] else '✗'
        if r['fits']: n_ok += 1
        else:          n_bad += 1
        print(f"0x{r['off']:06X}  {r['slot']:>4}  {r['src_len']:>3}  {r['target_len']:>3}  {mark}   {r['src']!r} -> {r['target']!r}")
    print(f"\nTotal: {len(rows)}  OK: {n_ok}  FAIL: {n_bad}")
    if n_bad:
        print("\nFailures (too long):")
        for r in rows:
            if not r['fits']:
                print(f"  0x{r['off']:X}: target {r['target_len']}B > budget {r['budget']}B: {r['target']!r}")
        sys.exit(2)

    # Emit patches.h (C header)
    with open('patches.h','w', encoding='utf-8') as f:
        f.write('/* Auto-generated from patches.py. Do not edit. */\n')
        f.write('#ifndef LUCKPROXY_PATCHES_H\n#define LUCKPROXY_PATCHES_H\n\n')
        f.write(f'#define PATCH_GAME_NAME "{PATCH_GAME_NAME}"\n')
        f.write(f'#define PATCH_VERSION   "{PATCH_VERSION}"\n\n')
        for i, r in enumerate(rows):
            rva = r['off'] + RVA_DELTA
            # Write length = max(len(src)+1, len(target)+1): overwrite at least full src
            # to erase stale bytes, and at least full target so no truncation.
            write_len = max(len(r['src_bytes']), len(r['target_bytes'])) + 1
            en_padded = list(r['src_bytes']) + [0] * (write_len - len(r['src_bytes']))
            fr_padded = list(r['target_bytes']) + [0] * (write_len - len(r['target_bytes']))
            assert len(en_padded) == write_len and len(fr_padded) == write_len
            en_arr = ','.join(f'0x{b:02X}' for b in en_padded)
            fr_arr = ','.join(f'0x{b:02X}' for b in fr_padded)
            f.write(f'static const BYTE s_src_{i:03d}[] = {{ {en_arr} }};\n')
            f.write(f'static const BYTE s_tgt_{i:03d}[] = {{ {fr_arr} }};\n')
        f.write('\nstatic const LuckPatch g_patches[] = {\n')
        for i, r in enumerate(rows):
            rva = r['off'] + RVA_DELTA
            write_len = max(len(r['src_bytes']), len(r['target_bytes'])) + 1
            # Escape newlines, backslashes, and double quotes for C string literal
            ctx = r['context'] + ': ' + r['src'][:30]
            ctx = ctx.replace('\\', '\\\\').replace('"', '\\"').replace('\n', '\\n').replace('\r', '\\r').replace('\t', '\\t')
            f.write(f'    {{ 0x{rva:06X}, {write_len:>4}, s_src_{i:03d}, s_tgt_{i:03d}, "{ctx}" }},\n')
        f.write('};\n\n#define N_PATCHES (sizeof(g_patches)/sizeof(g_patches[0]))\n')
        f.write('\n#endif\n')
    print(f"\nGenerated patches.h with {len(rows)} entries.")

    # Emit patches.csv
    with open('patches.csv','w', encoding='utf-8') as f:
        f.write('raw_offset,rva,slot,budget,src_len,target_len,fits,src,target,context,note\n')
        for r in rows:
            rva = r['off'] + RVA_DELTA
            # CSV-escape quotes in text
            def esc(s): return '"' + s.replace('"','""') + '"'
            f.write(f'0x{r["off"]:X},0x{rva:X},{r["slot"]},{r["budget"]},{r["src_len"]},{r["target_len"]},{r["fits"]},{esc(r["src"])},{esc(r["target"])},{esc(r["context"])},{esc(r["note"])}\n')
    print(f"Generated patches.csv with {len(rows)} entries.")


if __name__ == '__main__':
    main()
