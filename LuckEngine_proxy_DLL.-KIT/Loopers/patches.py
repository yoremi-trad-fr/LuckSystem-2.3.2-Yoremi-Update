#!/usr/bin/env python3
"""
LOOPERS (Luck Engine) - proxy DLL patch table.

This table was ported from the HarmoniaHD proxy table by matching the same
hardcoded Luck Engine UI strings inside LOOPERS.exe, then adjusting the few
LOOPERS-specific strings/case variants.

Current targets are French replacements for the hardcoded English strings.
Entries must fit inside the original string slot budget.

Generated files:
  - patches.h
  - patches.csv

Notes:
  - LOOPERS.exe imports VERSION.dll, so the shared version.dll proxy can load.
  - .rdata raw -> RVA delta is 0x1400 for the current Steam executable.
  - Ambiguous or unrelated strings are documented in port_report.md.
"""

from pathlib import Path

GAME_EXE = r'C:\Program Files (x86)\Steam\steamapps\common\LOOPERS\LOOPERS.exe'
RVA_DELTA = 0x1400
PATCH_GAME_NAME = 'Loopers'
PATCH_VERSION = '0.1'

# Each entry: (raw_offset, src_bytes, target_str, context, note)
# Ready-to-use state: target_str contains the French replacement.

PATCHES = [
    (0x40F9B8, b'English', 'FR', 'Language switcher', 'budget 7 hard cap'),
    (0x40F9D8, b'Yes', 'Oui', 'Global dialog button', 'budget 3 exact fit'),
    (0x40F9EC, b'No', 'Non', 'Global dialog button', 'budget 3 exact fit'),
    (0x40F9F4, b'Close', 'Fermer', 'Global close button', 'budget 7'),

    (0x4103E0, b'Do you wish to load this save?', 'Charger cette sauvegarde ?', 'Save/Load prompt', 'budget 31'),
    (0x4105B0, b'Do you wish to save?', 'Sauvegarder ?', 'Save prompt', 'budget 23'),
    (0x4105F8, b'Are you sure you wish to overwrite data?', 'Ecraser cette sauvegarde ?', 'Save/Load prompt', 'budget 47'),
    (0x4107E8, b'Save completed.', 'Sauvegarde OK', 'Save confirmation', 'budget 15'),
    (0x410820, b'Are you sure you wish to delete this save data?', 'Supprimer cette sauvegarde ?', 'Save delete prompt', 'budget 47'),
    (0x410880, b'This cannot be deleted.', 'Suppression impossible.', 'Save delete error', 'budget 23'),
    (0x4108D0, b'Return to the menu?', 'Retour au menu ?', 'Menu return prompt', 'budget 23'),
    (0x410928, b'Return to the title screen?', "Retour a l'ecran titre ?", 'Title return prompt', 'budget 31'),
    (0x411058, b'Delete', 'Suppr.', 'Save menu button', 'budget 7'),
    (0x411074, b'Latest', 'Dernier', 'Save menu button', 'budget 11'),
    (0x4110B8, b'Defaults', 'Defauts', 'Defaults button', 'budget 15'),
    (0x4113F0, b'Text preview.', 'Apercu texte.', 'Text preview label', 'budget 15'),
    (0x411508, b'Settings such as text speed are reflected.', 'Les reglages de texte sont appliques.', 'Text preview tooltip', 'budget 47'),

    (0x4386A0, b'$A1There is unsaved data.\n$A1Are you sure you wish to quit the game?', '$A1Donnees non sauvegardees.\n$A1Quitter le jeu ?', 'Quit prompt (unsaved)', 'budget 71, keep $A1 tags'),
    (0x438748, b'Are you sure you wish to quit the game?', 'Quitter le jeu ?', 'Quit prompt (saved)', 'budget 39'),

    (0x43EC60, b'Game Ver. ', 'Ver. jeu ', 'System/info', 'budget 15'),
    (0x43EC88, b'Previously Read Only', 'Deja lu seulement', 'Basic/Skip value', 'budget 23'),
    (0x43ECE0, b'Skip', 'Saut', 'Basic/Skip', 'budget 7'),
    (0x43ED58, b'Stop on New Message', 'Stop au nouveau msg', 'Basic/Voice value', 'budget 23'),
    (0x43ED70, b'      No Stops      ', '     Sans arret     ', 'Basic/Voice value', 'keep spaces for alignment, budget 23'),
    (0x43EDA8, b'Voice', 'Voix', 'Basic/Voice', 'budget 7'),
    (0x43EDD0, b'Positioned at \xe2\x9d\x9dYes\xe2\x9d\x9e', 'Place sur ❝Oui❞', 'Basic/Cursor value', 'budget 23'),
    (0x43EE48, b'Positioned at \xe2\x9d\x9dNo\xe2\x9d\x9e', 'Place sur ❝Non❞', 'Basic/Cursor value', 'budget 23'),
    (0x43EE78, b'Initial cursor position', 'Position initiale', 'Basic/Cursor', 'LOOPERS uses lowercase cursor, budget 23'),
    (0x43EE98, b'Disable', 'Arret', 'Basic/Rumble value', 'budget 7'),
    (0x43EEB8, b'Controller Rumble Function', 'Vibration manette', 'Basic/Rumble', 'budget 27'),

    (0x43F070, b'Voice: If you select No Stops, sound will continue to play even if you advance the text during voice playback.(It stops if there is sound on the next message.)', "Avec Sans arret, la voix continue si vous avancez le texte pendant sa lecture. Elle s'arrete si le message suivant a un son.", 'Basic tooltip', 'LOOPERS has no space before parenthesis, budget 159'),
    (0x43F110, b'Initial cursor position: Sets the initial position of the cursor when a Yes/No choice is available.', "Definit la position initiale du curseur lors d'un choix Oui/Non.", 'Basic tooltip', 'LOOPERS uses lowercase cursor, budget 111'),
    (0x43F370, b'Controller Rumble Function: Plug in a controller to use the controller Rumble function.', 'Branchez une manette pour utiliser la vibration.', 'Basic tooltip', 'budget 87'),

    (0x43F3C8, b'Language', 'Langue', 'Text1/Language', 'budget 11'),
    (0x43F45C, b'Mincho', 'Mincho', 'Text1/Font value', 'font family name, keep'),
    (0x43F478, b'Modern', 'Moderne', 'Text1/Font value', 'budget 7 tight'),
    (0x43F488, b'Font', 'Police', 'Text1/Font', 'budget 7'),
    (0x43F4D8, b'Green', 'Vert', 'Text1/Color value', 'budget 7'),
    (0x43F530, b'Window Frame Color', 'Couleur cadre', 'Text1/Window frame color', 'LOOPERS-specific option, budget 23'),
    (0x43F55C, b'Solid', 'Opaque', 'Text1/Window transparency value', 'budget 11'),
    (0x43F584, b'Clear', 'Transp', 'Text1/Window transparency value', 'budget 11'),
    (0x43F5A0, b'Window Transparency', 'Transp. fenetre', 'Text1/Window transparency', 'budget 23'),
    (0x43F5E0, b'Enable', 'Actif', 'Text1/Read text value', 'budget 7'),
    (0x43F600, b'Previously Read Text', 'Texte deja lu', 'Text1/Read target', 'budget 23'),
    (0x43F638, b'Color of', 'Couleur', 'Text1/Color label', 'budget 11'),
    (0x43F644, b'Yellow', 'Jaune', 'Text1/Color value', 'budget 7'),
    (0x43F658, b'Blue', 'Bleu', 'Text1/Color value', 'budget 7'),
    (0x43F670, b'Purple', 'Violet', 'Text1/Color value', 'budget 7'),
    (0x43F680, b'Read Text Color', 'Coul. texte lu', 'Text1/Read color', 'shortened for LOOPERS budget 15'),

    (0x43F850, b'Wait Time Per Character: In Auto Mode, this sets the wait time until the next message is displayed based on the number of characters in text.', 'En mode auto, definit le delai avant le message suivant selon le nombre de caracteres du texte.', 'Text2 tooltip', 'budget 143'),
    (0x43F8E0, b'Base Wait Time: You can set a Base Wait Time to add to\n\xe2\x9d\x9dWait Time Per Character\xe2\x9d\x9e.', "Ajoute un delai fixe a l'attente par caractere.", 'Text2 tooltip', 'budget 95'),
    (0x43FA30, b'Slow', 'Lent', 'Text2/Speed value', 'budget 7'),
    (0x43FA38, b'Fast', 'Rapide', 'Text2/Speed value', 'budget 7'),
    (0x43FA70, b'Text Speed', 'Vitesse texte', 'Text2/Speed', 'budget 15'),
    (0x43FA80, b'0 sec/char', '0 s/car.', 'Text2/Speed value', 'budget 15'),
    (0x43FAD0, b'0.1 sec/char', '0.1 s/car.', 'Text2/Speed value', 'budget 15'),
    (0x43FAF0, b'Wait Time Per Character', 'Attente par caractere', 'Text2/Wait', 'budget 23'),
    (0x43FB08, b'0 sec', '0 s', 'Text2/Wait value', 'budget 7'),
    (0x43FB18, b'1 sec', '1 s', 'Text2/Wait value', 'budget 7'),
    (0x43FB28, b'2 sec', '2 s', 'Text2/Wait value', 'budget 7'),
    (0x43FB38, b'3 sec', '3 s', 'Text2/Wait value', 'budget 7'),
    (0x43FB48, b'Base Wait Time', 'Attente base', 'Text2/Base', 'budget 15'),

    (0x43FBB0, b'Master Volume', 'Volume global', 'Sound/Master', 'budget 15'),
    (0x43FC00, b'System Sounds', 'Sons systeme', 'Sound/System', 'budget 15'),

    (0x43FC40, b'While Pressed', 'Maintenu', 'Keyboard/mode', 'budget 15'),
    (0x43FC98, b'Start/Stop', 'Marche/Arret', 'Keyboard/mode', 'budget 15'),
    (0x43FCA8, b'C (Skip)', 'C (Saut)', 'Keyboard/key label', 'budget 15'),
    (0x43FD08, b'Z (Rewind)', 'Z (Retour)', 'Keyboard/key label', 'budget 15'),
    (0x43FD18, b'  Disable  ', '  Inactif  ', 'Keyboard/value', 'retain spaces, budget 15'),
    (0x43FD70, b'Quick Save', 'Sauv. rap.', 'Keyboard/key label', 'budget 11'),
    (0x43FDB8, b'Quick Load', 'Charg. rapide', 'Keyboard/key label', 'budget 15'),
    (0x43FDC8, b'Switch Language', 'Changer langue', 'Keyboard/key label', 'budget 15'),
    (0x43FE30, b'Up Arrow', 'Fleche haut', 'Keyboard/key label', 'budget 15'),

    (0x43FE90, b'System Menu', 'Menu systeme', 'Mouse/target', 'budget 15'),
    (0x43FEA0, b'Hide Window', 'Cacher fenetre', 'Mouse/target', 'budget 15'),
    (0x43FEF8, b'Right Click', 'Clic droit', 'Mouse/binding', 'budget 15'),
    (0x43FF18, b'Left+Right Click', 'Clic gauche+droit', 'Mouse/binding', 'budget 23'),
    (0x43FFC0, b'Mouse Wheel Button', 'Bouton molette', 'Mouse/binding', 'budget 23'),
    (0x440065, b'Rewind Once', 'Retour x1', 'Mouse/target', 'inside bracketed [380] label, prefix kept'),
    (0x440078, b'Wheel up', 'Molette haut', 'Mouse/binding', 'LOOPERS uses lowercase up, budget 15'),
    (0x4400E5, b'Forward Once', 'Avance x1', 'Mouse/target', 'inside bracketed [380] label, prefix kept'),
    (0x4400F8, b'Wheel down', 'Molette bas', 'Mouse/binding', 'LOOPERS uses lowercase down, budget 15'),
    (0x440178, b'Jump and switch pages', 'Saut/changer page', 'Mouse/target', 'LOOPERS uses lowercase switch/pages, budget 23'),
    (0x440190, b'Return/Proceed Button', 'Bouton retour/avance', 'Mouse/target', 'budget 23'),
    (0x4401F8, b'Gestures', 'Gestes', 'Mouse/Gestures', 'budget 15'),
    (0x440238, b'Dialogue Only', 'Dialogue seul', 'Mouse/Gestures value', 'LOOPERS-specific value, budget 15'),
    (0x4402A8, b'Snap pointer', 'Aimant curseur', 'Mouse/Snap', 'LOOPERS uses lowercase pointer, budget 23'),
    (0x4402C0, b'Gestures: Moving the cursor while holding the left button works the same way as Touch controls.', 'Maintenir le bouton gauche et bouger le curseur agit comme les commandes tactiles.', 'Mouse tooltip', 'budget 95'),
    (0x440570, b'Left+Right Click: Hold left button then right click to switch between languages (English/Simplified Chinese/Japanese).', 'Maintenez le bouton gauche puis clic droit pour changer de langue.', 'Mouse tooltip', 'budget 119'),

    (0x4410B4, b'Window', 'Fenetre', 'System/Window', 'budget 11'),
    (0x441108, b'Full Screen', 'Plein ecran', 'System/FullScreen', 'budget 15'),
    (0x441128, b'Screen Mode', 'Mode ecran', 'System/ScreenMode', 'budget 11'),
    (0x441134, b'Auto', 'Auto', 'System/value', 'same word'),
    (0x441180, b'Window Size', 'Taille fen.', 'System/WindowSize', 'budget 11'),
    (0x441420, 'Changing the setting to ❝Auto❞ will open the window at a scale based on the Windows ❝Display❞ setting.\nChanging ❝%%❞ will scale the display, with the default resolution being %d×%d pixels.\n (Scale cannot be increased beyond the maximum resolution of your display.)'.encode('utf-8'), "Le mode ❝Auto❞ ouvre la fenetre selon le reglage Windows ❝Affichage❞.\nModifier ❝%%❞ redimensionne l'affichage, resolution par defaut : %d×%d pixels.\n(Le zoom ne peut pas depasser la resolution maximale de l'ecran.)", 'System tab tooltip', 'budget 279, keep %%/%d×%d intact'),
    (0x441538, b'Low resolution', 'Basse res.', 'System/Movie quality value', 'LOOPERS-specific option, budget 15'),
    (0x441588, b'Movie Quality', 'Qualite video', 'System/Movie quality', 'LOOPERS-specific option, budget 15'),
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
