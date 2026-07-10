package font

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	"github.com/go-restruct/restruct"
)

func TestBuildFontAliasEntryKeepsTargetInfoGeometry(t *testing.T) {
	restruct.EnableExprBeta()
	source := aliasTestInfo(30, 31, 3)
	source.DrawSize[1] = DrawSize{X: 2, W: 17, Y: 4}
	source.UnicodeIndex['A'] = 1
	source.UnicodeSize['A'] = CharSize{W: 16}
	target := aliasTestInfo(32, 33, 2)

	sourceData := writeAliasTestInfo(t, source)
	targetData := writeAliasTestInfo(t, target)
	result, err := BuildFontAliasEntry("info30", "info32", sourceData, targetData)
	if err != nil {
		t.Fatal(err)
	}

	got := LoadFontInfo(result)
	if got.FontSize != 32 || got.BlockSize != 33 {
		t.Fatalf("target geometry not preserved: font=%d block=%d", got.FontSize, got.BlockSize)
	}
	if got.CharNum != source.CharNum {
		t.Fatalf("source character map not preserved: got %d, want %d", got.CharNum, source.CharNum)
	}
	if got.UnicodeIndex['A'] != 1 || got.DrawSize[1] != source.DrawSize[1] || got.UnicodeSize['A'] != source.UnicodeSize['A'] {
		t.Fatal("source glyph mapping or metrics were not preserved")
	}
}

func TestRegridFontAtlasPreservesGlyphPixelsAndTargetWidth(t *testing.T) {
	const (
		sourceBlock = 3
		targetBlock = 5
		targetWidth = fontAtlasColumns*targetBlock + 12
	)
	source := image.NewNRGBA(image.Rect(0, 0, fontAtlasColumns*sourceBlock+4, sourceBlock*2))
	marked := color.NRGBA{R: 10, G: 20, B: 30, A: 200}
	// Mark the last pixel of cell 1 on row 1. Regridding must move the
	// cell origin from x=3/y=3 to x=5/y=5 without scaling the glyph.
	source.SetNRGBA(sourceBlock+2, sourceBlock+1, marked)

	got, err := regridFontAtlas(source, sourceBlock, targetBlock, targetWidth)
	if err != nil {
		t.Fatal(err)
	}
	if got.Bounds().Dx() != targetWidth || got.Bounds().Dy() != targetBlock*2 {
		t.Fatalf("unexpected target atlas dimensions: %v", got.Bounds())
	}
	if pixel := got.NRGBAAt(targetBlock+2, targetBlock+1); pixel != marked {
		t.Fatalf("glyph pixel moved incorrectly: got %#v, want %#v", pixel, marked)
	}
	if pixel := got.NRGBAAt(targetBlock+sourceBlock, targetBlock+1); pixel.A != 0 {
		t.Fatalf("target-only cell padding is not empty: %#v", pixel)
	}
	if pixel := got.NRGBAAt(targetWidth-1, 0); pixel.A != 0 {
		t.Fatalf("target texture alignment padding is not empty: %#v", pixel)
	}
}

func TestBuildFontAliasEntryRejectsMixedEntryKinds(t *testing.T) {
	_, err := BuildFontAliasEntry("info30", "明朝32", []byte{0, 0, 0, 0, 0, 0}, make([]byte, 15))
	if err == nil {
		t.Fatal("expected mixed info/image alias to fail")
	}
}

func TestBuildFontAliasEntryPreservesLongerTargetEntryLength(t *testing.T) {
	restruct.EnableExprBeta()
	sourceData := writeAliasTestInfo(t, aliasTestInfo(30, 31, 2))
	targetData := writeAliasTestInfo(t, aliasTestInfo(32, 33, 5))

	result, err := BuildFontAliasEntry("info30", "info32", sourceData, targetData)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != len(targetData) {
		t.Fatalf("alias entry length = %d, want target length %d", len(result), len(targetData))
	}
	got := LoadFontInfo(result)
	if got.CharNum != 2 || got.FontSize != 32 || got.BlockSize != 33 {
		t.Fatalf("unexpected padded alias info: chars=%d font=%d block=%d", got.CharNum, got.FontSize, got.BlockSize)
	}
}

func aliasTestInfo(fontSize, blockSize, charNum int) *Info {
	return &Info{
		FontSize:     uint16(fontSize),
		BlockSize:    uint16(blockSize),
		CharNum:      uint16(charNum),
		DrawSize:     make([]DrawSize, charNum),
		UnicodeIndex: make([]uint16, 65536),
		UnicodeSize:  make([]CharSize, 65536),
	}
}

func writeAliasTestInfo(t *testing.T, info *Info) []byte {
	t.Helper()
	var out bytes.Buffer
	if err := info.Write(&out); err != nil {
		t.Fatal(err)
	}
	return out.Bytes()
}
