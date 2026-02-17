# LuckSystem 2.3.2 — Patch Yoremi v4

Fork de [LuckSystem 2.3.2](https://github.com/wetor/LuckSystem) avec corrections pour le support de la traduction de visual novels utilisant le moteur ProtoDB/SiglusEngine (AIR, CLANNAD, Kanon, Summer Pockets, Harmonia, etc.).

## Corrections apportées

### Patch 1 — Import de scripts à longueur variable
**Fichier :** `script/script.go`

L'import de scripts traduits échouait avec un panic quand la traduction avait une longueur différente de l'original. Le code vérifiait strictement `len(paramList) == len(code.Params)`, ce qui bloquait toute traduction plus longue ou plus courte.

- Suppression de la vérification stricte du nombre de paramètres
- Ajout de bounds checking dans la boucle de conversion et le merge des paramètres
- Les offsets de jump (GOTO, IFN, IFY…) sont recalculés automatiquement

### Patch 2 — Correction du pipeline CZ3 (export/import PNG)
**Fichiers :** `czimage/cz3.go`, `czimage/imagefix.go`

L'export et l'import de CZ3 corrompaient silencieusement les données pixels à cause de bugs dans la gestion mémoire et le format d'image.

- **Magic byte** : `Write()` écrasait le magic "CZ3" → "CZ0", rendant le fichier illisible par le jeu
- **Format NRGBA** : Conversion automatique de tout format PNG (RGB, RGBA, paletted) en NRGBA 32-bit avant encodage
- **Buffer aliasing** : `DiffLine()` et `LineDiff()` partageaient des slices au lieu de copier, provoquant une corruption des données delta lors de l'écriture

### Patch 3 — Corruption mémoire dans le décompresseur LZW
**Fichier :** `czimage/lzw.go`

Le décompresseur LZW (`decompressLZW` et `decompressLZW2`) ajoutait des entrées dictionnaire qui référençaient directement le slice `w` au lieu d'en faire une copie. Quand `w` était réassigné, les anciennes entrées du dictionnaire pointaient vers des données corrompues.

- Allocation explicite de `newEntry` avec copie de `w` avant ajout au dictionnaire

### Patch 4 — RawSize incorrect dans la table de blocs CZ
**Fichier :** `czimage/util.go`

Bug critique causant la corruption visuelle des CG en jeu (artefacts colorés dans la moitié inférieure de l'image). Les fonctions `Compress()` et `Compress2()` calculaient un `RawSize` erroné pour chaque bloc LZW, à cause de deux problèmes :

1. **Carry-over non compensé** : Quand la compression LZW d'un bloc s'arrête, le dernier élément en cours (`lastElement`) est reporté au bloc suivant. Mais `count` (bytes lus) incluait ces bytes, alors qu'ils ne seraient décompressés que dans le bloc suivant. Résultat : le premier bloc est trop grand, le dernier trop petit.

2. **Encodage UTF-8 de Go** : En Go, `string(byte(200))` produit une chaîne UTF-8 de 2 bytes (`\xc3\x88`), mais représente un seul byte de données. L'utilisation de `len(last)` pour compter les bytes de données donnait 2 au lieu de 1 pour les octets > 127, causant des erreurs ±1 sur les blocs intermédiaires.

Ce bug affecte tous les jeux utilisant le format CZ3 avec des images multi-blocs. Il correspond à [l'issue #X reportée par masagrator](https://github.com/wetor/LuckSystem/issues/) pour Summer Pockets.

## Installation

```bash
# Remplacer les fichiers dans les dossiers correspondants :
# czimage/cz3.go, czimage/imagefix.go, czimage/lzw.go, czimage/util.go
# pak/pak.go
# script/script.go

# Puis compiler :
go clean
go build -o lucksystem.exe
```

## Jeux testés

- AIR (Steam) — traduction française complète
- Summer Pockets — fix RawSize confirmé (rapport masagrator)

## Crédits

- **wetor** — LuckSystem original
- **masagrator** — identification du bug RawSize (CZ3 layers)
- **Yoremi** — patches 1-4, traduction française d'AIR
