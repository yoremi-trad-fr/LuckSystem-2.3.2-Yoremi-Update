#!/usr/bin/env python3
"""
Kanon Steam — full FR patch table.
Single source of truth for all EN -> FR string patches in Kanon.exe.
Generates:
  - patches.h (C array to include in version.c)
  - patches.csv (human-readable review table)

Rules:
  - Each entry: (raw_offset_in_exe, en_bytes, fr_utf8_bytes, comment)
  - raw offset is converted to RVA on the fly: rva = raw + 0xC00
  - fr_bytes length MUST be <= slot_size - 1 (i.e., budget shown in the CSV)
  - A trailing \\0 is appended automatically; remaining padding stays \\0
  - If FR len equals EN len exactly, that's perfect.
  - If FR len is shorter, the extra bytes inside the slot stay \\0 (safe).

Conventions used for FR (aligned with Jeremy's ONE translation):
  - Present tense, "tu/toi" style (but not applicable to UI labels — formal infinitive)
  - UI verbs in infinitive form ("Activer", "Désactiver")
  - No bracketed qualifiers dropped
  - Accents kept where budget allows; dropped otherwise with ASCII fallback
  - 3-letter abbreviations (Min/Max) kept as-is since they're universal in FR
"""

# (raw_offset, EN original, FR proposed, context/tab, note)
# raw_offset is where the EN string starts in Kanon.exe file bytes
PATCHES = [
    # -- OPTIONS screen header / bottom buttons --
    (0x4874A4, b'Close',                       'Fermer',                'bottom-right button', 'slot 8, budget 7 -> "Fermer" 6B ok'),
    (0x488B98, b'Defaults',                    'Défauts',               'bottom-left button',  '"Défauts" é=2B so 7B, budget 15 ok'),

    # -- Basic tab --
    (0x499104, b'Basic',                       'Base',                  'Basic tab', 'already done via PAK? keep anyway'),
    (0x4991A0, b'Shortcut Menu',               'Raccourcis',            'Basic/Shortcut', 'budget 15'),
    (0x499144, b'Hide',                        'Cacher',                'Basic/Shortcut value', 'budget 11'),
    (0x499150, b'Display',                     'Visible',               'Basic/Shortcut value', 'budget 7 -> "Visible" 7B exact'),
    (0x499208, b'Skip',                        'Passer',                'Basic/Skip', 'budget 7'),
    (0x4991B0, b'Previously Read Only',        'Déjà lu seulement',     'Basic/Skip value', '"Déjà lu seulement" = 19B, budget 23 ok'),
    # "All" has slot 4 / budget 3 -> NO French word fits. Leave as-is with comment.
    # (0x499204, b'All',                         'All',                   'Basic/Skip value', 'SKIP: no 3-byte FR alt'),
    (0x499248, b'Position of Choices',         'Position des choix',    'Basic/Position',      '19B, budget 23 ok'),
    (0x499230, b'Bottom',                      'Bas',                   'Basic/Position value', 'budget 7'),
    (0x499238, b'Center',                      'Centre',                'Basic/Position value', 'budget 7'),
    (0x499320, b'Voice',                       'Voix',                  'Basic/Voice',          'budget 7'),
    (0x4992D0, b'Stop on New Message',         'Arrêt au message',      'Basic/Voice value',    '"Arrêt au message" = 17B (ê=2), budget 23'),
    (0x4992E8, b'      No Stops      ',        '     Sans arrêt     ',  'Basic/Voice value',    'keep spaces for alignment, FR 21B, budget 23'),
    (0x499398, b'Display Date',                'Afficher date',         'Basic/Date',           '13B, budget 15'),
    # Display Date: Off / On are shared with other "On/Off" toggles - risky to repatch globally
    # (0x4987E8, b'Off', 'Off', ...) - skip
    # (0x499358, b'On',  'On',  ...) - skip (too short to replace anyway: "Oui"/"Non" are 3B; budget 7 for "On" slot. "Marche"/"Arrêt" too long)

    (0x499450, b'Initial Cursor Position',     'Position du curseur',   'Basic/Cursor',         '20B, budget 23'),
    (0x4993A8, b'Positioned at \xe2\x9d\x9dYes\xe2\x9d\x9e', 'Sur « Oui »', 'Basic/Cursor value', 'Keep stylish chevrons if fits; budget 23'),
    (0x499420, b'Positioned at \xe2\x9d\x9dNo\xe2\x9d\x9e',  'Sur « Non »', 'Basic/Cursor value', 'budget 23'),
    (0x499490, b'Controller Rumble Function',  'Vibration manette',     'Basic/Rumble',         '17B, budget 27'),
    (0x499470, b'Disable',                     'Aucun',                 'Basic/Rumble value',   '"Aucun" 5B, budget 7 (rumble intensity: Aucun/Min/Mid/Max)'),
    # Min / Mid / Max (slot 4, budget 3) — keep as-is, universal
    # (0x499478, Min) - SKIP
    # (0x499480, Mid) - SKIP, or "Moy" 3B but awkward
    # (0x499488, Max) - SKIP

    # -- Text1 tab (fonts, colors, window) --
    (0x49A4D0, b'Text1',                       'Texte1',                'Text1 tab',            '6B, budget 7'),
    (0x49A4D8, b'Language',                    'Langue',                'Text1/Language',       '6B, budget 11'),
    (0x49A598, b'Font',                        'Police',                'Text1/Font',           '6B, budget 7'),
    (0x49A56C, b'Mincho',                      'Mincho',                'Text1/Font value',     'japanese font family name, keep'),
    (0x49A588, b'Modern',                      'Moderne',               'Text1/Font value',     '7B, budget 7 tight!'),
    (0x49A5D4, b'Solid',                       'Opaque',                'Text1/Window Transp value','6B, budget 11'),
    (0x49A5FC, b'Clear',                       'Clair',                 'Text1/Window Transp value','5B, budget 11'),
    (0x49A618, b'Window Transparency',         'Transparence',          'Text1/Window Transp',     '12B, budget 23'),
    (0x49A630, b'Only Choices',                'Choix seulement',       'Text1/Transp target',     '15B, budget 15 exact'),
    (0x49A658, b'Previously Read Text',        'Texte déjà lu',         'Text1/Read target',       '14B (é=2), budget 23'),
    (0x49A688, b'Color of',                    'Couleur',               'Text1/Color label',       '7B, budget 15'),
    (0x49A6D0, b'Green',                       'Vert',                  'Text1/Color value',       '4B, budget 7'),
    (0x49A6D8, b'Blue',                        'Bleu',                  'Text1/Color value',       '4B, budget 7'),
    (0x49A6F0, b'Purple',                      'Violet',                'Text1/Color value',       '6B, budget 7'),
    # (0x49A6FC, 'Red' -> 'Rouge' 5B but slot 4 / budget 3 TOO LONG -> keep "Red" or use "R.")
    (0x49A494, b'Orange',                      'Orange',                'Text1/Color value',       'same word in FR'),
    (0x49A4C0, b'Yellow',                      'Jaune',                 'Text1/Color value',       '5B, budget 7'),
    (0x49A700, b'Read Text Color',             'Coul. texte lu',        'Text1/Read color',        '14B, budget 15'),

    # -- Text2 tab --
    (0x49AABC, b'Text2',                       'Texte2',                'Text2 tab',            '6B, budget 11'),
    (0x49AB30, b'Text Speed',                  'Vitesse texte',         'Text2/Speed',          '13B, budget 15'),
    (0x49AAEC, b'Slow',                        'Lent',                  'Text2/Speed value',    '4B, budget 7'),
    (0x49AAF4, b'Fast',                        'Rapide',                'Text2/Speed value',    '6B, budget 7'),
    (0x49AB40, b'0 sec/char',                  '0 s/car',               'Text2/Speed value',    '7B, budget 15'),
    (0x49AB90, b'0.1 sec/char',                '0,1 s/car',             'Text2/Speed value',    '9B, budget 15 (comma for FR decimal)'),
    (0x49ABB0, b'Wait Time Per Character',     'Attente par caractère', 'Text2/Wait',           '22B (è=2), budget 23'),
    (0x49ABC8, b'0 sec',                       '0 s',                   'Text2/Wait value',     '3B, budget 7'),
    (0x49ABD8, b'1 sec',                       '1 s',                   'Text2/Wait value',     '3B, budget 7'),
    (0x49ABE8, b'2 sec',                       '2 s',                   'Text2/Wait value',     '3B, budget 7'),
    (0x49ABF8, b'3 sec',                       '3 s',                   'Text2/Wait value',     '3B, budget 7'),
    (0x49AC08, b'Base Wait Time',              'Attente de base',       'Text2/Base',           '15B, budget 15 exact'),

    # -- Sound tab --
    (0x49AC5C, b'Sound',                       'Son',                   'Sound tab',            '3B, budget 11'),
    (0x49AC68, b'Master Volume',               'Volume maître',         'Sound/Master',         '"Volume maître" 14B (î=2), budget 15'),
    # SFX: slot 8 / budget 7. "Effets" 6B, but SFX universal. Keep "SFX"
    (0x49AD08, b'System Sounds',               'Sons système',          'Sound/System',         '13B (è=2), budget 15'),

    # -- Voice tab --
    # (0x49AD30 = \\t(Young)\\t → voice sample label, budget 15, "\\t(Jeune)\\t" = 11B if FR. But it's user-facing)
    (0x49AD30, b'\t(Young)\t',                 '\t(Jeune)\t',          'Voice tab young label',   'Yuichi child voice. 9B, budget 15'),
    # Yuichi Aizawa — character name, DO NOT translate
    # PARTS/voice_icon — asset path, DO NOT translate

    # -- Keyboard tab --
    # keyboard (0x49ADA0) likely asset label, keep
    (0x49ADF0, b'While Pressed',               'Maintenu',              'Keyboard/mode',        '8B, budget 15'),
    (0x49AE00, b'Start/Stop',                  'Marche/Arrêt',          'Keyboard/mode',        '13B (ê=2), budget 15'),
    (0x49AE58, b'C (Skip)',                    'C (Passer)',            'Keyboard/key label',   '10B, budget 15'),
    (0x49AE68, b'Z (Rewind)',                  'Z (Retour)',            'Keyboard/key label',   '10B, budget 15'),
    (0x49AEC0, b'  Disable  ',                 '  Désactivé ',          'Keyboard/value',       'retain spaces; 13B (é=2), budget 15'),
    (0x49AED0, b'Quick Save',                  'Sauv. rapide',          'Keyboard/key label',   '12B, budget 15'),
    (0x49AF20, b'Quick Load',                  'Chrg. rapide',          'Keyboard/key label',   '12B, budget 15'),
    (0x49AF80, b'Switch Language',             'Changer langue',        'Keyboard/key label',   '14B, budget 15'),
    (0x49AFE0, b'Up Arrow',                    'Flèche haut',           'Keyboard/key label',   '12B (è=2), budget 15'),

    # -- Mouse tab --
    (0x49B02C, b'Mouse',                       'Souris',                'Mouse tab',            '6B, budget 11'),
    (0x49B080, b'System Menu',                 'Menu système',          'Mouse/target',         '13B (è=2), budget 15'),
    (0x49B090, b'Hide Window',                 'Cacher fenêtre',        'Mouse/target',         '15B (ê=2), budget 15 exact'),
    (0x49B0E8, b'Right Click',                 'Clic droit',            'Mouse/binding',        '10B, budget 15'),
    (0x49B108, b'Left+Right Click',            'Clic gauche+droit',     'Mouse/binding',        '18B, budget 23'),
    (0x49B120, b'Mouse Wheel Button',          'Bouton molette',        'Mouse/binding',        '14B, budget 23'),
    (0x49B180, b'Rewind Once',                 'Retour un pas',         'Mouse/target',         '13B, budget 15'),
    (0x49B190, b'Wheel Up',                    'Molette haut',          'Mouse/binding',        '12B, budget 15'),
    (0x49B1F0, b'Forward Once',                'Avance un pas',         'Mouse/target',         '14B, budget 15'),
    (0x49B200, b'Wheel Down',                  'Molette bas',           'Mouse/binding',        '11B, budget 15'),
    (0x49B280, b'Jump and Switch Pages',       'Saut et pages',         'Mouse/target',         '14B, budget 23'),
    (0x49B298, b'Return/Proceed Button',       'Retour/Continuer',      'Mouse/target',         '16B, budget 23'),
    (0x49B2F4, b'Enable',                      'Activé',                'Mouse/Gestures value', '7B (é=2), budget 11'),
    (0x49B310, b'Gestures',                    'Gestes',                'Mouse/Gestures',       '6B, budget 15'),
    (0x49B350, b'Dialog and Choices',          'Dialogue et choix',     'Mouse/Snap target',    '18B, budget 23'),
    (0x49B3D8, b'Snap Pointer',                'Aimanter curseur',      'Mouse/Snap',           '16B, budget 23'),

    # -- System tab --
    (0x49B748, b'Game Ver. ',                  'Version jeu ',          'System/info',          '12B, budget 15'),
    # "Steam Deck: Proton " keep as-is — technical info, same in FR
    (0x49B76C, b'System',                      'Système',               'System tab',           '8B (è=2), budget 11'),
    (0x49B7B0, b'Window',                      'Fenetre',               'System/Window',        '"Fenetre" no accent, 7B, budget 7 exact'),
    (0x49B7B8, b'Full Screen',                 'Plein écran',           'System/FullScreen',    '12B (é=2), budget 15'),
    (0x49B7E0, b'Screen Mode',                 'Mode écran',            'System/ScreenMode',    '11B (é=2), budget 15'),
    (0x49B820, b'Auto',                        'Auto',                  'System/value',         'same in FR'),
    (0x49B870, b'Window Size',                 'Taille fenêtre',        'System/WindowSize',    '15B (ê=2), budget 15 exact'),

    # -- Long tooltips (lots of budget) --
    (0x49A8D0, b'Wait Time Per Character: In Auto Mode, this sets the wait time until the next message is displayed based on the number of characters in text.',
                                               'Attente par caractère : en mode Auto, règle le temps avant le message suivant selon le nombre de caractères du texte.',
                                               'Texte2 tooltip', 'budget 143'),
    (0x49A960, b'Base Wait Time: You can set a Base Wait Time to add to\n\xe2\x9d\x9dWait Time Per Character\xe2\x9d\x9e.',
                                               'Attente de base : définit un temps fixe à ajouter à\n❝Attente par caractère❞.',
                                               'Texte2 tooltip', 'budget 95, keep stylized quotes'),
    (0x4996F0, b'Initial Cursor Position: Sets the initial position of the cursor when a Yes/No choice is available.',
                                               'Position du curseur : définit la position initiale du curseur pour un choix Oui/Non.',
                                               'Basic tooltip', 'budget 111'),
    (0x499960, b'Controller Rumble Function: Plug in a controller to use the controller Rumble function.',
                                               'Vibration manette : branchez une manette pour utiliser la vibration.',
                                               'Basic tooltip', 'budget 87'),
    (0x49B3F0, b'Gestures: Moving the cursor while holding the left button works the same way as Touch controls.',
                                               'Gestes : maintenir le bouton gauche en déplaçant le curseur agit comme en mode tactile.',
                                               'Mouse tooltip', 'budget 95'),
    (0x49B6A0, b'Left+Right Click: Hold left button then right click to switch between languages (English/Simplified Chinese/Japanese).',
                                               'Clic gauche+droit : maintenir le clic gauche puis clic droit pour changer de langue (anglais/chinois/japonais).',
                                               'Mouse tooltip', 'budget 119'),
    (0x499640, b'Voice: If you select No Stops, sound will continue to play even if you advance the text during voice playback. (It stops if there is sound on the next message.)',
                                               'Voix : avec Sans arrêt, le son continue même si vous avancez le texte durant la lecture. (S\'arrête si le message suivant a du son.)',
                                               'Basic tooltip', 'budget 175'),

    # -- v0.3 ADDITIONS: Save/Load prompts --
    (0x487EB0, b'Do you wish to load this save?',
                                               'Charger cette sauvegarde ?',
                                               'Save/Load prompt', 'budget 31'),
    (0x4880C8, b'Are you sure you wish to overwrite data?',
                                               'Écraser la sauvegarde ?',
                                               'Save/Load prompt', 'budget 47'),
    (0x488080, b'Do you wish to save?',
                                               'Sauvegarder ?',
                                               'Save prompt', 'budget 23, bonus'),
    (0x4882F0, b'Are you sure you wish to delete this save data?',
                                               'Supprimer cette sauvegarde ?',
                                               'Save delete prompt', 'budget 47'),
    (0x488350, b'This cannot be deleted.',
                                               'Suppression impossible.',
                                               'Save delete error', 'budget 23'),
    (0x488B2C, b'Delete',
                                               'Suppr.',
                                               'Save menu button', 'budget 11'),
    (0x488B54, b'Latest',
                                               'Récent',
                                               'Save menu button', 'budget 11 (é=2)'),

    # -- v0.3 ADDITIONS: Text1 tab extras --
    (0x488EE8, b'Text preview.',
                                               'Aperçu.',
                                               'Text preview label', 'budget 15'),
    (0x489008, b'Settings such as text speed are reflected.',
                                               'Les réglages comme la vitesse sont appliqués.',
                                               'Text preview tooltip', 'budget 47'),
    (0x487458, b'English',
                                               'FR',
                                               'Language switcher', 'budget 7 hard cap (neighbor is 简体中文). Using ISO code "FR" (2B) for cleanness. Alternatives in 7B: Franç./Franc./Francai'),

    # -- v0.3 ADDITIONS: Read Text Color On/Off toggles --
    (0x4987E8, b'Off',
                                               'Non',
                                               'Read text color toggle', 'budget 7, "Non" 3B'),
    (0x499358, b'On',
                                               'Oui',
                                               'Read text color toggle', 'budget 7, "Oui" 3B'),

    # -- v0.3 ADDITIONS: Right-click / system menu prompts --
    (0x4883F8, b'Return to the title screen?',
                                               'Retour au titre ?',
                                               'Title return prompt', 'budget 31'),
    (0x4883A0, b'Return to the menu?',
                                               'Retour au menu ?',
                                               'Menu return prompt', 'budget 23, bonus'),
    (0x495B30, b'$A1There is unsaved data.\n$A1Are you sure you wish to quit the game?',
                                               '$A1Données non sauvegardées.\n$A1Quitter le jeu ?',
                                               'Quit prompt (unsaved)', 'budget 79, keep $A1 tags'),
    # -- v0.4 ADDITIONS: Global Yes/No dialog buttons (used by all prompts) --
    (0x487478, b'Yes',
                                               'Oui',
                                               'Global dialog button', 'budget 3, exact fit. Used by ALL Yes/No prompts'),
    (0x48748C, b'No',
                                               'Non',
                                               'Global dialog button', 'budget 3, exact fit. Used by ALL Yes/No prompts'),

    # -- v0.5 ADDITIONS: leftover strings spotted in-game --
    (0x4882B8, b'Save completed.',
                                               'Sauvegardé.',
                                               'Save confirmation', '"Sauvegardé." 12B (é=2), budget 15'),
    (0x495C58, b'Are you sure you wish to quit the game?',
                                               'Quitter le jeu ?',
                                               'Quit prompt (saved)', 'budget 39 — the short variant, used when no unsaved data'),
    (0x49B880, b'Changing the setting to \xe2\x9d\x9dAuto\xe2\x9d\x9e will open the window at a scale based on the Windows \xe2\x9d\x9dDisplay\xe2\x9d\x9e setting.\nChanging \xe2\x9d\x9d%%\xe2\x9d\x9e will scale the display, with the default resolution being %d\xc3\x97%d pixels.\n (Scale cannot be increased beyond the maximum resolution of your display.)',
                                               'Le réglage ❝Auto❞ ouvre la fenêtre à une échelle basée sur les paramètres ❝Affichage❞ de Windows.\nLe réglage ❝%%❞ ajuste l\'échelle, la résolution par défaut étant %d×%d pixels.\n (L\'échelle ne peut dépasser la résolution maximale de votre écran.)',
                                               'System tab tooltip', 'budget 287, keep %%/%d×%d intact'),
]


def main():
    import sys
    data = open('Kanon.exe','rb').read()

    def slot_size(start):
        i = start
        while data[i] != 0: i += 1
        while i < len(data) and data[i] == 0: i += 1
        return i - start

    rows = []
    errors = []
    for off, en, fr, context, note in PATCHES:
        actual = data[off:off+len(en)]
        if actual != en:
            errors.append(f"0x{off:X}: expected {en!r}, got {actual!r}")
            continue
        fr_bytes = fr.encode('utf-8')
        slot = slot_size(off)
        budget = slot - 1
        fits = len(fr_bytes) <= budget
        rows.append({
            'off': off,
            'en': en.decode('utf-8', errors='replace'),
            'fr': fr,
            'en_len': len(en),
            'fr_len': len(fr_bytes),
            'slot': slot,
            'budget': budget,
            'fits': fits,
            'context': context,
            'note': note,
            'en_bytes': en,
            'fr_bytes': fr_bytes,
        })

    if errors:
        print("=== OFFSET MISMATCH (aborting) ===", file=sys.stderr)
        for e in errors:
            print("  " + e, file=sys.stderr)
        sys.exit(1)

    # Print status
    print(f"{'off':>8}  {'slot':>4}  {'EN':>3}  {'FR':>3}  {'fit':3}  EN -> FR")
    print('-' * 100)
    n_ok = n_bad = 0
    for r in rows:
        mark = '✓' if r['fits'] else '✗'
        if r['fits']: n_ok += 1
        else:          n_bad += 1
        print(f"0x{r['off']:06X}  {r['slot']:>4}  {r['en_len']:>3}  {r['fr_len']:>3}  {mark}   {r['en']!r} -> {r['fr']!r}")
    print(f"\nTotal: {len(rows)}  OK: {n_ok}  FAIL: {n_bad}")
    if n_bad:
        print("\nFailures (too long):")
        for r in rows:
            if not r['fits']:
                print(f"  0x{r['off']:X}: FR {r['fr_len']}B > budget {r['budget']}B: {r['fr']!r}")
        sys.exit(2)

    # Emit patches.h (C header)
    with open('patches.h','w', encoding='utf-8') as f:
        f.write('/* Auto-generated from patches.py. Do not edit. */\n')
        f.write('#ifndef KANON_FR_PATCHES_H\n#define KANON_FR_PATCHES_H\n\n')
        for i, r in enumerate(rows):
            rva = r['off'] + 0xC00
            # Write length = max(len(EN)+1, len(FR)+1): overwrite at least full EN
            # to erase stale bytes, and at least full FR so no truncation.
            write_len = max(len(r['en_bytes']), len(r['fr_bytes'])) + 1
            en_padded = list(r['en_bytes']) + [0] * (write_len - len(r['en_bytes']))
            fr_padded = list(r['fr_bytes']) + [0] * (write_len - len(r['fr_bytes']))
            assert len(en_padded) == write_len and len(fr_padded) == write_len
            en_arr = ','.join(f'0x{b:02X}' for b in en_padded)
            fr_arr = ','.join(f'0x{b:02X}' for b in fr_padded)
            f.write(f'static const BYTE s_en_{i:03d}[] = {{ {en_arr} }};\n')
            f.write(f'static const BYTE s_fr_{i:03d}[] = {{ {fr_arr} }};\n')
        f.write('\nstatic const KanonPatch g_patches[] = {\n')
        for i, r in enumerate(rows):
            rva = r['off'] + 0xC00
            write_len = max(len(r['en_bytes']), len(r['fr_bytes'])) + 1
            # Escape newlines, backslashes, and double quotes for C string literal
            ctx = r['context'] + ': ' + r['en'][:30]
            ctx = ctx.replace('\\', '\\\\').replace('"', '\\"').replace('\n', '\\n').replace('\r', '\\r').replace('\t', '\\t')
            f.write(f'    {{ 0x{rva:06X}, {write_len:>4}, s_en_{i:03d}, s_fr_{i:03d}, "{ctx}" }},\n')
        f.write('};\n\n#define N_PATCHES (sizeof(g_patches)/sizeof(g_patches[0]))\n')
        f.write('\n#endif\n')
    print(f"\nGenerated patches.h with {len(rows)} entries.")

    # Emit patches.csv
    with open('patches.csv','w', encoding='utf-8') as f:
        f.write('raw_offset,rva,slot,budget,en_len,fr_len,fits,en,fr,context,note\n')
        for r in rows:
            rva = r['off'] + 0xC00
            # CSV-escape quotes in text
            def esc(s): return '"' + s.replace('"','""') + '"'
            f.write(f'0x{r["off"]:X},0x{rva:X},{r["slot"]},{r["budget"]},{r["en_len"]},{r["fr_len"]},{r["fits"]},{esc(r["en"])},{esc(r["fr"])},{esc(r["context"])},{esc(r["note"])}\n')
    print(f"Generated patches.csv with {len(rows)} entries.")


if __name__ == '__main__':
    main()
