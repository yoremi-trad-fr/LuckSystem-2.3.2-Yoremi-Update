package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-restruct/restruct"

	"lucksystem/charset"
	"lucksystem/font"
	"lucksystem/pak"
)

func main() {
	restruct.EnableExprBeta()
	if len(os.Args) != 4 {
		fatalf("usage: fontdiag <font-root> <output-dir> <family-pak-name>")
	}

	root, err := filepath.Abs(os.Args[1])
	check(err)
	outDir, err := filepath.Abs(os.Args[2])
	check(err)
	familyName := os.Args[3]
	check(os.MkdirAll(outDir, 0755))

	infoName := "FONT__INFO.PAK"
	infoDir := filepath.Join(root, "font_win32_1280")
	familyPath := filepath.Join(infoDir, familyName)
	if _, err := os.Stat(familyPath); err != nil {
		infoName = "FONTZC__INFO.PAK"
		infoDir = filepath.Join(root, "fontzc_win32_1280")
		familyPath = filepath.Join(infoDir, familyName)
	}

	infoPak := pak.LoadPak(filepath.Join(infoDir, infoName), charset.UTF_8)
	infoPak.ReadAll()
	familyPak := pak.LoadPak(familyPath, charset.UTF_8)
	familyPak.ReadAll()
	if len(infoPak.Files) != len(familyPak.Files) {
		fatalf("file count mismatch: %d info files, %d glyph files", len(infoPak.Files), len(familyPak.Files))
	}

	for i, glyphEntry := range familyPak.Files {
		lucaFont := font.LoadLucaFont(infoPak.Files[i].Data, glyphEntry.Data)
		var glyphOut bytes.Buffer
		if err := lucaFont.Write(&glyphOut, nil); err != nil {
			fatalf("%s: %v", glyphEntry.Name, err)
		}
		check(familyPak.Set(glyphEntry.Name, bytes.NewReader(glyphOut.Bytes())))
	}
	familyPak.Rebuild = true

	outName := filepath.Join(outDir, "ROUNDTRIP_"+familyName)
	out, err := os.Create(outName)
	check(err)
	check(familyPak.Write(out))
	check(out.Close())
	fmt.Println(outName)
}

func check(err error) {
	if err != nil {
		fatalf("%v", err)
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
