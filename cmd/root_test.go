package cmd

import "testing"

func TestForkVersion(t *testing.T) {
	const want = "2.3.2-yoremi.3.27"
	if rootCmd.Version != want {
		t.Fatalf("root version = %q, want %q", rootCmd.Version, want)
	}
}
