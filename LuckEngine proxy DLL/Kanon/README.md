# Kanon FR proxy DLL — v0.5

Traduction française en mémoire des strings hardcodées dans `Kanon.exe`
(Steam), via proxy de `VERSION.dll`. Zéro modification de l'exe sur disque
→ SteamStub DRM intact → Steam content.

## Nouveautés v0.5 (3 nouvelles strings, 112 au total)

- `Save completed.` → `Sauvegardé.`
- `Are you sure you wish to quit the game?` (variante courte, quand tout est sauvegardé) → `Quitter le jeu ?`
- Tooltip de l'onglet **Système** (bloc de 277 octets sur "Mode écran / Taille fenêtre") → traduction complète :
  > « Le réglage ❝Auto❞ ouvre la fenêtre à une échelle basée sur les paramètres ❝Affichage❞ de Windows.
  > Le réglage ❝%%❞ ajuste l'échelle, la résolution par défaut étant %d×%d pixels.
  > (L'échelle ne peut dépasser la résolution maximale de votre écran.) »

Les tokens `%%` (affichage du caractère `%`) et `%d×%d` (paramètres de
formatage runtime) sont préservés à l'identique.

## Install

1. Backup l'ancienne `version.dll` si tu en as une.
2. Copier **la nouvelle `version.dll`** dans le dossier du jeu
   (celui qui contient `Kanon.exe`).
3. Lancer via **Steam**.

## Activer les logs (si souci)

Steam → Kanon → Propriétés → Options de lancement :
```
KANON_FR_LOG=1 %command%
```

Log v0.5 normal :
```
[HH:MM:SS.mmm] DLL_PROCESS_ATTACH (Kanon FR proxy v0.5, 112 patches)
[HH:MM:SS.mmm] Loaded real version.dll from C:\WINDOWS\system32\version.dll
[HH:MM:SS.mmm] Sentinel ready, applying 112 patch(es)
[HH:MM:SS.mmm] ... (112 lignes de patches)
[HH:MM:SS.mmm] Patch thread done.
```

## Strings intentionnellement non traduites

| Offset raw | EN | Pourquoi |
|---|---|---|
| 0x499204 | `All` | slot 4B, pas d'équivalent FR en 3B |
| 0x49A6FC | `Red` | slot 4B, "Rouge" = 5B |
| 0x499478/80/88 | `Min`/`Mid`/`Max` | universels, slot 4B |
| 0x49ACC0 | `SFX` | terme universel |
| 0x49AD58 | `Yuichi Aizawa` | nom propre |
| 0x49AD88 | `PARTS/voice_icon` | asset path |
| 0x49ADA0 | `keyboard` | preset interne non user-facing |
| 0x49B758 | `Steam Deck: Proton ` | étiquette technique invariante |

## Itérer

1. Éditer `patches.py` (seul fichier à modifier pour les textes)
2. `python3 patches.py` → régénère `patches.h` + `patches.csv`
3. `make` → recompile `version.dll`

## Changelog

- **v0.5** (112 patches) : tooltip System tab (Mode écran), Save completed, 2ᵉ variante Quit prompt
- **v0.4** (109 patches) : boutons globaux `Yes`/`No` → `Oui`/`Non`
- **v0.3** (108 patches) : menu Sauvegardes, menu clic droit, Aperçu Texte1, toggle Oui/Non du Read Text Color, switcher langue → FR
- **v0.2** (91 patches)  : tous les onglets Options complets
- **v0.1** (1 patch)     : POC `Slow` → `Lent`
