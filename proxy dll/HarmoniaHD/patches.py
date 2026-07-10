#!/usr/bin/env python3
"""
Harmonia Full HD Edition (Luck Engine) — proxy DLL patch table.

This table was ported from the Kanon Steam proxy table by matching the same
hardcoded Luck Engine UI strings inside HarmoniaFHD.exe.

Current targets are French replacements for the hardcoded English strings.
Entries must fit inside the original string slot budget.

Generated files:
  - patches.h
  - patches.csv

Notes:
  - HarmoniaFHD.exe imports VERSION.dll, so the shared version.dll proxy can load.
  - .rdata raw -> RVA delta is 0xE00 for the current Steam executable.
  - Ambiguous strings were only accepted when they landed near the surrounding
    unique strings. Skipped entries are documented in port_report.md.
"""

from pathlib import Path

GAME_EXE = r'C:\Program Files (x86)\Steam\steamapps\common\Harmonia Full HD Edition\HarmoniaFHD.exe'
RVA_DELTA = 0xE00
PATCH_GAME_NAME = 'HarmoniaHD'
PATCH_VERSION = '0.1'

# Each entry: (raw_offset, src_bytes, target_str, context, note)
# Ready-to-use state: target_str contains the French replacement.

PATCHES = [
    (0x48DBC8, b'Defaults', 'Defauts', 'bottom-left button', 'budget 15 ok'),
    (0x48C4D4, b'Close', 'Fermer', 'bottom-right button', 'slot 8, budget 7; AUTO-SELECTED multiple/6 nearest-local dist=0x0'),
    (0x49DB60, b'Basic', 'Base', 'Basic tab', 'already done via PAK? keep anyway'),
    (0x49DBC0, b'Skip', 'Saut', 'Basic/Skip', 'budget 7; AUTO-SELECTED multiple/5 nearest-local dist=0x18'),
    (0x49DB68, b'Previously Read Only', 'Deja lu seulement', 'Basic/Skip value', 'budget 23'),
    (0x49DC88, b'Voice', 'Voix', 'Basic/Voice', 'budget 7; AUTO-SELECTED multiple/18 nearest-local dist=0x1C'),
    (0x49DC38, b'Stop on New Message', 'Stop au nouveau msg', 'Basic/Voice value', 'budget 23'),
    (0x49DC50, b'      No Stops      ', '     Sans arret     ', 'Basic/Voice value', 'keep spaces for alignment, budget 23'),
    (0x49DD58, b'Initial Cursor Position', 'Position initiale', 'Basic/Cursor', '20B, budget 23; AUTO-SELECTED multiple/2 nearest-local dist=0x1'),
    (0x49DCB0, b'Positioned at \xe2\x9d\x9dYes\xe2\x9d\x9e', 'Place sur ❝Oui❞', 'Basic/Cursor value', 'Keep stylish chevrons if fits; budget 23'),
    (0x49DD28, b'Positioned at \xe2\x9d\x9dNo\xe2\x9d\x9e', 'Place sur ❝Non❞', 'Basic/Cursor value', 'budget 23'),
    (0x49DD98, b'Controller Rumble Function', 'Vibration manette', 'Basic/Rumble', '17B, budget 27; AUTO-SELECTED multiple/2 nearest-local dist=0x2'),
    (0x49DD78, b'Disable', 'Arret', 'Basic/Rumble value', 'budget 7 (rumble intensity: None/Min/Mid/Max); AUTO-SELECTED multiple/3 nearest-local dist=0x1'),
    (0x49EDCC, b'Text1', 'Texte1', 'Text1 tab', '6B, budget 7'),
    (0x49EE00, b'Language', 'Langue', 'Text1/Language', '6B, budget 11; AUTO-SELECTED multiple/9 nearest-local dist=0x2C'),
    (0x49EE94, b'Font', 'Police', 'Text1/Font', '6B, budget 7; AUTO-SELECTED multiple/3 nearest-local dist=0x5'),
    (0x49EE68, b'Mincho', 'Mincho', 'Text1/Font value', 'japanese font family name, keep'),
    (0x49EE7C, b'Modern', 'Moderne', 'Text1/Font value', '7B, budget 7 tight!'),
    (0x49EED4, b'Solid', 'Opaque', 'Text1/Window Transp value', '6B, budget 11'),
    (0x49EEEC, b'Clear', 'Transp', 'Text1/Window Transp value', '5B, budget 11; AUTO-SELECTED multiple/5 nearest-local dist=0x2'),
    (0x49EF00, b'Window Transparency', 'Transp. fenetre', 'Text1/Window Transp', '12B, budget 23'),
    (0x49EF90, b'Previously Read Text', 'Texte deja lu', 'Text1/Read target', 'budget 23'),
    (0x49EFC8, b'Color of', 'Couleur', 'Text1/Color label', '7B, budget 15'),
    (0x49EFD4, b'Green', 'Vert', 'Text1/Color value', '4B, budget 7'),
    (0x49EFEC, b'Blue', 'Bleu', 'Text1/Color value', '4B, budget 7; AUTO-SELECTED multiple/2 nearest-local dist=0x10'),
    (0x49EFF4, b'Purple', 'Violet', 'Text1/Color value', '6B, budget 7'),
    (0x49EDB8, b'Orange', 'Orange', 'Text1/Color value', 'same word'),
    (0x49EDC0, b'Yellow', 'Jaune', 'Text1/Color value', '5B, budget 7'),
    (0x49F028, b'Read Text Color', 'Couleur texte lu', 'Text1/Read color', '14B, budget 15'),
    (0x49F3E8, b'Text2', 'Texte2', 'Text2 tab', '6B, budget 11'),
    (0x49F418, b'Text Speed', 'Vitesse texte', 'Text2/Speed', '13B, budget 15'),
    (0x49F3F0, b'Slow', 'Lent', 'Text2/Speed value', '4B, budget 7'),
    (0x49F410, b'Fast', 'Rapide', 'Text2/Speed value', '6B, budget 7'),
    (0x49F470, b'0 sec/char', '0 s/car.', 'Text2/Speed value', '7B, budget 15'),
    (0x49F480, b'0.1 sec/char', '0.1 s/car.', 'Text2/Speed value', '9B, budget 15'),
    (0x49F4A0, b'Wait Time Per Character', 'Attente par caractere', 'Text2/Wait', 'budget 23; AUTO-SELECTED multiple/3 nearest-local dist=0x11'),
    (0x49F4E8, b'0 sec', '0 s', 'Text2/Wait value', '3B, budget 7; AUTO-SELECTED multiple/2 nearest-local dist=0x11'),
    (0x49F4F8, b'1 sec', '1 s', 'Text2/Wait value', '3B, budget 7; AUTO-SELECTED multiple/2 nearest-local dist=0x9'),
    (0x49F508, b'2 sec', '2 s', 'Text2/Wait value', '3B, budget 7'),
    (0x49F524, b'3 sec', '3 s', 'Text2/Wait value', '3B, budget 7'),
    (0x49F548, b'Base Wait Time', 'Attente base', 'Text2/Base', '15B, budget 15 exact; AUTO-SELECTED multiple/3 nearest-local dist=0x12'),
    (0x49F558, b'Sound', 'Son', 'Sound tab', '3B, budget 11; AUTO-SELECTED multiple/23 nearest-local dist=0x3B'),
    (0x49F5A0, b'Master Volume', 'Volume global', 'Sound/Master', 'budget 15'),
    (0x49F5F0, b'System Sounds', 'Sons systeme', 'Sound/System', 'budget 15'),
    (0x49F678, b'While Pressed', 'Maintenu', 'Keyboard/mode', '8B, budget 15'),
    (0x49F6D0, b'Start/Stop', 'Marche/Arret', 'Keyboard/mode', 'budget 15'),
    (0x49F6E0, b'C (Skip)', 'C (Saut)', 'Keyboard/key label', '10B, budget 15'),
    (0x49F740, b'Z (Rewind)', 'Z (Retour)', 'Keyboard/key label', '10B, budget 15'),
    (0x49F750, b'  Disable  ', '  Inactif  ', 'Keyboard/value', 'retain spaces, budget 15'),
    (0x49F7A8, b'Quick Save', 'Sauv. rap.', 'Keyboard/key label', '12B, budget 15'),
    (0x49F7F0, b'Quick Load', 'Charg. rapide', 'Keyboard/key label', '12B, budget 15'),
    (0x49F800, b'Switch Language', 'Changer langue', 'Keyboard/key label', '14B, budget 15'),
    (0x49F868, b'Up Arrow', 'Fleche haut', 'Keyboard/key label', 'budget 15'),
    (0x49F8EC, b'Mouse', 'Souris', 'Mouse tab', '6B, budget 11; AUTO-SELECTED multiple/2 nearest-local dist=0xA'),
    (0x49F968, b'System Menu', 'Menu systeme', 'Mouse/target', 'budget 15'),
    (0x49F8F8, b'Hide Window', 'Cacher fenetre', 'Mouse/target', 'budget 15 exact'),
    (0x49F99D, b'Right Click', 'Clic droit', 'Mouse/binding', '10B, budget 15; AUTO-SELECTED multiple/3 nearest-local dist=0xD'),
    (0x49F998, b'Left+Right Click', 'Clic gauche+droit', 'Mouse/binding', '18B, budget 23; AUTO-SELECTED multiple/2 nearest-local dist=0x2F'),
    (0x49F9F0, b'Mouse Wheel Button', 'Bouton molette', 'Mouse/binding', '14B, budget 23'),
    (0x49FA08, b'Rewind Once', 'Retour x1', 'Mouse/target', '13B, budget 15'),
    (0x49FA68, b'Wheel Up', 'Molette haut', 'Mouse/binding', '12B, budget 15'),
    (0x49FA78, b'Forward Once', 'Avance x1', 'Mouse/target', '14B, budget 15'),
    (0x49FAD8, b'Wheel Down', 'Molette bas', 'Mouse/binding', '11B, budget 15'),
    (0x49FAE8, b'Jump and Switch Pages', 'Saut/changer page', 'Mouse/target', '14B, budget 23'),
    (0x49FB78, b'Return/Proceed Button', 'Bouton retour/avance', 'Mouse/target', '16B, budget 23'),
    (0x49FB90, b'Enable', 'Actif', 'Mouse/Gestures value', 'budget 11'),
    (0x49FBA0, b'Gestures', 'Gestes', 'Mouse/Gestures', '6B, budget 15; AUTO-SELECTED multiple/2 nearest-local dist=0x4'),
    (0x49FC30, b'Snap Pointer', 'Aimant curseur', 'Mouse/Snap', '16B, budget 23'),
    (0x4A0008, b'Game Ver. ', 'Ver. jeu ', 'System/info', '12B, budget 15'),
    (0x4A001C, b'System', 'Systeme', 'System tab', 'budget 11; AUTO-SELECTED multiple/23 nearest-local dist=0xE'),
    (0x4A0024, b'Window', 'Fenetre', 'System/Window', '7B, budget 7 exact; AUTO-SELECTED multiple/29 nearest-local dist=0x47'),
    (0x4A0078, b'Full Screen', 'Plein ecran', 'System/FullScreen', 'budget 15; AUTO-SELECTED multiple/2 nearest-local dist=0x6'),
    (0x4A0098, b'Screen Mode', 'Mode ecran', 'System/ScreenMode', 'budget 15'),
    (0x4A00A4, b'Auto', 'Auto', 'System/value', 'same; AUTO-SELECTED multiple/7 nearest-local dist=0x1B'),
    (0x4A00F0, b'Window Size', 'Taille fen.', 'System/WindowSize', 'budget 15 exact'),
    (0x49F040, b'Wait Time Per Character: In Auto Mode, this sets the wait time until the next message is displayed based on the number of characters in text.', 'En mode auto, definit le delai avant le message suivant selon le nombre de caracteres du texte.', 'Texte2 tooltip', 'budget 143'),
    (0x49F360, b'Base Wait Time: You can set a Base Wait Time to add to\n\xe2\x9d\x9dWait Time Per Character\xe2\x9d\x9e.', "Ajoute un delai fixe a l'attente par caractere.", 'Texte2 tooltip', 'budget 95, keep stylized quotes'),
    (0x49E000, b'Initial Cursor Position: Sets the initial position of the cursor when a Yes/No choice is available.', "Definit la position initiale du curseur lors d'un choix Oui/Non.", 'Basic tooltip', 'budget 111'),
    (0x49E270, b'Controller Rumble Function: Plug in a controller to use the controller Rumble function.', 'Branchez une manette pour utiliser la vibration.', 'Basic tooltip', 'budget 87'),
    (0x49FDA0, b'Gestures: Moving the cursor while holding the left button works the same way as Touch controls.', 'Maintenir le bouton gauche et bouger le curseur agit comme les commandes tactiles.', 'Mouse tooltip', 'budget 95'),
    (0x49FE00, b'Left+Right Click: Hold left button then right click to switch between languages (English/Simplified Chinese/Japanese).', 'Maintenez le bouton gauche puis clic droit pour changer de langue.', 'Mouse tooltip', 'budget 119'),
    (0x49DF50, b'Voice: If you select No Stops, sound will continue to play even if you advance the text during voice playback. (It stops if there is sound on the next message.)', "Avec Sans arret, la voix continue si vous avancez le texte pendant sa lecture. Elle s'arrete si le message suivant a un son.", 'Basic tooltip', 'budget 175'),
    (0x48CEE0, b'Do you wish to load this save?', 'Charger cette sauvegarde ?', 'Save/Load prompt', 'budget 31'),
    (0x48D0F8, b'Are you sure you wish to overwrite data?', 'Ecraser cette sauvegarde ?', 'Save/Load prompt', 'budget 47'),
    (0x48D0B0, b'Do you wish to save?', 'Sauvegarder ?', 'Save prompt', 'budget 23, bonus'),
    (0x48D320, b'Are you sure you wish to delete this save data?', 'Supprimer cette sauvegarde ?', 'Save delete prompt', 'budget 47'),
    (0x48D380, b'This cannot be deleted.', 'Suppression impossible.', 'Save delete error', 'budget 23'),
    (0x48DB5C, b'Delete', 'Suppr.', 'Save menu button', 'budget 11; AUTO-SELECTED multiple/10 nearest-local dist=0x0'),
    (0x48DB84, b'Latest', 'Dernier', 'Save menu button', 'budget 11'),
    (0x48DF18, b'Text preview.', 'Apercu texte.', 'Text preview label', 'budget 15'),
    (0x48E038, b'Settings such as text speed are reflected.', 'Les reglages de texte sont appliques.', 'Text preview tooltip', 'budget 47'),
    (0x48C488, b'English', 'FR', 'Language switcher', 'budget 7 hard cap (neighbor is 简体中文). Fill in your ISO language code.; AUTO-SELECTED multiple/2 nearest-local dist=0x0'),
    (0x49CEB8, b'Off', 'Non', 'Read text color toggle', 'budget 7'),
    (0x48D428, b'Return to the title screen?', "Retour a l'ecran titre ?", 'Title return prompt', 'budget 31'),
    (0x48D3D0, b'Return to the menu?', 'Retour au menu ?', 'Menu return prompt', 'budget 23, bonus'),
    (0x49A780, b'$A1There is unsaved data.\n$A1Are you sure you wish to quit the game?', '$A1Donnees non sauvegardees.\n$A1Quitter le jeu ?', 'Quit prompt (unsaved)', 'budget 79, keep $A1 tags'),
    (0x48C4A8, b'Yes', 'Oui', 'Global dialog button', 'budget 3, exact fit. Used by ALL Yes/No prompts; AUTO-SELECTED multiple/3 nearest-local dist=0x0'),
    (0x48C4BC, b'No', 'Non', 'Global dialog button', 'budget 3, exact fit. Used by ALL Yes/No prompts; AUTO-SELECTED multiple/96 nearest-local dist=0x0'),
    (0x48D2E8, b'Save completed.', 'Sauvegarde OK', 'Save confirmation', 'budget 15'),
    (0x49A7C8, b'Are you sure you wish to quit the game?', 'Quitter le jeu ?', 'Quit prompt (saved)', 'budget 39 — the short variant, used when no unsaved data; AUTO-SELECTED multiple/2 nearest-local dist=0xBC'),
    (0x4A0390, b'Changing the setting to \xe2\x9d\x9dAuto\xe2\x9d\x9e will open the window at a scale based on the Windows \xe2\x9d\x9dDisplay\xe2\x9d\x9e setting.\nChanging \xe2\x9d\x9d%%\xe2\x9d\x9e will scale the display, with the default resolution being %d\xc3\x97%d pixels.\n (Scale cannot be increased beyond the maximum resolution of your display.)', "Le mode ❝Auto❞ ouvre la fenetre selon le reglage Windows ❝Affichage❞.\nModifier ❝%%❞ redimensionne l'affichage, resolution par defaut : %d×%d pixels.\n(Le zoom ne peut pas depasser la resolution maximale de l'ecran.)", 'System tab tooltip', 'budget 287, keep %%/%d×%d intact'),
]


def main():
    import sys
    if hasattr(sys.stdout, 'reconfigure'):
        sys.stdout.reconfigure(encoding='utf-8')
        sys.stderr.reconfigure(encoding='utf-8')

    data = Path(GAME_EXE).read_bytes()

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
            'off': off,
            'src': src.decode('utf-8', errors='replace'),
            'target': target,
            'src_len': len(src),
            'target_len': len(target_bytes),
            'slot': slot,
            'budget': budget,
            'fits': fits,
            'context': context,
            'note': note,
            'src_bytes': src,
            'target_bytes': target_bytes,
        })

    if errors:
        print('=== OFFSET MISMATCH (aborting) ===', file=sys.stderr)
        for e in errors:
            print('  ' + e, file=sys.stderr)
        sys.exit(1)

    print(f"{'off':>8}  {'slot':>4}  {'src':>3}  {'tgt':>3}  {'fit':3}  src -> target")
    print('-' * 100)
    n_ok = n_bad = 0
    for r in rows:
        mark = 'OK' if r['fits'] else 'NO'
        if r['fits']:
            n_ok += 1
        else:
            n_bad += 1
        print(f"0x{r['off']:06X}  {r['slot']:>4}  {r['src_len']:>3}  {r['target_len']:>3}  {mark:3}  {r['src']!r} -> {r['target']!r}")
    print(f"\nTotal: {len(rows)}  OK: {n_ok}  FAIL: {n_bad}")
    if n_bad:
        print('\nFailures (target too long):')
        for r in rows:
            if not r['fits']:
                print(f"  0x{r['off']:X}: target {r['target_len']}B > budget {r['budget']}B: {r['target']!r}")
        sys.exit(2)

    with open('patches.h', 'w', encoding='utf-8') as f:
        f.write('/* Auto-generated from patches.py. Do not edit. */\n')
        f.write('#ifndef LUCKPROXY_PATCHES_H\n#define LUCKPROXY_PATCHES_H\n\n')
        f.write(f'#define PATCH_GAME_NAME "{PATCH_GAME_NAME}"\n')
        f.write(f'#define PATCH_VERSION   "{PATCH_VERSION}"\n\n')
        for i, r in enumerate(rows):
            rva = r['off'] + RVA_DELTA
            write_len = max(len(r['src_bytes']), len(r['target_bytes'])) + 1
            src_padded = list(r['src_bytes']) + [0] * (write_len - len(r['src_bytes']))
            target_padded = list(r['target_bytes']) + [0] * (write_len - len(r['target_bytes']))
            src_arr = ','.join(f'0x{b:02X}' for b in src_padded)
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

    with open('patches.csv', 'w', encoding='utf-8') as f:
        f.write('raw_offset,rva,slot,budget,src_len,target_len,fits,src,target,context,note\n')
        for r in rows:
            rva = r['off'] + RVA_DELTA
            def esc(s):
                return '"' + s.replace('"', '""') + '"'
            f.write(
                f'0x{r["off"]:X},0x{rva:X},{r["slot"]},{r["budget"]},'
                f'{r["src_len"]},{r["target_len"]},{r["fits"]},'
                f'{esc(r["src"])},{esc(r["target"])},{esc(r["context"])},{esc(r["note"])}\n'
            )
    print(f"Generated patches.csv with {len(rows)} entries.")


if __name__ == '__main__':
    main()
