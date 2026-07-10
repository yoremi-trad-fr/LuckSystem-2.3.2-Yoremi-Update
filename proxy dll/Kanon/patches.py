#!/usr/bin/env python3
"""
Kanon Steam (Luck Engine) — string patch table.
Single source of truth for all EN -> target string patches in Kanon.exe.
Generates:
  - patches.h  (C array included by version.c)
  - patches.csv (human-readable review table)

Configuration:
  GAME_EXE       — exe to read offsets from
  RVA_DELTA      — raw_offset + RVA_DELTA = RVA  (0xC00 for Kanon Steam v1.5.0.6)
  PATCH_GAME_NAME / PATCH_VERSION — appear in DLL log output

Rules:
  - Each entry: (raw_offset, src_bytes, target_str, context, note)
  - target length (UTF-8) MUST be <= slot_size - 1  (budget shown in CSV)
  - A trailing \\0 is appended automatically
  - If target len == src len, that is ideal (no wasted bytes)
  - If target len < src len, remaining slot bytes stay \\0 (safe)
"""

GAME_EXE        = 'Kanon.exe'
RVA_DELTA       = 0xC00
PATCH_GAME_NAME = 'Kanon'
PATCH_VERSION   = '0.1'

# (raw_offset, src_bytes, target_str, context, note)

PATCHES = [

    # -- OPTIONS screen header / bottom buttons --
    (0x4874A4, b'Close', '', 'bottom-right button', 'slot 8, budget 7'),
    (0x488B98, b'Defaults', '', 'bottom-left button', 'budget 15 ok'),

    # -- Basic tab --
    (0x499104, b'Basic', '', 'Basic tab', 'already done via PAK? keep anyway'),
    (0x4991A0, b'Shortcut Menu', '', 'Basic/Shortcut', 'budget 15'),
    (0x499144, b'Hide', '', 'Basic/Shortcut value', 'budget 11'),
    (0x499150, b'Display', '', 'Basic/Shortcut value', 'budget 7exact'),
    (0x499208, b'Skip', '', 'Basic/Skip', 'budget 7'),
    (0x4991B0, b'Previously Read Only', '', 'Basic/Skip value', 'budget 23'),
    # "All" has slot 4 / budget 3 -> NO French word fits. Leave as-is with comment.
    # (0x499204, b'All',                         'All',                   'Basic/Skip value', 'SKIP: no short equivalent'),
    (0x499248, b'Position of Choices', '', 'Basic/Position', '19B, budget 23 ok'),
    (0x499230, b'Bottom', '', 'Basic/Position value', 'budget 7'),
    (0x499238, b'Center', '', 'Basic/Position value', 'budget 7'),
    (0x499320, b'Voice', '', 'Basic/Voice', 'budget 7'),
    (0x4992D0, b'Stop on New Message', '', 'Basic/Voice value', 'budget 23'),
    (0x4992E8, b'      No Stops      ', '', 'Basic/Voice value', 'keep spaces for alignment, budget 23'),
    (0x499398, b'Display Date', '', 'Basic/Date', '13B, budget 15'),
    # Display Date: Off / On are shared with other "On/Off" toggles - risky to repatch globally
    # (0x4987E8, b'Off', 'Off', ...) - skip
    # (0x499358, b'On',  'On',  ...) - skip (too short to replace anyway: budget 7 for "On" slot. too long)

    (0x499450, b'Initial Cursor Position', '', 'Basic/Cursor', '20B, budget 23'),
    (0x4993A8, b'Positioned at \xe2\x9d\x9dYes\xe2\x9d\x9e', '', 'Basic/Cursor value', 'Keep stylish chevrons if fits; budget 23'),
    (0x499420, b'Positioned at \xe2\x9d\x9dNo\xe2\x9d\x9e', '', 'Basic/Cursor value', 'budget 23'),
    (0x499490, b'Controller Rumble Function', '', 'Basic/Rumble', '17B, budget 27'),
    (0x499470, b'Disable', '', 'Basic/Rumble value', 'budget 7 (rumble intensity: None/Min/Mid/Max)'),
    # Min / Mid / Max (slot 4, budget 3) — keep as-is, universal
    # (0x499478, Min) - SKIP
    # (0x499480, Mid) - SKIP
    # (0x499488, Max) - SKIP

    # -- Text1 tab (fonts, colors, window) --
    (0x49A4D0, b'Text1', '', 'Text1 tab', '6B, budget 7'),
    (0x49A4D8, b'Language', '', 'Text1/Language', '6B, budget 11'),
    (0x49A598, b'Font', '', 'Text1/Font', '6B, budget 7'),
    (0x49A56C, b'Mincho', '', 'Text1/Font value', 'japanese font family name, keep'),
    (0x49A588, b'Modern', '', 'Text1/Font value', '7B, budget 7 tight!'),
    (0x49A5D4, b'Solid', '', 'Text1/Window Transp value', '6B, budget 11'),
    (0x49A5FC, b'Clear', '', 'Text1/Window Transp value', '5B, budget 11'),
    (0x49A618, b'Window Transparency', '', 'Text1/Window Transp', '12B, budget 23'),
    (0x49A630, b'Only Choices', '', 'Text1/Transp target', '15B, budget 15 exact'),
    (0x49A658, b'Previously Read Text', '', 'Text1/Read target', 'budget 23'),
    (0x49A688, b'Color of', '', 'Text1/Color label', '7B, budget 15'),
    (0x49A6D0, b'Green', '', 'Text1/Color value', '4B, budget 7'),
    (0x49A6D8, b'Blue', '', 'Text1/Color value', '4B, budget 7'),
    (0x49A6F0, b'Purple', '', 'Text1/Color value', '6B, budget 7'),
    # (0x49A6FC, 'Red' — too long: slot 4 / budget 3 TOO LONG -> keep "Red" or use "R.")
    (0x49A494, b'Orange', '', 'Text1/Color value', 'same word'),
    (0x49A4C0, b'Yellow', '', 'Text1/Color value', '5B, budget 7'),
    (0x49A700, b'Read Text Color', '', 'Text1/Read color', '14B, budget 15'),

    # -- Text2 tab --
    (0x49AABC, b'Text2', '', 'Text2 tab', '6B, budget 11'),
    (0x49AB30, b'Text Speed', '', 'Text2/Speed', '13B, budget 15'),
    (0x49AAEC, b'Slow', '', 'Text2/Speed value', '4B, budget 7'),
    (0x49AAF4, b'Fast', '', 'Text2/Speed value', '6B, budget 7'),
    (0x49AB40, b'0 sec/char', '', 'Text2/Speed value', '7B, budget 15'),
    (0x49AB90, b'0.1 sec/char', '', 'Text2/Speed value', '9B, budget 15'),
    (0x49ABB0, b'Wait Time Per Character', '', 'Text2/Wait', 'budget 23'),
    (0x49ABC8, b'0 sec', '', 'Text2/Wait value', '3B, budget 7'),
    (0x49ABD8, b'1 sec', '', 'Text2/Wait value', '3B, budget 7'),
    (0x49ABE8, b'2 sec', '', 'Text2/Wait value', '3B, budget 7'),
    (0x49ABF8, b'3 sec', '', 'Text2/Wait value', '3B, budget 7'),
    (0x49AC08, b'Base Wait Time', '', 'Text2/Base', '15B, budget 15 exact'),

    # -- Sound tab --
    (0x49AC5C, b'Sound', '', 'Sound tab', '3B, budget 11'),
    (0x49AC68, b'Master Volume', '', 'Sound/Master', 'budget 15'),
    # SFX: slot 8 / budget 7. "Effets" 6B, but SFX universal. Keep "SFX"
    (0x49AD08, b'System Sounds', '', 'Sound/System', 'budget 15'),

    # -- Voice tab --
    # (0x49AD30 = \\t(Young)\\t → voice sample label, budget 15, "\\t(Jeune)\\t" = 11B if target lang. But it's user-facing)
    (0x49AD30, b'\t(Young)\t', '', 'Voice tab young label', 'Yuichi child voice. 9B, budget 15'),
    # Yuichi Aizawa — character name, DO NOT translate
    # PARTS/voice_icon — asset path, DO NOT translate

    # -- Keyboard tab --
    # keyboard (0x49ADA0) likely asset label, keep
    (0x49ADF0, b'While Pressed', '', 'Keyboard/mode', '8B, budget 15'),
    (0x49AE00, b'Start/Stop', '', 'Keyboard/mode', 'budget 15'),
    (0x49AE58, b'C (Skip)', '', 'Keyboard/key label', '10B, budget 15'),
    (0x49AE68, b'Z (Rewind)', '', 'Keyboard/key label', '10B, budget 15'),
    (0x49AEC0, b'  Disable  ', '', 'Keyboard/value', 'retain spaces, budget 15'),
    (0x49AED0, b'Quick Save', '', 'Keyboard/key label', '12B, budget 15'),
    (0x49AF20, b'Quick Load', '', 'Keyboard/key label', '12B, budget 15'),
    (0x49AF80, b'Switch Language', '', 'Keyboard/key label', '14B, budget 15'),
    (0x49AFE0, b'Up Arrow', '', 'Keyboard/key label', 'budget 15'),

    # -- Mouse tab --
    (0x49B02C, b'Mouse', '', 'Mouse tab', '6B, budget 11'),
    (0x49B080, b'System Menu', '', 'Mouse/target', 'budget 15'),
    (0x49B090, b'Hide Window', '', 'Mouse/target', 'budget 15 exact'),
    (0x49B0E8, b'Right Click', '', 'Mouse/binding', '10B, budget 15'),
    (0x49B108, b'Left+Right Click', '', 'Mouse/binding', '18B, budget 23'),
    (0x49B120, b'Mouse Wheel Button', '', 'Mouse/binding', '14B, budget 23'),
    (0x49B180, b'Rewind Once', '', 'Mouse/target', '13B, budget 15'),
    (0x49B190, b'Wheel Up', '', 'Mouse/binding', '12B, budget 15'),
    (0x49B1F0, b'Forward Once', '', 'Mouse/target', '14B, budget 15'),
    (0x49B200, b'Wheel Down', '', 'Mouse/binding', '11B, budget 15'),
    (0x49B280, b'Jump and Switch Pages', '', 'Mouse/target', '14B, budget 23'),
    (0x49B298, b'Return/Proceed Button', '', 'Mouse/target', '16B, budget 23'),
    (0x49B2F4, b'Enable', '', 'Mouse/Gestures value', 'budget 11'),
    (0x49B310, b'Gestures', '', 'Mouse/Gestures', '6B, budget 15'),
    (0x49B350, b'Dialog and Choices', '', 'Mouse/Snap target', '18B, budget 23'),
    (0x49B3D8, b'Snap Pointer', '', 'Mouse/Snap', '16B, budget 23'),

    # -- System tab --
    (0x49B748, b'Game Ver. ', '', 'System/info', '12B, budget 15'),
    # "Steam Deck: Proton " keep as-is — technical info, same
    (0x49B76C, b'System', '', 'System tab', 'budget 11'),
    (0x49B7B0, b'Window', '', 'System/Window', '7B, budget 7 exact'),
    (0x49B7B8, b'Full Screen', '', 'System/FullScreen', 'budget 15'),
    (0x49B7E0, b'Screen Mode', '', 'System/ScreenMode', 'budget 15'),
    (0x49B820, b'Auto', '', 'System/value', 'same'),
    (0x49B870, b'Window Size', '', 'System/WindowSize', 'budget 15 exact'),

    # -- Long tooltips (lots of budget) --
    (0x49A8D0, b'Wait Time Per Character: In Auto Mode, this sets the wait time until the next message is displayed based on the number of characters in text.',
                                               '',
                                               'Texte2 tooltip', 'budget 143'),
    (0x49A960, b'Base Wait Time: You can set a Base Wait Time to add to\n\xe2\x9d\x9dWait Time Per Character\xe2\x9d\x9e.',
                                               '',
                                               'Texte2 tooltip', 'budget 95, keep stylized quotes'),
    (0x4996F0, b'Initial Cursor Position: Sets the initial position of the cursor when a Yes/No choice is available.',
                                               '',
                                               'Basic tooltip', 'budget 111'),
    (0x499960, b'Controller Rumble Function: Plug in a controller to use the controller Rumble function.',
                                               '',
                                               'Basic tooltip', 'budget 87'),
    (0x49B3F0, b'Gestures: Moving the cursor while holding the left button works the same way as Touch controls.',
                                               '',
                                               'Mouse tooltip', 'budget 95'),
    (0x49B6A0, b'Left+Right Click: Hold left button then right click to switch between languages (English/Simplified Chinese/Japanese).',
                                               '',
                                               'Mouse tooltip', 'budget 119'),
    (0x499640, b'Voice: If you select No Stops, sound will continue to play even if you advance the text during voice playback. (It stops if there is sound on the next message.)',
                                               '',
                                               'Basic tooltip', 'budget 175'),

    # -- v0.3 ADDITIONS: Save/Load prompts --
    (0x487EB0, b'Do you wish to load this save?',
                                               '',
                                               'Save/Load prompt', 'budget 31'),
    (0x4880C8, b'Are you sure you wish to overwrite data?',
                                               '',
                                               'Save/Load prompt', 'budget 47'),
    (0x488080, b'Do you wish to save?',
                                               '',
                                               'Save prompt', 'budget 23, bonus'),
    (0x4882F0, b'Are you sure you wish to delete this save data?',
                                               '',
                                               'Save delete prompt', 'budget 47'),
    (0x488350, b'This cannot be deleted.',
                                               '',
                                               'Save delete error', 'budget 23'),
    (0x488B2C, b'Delete',
                                               '',
                                               'Save menu button', 'budget 11'),
    (0x488B54, b'Latest',
                                               '',
                                               'Save menu button', 'budget 11'),

    # -- v0.3 ADDITIONS: Text1 tab extras --
    (0x488EE8, b'Text preview.',
                                               '',
                                               'Text preview label', 'budget 15'),
    (0x489008, b'Settings such as text speed are reflected.',
                                               '',
                                               'Text preview tooltip', 'budget 47'),
    (0x487458, b'English',
                                               '',
                                               'Language switcher', 'budget 7 hard cap (neighbor is 简体中文). Fill in your ISO language code.'),

    # -- v0.3 ADDITIONS: Read Text Color On/Off toggles --
    (0x4987E8, b'Off',
                                               '',
                                               'Read text color toggle', 'budget 7'),
    (0x499358, b'On',
                                               '',
                                               'Read text color toggle', 'budget 7'),

    # -- v0.3 ADDITIONS: Right-click / system menu prompts --
    (0x4883F8, b'Return to the title screen?',
                                               '',
                                               'Title return prompt', 'budget 31'),
    (0x4883A0, b'Return to the menu?',
                                               '',
                                               'Menu return prompt', 'budget 23, bonus'),
    (0x495B30, b'$A1There is unsaved data.\n$A1Are you sure you wish to quit the game?',
                                               '',
                                               'Quit prompt (unsaved)', 'budget 79, keep $A1 tags'),
    # -- v0.4 ADDITIONS: Global Yes/No dialog buttons (used by all prompts) --
    (0x487478, b'Yes',
                                               '',
                                               'Global dialog button', 'budget 3, exact fit. Used by ALL Yes/No prompts'),
    (0x48748C, b'No',
                                               '',
                                               'Global dialog button', 'budget 3, exact fit. Used by ALL Yes/No prompts'),

    # -- v0.5 ADDITIONS: leftover strings spotted in-game --
    (0x4882B8, b'Save completed.',
                                               '',
                                               'Save confirmation', 'budget 15'),
    (0x495C58, b'Are you sure you wish to quit the game?',
                                               '',
                                               'Quit prompt (saved)', 'budget 39 — the short variant, used when no unsaved data'),
    (0x49B880, b'Changing the setting to \xe2\x9d\x9dAuto\xe2\x9d\x9e will open the window at a scale based on the Windows \xe2\x9d\x9dDisplay\xe2\x9d\x9e setting.\nChanging \xe2\x9d\x9d%%\xe2\x9d\x9e will scale the display, with the default resolution being %d\xc3\x97%d pixels.\n (Scale cannot be increased beyond the maximum resolution of your display.)',
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
