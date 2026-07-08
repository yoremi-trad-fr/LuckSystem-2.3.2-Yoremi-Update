package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveImageBatchImportSourceAcceptsCzSuffixPngForExtensionlessPakEntry(t *testing.T) {
	dir := t.TempDir()
	sourceName := "ET_YK00_MOJI01_EN"
	if err := os.WriteFile(filepath.Join(dir, sourceName), []byte("CZ3\x00"), 0644); err != nil {
		t.Fatal(err)
	}

	sourceCz, outputName, _, ok := resolveImageBatchImportSource(dir, sourceName+".cz3.png")
	if !ok {
		t.Fatal("expected CZ source to be found")
	}
	if filepath.Base(sourceCz) != sourceName {
		t.Fatalf("sourceCz = %q, want %q", filepath.Base(sourceCz), sourceName)
	}
	if outputName != sourceName {
		t.Fatalf("outputName = %q, want %q", outputName, sourceName)
	}
}

func TestResolveImageBatchImportSourcePrefersExactCzName(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"TITLE", "TITLE.cz3"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("CZ3\x00"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	sourceCz, outputName, _, ok := resolveImageBatchImportSource(dir, "TITLE.cz3.png")
	if !ok {
		t.Fatal("expected CZ source to be found")
	}
	if filepath.Base(sourceCz) != "TITLE.cz3" {
		t.Fatalf("sourceCz = %q, want TITLE.cz3", filepath.Base(sourceCz))
	}
	if outputName != "TITLE.cz3" {
		t.Fatalf("outputName = %q, want TITLE.cz3", outputName)
	}
}
