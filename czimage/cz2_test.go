package czimage

import (
	"bytes"
	"image/color"
	"testing"
)

func TestCz2FindClosestPaletteEntryReversesExportedColor(t *testing.T) {
	cz := &Cz2Image{
		ColorPanel: []color.NRGBA{
			{A: 0},
			{A: 80},
			{A: 160},
			{A: 255},
		},
	}

	if got := cz.findClosestPaletteEntry(color.NRGBA{A: 160}); got != 2 {
		t.Fatalf("exact palette color mapped to index %d, want 2", got)
	}
	if got := cz.findClosestPaletteEntry(color.NRGBA{A: 150}); got != 2 {
		t.Fatalf("nearest palette color mapped to index %d, want 2", got)
	}
}

func TestCz2WritePreservesOriginalEntryLength(t *testing.T) {
	cz := &Cz2Image{
		CzHeader: CzHeader{
			Magic:        []byte{'C', 'Z', '2', 0},
			HeaderLength: 18,
		},
		OriginalLength: 64,
		CzData: CzData{
			Raw: []byte{1, 2, 3},
			OutputInfo: &CzOutputInfo{
				FileCount: 1,
				BlockInfo: []CzBlockInfo{{CompressedSize: 3, RawSize: 3}},
			},
		},
	}

	var out bytes.Buffer
	if err := cz.Write(&out); err != nil {
		t.Fatal(err)
	}
	if out.Len() != cz.OriginalLength {
		t.Fatalf("CZ2 length = %d, want preserved length %d", out.Len(), cz.OriginalLength)
	}
}
