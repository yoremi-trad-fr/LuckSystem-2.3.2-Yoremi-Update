package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/go-restruct/restruct"
	"golang.org/x/image/math/fixed"

	"lucksystem/charset"
	"lucksystem/font"
	"lucksystem/pak"
)

type vietnameseFontSet struct {
	Name       string
	InfoPak    string
	FamilyPaks []string
}

// VietnameseFontPatch generates AIR/Planetarian SG Vietnamese font PAKs from the GUI.
// It embeds the tools/vietfontpatch workflow so users do not need a separate command-line tool.
func (a *App) VietnameseFontPatch(fontRoot, charsetFile, ttfFile, outputDir, slot, family string, yOffsets []int, redrawLatin bool) string {
	if fontRoot == "" || charsetFile == "" || ttfFile == "" || outputDir == "" {
		a.logError("Font root, charset file, TTF file, and output folder are required")
		return "ERROR"
	}
	if len(yOffsets) == 0 {
		a.logError("Select at least one Y offset")
		return "ERROR"
	}
	if slot == "" {
		slot = "en"
	}
	if family == "" {
		family = "GOTHIC1"
	}

	restruct.EnableExprBeta()

	fontRoot, err := filepath.Abs(fontRoot)
	if err != nil {
		a.logError(err.Error())
		return "ERROR"
	}
	charsetFile, err = filepath.Abs(charsetFile)
	if err != nil {
		a.logError(err.Error())
		return "ERROR"
	}
	ttfFile, err = filepath.Abs(ttfFile)
	if err != nil {
		a.logError(err.Error())
		return "ERROR"
	}
	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		a.logError(err.Error())
		return "ERROR"
	}

	charsBytes, err := os.ReadFile(charsetFile)
	if err != nil {
		a.logError(fmt.Sprintf("Cannot read charset file: %v", err))
		return "ERROR"
	}
	chars := strings.TrimPrefix(string(charsBytes), "\ufeff")
	chars = strings.TrimRight(chars, "\r\n")
	if chars == "" {
		a.logError("Charset file is empty")
		return "ERROR"
	}

	ttfBytes, err := os.ReadFile(ttfFile)
	if err != nil {
		a.logError(fmt.Sprintf("Cannot read TTF file: %v", err))
		return "ERROR"
	}

	sort.Ints(yOffsets)
	yOffsets = uniqueInts(yOffsets)

	a.log("════════════════════════════════════════")
	a.log("  AIR / SG VIETNAMESE FONT PATCH")
	a.log("════════════════════════════════════════")
	a.log(fmt.Sprintf("Font root: %s", fontRoot))
	a.log(fmt.Sprintf("Charset:   %s", charsetFile))
	a.log(fmt.Sprintf("TTF:       %s", ttfFile))
	a.log(fmt.Sprintf("Slot:      %s", slot))
	a.log(fmt.Sprintf("Family:    %s", family))
	a.log(fmt.Sprintf("Y offsets: %s", formatYOffsets(yOffsets)))
	if redrawLatin {
		a.log("Mode:      experimental Latin redraw")
	} else {
		a.log("Mode:      safe missing-glyph injection")
	}

	for _, yOffset := range yOffsets {
		runDir := filepath.Join(outputDir, buildVietnameseOutputName(ttfFile, slot, family, yOffset, redrawLatin))
		if err := os.MkdirAll(runDir, 0755); err != nil {
			a.logError(fmt.Sprintf("Cannot create output folder: %v", err))
			return "ERROR"
		}
		a.log("────────────────────────────────────────")
		a.log(fmt.Sprintf("Generating Y%+d -> %s", yOffset, runDir))
		if err := a.patchVietnameseFontOnce(fontRoot, chars, ttfBytes, runDir, slot, family, yOffset, redrawLatin); err != nil {
			a.logError(err.Error())
			return "ERROR"
		}
		a.logOK(fmt.Sprintf("Generated Y%+d in %s", yOffset, runDir))
	}

	a.log("════════════════════════════════════════")
	a.logOK("Vietnamese font patch completed")
	return "OK"
}

func (a *App) patchVietnameseFontOnce(fontRoot, chars string, ttfBytes []byte, outputDir, slot, family string, yOffset int, redrawLatin bool) error {
	sets := []vietnameseFontSet{
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
		if slot != "all" && slot != set.Name && !(slot == "en" && set.Name == "jp") {
			continue
		}
		set.FamilyPaks = filterVietnameseFamilyPaks(set.FamilyPaks, family)
		if len(set.FamilyPaks) == 0 {
			continue
		}
		if err := validateVietnameseSetFiles(set); err != nil {
			return err
		}
		if err := a.patchVietnameseSet(set, chars, ttfBytes, outputDir, yOffset, redrawLatin); err != nil {
			return err
		}
		patched = true
	}
	if !patched {
		return fmt.Errorf("no font family matched slot=%s family=%s", slot, family)
	}
	return nil
}

func validateVietnameseSetFiles(set vietnameseFontSet) error {
	if _, err := os.Stat(set.InfoPak); err != nil {
		return fmt.Errorf("missing info PAK for %s slot: %s", set.Name, set.InfoPak)
	}
	for _, familyPak := range set.FamilyPaks {
		if _, err := os.Stat(familyPak); err != nil {
			return fmt.Errorf("missing family PAK for %s slot: %s", set.Name, familyPak)
		}
	}
	return nil
}

func (a *App) patchVietnameseSet(set vietnameseFontSet, chars string, ttfBytes []byte, outputDir string, yOffset int, redrawLatin bool) error {
	infoPak := pak.LoadPak(set.InfoPak, charset.UTF_8)
	infoPak.ReadAll()
	if len(infoPak.Files) == 0 {
		return fmt.Errorf("empty info PAK: %s", set.InfoPak)
	}

	requestedRunes := []rune(chars)
	referenceInfo := font.LoadFontInfo(infoPak.Files[0].Data)
	patchRunes := missingVietnameseRunes(referenceInfo, requestedRunes)
	if len(patchRunes) == 0 {
		return fmt.Errorf("%s already contains all requested characters", set.InfoPak)
	}
	patchChars := string(patchRunes)
	redrawRunes := []rune(nil)
	if redrawLatin {
		redrawRunes = latinExperimentRunes(referenceInfo, requestedRunes)
		if len(redrawRunes) == 0 {
			return fmt.Errorf("%s has no mapped Latin/Vietnamese glyphs to redraw", set.InfoPak)
		}
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

	a.log(fmt.Sprintf("Patching %s slot", set.Name))
	a.log(fmt.Sprintf("  requested: %d, already present: %d, injected: %d",
		len(requestedRunes), len(requestedRunes)-len(patchRunes), len(patchRunes)))
	a.log(fmt.Sprintf("  base cells: %d, replace index: %d", baseCount, startIndex))
	if redrawLatin {
		a.log(fmt.Sprintf("  experimental redraw: %d existing Latin/Vietnamese glyphs", len(redrawRunes)))
	}

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
			lowerY := vietnameseReferenceY(lucaFont.Info, []rune{'ó', 'á', 'a', 'o'})
			upperY := vietnameseReferenceY(lucaFont.Info, []rune{'Á', 'Â', 'A', 'O'})
			lucaFont.ReplaceChars(bytes.NewReader(ttfBytes), patchChars, startIndex, false)
			normalizeRunesVerticalMetrics(lucaFont.Info, patchRunes, yOffset, lowerY, upperY)
			if redrawLatin {
				if err := redrawMappedRunesFromCurrentFace(lucaFont, redrawRunes); err != nil {
					return fmt.Errorf("%s/%s: %w", familyPakName, glyphEntry.Name, err)
				}
				normalizeRunesVerticalMetrics(lucaFont.Info, redrawRunes, yOffset, lowerY, upperY)
			}

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
		if err := outFile.Close(); err != nil {
			return err
		}
		a.logOK(fmt.Sprintf("Wrote %s", outName))
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
	if err := outInfo.Close(); err != nil {
		return err
	}
	a.logOK(fmt.Sprintf("Wrote %s", outInfoName))
	return nil
}

func latinExperimentRunes(info *font.Info, requested []rune) []rune {
	candidates := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	candidates = append(candidates, requested...)
	return mappedUniqueRunes(info, candidates)
}

func mappedUniqueRunes(info *font.Info, chars []rune) []rune {
	seen := make(map[rune]bool, len(chars))
	out := make([]rune, 0, len(chars))
	for _, char := range chars {
		if seen[char] {
			continue
		}
		seen[char] = true
		if !hasVietnameseRune(info, char) {
			continue
		}
		out = append(out, char)
	}
	return out
}

func redrawMappedRunesFromCurrentFace(lucaFont *font.LucaFont, chars []rune) error {
	if lucaFont.Info == nil || lucaFont.Info.FontFace == nil {
		return fmt.Errorf("font face is not loaded")
	}
	if lucaFont.Image == nil {
		return fmt.Errorf("font atlas image is not loaded")
	}

	size := int(lucaFont.Info.BlockSize)
	alphaMask := image.NewAlpha(image.Rect(0, 0, size, size))
	for _, char := range chars {
		if char < 0 || int(char) >= len(lucaFont.Info.UnicodeIndex) {
			continue
		}
		fontIndex := lucaFont.Info.UnicodeIndex[int(char)]
		if fontIndex == 0 && char != ' ' {
			continue
		}
		index := int(fontIndex)
		bounds, advance, ok := lucaFont.Info.FontFace.GlyphBounds(char)
		if !ok {
			return fmt.Errorf("selected TTF does not contain %q (U+%04X)", char, char)
		}

		w := uint8(advance.Ceil())
		if char == 32 || w == 0 {
			w = uint8(lucaFont.Info.FontSize)
		}
		lucaFont.Info.DrawSize[index].X = uint8(bounds.Min.X.Floor())
		lucaFont.Info.DrawSize[index].W = w
		lucaFont.Info.DrawSize[index].Y = uint8(bounds.Min.Y.Floor())
		lucaFont.Info.UnicodeSize[char].W = w

		y := index / 100
		x := index % 100
		point := fixed.Point26_6{
			X: fixed.Int26_6(x * 64),
			Y: fixed.Int26_6(y * 64),
		}
		_, img, _, _, ok := lucaFont.Info.FontFace.Glyph(point, char)
		if !ok {
			return fmt.Errorf("selected TTF cannot draw %q (U+%04X)", char, char)
		}
		cellRect := image.Rect(x*size, y*size, (x+1)*size, (y+1)*size).Intersect(lucaFont.Image.Bounds())
		if cellRect.Empty() {
			return fmt.Errorf("mapped glyph %q points outside the atlas at index %d", char, index)
		}
		draw.Draw(lucaFont.Image, cellRect, alphaMask, alphaMask.Bounds().Min, draw.Src)
		draw.Draw(lucaFont.Image, cellRect, img, img.Bounds().Min, draw.Src)
	}
	return nil
}

func filterVietnameseFamilyPaks(familyPaks []string, family string) []string {
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

func missingVietnameseRunes(info *font.Info, chars []rune) []rune {
	seen := make(map[rune]bool, len(chars))
	missing := make([]rune, 0, len(chars))
	for _, char := range chars {
		if seen[char] {
			continue
		}
		seen[char] = true
		if hasVietnameseRune(info, char) {
			continue
		}
		missing = append(missing, char)
	}
	return missing
}

func hasVietnameseRune(info *font.Info, char rune) bool {
	if char < 0 || int(char) >= len(info.UnicodeIndex) {
		return false
	}
	return char == ' ' || info.UnicodeIndex[int(char)] != 0
}

func normalizeVietnameseVerticalMetrics(info *font.Info, chars []rune, yOffset int) {
	lowerY := vietnameseReferenceY(info, []rune{'ó', 'á', 'a', 'o'})
	upperY := vietnameseReferenceY(info, []rune{'Á', 'Â', 'A', 'O'})
	normalizeRunesVerticalMetrics(info, chars, yOffset, lowerY, upperY)
}

func normalizeRunesVerticalMetrics(info *font.Info, chars []rune, yOffset int, lowerY, upperY uint8) {
	for _, char := range chars {
		if char < 0 || int(char) >= len(info.UnicodeIndex) {
			continue
		}
		index := info.UnicodeIndex[int(char)]
		if index == 0 && char != ' ' {
			continue
		}
		if unicode.IsUpper(char) {
			info.DrawSize[index].Y = addSignedVietnameseYOffset(upperY, yOffset)
		} else {
			info.DrawSize[index].Y = addSignedVietnameseYOffset(lowerY, yOffset)
		}
	}
}

func vietnameseReferenceY(info *font.Info, candidates []rune) uint8 {
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

func addSignedVietnameseYOffset(raw uint8, offset int) uint8 {
	value := int(int8(raw)) + offset
	if value < -128 {
		value = -128
	}
	if value > 127 {
		value = 127
	}
	return uint8(int8(value))
}

func uniqueInts(values []int) []int {
	if len(values) == 0 {
		return values
	}
	out := values[:1]
	for _, value := range values[1:] {
		if value != out[len(out)-1] {
			out = append(out, value)
		}
	}
	return out
}

func formatYOffsets(values []int) string {
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = fmt.Sprintf("Y%+d", value)
	}
	return strings.Join(parts, ", ")
}

func buildVietnameseOutputName(ttfFile, slot, family string, yOffset int, redrawLatin bool) string {
	base := strings.TrimSuffix(filepath.Base(ttfFile), filepath.Ext(ttfFile))
	base = sanitizeVietnamesePathPart(base)
	slot = sanitizeVietnamesePathPart(slot)
	family = sanitizeVietnamesePathPart(strings.ToUpper(family))
	mode := ""
	if redrawLatin {
		mode = "_LATIN"
	}
	return fmt.Sprintf("%s_%s_%s%s_Y%+d", base, slot, family, mode, yOffset)
}

func sanitizeVietnamesePathPart(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "font"
	}
	re := regexp.MustCompile(`[^A-Za-z0-9._+-]+`)
	value = re.ReplaceAllString(value, "_")
	value = strings.Trim(value, "._")
	if value == "" {
		return "font"
	}
	return value
}
