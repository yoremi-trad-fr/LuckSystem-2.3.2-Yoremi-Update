## Key Steam Proxy DLL — Kanon & AIR (FR)

English Version

This project provides an in-memory French translation for hardcoded strings in the Steam versions of Kanon and AIR. It uses a VERSION.dll proxy method, meaning:

    Zero disk modifications: The original .exe remains untouched.

    Steam-Safe: SteamStub DRM and file integrity remain intact.

    Dynamic Patching: Translations are applied in RAM at runtime.

## Translated Content

    Kanon (v0.5): 112 patches covering Options, Save/Load menus, Right-click menus, Global Yes/No buttons, and long tooltips (System tab).

    AIR (v0.1): 111 patches automatically derived from the Kanon table. Matches all UI strings except for Kanon-specific voice samples (e.g., Yuichi's voice icon).

## Installation

    Backup any existing version.dll in your game folder.

    Copy the version.dll corresponding to your game into the root folder (where Kanon.exe or AIR.exe is located).

    Launch the game via Steam.
	
	--------------------------------------------------------------------------------------
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

# AIR FR proxy DLL — v0.1

Traduction française en mémoire des strings hardcodées dans `AIR.exe`
(Steam), via proxy de `VERSION.dll`. Même technique que le patch Kanon :
zéro modification de l'exe sur disque, SteamStub DRM intact, Steam content.

**111 patches** automatiquement dérivés de la table Kanon v0.5 par
correspondance de chaînes EN exactes (Luck Engine a la même base de code
entre Kanon et AIR, seuls les contenus spécifiques au jeu diffèrent).

## Contenu traduit

Toutes les strings du menu Options, du menu Sauvegardes, du menu clic
droit, des boutons globaux Oui/Non, et des tooltips longs — identiques
à Kanon v0.5 sauf pour 1 string :

- `\t(Young)\t` (sample vocal de Yuichi Aizawa dans Kanon) n'existe pas
  dans AIR car les personnages sont différents. Non patché.

## Install

1. Backup toute `version.dll` existante dans le dossier d'AIR.
2. Copier **`version.dll`** dans le dossier du jeu (celui qui contient
   `AIR.exe`).
3. Lancer via **Steam**.

## Activer les logs (optionnel)

Steam → AIR → Propriétés → Options de lancement :
```
KANON_FR_LOG=1 %command%
```

(Le nom de la variable commence par `KANON_FR_LOG` — c'est historique,
le DLL a été prototypé sur Kanon. Le log s'appelle `kanon_fr.log` même
pour AIR. Non gênant en pratique.)

Log attendu :
```
[HH:MM:SS.mmm] DLL_PROCESS_ATTACH (AIR FR proxy v0.1, 111 patches)
[HH:MM:SS.mmm] Loaded real version.dll from C:\WINDOWS\system32\version.dll
[HH:MM:SS.mmm] Sentinel ready, applying 111 patch(es)
[HH:MM:SS.mmm] ... (111 lignes)
[HH:MM:SS.mmm] Patch thread done.
```

## Différences techniques avec Kanon

| Aspect | Kanon | AIR |
|---|---|---|
| Taille exe | 13.02 MB | 12.99 MB |
| Version jeu | 1.5.0.6 | à vérifier en-jeu |
| `.rdata` raw offset | 0x459400 | 0x45D800 |
| `.rdata` virtual addr | 0x45A000 | 0x45F000 |
| Delta raw→RVA | **0xC00** | **0x1800** |
| SteamStub (`.bind`) | oui | oui |
| Triplets JP/EN/ZH | UTF-8 | UTF-8 (identique) |
| Total patches | 112 | 111 |

Le code C du DLL est strictement le même (logique de poll sentinel + patch
runtime via `VirtualProtect`). Seule la table des offsets dans `patches.h`
diffère — auto-générée par `patches.py` à partir de `AIR.exe`.

## Itérer (corrections de traduction)

1. Éditer `patches.py`
2. `python3 patches.py` → régénère `patches.h` + `patches.csv`
3. `make` → recompile `version.dll`


## Changelog

- **v0.1** (111 patches) : portage automatique depuis Kanon v0.5
