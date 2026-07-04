package font

import (
	"image"
	"image/color"
	"testing"
)

func TestAdjustMetricsSignedOffsets(t *testing.T) {
	char := 'ﺑ'
	info := metricTestInfo(char, DrawSize{
		X: signedMetricByte(-1),
		W: 9,
		Y: signedMetricByte(-15),
	})

	info.AdjustMetrics([]rune{char}, MetricAdjustOptions{
		SetY:    true,
		Y:       3,
		XOffset: 2,
		WOffset: -1,
	})

	draw := info.DrawSize[1]
	if got := int(int8(draw.X)); got != 1 {
		t.Fatalf("draw_x = %d, want 1", got)
	}
	if got := draw.W; got != 9 {
		t.Fatalf("draw_w = %d, want 9", got)
	}
	if got := int(int8(draw.Y)); got != 3 {
		t.Fatalf("draw_y = %d, want 3", got)
	}
	if got := info.UnicodeSize[char].W; got != 8 {
		t.Fatalf("unicode width = %d, want 8", got)
	}
}

func TestAdjustMetricsClampsValues(t *testing.T) {
	char := 'ﺑ'
	info := metricTestInfo(char, DrawSize{
		X: signedMetricByte(120),
		W: 2,
		Y: signedMetricByte(-120),
	})

	info.AdjustMetrics([]rune{char}, MetricAdjustOptions{
		YOffset: -20,
		XOffset: 20,
		WOffset: -20,
	})

	draw := info.DrawSize[1]
	if got := int(int8(draw.X)); got != 127 {
		t.Fatalf("draw_x = %d, want 127", got)
	}
	if got := draw.W; got != 2 {
		t.Fatalf("draw_w = %d, want 2", got)
	}
	if got := int(int8(draw.Y)); got != -128 {
		t.Fatalf("draw_y = %d, want -128", got)
	}
	if got := info.UnicodeSize[char].W; got != 1 {
		t.Fatalf("unicode width = %d, want 1", got)
	}
}

func TestAdjustMetricsSkipsDuplicateGlyphs(t *testing.T) {
	char := 'ﺑ'
	info := metricTestInfo(char, DrawSize{W: 9})

	info.AdjustMetrics([]rune{char, char}, MetricAdjustOptions{WOffset: -1})

	if got := info.DrawSize[1].W; got != 9 {
		t.Fatalf("draw_w = %d, want 9", got)
	}
	if got := info.UnicodeSize[char].W; got != 8 {
		t.Fatalf("unicode width = %d, want 8", got)
	}
}

func TestAdjustMetricsKeepsSharedGlyphAdvancesSeparate(t *testing.T) {
	first := 'ﺑ'
	second := 'ﺒ'
	info := metricTestInfo(first, DrawSize{W: 9})
	info.UnicodeIndex[second] = 1
	info.UnicodeSize[second].W = 11

	info.AdjustMetrics([]rune{first, second}, MetricAdjustOptions{WOffset: -1})

	if got := info.DrawSize[1].W; got != 9 {
		t.Fatalf("draw_w = %d, want 9", got)
	}
	if got := info.UnicodeSize[first].W; got != 8 {
		t.Fatalf("first unicode width = %d, want 8", got)
	}
	if got := info.UnicodeSize[second].W; got != 10 {
		t.Fatalf("second unicode width = %d, want 10", got)
	}
}

func TestGetStringImageUsesUnicodeAdvance(t *testing.T) {
	char := 'A'
	info := &Info{
		BlockSize:    8,
		DrawSize:     make([]DrawSize, 2),
		UnicodeIndex: make([]uint16, 65536),
		UnicodeSize:  make([]CharSize, 65536),
	}
	info.DrawSize[1] = DrawSize{W: 10}
	info.UnicodeIndex[char] = 1
	info.UnicodeSize[char].W = 4

	atlas := image.NewNRGBA(image.Rect(0, 0, 16, 8))
	atlas.SetNRGBA(8, 0, color.NRGBA{R: 255, A: 255})
	lucaFont := &LucaFont{
		Info:  info,
		Image: atlas,
	}

	out := lucaFont.GetStringImage("AA").(*image.NRGBA)

	if got := out.NRGBAAt(4, 0).A; got != 255 {
		t.Fatalf("second glyph alpha at unicode advance = %d, want 255", got)
	}
	if got := out.NRGBAAt(10, 0).A; got != 0 {
		t.Fatalf("second glyph alpha at draw width = %d, want 0", got)
	}
}

func TestBleedGlyphEdgesExtendsSelectedGlyphPixels(t *testing.T) {
	char := 'ﺒ'
	info := &Info{
		BlockSize:    10,
		DrawSize:     make([]DrawSize, 2),
		UnicodeIndex: make([]uint16, 65536),
	}
	info.DrawSize[1] = DrawSize{W: 10}
	info.UnicodeIndex[char] = 1
	atlas := image.NewNRGBA(image.Rect(0, 0, 20, 10))
	for y := 3; y <= 5; y++ {
		for x := 13; x <= 17; x++ {
			atlas.SetNRGBA(x, y, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}
	lucaFont := &LucaFont{
		Info:  info,
		Image: atlas,
	}

	lucaFont.BleedGlyphEdges([]rune{char}, 2)

	for _, x := range []int{12, 18} {
		if got := atlas.NRGBAAt(x, 4).A; got != 255 {
			t.Fatalf("alpha at x=%d = %d, want 255", x, got)
		}
	}
	for _, x := range []int{11, 19} {
		if got := atlas.NRGBAAt(x, 4).A; got != 255 {
			t.Fatalf("core alpha at x=%d = %d, want 255", x, got)
		}
	}
	if got := atlas.NRGBAAt(11, 3).A; got != 0 {
		t.Fatalf("non-core second-step alpha = %d, want 0", got)
	}
	if got := atlas.NRGBAAt(10, 4).A; got != 0 {
		t.Fatalf("alpha outside bleed = %d, want 0", got)
	}
	if got := info.DrawSize[1].W; got != 10 {
		t.Fatalf("draw_w = %d, want unchanged 10", got)
	}
}

func TestBleedGlyphEdgesDoesNotDrawPastGlyphCropWidth(t *testing.T) {
	char := 'ﺒ'
	info := &Info{
		BlockSize:    10,
		DrawSize:     make([]DrawSize, 2),
		UnicodeIndex: make([]uint16, 65536),
	}
	info.DrawSize[1] = DrawSize{W: 8}
	info.UnicodeIndex[char] = 1
	atlas := image.NewNRGBA(image.Rect(0, 0, 20, 10))
	for y := 3; y <= 5; y++ {
		for x := 13; x <= 17; x++ {
			atlas.SetNRGBA(x, y, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}
	lucaFont := &LucaFont{
		Info:  info,
		Image: atlas,
	}

	lucaFont.BleedGlyphEdges([]rune{char}, 2)

	if got := atlas.NRGBAAt(12, 4).A; got != 255 {
		t.Fatalf("left crop-safe bleed alpha at x=12 = %d, want 255", got)
	}
	if got := atlas.NRGBAAt(11, 4).A; got != 255 {
		t.Fatalf("left crop-safe second-step alpha at x=11 = %d, want 255", got)
	}
	for _, x := range []int{18, 19} {
		if got := atlas.NRGBAAt(x, 4).A; got != 0 {
			t.Fatalf("right alpha beyond draw_w at x=%d = %d, want 0", x, got)
		}
	}
	if got := info.DrawSize[1].W; got != 8 {
		t.Fatalf("draw_w = %d, want unchanged 8", got)
	}
}

func TestBleedGlyphEdgesKeepsOnePixelBleedOpaque(t *testing.T) {
	char := 'ﺒ'
	info := &Info{
		BlockSize:    10,
		DrawSize:     make([]DrawSize, 2),
		UnicodeIndex: make([]uint16, 65536),
	}
	info.DrawSize[1] = DrawSize{W: 10}
	info.UnicodeIndex[char] = 1
	atlas := image.NewNRGBA(image.Rect(0, 0, 20, 10))
	for y := 3; y <= 5; y++ {
		for x := 13; x <= 17; x++ {
			atlas.SetNRGBA(x, y, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}
	lucaFont := &LucaFont{
		Info:  info,
		Image: atlas,
	}

	lucaFont.BleedGlyphEdges([]rune{char}, 1)

	for _, x := range []int{12, 18} {
		if got := atlas.NRGBAAt(x, 4).A; got != 255 {
			t.Fatalf("one-pixel bleed alpha at x=%d = %d, want 255", x, got)
		}
	}
	for _, x := range []int{11, 19} {
		if got := atlas.NRGBAAt(x, 4).A; got != 0 {
			t.Fatalf("outside one-pixel bleed alpha at x=%d = %d, want 0", x, got)
		}
	}
	if got := info.DrawSize[1].W; got != 10 {
		t.Fatalf("draw_w = %d, want unchanged 10", got)
	}
}

func TestBleedGlyphEdgesUsesArabicJoiningSides(t *testing.T) {
	char := 'ﺑ'
	info := &Info{
		BlockSize:    10,
		DrawSize:     make([]DrawSize, 2),
		UnicodeIndex: make([]uint16, 65536),
	}
	info.DrawSize[1] = DrawSize{W: 8}
	info.UnicodeIndex[char] = 1
	atlas := image.NewNRGBA(image.Rect(0, 0, 20, 10))
	for y := 3; y <= 5; y++ {
		for x := 13; x <= 17; x++ {
			atlas.SetNRGBA(x, y, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}
	lucaFont := &LucaFont{
		Info:  info,
		Image: atlas,
	}

	lucaFont.BleedGlyphEdges([]rune{char}, 2)

	if got := atlas.NRGBAAt(12, 4).A; got != 255 {
		t.Fatalf("left connector alpha at x=12 = %d, want 255", got)
	}
	if got := atlas.NRGBAAt(11, 4).A; got != 255 {
		t.Fatalf("left connector second-step alpha at x=11 = %d, want 255", got)
	}
	for _, x := range []int{18, 19} {
		if got := atlas.NRGBAAt(x, 4).A; got != 0 {
			t.Fatalf("right non-connector alpha at x=%d = %d, want 0", x, got)
		}
	}
	if got := info.DrawSize[1].W; got != 8 {
		t.Fatalf("draw_w = %d, want unchanged 8", got)
	}
}

func TestBleedGlyphEdgesIgnoresSoftAntialiasRows(t *testing.T) {
	char := 'ﺒ'
	info := &Info{
		BlockSize:    10,
		DrawSize:     make([]DrawSize, 2),
		UnicodeIndex: make([]uint16, 65536),
	}
	info.DrawSize[1] = DrawSize{W: 8}
	info.UnicodeIndex[char] = 1
	atlas := image.NewNRGBA(image.Rect(0, 0, 20, 10))
	for x := 13; x <= 17; x++ {
		atlas.SetNRGBA(x, 4, color.NRGBA{R: 64, G: 64, B: 64, A: 64})
	}
	lucaFont := &LucaFont{
		Info:  info,
		Image: atlas,
	}

	lucaFont.BleedGlyphEdges([]rune{char}, 2)

	for _, x := range []int{11, 12, 18, 19} {
		if got := atlas.NRGBAAt(x, 4).A; got != 0 {
			t.Fatalf("soft edge bleed alpha at x=%d = %d, want 0", x, got)
		}
	}
	if got := info.DrawSize[1].W; got != 8 {
		t.Fatalf("draw_w = %d, want unchanged 8", got)
	}
}

func TestBleedGlyphEdgesIgnoresIsolatedHardRows(t *testing.T) {
	char := 'ﺒ'
	info := &Info{
		BlockSize:    10,
		DrawSize:     make([]DrawSize, 2),
		UnicodeIndex: make([]uint16, 65536),
	}
	info.DrawSize[1] = DrawSize{W: 10}
	info.UnicodeIndex[char] = 1
	atlas := image.NewNRGBA(image.Rect(0, 0, 20, 10))
	for x := 13; x <= 17; x++ {
		atlas.SetNRGBA(x, 4, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	}
	lucaFont := &LucaFont{
		Info:  info,
		Image: atlas,
	}

	lucaFont.BleedGlyphEdges([]rune{char}, 2)

	for _, x := range []int{11, 12, 18, 19} {
		if got := atlas.NRGBAAt(x, 4).A; got != 0 {
			t.Fatalf("isolated row bleed alpha at x=%d = %d, want 0", x, got)
		}
	}
	if got := info.DrawSize[1].W; got != 10 {
		t.Fatalf("draw_w = %d, want unchanged 10", got)
	}
}

func TestArabicPresentationFormBleedSidesRepeatAcrossGroups(t *testing.T) {
	if got := glyphBleedSidesForRune('ﻬ'); got != glyphBleedBoth {
		t.Fatalf("HEH medial bleed sides = %d, want both sides", got)
	}
	if got := glyphBleedSidesForRune('ﻫ'); got != glyphBleedLeft {
		t.Fatalf("HEH initial bleed sides = %d, want left side", got)
	}
	if got := glyphBleedSidesForRune('ﻪ'); got != glyphBleedRight {
		t.Fatalf("HEH final bleed sides = %d, want right side", got)
	}
}

func metricTestInfo(char rune, draw DrawSize) *Info {
	info := &Info{
		DrawSize:     make([]DrawSize, 2),
		UnicodeIndex: make([]uint16, 65536),
		UnicodeSize:  make([]CharSize, 65536),
	}
	info.DrawSize[1] = draw
	info.UnicodeIndex[char] = 1
	info.UnicodeSize[char].W = draw.W
	return info
}
