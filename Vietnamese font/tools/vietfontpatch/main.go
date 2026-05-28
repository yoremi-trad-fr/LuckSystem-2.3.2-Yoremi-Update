package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/go-restruct/restruct"

	"lucksystem/charset"
	"lucksystem/font"
	"lucksystem/pak"
)

type fontSet struct {
	Name       string
	InfoPak    string
	FamilyPaks []string
}

func main() {
	restruct.EnableExprBeta()

	slot := flag.String("slot", "all", "font slot to patch: all, en, zc")
	family := flag.String("family", "all", "font family to patch: all, GOTHIC1, GOTHIC2, GOTHIC3, MINCHO, MODERN")
	yOffset := flag.Int("yoffset", 0, "extra signed vertical offset for injected characters")
	flag.Parse()

	if flag.NArg() != 4 {
		fatalf("usage: vietfontpatch [-slot all|en|zc] [-family all|GOTHIC1|...] [-yoffset N] <font-root> <charset-file> <ttf-file> <output-dir>")
	}

	fontRoot, err := filepath.Abs(flag.Arg(0))
	check(err)
	charsetFile, err := filepath.Abs(flag.Arg(1))
	check(err)
	ttfFile, err := filepath.Abs(flag.Arg(2))
	check(err)
	outputDir, err := filepath.Abs(flag.Arg(3))
	check(err)
	check(os.MkdirAll(outputDir, 0755))

	charsBytes, err := os.ReadFile(charsetFile)
	check(err)
	chars := strings.TrimPrefix(string(charsBytes), "\ufeff")
	chars = strings.TrimRight(chars, "\r\n")
	if chars == "" {
		fatalf("charset is empty: %s", charsetFile)
	}

	ttfBytes, err := os.ReadFile(ttfFile)
	check(err)

	sets := []fontSet{
		{
			Name:    "jp",
			InfoPak: filepath.Join(fontRoot, "font_win32_1280", "FONT__INFO.PAK"),
			FamilyPaks: []string{
				filepath.Join(fontRoot, "font_win32_1280", "FONT_GOTHIC1.PAK"),
				filepath.Join(fontRoot, "font_win32_1280", "FONT_GOTHIC2.PAK"),
				filepath.Join(fontRoot, "font_win32_1280", "FONT_GOTHIC3.PAK"),
				filepath.Join(fontRoot, "font_win32_1280", "FONT_MINCHO.PAK"),
				filepath.Join(fontRoot, "font_win32_1280", "FONT_MODERN.PAK"),
			},
		},
		{
			Name:    "zc",
			InfoPak: filepath.Join(fontRoot, "fontzc_win32_1280", "FONTZC__INFO.PAK"),
			FamilyPaks: []string{
				filepath.Join(fontRoot, "fontzc_win32_1280", "FONTZC_GOTHIC1.PAK"),
				filepath.Join(fontRoot, "fontzc_win32_1280", "FONTZC_GOTHIC2.PAK"),
				filepath.Join(fontRoot, "fontzc_win32_1280", "FONTZC_MINCHO.PAK"),
			},
		},
	}

	patched := false
	for _, set := range sets {
		if *slot != "all" && *slot != set.Name && !(*slot == "en" && set.Name == "jp") {
			continue
		}
		set.FamilyPaks = filterFamilyPaks(set.FamilyPaks, *family)
		if len(set.FamilyPaks) == 0 {
			continue
		}
		fmt.Printf("patching %s\n", set.Name)
		check(patchSet(set, chars, ttfBytes, outputDir, *yOffset))
		patched = true
	}
	if !patched {
		fatalf("no font family matched slot=%s family=%s", *slot, *family)
	}
}

func filterFamilyPaks(familyPaks []string, family string) []string {
	filter := strings.ToUpper(strings.TrimSuffix(family, ".PAK"))
	if filter == "" || filter == "ALL" {
		return familyPaks
	}
	selected := make([]string, 0, len(familyPaks))
	for _, familyPak := range familyPaks {
		base := strings.ToUpper(strings.TrimSuffix(filepath.Base(familyPak), ".PAK"))
		if base == filter || strings.HasSuffix(base, "_"+filter) || strings.Contains(base, filter) {
			selected = append(selected, familyPak)
		}
	}
	return selected
}

func patchSet(set fontSet, chars string, ttfBytes []byte, outputDir string, yOffset int) error {
	infoPak := pak.LoadPak(set.InfoPak, charset.UTF_8)
	infoPak.ReadAll()
	if len(infoPak.Files) == 0 {
		return fmt.Errorf("empty info pak: %s", set.InfoPak)
	}
	requestedRunes := []rune(chars)
	referenceInfo := font.LoadFontInfo(infoPak.Files[0].Data)
	patchRunes := missingRunes(referenceInfo, requestedRunes)
	patchChars := string(patchRunes)
	if len(patchRunes) == 0 {
		return fmt.Errorf("%s already contains all requested characters", set.InfoPak)
	}

	baseCount := uint16(0)
	for _, entry := range infoPak.Files {
		info := font.LoadFontInfo(entry.Data)
		if int(info.CharNum) >= len(patchRunes) && (baseCount == 0 || info.CharNum < baseCount) {
			baseCount = info.CharNum
		}
	}
	if baseCount == 0 {
		return fmt.Errorf("%s has no info table large enough for %d chars", set.InfoPak, len(patchRunes))
	}
	startIndex := int(baseCount) - len(patchRunes)
	fmt.Printf("  requested: %d, already present: %d, injected: %d\n",
		len(requestedRunes), len(requestedRunes)-len(patchRunes), len(patchRunes))
	fmt.Printf("  base cells: %d, replace index: %d\n", baseCount, startIndex)

	var patchedInfos [][]byte
	for _, familyPakName := range set.FamilyPaks {
		familyPak := pak.LoadPak(familyPakName, charset.UTF_8)
		familyPak.ReadAll()
		if len(familyPak.Files) != len(infoPak.Files) {
			return fmt.Errorf("file count mismatch: %s has %d, %s has %d",
				familyPakName, len(familyPak.Files), set.InfoPak, len(infoPak.Files))
		}

		for index, glyphEntry := range familyPak.Files {
			infoEntry := infoPak.Files[index]
			lucaFont := font.LoadLucaFont(infoEntry.Data, glyphEntry.Data)
			lucaFont.ReplaceChars(bytes.NewReader(ttfBytes), patchChars, startIndex, false)
			normalizeVerticalMetrics(lucaFont.Info, patchRunes, yOffset)

			var glyphOut bytes.Buffer
			var infoOut bytes.Buffer
			if err := lucaFont.Write(&glyphOut, &infoOut); err != nil {
				return fmt.Errorf("%s/%s: %w", familyPakName, glyphEntry.Name, err)
			}
			if err := familyPak.Set(glyphEntry.Name, bytes.NewReader(glyphOut.Bytes())); err != nil {
				return err
			}
			if patchedInfos == nil {
				patchedInfos = make([][]byte, len(infoPak.Files))
			}
			if patchedInfos[index] == nil {
				patchedInfos[index] = append([]byte(nil), infoOut.Bytes()...)
			}
		}
		familyPak.Rebuild = true

		outName := filepath.Join(outputDir, filepath.Base(familyPakName))
		outFile, err := os.Create(outName)
		if err != nil {
			return err
		}
		if err := familyPak.Write(outFile); err != nil {
			_ = outFile.Close()
			return err
		}
		check(outFile.Close())
		fmt.Printf("  wrote %s\n", outName)
	}

	for index, infoBytes := range patchedInfos {
		if infoBytes == nil {
			return fmt.Errorf("missing patched info at index %d", index)
		}
		if err := infoPak.Set(infoPak.Files[index].Name, bytes.NewReader(infoBytes)); err != nil {
			return err
		}
	}
	outInfoName := filepath.Join(outputDir, filepath.Base(set.InfoPak))
	outInfo, err := os.Create(outInfoName)
	if err != nil {
		return err
	}
	if err := infoPak.Write(outInfo); err != nil {
		_ = outInfo.Close()
		return err
	}
	check(outInfo.Close())
	fmt.Printf("  wrote %s\n", outInfoName)
	return nil
}

func missingRunes(info *font.Info, chars []rune) []rune {
	seen := make(map[rune]bool, len(chars))
	missing := make([]rune, 0, len(chars))
	for _, char := range chars {
		if seen[char] {
			continue
		}
		seen[char] = true
		if hasRune(info, char) {
			continue
		}
		missing = append(missing, char)
	}
	return missing
}

func hasRune(info *font.Info, char rune) bool {
	if char < 0 || int(char) >= len(info.UnicodeIndex) {
		return false
	}
	return char == ' ' || info.UnicodeIndex[int(char)] != 0
}

func normalizeVerticalMetrics(info *font.Info, chars []rune, yOffset int) {
	lowerY := referenceY(info, []rune{'ó', 'á', 'a', 'o'})
	upperY := referenceY(info, []rune{'Á', 'Â', 'A', 'O'})
	for _, char := range chars {
		if char < 0 || int(char) >= len(info.UnicodeIndex) {
			continue
		}
		index := info.UnicodeIndex[int(char)]
		if index == 0 && char != ' ' {
			continue
		}
		if unicode.IsUpper(char) {
			info.DrawSize[index].Y = addSignedYOffset(upperY, yOffset)
		} else {
			info.DrawSize[index].Y = addSignedYOffset(lowerY, yOffset)
		}
	}
}

func addSignedYOffset(raw uint8, offset int) uint8 {
	value := int(int8(raw)) + offset
	if value < -128 {
		value = -128
	}
	if value > 127 {
		value = 127
	}
	return uint8(int8(value))
}

func referenceY(info *font.Info, candidates []rune) uint8 {
	for _, char := range candidates {
		if char < 0 || int(char) >= len(info.UnicodeIndex) {
			continue
		}
		index := info.UnicodeIndex[int(char)]
		if index == 0 && char != ' ' {
			continue
		}
		return info.DrawSize[index].Y
	}
	return 0
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
