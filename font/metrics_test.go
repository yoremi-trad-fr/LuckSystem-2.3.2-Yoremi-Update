package font

import "testing"

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
	if got := draw.W; got != 8 {
		t.Fatalf("draw_w = %d, want 8", got)
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
	if got := draw.W; got != 1 {
		t.Fatalf("draw_w = %d, want 1", got)
	}
	if got := int(int8(draw.Y)); got != -128 {
		t.Fatalf("draw_y = %d, want -128", got)
	}
}

func TestAdjustMetricsSkipsDuplicateGlyphs(t *testing.T) {
	char := 'ﺑ'
	info := metricTestInfo(char, DrawSize{W: 9})

	info.AdjustMetrics([]rune{char, char}, MetricAdjustOptions{WOffset: -1})

	if got := info.DrawSize[1].W; got != 8 {
		t.Fatalf("draw_w = %d, want 8", got)
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
