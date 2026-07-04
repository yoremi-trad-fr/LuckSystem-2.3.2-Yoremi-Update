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
	char := 'ﺑ'
	info := &Info{
		BlockSize:    10,
		DrawSize:     make([]DrawSize, 2),
		UnicodeIndex: make([]uint16, 65536),
	}
	info.DrawSize[1] = DrawSize{W: 8}
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
		if got := atlas.NRGBAAt(x, 4).A; got != 255 {
			t.Fatalf("alpha at x=%d = %d, want 255", x, got)
		}
	}
	if got := atlas.NRGBAAt(10, 4).A; got != 0 {
		t.Fatalf("alpha outside bleed = %d, want 0", got)
	}
	if got := info.DrawSize[1].W; got != 10 {
		t.Fatalf("draw_w = %d, want 10", got)
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
