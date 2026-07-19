package cmd

import "testing"

func TestForkVersion(t *testing.T) {
	const want = "2.3.2-yoremi.3.31"
	if rootCmd.Version != want {
		t.Fatalf("root version = %q, want %q", rootCmd.Version, want)
	}
}

func TestLoggingFlagsAreAvailableToSubcommands(t *testing.T) {
	for _, name := range []string{"log", "log_level", "log_dir"} {
		if rootCmd.PersistentFlags().Lookup(name) == nil {
			t.Errorf("logging flag %q is not persistent", name)
		}
	}
}
