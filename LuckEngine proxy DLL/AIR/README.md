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

## Possibles strings AIR-spécifiques restantes

Si tu découvres en jouant des strings encore en EN qui sont propres à
l'univers d'AIR (noms de personnages, BGM, lieux, ou prompts que Kanon
n'avait pas) :

1. Signale-les avec capture d'écran
2. Je scanne AIR.exe pour localiser l'offset
3. On ajoute à `patches.py` et rebuild

La méthode de scan est la même qu'avec Kanon.

## Changelog

- **v0.1** (111 patches) : portage automatique depuis Kanon v0.5
