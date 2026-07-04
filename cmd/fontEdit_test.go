package cmd

import "testing"

func TestFontEditEffectiveArabicConnectorBleed(t *testing.T) {
	tests := []struct {
		name          string
		arabicMetrics bool
		bleed         int
		bleedChanged  bool
		want          int
	}{
		{name: "off by default", want: 0},
		{name: "manual without preset", bleed: 1, bleedChanged: true, want: 1},
		{name: "preset keeps connector bleed disabled by default", arabicMetrics: true, want: defaultArabicConnectorBleed},
		{name: "preset allows explicit disable", arabicMetrics: true, bleedChanged: true, want: 0},
		{name: "preset allows explicit override", arabicMetrics: true, bleed: 1, bleedChanged: true, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fontEditEffectiveArabicConnectorBleed(tt.arabicMetrics, tt.bleed, tt.bleedChanged)
			if got != tt.want {
				t.Fatalf("bleed = %d, want %d", got, tt.want)
			}
		})
	}
}
