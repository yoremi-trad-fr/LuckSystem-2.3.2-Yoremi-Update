# Analyse technique — LuckSystem Yoremi update

Document technique détaillant les 4 corrections appliquées à LuckSystem 2.3.2 pour le support des traductions de visual novels.

---

## Patch 1 — Import de scripts à longueur variable

### Fichier modifié
`script/script.go` — lignes 172-243

### Problème
La fonction `VMRun()` en mode import vérifiait strictement que le nombre de paramètres importés (`code.Params`) correspondait au nombre attendu (`expectedExportCount`). Cette vérification échouait systématiquement avec des traductions de longueur différente de l'original, car les `StringParam` de taille variable n'étaient pas correctement comptabilisés.

```go
// AVANT : panic si les longueurs diffèrent
if expectedExportCount != len(code.Params) {
    panic("导入参数数量不匹配...")
}
```

### Correction
- Suppression du bloc de vérification stricte (lignes 175-194 de l'original)
- Remplacement de la boucle `for i := 0; i < len(paramList)` par `for i := 0; i < maxLen` avec `maxLen = min(len(paramList), len(code.Params))`
- Ajout de bounds checking (`if pi < len(code.Params)`) dans le merge des `StringParam`, `JumpParam` et `[]uint16`

### Impact
Les traductions peuvent désormais être plus longues ou plus courtes que l'original. Les offsets de jump sont automatiquement recalculés par le code existant en aval.

---

## Patch 2 — Correction du pipeline CZ3

### Fichiers modifiés
`czimage/cz3.go`, `czimage/imagefix.go`

### Problème 1 — Magic byte écrasé (cz3.go, ligne 185)
La fonction `Write()` laissait le champ `CzHeader.Magic` se corrompre en "CZ0" au lieu de "CZ3". Le jeu ne reconnaissait plus le format du fichier.

```go
// FIX : forcer le magic avant écriture
cz.CzHeader.Magic = []byte{'C', 'Z', '3', 0}
```

### Problème 2 — Format NRGBA non garanti (cz3.go, lignes 84-99, 120-125)
Le format CZ3 encode les pixels en BGRA 32-bit. La bibliothèque PNG de Go peut décoder en RGB, RGBA, NRGBA ou paletted selon le fichier source. Sans conversion forcée, un PNG RGB 24-bit est traité comme 32-bit, décalant tous les pixels.

```go
// FIX : conversion systématique en NRGBA 32-bit
pic = ImageToNRGBA(cz.PngImage)
```

### Problème 3 — Buffer aliasing dans DiffLine/LineDiff (imagefix.go)
Le code original créait un alias de slice (`currLine = pic.Pix[i:...]`) au lieu d'une copie. L'opération de delta (`currLine[x] -= preLine[x]`) modifiait directement `pic.Pix`, corrompant les données source pour les lignes suivantes.

```go
// AVANT (buggé) : alias, modifie pic.Pix
currLine = pic.Pix[i : i+lineByteCount]

// APRÈS (corrigé) : copie dans un buffer séparé
copy(currLine, pic.Pix[i:i+lineByteCount])
```

Même problème dans `LineDiff()` : `preLine = currLine` créait un alias au lieu d'une copie.

---

## Patch 3 — Corruption mémoire dans le décompresseur LZW

### Fichier modifié
`czimage/lzw.go` — fonctions `decompressLZW()` et `decompressLZW2()`

### Problème
Le dictionnaire LZW ajoutait des entrées en faisant `dictionary[dictionaryCount] = append(w, entry[0])`. En Go, `append()` peut retourner le même slice sous-jacent si la capacité le permet. Quand `w` était ensuite réassigné (`w = entry`), l'ancienne entrée du dictionnaire pouvait pointer vers des données modifiées.

### Correction
Allocation explicite d'un nouveau slice avant ajout au dictionnaire :

```go
// AVANT (buggé) :
dictionary[dictionaryCount] = append(w, entry[0])

// APRÈS (corrigé) :
newEntry := make([]byte, len(w)+1)
copy(newEntry, w)
newEntry[len(w)] = entry[0]
dictionary[dictionaryCount] = newEntry
```

---

## Patch 4 — RawSize incorrect dans la table de blocs CZ

### Fichier modifié
`czimage/util.go` — fonctions `Compress()` et `Compress2()`

### Contexte
Le format CZ3 stocke les pixels sous forme de blocs LZW compressés. Chaque bloc déclare un `CompressedSize` et un `RawSize` dans une table en-tête. Le moteur de jeu décompresse chaque bloc en se fiant strictement au `RawSize` déclaré pour allouer les buffers et positionner les données.

### Problème 1 — Carry-over LZW non compensé

La fonction `compressLZW()` opère par blocs de `size` codes maximum. Quand elle atteint la limite, elle conserve un `lastElement` (l'élément en cours de construction dans le dictionnaire) qui sera reporté au bloc suivant. Le compteur `count` retourné inclut les bytes de cet élément :

```
Bloc N: lit 502 bytes, produit 500 codes, garde 1 byte en carry-over
  → count = 502, lastElement = 1 byte
  → RawSize DEVRAIT être 501 (les 502 lus - le 1 reporté)
  → RawSize ÉTAIT 502 (buggé)
```

Effet : le premier bloc déclare un `RawSize` trop grand de 1, le dernier bloc trop petit de 1. Le moteur de jeu lit 1 byte de trop dans le premier bloc et 1 byte de moins dans le dernier, provoquant un décalage qui se propage en cascade.

### Problème 2 — Encodage UTF-8 de Go

LuckSystem utilise des `string` Go comme clés du dictionnaire LZW. L'élément carry-over est construit par `element = string(c)` où `c` est un `byte` (0-255).

En Go, `string(byte(c))` effectue une conversion `byte → rune → UTF-8`. Pour les bytes 0-127, le résultat a une longueur de 1. Pour les bytes 128-255, Go produit une chaîne UTF-8 de **2 bytes** :

```go
string(byte(127)) // len = 1 (ASCII)
string(byte(128)) // len = 2 (UTF-8: 0xC2 0x80)
string(byte(255)) // len = 2 (UTF-8: 0xC3 0xBF)
```

La première tentative de fix utilisait `len(last)` pour compter les bytes de données en carry-over. Pour un carry de valeur 200, `len(last) = 2` en Go alors qu'il représente 1 seul byte de données. Cela causait des erreurs ±1 sur les blocs dont le carry-over tombait sur un octet > 127.

### Correction finale

Le carry-over de `compressLZW()` est **toujours 0 ou 1 byte de données**, quel que soit `len(last)` en Go. La construction du dictionnaire (`element = string(c)`) assigne toujours un seul byte source à `element` quand il y a un carry.

```go
// AVANT (buggé) — version originale :
RawSize: uint32(count) // inclut le carry-over

// TENTATIVE 1 (partiellement buggé) :
rawSize := prevCarryLen + count - len(last) // len(last) ≠ 1 pour bytes > 127

// APRÈS (corrigé) :
carry := 0
if len(last) > 0 {
    carry = 1  // toujours 1 DATA byte, peu importe len(last) en Go
}
rawSize := prevCarry + count - carry
```

### Vérification
Test de round-trip sur le CZ3 original d'AIR (`title1a`, 1280×720, 32-bit, 10 blocs) : les 10 `RawSize` produits par la version corrigée correspondent **exactement** à ceux du fichier original créé par les outils de Visual Art's.

```
Block | Visual Art's | LuckSystem corrigé | Diff
   0  |      447246  |           447246   |   0 ✅
   1  |      471332  |           471332   |   0 ✅
   2  |      612039  |           612039   |   0 ✅
  ...       ...               ...           ...
   9  |      271697  |           271697   |   0 ✅
```

### Portée
Ce bug affecte **tous les jeux** supportés par LuckSystem qui utilisent des images CZ multi-blocs (c'est-à-dire toute image dont les données compressées dépassent 0xFEFD codes LZW, soit la grande majorité des CG). Les CZ mono-bloc (petites images UI) ne sont pas affectés car il n'y a pas de carry-over.

---

## Correction annexe — Alignement PAK

### Fichier modifié
`pak/pak.go` — fonction `Write()`

### Correction
Ajout de padding (bytes nuls) en fin de fichier PAK pour aligner la taille totale sur `BlockSize`. Certains moteurs de jeu vérifient cet alignement lors du chargement.

---

## Fichiers modifiés (résumé)

| Fichier | Patch | Description |
|---------|-------|-------------|
| `script/script.go` | 1 | Import de scripts à longueur variable |
| `czimage/cz3.go` | 2 | Magic byte, conversion NRGBA, logs |
| `czimage/imagefix.go` | 2 | Buffer aliasing DiffLine/LineDiff |
| `czimage/lzw.go` | 3 | Corruption mémoire dictionnaire LZW |
| `czimage/util.go` | 4 | RawSize carry-over + UTF-8 |
| `pak/pak.go` | annexe | Alignement PAK |
