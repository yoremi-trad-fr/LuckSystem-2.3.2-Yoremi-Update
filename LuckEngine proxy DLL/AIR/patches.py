#!/usr/bin/env python3
"""
AIR Steam — FR patch table, auto-derived from Kanon's patches.py by offset remap.
Match criteria: byte-exact EN string present exactly once in AIR.exe.
Strings that didn't port (multi-match or missing) are listed below.
"""

# (raw_offset_in_AIR, EN original, FR proposed, context, note)
PATCHES = [
    (0x48C534, b'Close',
                                               'Fermer',
                                               'bottom-right button', 'slot 8, budget 7 -> "Fermer" 6B ok'),
    (0x48DC28, b'Defaults',
                                               'Défauts',
                                               'bottom-left button', '"Défauts" é=2B so 7B, budget 15 ok'),
    (0x49E490, b'Basic',
                                               'Base',
                                               'Basic tab', 'already done via PAK? keep anyway'),
    (0x49E528, b'Shortcut Menu',
                                               'Raccourcis',
                                               'Basic/Shortcut', 'budget 15'),
    (0x49E4CC, b'Hide',
                                               'Cacher',
                                               'Basic/Shortcut value', 'budget 11'),
    (0x49E4D8, b'Display',
                                               'Visible',
                                               'Basic/Shortcut value', 'budget 7 -> "Visible" 7B exact'),
    (0x49E590, b'Skip',
                                               'Passer',
                                               'Basic/Skip', 'budget 7'),
    (0x49E538, b'Previously Read Only',
                                               'Déjà lu seulement',
                                               'Basic/Skip value', '"Déjà lu seulement" = 19B, budget 23 ok'),
    (0x49E5D0, b'Position of Choices',
                                               'Position des choix',
                                               'Basic/Position', '19B, budget 23 ok'),
    (0x49E5B8, b'Bottom',
                                               'Bas',
                                               'Basic/Position value', 'budget 7'),
    (0x49E5C0, b'Center',
                                               'Centre',
                                               'Basic/Position value', 'budget 7'),
    (0x49E6A8, b'Voice',
                                               'Voix',
                                               'Basic/Voice', 'budget 7'),
    (0x49E658, b'Stop on New Message',
                                               'Arrêt au message',
                                               'Basic/Voice value', '"Arrêt au message" = 17B (ê=2), budget 23'),
    (0x49E670, b'      No Stops      ',
                                               '     Sans arrêt     ',
                                               'Basic/Voice value', 'keep spaces for alignment, FR 21B, budget 23'),
    (0x49E720, b'Display Date',
                                               'Afficher date',
                                               'Basic/Date', '13B, budget 15'),
    (0x49E7D8, b'Initial Cursor Position',
                                               'Position du curseur',
                                               'Basic/Cursor', '20B, budget 23'),
    (0x49E730, b'Positioned at \xe2\x9d\x9dYes\xe2\x9d\x9e',
                                               'Sur « Oui »',
                                               'Basic/Cursor value', 'Keep stylish chevrons if fits; budget 23'),
    (0x49E7A8, b'Positioned at \xe2\x9d\x9dNo\xe2\x9d\x9e',
                                               'Sur « Non »',
                                               'Basic/Cursor value', 'budget 23'),
    (0x49E818, b'Controller Rumble Function',
                                               'Vibration manette',
                                               'Basic/Rumble', '17B, budget 27'),
    (0x49E7F8, b'Disable',
                                               'Aucun',
                                               'Basic/Rumble value', '"Aucun" 5B, budget 7 (rumble intensity: Aucun/Min/Mid/Max)'),
    (0x49F840, b'Text1',
                                               'Texte1',
                                               'Text1 tab', '6B, budget 7'),
    (0x49F848, b'Language',
                                               'Langue',
                                               'Text1/Language', '6B, budget 11'),
    (0x49F908, b'Font',
                                               'Police',
                                               'Text1/Font', '6B, budget 7'),
    (0x49F8DC, b'Mincho',
                                               'Mincho',
                                               'Text1/Font value', 'japanese font family name, keep'),
    (0x49F8F8, b'Modern',
                                               'Moderne',
                                               'Text1/Font value', '7B, budget 7 tight!'),
    (0x49F944, b'Solid',
                                               'Opaque',
                                               'Text1/Window Transp value', '6B, budget 11'),
    (0x49F96C, b'Clear',
                                               'Clair',
                                               'Text1/Window Transp value', '5B, budget 11'),
    (0x49F988, b'Window Transparency',
                                               'Transparence',
                                               'Text1/Window Transp', '12B, budget 23'),
    (0x49F9A0, b'Only Choices',
                                               'Choix seulement',
                                               'Text1/Transp target', '15B, budget 15 exact'),
    (0x49F9C8, b'Previously Read Text',
                                               'Texte déjà lu',
                                               'Text1/Read target', '14B (é=2), budget 23'),
    (0x49F9F8, b'Color of',
                                               'Couleur',
                                               'Text1/Color label', '7B, budget 15'),
    (0x49FA40, b'Green',
                                               'Vert',
                                               'Text1/Color value', '4B, budget 7'),
    (0x49FA48, b'Blue',
                                               'Bleu',
                                               'Text1/Color value', '4B, budget 7'),
    (0x49FA60, b'Purple',
                                               'Violet',
                                               'Text1/Color value', '6B, budget 7'),
    (0x49F804, b'Orange',
                                               'Orange',
                                               'Text1/Color value', 'same word in FR'),
    (0x49F830, b'Yellow',
                                               'Jaune',
                                               'Text1/Color value', '5B, budget 7'),
    (0x49FA70, b'Read Text Color',
                                               'Coul. texte lu',
                                               'Text1/Read color', '14B, budget 15'),
    (0x49FE2C, b'Text2',
                                               'Texte2',
                                               'Text2 tab', '6B, budget 11'),
    (0x49FEA0, b'Text Speed',
                                               'Vitesse texte',
                                               'Text2/Speed', '13B, budget 15'),
    (0x49FE5C, b'Slow',
                                               'Lent',
                                               'Text2/Speed value', '4B, budget 7'),
    (0x49FE64, b'Fast',
                                               'Rapide',
                                               'Text2/Speed value', '6B, budget 7'),
    (0x49FEB0, b'0 sec/char',
                                               '0 s/car',
                                               'Text2/Speed value', '7B, budget 15'),
    (0x49FF00, b'0.1 sec/char',
                                               '0,1 s/car',
                                               'Text2/Speed value', '9B, budget 15 (comma for FR decimal)'),
    (0x49FF20, b'Wait Time Per Character',
                                               'Attente par caractère',
                                               'Text2/Wait', '22B (è=2), budget 23'),
    (0x49FF38, b'0 sec',
                                               '0 s',
                                               'Text2/Wait value', '3B, budget 7'),
    (0x49FF48, b'1 sec',
                                               '1 s',
                                               'Text2/Wait value', '3B, budget 7'),
    (0x49FF58, b'2 sec',
                                               '2 s',
                                               'Text2/Wait value', '3B, budget 7'),
    (0x49FF68, b'3 sec',
                                               '3 s',
                                               'Text2/Wait value', '3B, budget 7'),
    (0x49FF78, b'Base Wait Time',
                                               'Attente de base',
                                               'Text2/Base', '15B, budget 15 exact'),
    (0x49FFCC, b'Sound',
                                               'Son',
                                               'Sound tab', '3B, budget 11'),
    (0x49FFD8, b'Master Volume',
                                               'Volume maître',
                                               'Sound/Master', '"Volume maître" 14B (î=2), budget 15'),
    (0x4A0078, b'System Sounds',
                                               'Sons système',
                                               'Sound/System', '13B (è=2), budget 15'),
    (0x4A0110, b'While Pressed',
                                               'Maintenu',
                                               'Keyboard/mode', '8B, budget 15'),
    (0x4A0120, b'Start/Stop',
                                               'Marche/Arrêt',
                                               'Keyboard/mode', '13B (ê=2), budget 15'),
    (0x4A0178, b'C (Skip)',
                                               'C (Passer)',
                                               'Keyboard/key label', '10B, budget 15'),
    (0x4A0188, b'Z (Rewind)',
                                               'Z (Retour)',
                                               'Keyboard/key label', '10B, budget 15'),
    (0x4A01E0, b'  Disable  ',
                                               '  Désactivé ',
                                               'Keyboard/value', 'retain spaces; 13B (é=2), budget 15'),
    (0x4A01F0, b'Quick Save',
                                               'Sauv. rapide',
                                               'Keyboard/key label', '12B, budget 15'),
    (0x4A0240, b'Quick Load',
                                               'Chrg. rapide',
                                               'Keyboard/key label', '12B, budget 15'),
    (0x4A02A0, b'Switch Language',
                                               'Changer langue',
                                               'Keyboard/key label', '14B, budget 15'),
    (0x4A0300, b'Up Arrow',
                                               'Flèche haut',
                                               'Keyboard/key label', '12B (è=2), budget 15'),
    (0x4A034C, b'Mouse',
                                               'Souris',
                                               'Mouse tab', '6B, budget 11'),
    (0x4A03A0, b'System Menu',
                                               'Menu système',
                                               'Mouse/target', '13B (è=2), budget 15'),
    (0x4A03B0, b'Hide Window',
                                               'Cacher fenêtre',
                                               'Mouse/target', '15B (ê=2), budget 15 exact'),
    (0x4A0408, b'Right Click',
                                               'Clic droit',
                                               'Mouse/binding', '10B, budget 15'),
    (0x4A0428, b'Left+Right Click',
                                               'Clic gauche+droit',
                                               'Mouse/binding', '18B, budget 23'),
    (0x4A0440, b'Mouse Wheel Button',
                                               'Bouton molette',
                                               'Mouse/binding', '14B, budget 23'),
    (0x4A04A0, b'Rewind Once',
                                               'Retour un pas',
                                               'Mouse/target', '13B, budget 15'),
    (0x4A04B0, b'Wheel Up',
                                               'Molette haut',
                                               'Mouse/binding', '12B, budget 15'),
    (0x4A0510, b'Forward Once',
                                               'Avance un pas',
                                               'Mouse/target', '14B, budget 15'),
    (0x4A0520, b'Wheel Down',
                                               'Molette bas',
                                               'Mouse/binding', '11B, budget 15'),
    (0x4A05A0, b'Jump and Switch Pages',
                                               'Saut et pages',
                                               'Mouse/target', '14B, budget 23'),
    (0x4A05B8, b'Return/Proceed Button',
                                               'Retour/Continuer',
                                               'Mouse/target', '16B, budget 23'),
    (0x4A0614, b'Enable',
                                               'Activé',
                                               'Mouse/Gestures value', '7B (é=2), budget 11'),
    (0x4A0630, b'Gestures',
                                               'Gestes',
                                               'Mouse/Gestures', '6B, budget 15'),
    (0x4A0670, b'Dialog and Choices',
                                               'Dialogue et choix',
                                               'Mouse/Snap target', '18B, budget 23'),
    (0x4A06F8, b'Snap Pointer',
                                               'Aimanter curseur',
                                               'Mouse/Snap', '16B, budget 23'),
    (0x4A0A68, b'Game Ver. ',
                                               'Version jeu ',
                                               'System/info', '12B, budget 15'),
    (0x4A0A8C, b'System',
                                               'Système',
                                               'System tab', '8B (è=2), budget 11'),
    (0x4A0AD0, b'Window',
                                               'Fenetre',
                                               'System/Window', '"Fenetre" no accent, 7B, budget 7 exact'),
    (0x4A0AD8, b'Full Screen',
                                               'Plein écran',
                                               'System/FullScreen', '12B (é=2), budget 15'),
    (0x4A0B00, b'Screen Mode',
                                               'Mode écran',
                                               'System/ScreenMode', '11B (é=2), budget 15'),
    (0x4A0B40, b'Auto',
                                               'Auto',
                                               'System/value', 'same in FR'),
    (0x4A0B90, b'Window Size',
                                               'Taille fenêtre',
                                               'System/WindowSize', '15B (ê=2), budget 15 exact'),
    (0x49FC40, b'Wait Time Per Character: In Auto Mode, this sets the wait time until the next message is displayed based on the number of characters in text.',
                                               'Attente par caractère : en mode Auto, règle le temps avant le message suivant selon le nombre de caractères du texte.',
                                               'Texte2 tooltip', 'budget 143'),
    (0x49FCD0, b'Base Wait Time: You can set a Base Wait Time to add to\n\xe2\x9d\x9dWait Time Per Character\xe2\x9d\x9e.',
                                               'Attente de base : définit un temps fixe à ajouter à\n❝Attente par caractère❞.',
                                               'Texte2 tooltip', 'budget 95, keep stylized quotes'),
    (0x49EA80, b'Initial Cursor Position: Sets the initial position of the cursor when a Yes/No choice is available.',
                                               'Position du curseur : définit la position initiale du curseur pour un choix Oui/Non.',
                                               'Basic tooltip', 'budget 111'),
    (0x49ECF0, b'Controller Rumble Function: Plug in a controller to use the controller Rumble function.',
                                               'Vibration manette : branchez une manette pour utiliser la vibration.',
                                               'Basic tooltip', 'budget 87'),
    (0x4A0710, b'Gestures: Moving the cursor while holding the left button works the same way as Touch controls.',
                                               'Gestes : maintenir le bouton gauche en déplaçant le curseur agit comme en mode tactile.',
                                               'Mouse tooltip', 'budget 95'),
    (0x4A09C0, b'Left+Right Click: Hold left button then right click to switch between languages (English/Simplified Chinese/Japanese).',
                                               'Clic gauche+droit : maintenir le clic gauche puis clic droit pour changer de langue (anglais/chinois/japonais).',
                                               'Mouse tooltip', 'budget 119'),
    (0x49E9D0, b'Voice: If you select No Stops, sound will continue to play even if you advance the text during voice playback. (It stops if there is sound on the next message.)',
                                               "Voix : avec Sans arrêt, le son continue même si vous avancez le texte durant la lecture. (S'arrête si le message suivant a du son.)",
                                               'Basic tooltip', 'budget 175'),
    (0x48CF40, b'Do you wish to load this save?',
                                               'Charger cette sauvegarde ?',
                                               'Save/Load prompt', 'budget 31'),
    (0x48D158, b'Are you sure you wish to overwrite data?',
                                               'Écraser la sauvegarde ?',
                                               'Save/Load prompt', 'budget 47'),
    (0x48D110, b'Do you wish to save?',
                                               'Sauvegarder ?',
                                               'Save prompt', 'budget 23, bonus'),
    (0x48D380, b'Are you sure you wish to delete this save data?',
                                               'Supprimer cette sauvegarde ?',
                                               'Save delete prompt', 'budget 47'),
    (0x48D3E0, b'This cannot be deleted.',
                                               'Suppression impossible.',
                                               'Save delete error', 'budget 23'),
    (0x48DBBC, b'Delete',
                                               'Suppr.',
                                               'Save menu button', 'budget 11'),
    (0x48DBE4, b'Latest',
                                               'Récent',
                                               'Save menu button', 'budget 11 (é=2)'),
    (0x48DF78, b'Text preview.',
                                               'Aperçu.',
                                               'Text preview label', 'budget 15'),
    (0x48E098, b'Settings such as text speed are reflected.',
                                               'Les réglages comme la vitesse sont appliqués.',
                                               'Text preview tooltip', 'budget 47'),
    (0x48C4E8, b'English',
                                               'FR',
                                               'Language switcher', 'budget 7 hard cap (neighbor is 简体中文). Using ISO code "FR" (2B) for cleanness. Alternatives in 7B: Franç./Franc./Francai'),
    (0x49DA78, b'Off',
                                               'Non',
                                               'Read text color toggle', 'budget 7, "Non" 3B'),
    (0x49E6E0, b'On',
                                               'Oui',
                                               'Read text color toggle', 'budget 7, "Oui" 3B'),
    (0x48D488, b'Return to the title screen?',
                                               'Retour au titre ?',
                                               'Title return prompt', 'budget 31'),
    (0x48D430, b'Return to the menu?',
                                               'Retour au menu ?',
                                               'Menu return prompt', 'budget 23, bonus'),
    (0x49AE80, b'$A1There is unsaved data.\n$A1Are you sure you wish to quit the game?',
                                               '$A1Données non sauvegardées.\n$A1Quitter le jeu ?',
                                               'Quit prompt (unsaved)', 'budget 79, keep $A1 tags'),
    (0x48C508, b'Yes',
                                               'Oui',
                                               'Global dialog button', 'budget 3, exact fit. Used by ALL Yes/No prompts'),
    (0x48C51C, b'No',
                                               'Non',
                                               'Global dialog button', 'budget 3, exact fit. Used by ALL Yes/No prompts'),
    (0x48D348, b'Save completed.',
                                               'Sauvegardé.',
                                               'Save confirmation', '"Sauvegardé." 12B (é=2), budget 15'),
    (0x49AEC8, b'Are you sure you wish to quit the game?',
                                               'Quitter le jeu ?',
                                               'Quit prompt (saved)', 'budget 39 — the short variant, used when no unsaved data'),
    (0x4A0BA0, b'Changing the setting to \xe2\x9d\x9dAuto\xe2\x9d\x9e will open the window at a scale based on the Windows \xe2\x9d\x9dDisplay\xe2\x9d\x9e setting.\nChanging \xe2\x9d\x9d%%\xe2\x9d\x9e will scale the display, with the default resolution being %d\xc3\x97%d pixels.\n (Scale cannot be increased beyond the maximum resolution of your display.)',
                                               "Le réglage ❝Auto❞ ouvre la fenêtre à une échelle basée sur les paramètres ❝Affichage❞ de Windows.\nLe réglage ❝%%❞ ajuste l'échelle, la résolution par défaut étant %d×%d pixels.\n (L'échelle ne peut dépasser la résolution maximale de votre écran.)",
                                               'System tab tooltip', 'budget 287, keep %%/%d×%d intact'),
]


def main():
    import sys, os
    here = os.path.dirname(os.path.abspath(__file__))
    data = open(os.path.join(here,'AIR.exe'),'rb').read()

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
        rows.append({'off':off,'en':en.decode('utf-8','replace'),'fr':fr,
                     'en_len':len(en),'fr_len':len(fr_bytes),'slot':slot,
                     'budget':budget,'fits':fits,'context':context,'note':note,
                     'en_bytes':en,'fr_bytes':fr_bytes})

    if errors:
        print("OFFSET MISMATCH:", file=sys.stderr)
        for e in errors: print("  "+e, file=sys.stderr)
        sys.exit(1)

    print(f"{'off':>8}  {'slot':>4}  {'EN':>3}  {'FR':>3}  {'fit':3}  EN -> FR")
    print('-'*100)
    n_ok=n_bad=0
    for r in rows:
        mark='✓' if r['fits'] else '✗'
        if r['fits']: n_ok+=1
        else: n_bad+=1
        print(f"0x{r['off']:06X}  {r['slot']:>4}  {r['en_len']:>3}  {r['fr_len']:>3}  {mark}   {r['en']!r} -> {r['fr']!r}")
    print(f"\nTotal: {len(rows)}  OK: {n_ok}  FAIL: {n_bad}")
    if n_bad: sys.exit(2)

    with open('patches.h','w',encoding='utf-8') as f:
        f.write('/* Auto-generated from patches_air.py. Do not edit. */\n')
        f.write('#ifndef KANON_FR_PATCHES_H\n#define KANON_FR_PATCHES_H\n\n')
        for i,r in enumerate(rows):
            write_len = max(len(r['en_bytes']), len(r['fr_bytes'])) + 1
            en_padded = list(r['en_bytes']) + [0]*(write_len-len(r['en_bytes']))
            fr_padded = list(r['fr_bytes']) + [0]*(write_len-len(r['fr_bytes']))
            en_arr=','.join(f'0x{b:02X}' for b in en_padded)
            fr_arr=','.join(f'0x{b:02X}' for b in fr_padded)
            f.write(f'static const BYTE s_en_{i:03d}[] = {{ {en_arr} }};\n')
            f.write(f'static const BYTE s_fr_{i:03d}[] = {{ {fr_arr} }};\n')
        f.write('\nstatic const KanonPatch g_patches[] = {\n')
        for i,r in enumerate(rows):
            # AIR raw .rdata starts at 0x45D800, vaddr 0x45F000 -> RVA = raw + 0x1800
            rva = r['off'] + 0x1800
            write_len = max(len(r['en_bytes']), len(r['fr_bytes'])) + 1
            ctx = (r['context']+': '+r['en'][:30]).replace('\\','\\\\').replace('"','\\"').replace('\n','\\n').replace('\r','\\r').replace('\t','\\t')
            f.write(f'    {{ 0x{rva:06X}, {write_len:>4}, s_en_{i:03d}, s_fr_{i:03d}, "{ctx}" }},\n')
        f.write('};\n\n#define N_PATCHES (sizeof(g_patches)/sizeof(g_patches[0]))\n\n#endif\n')
    print(f"\nGenerated patches.h with {len(rows)} entries.")

    with open('patches.csv','w',encoding='utf-8') as f:
        f.write('raw_offset,rva,slot,budget,en_len,fr_len,fits,en,fr,context,note\n')
        def esc(s): return '"'+s.replace('"','""')+'"'
        for r in rows:
            rva=r['off']+0x1800
            f.write(f'0x{r["off"]:X},0x{rva:X},{r["slot"]},{r["budget"]},{r["en_len"]},{r["fr_len"]},{r["fits"]},{esc(r["en"])},{esc(r["fr"])},{esc(r["context"])},{esc(r["note"])}\n')
    print(f"Generated patches.csv with {len(rows)} entries.")


if __name__ == '__main__':
    main()
