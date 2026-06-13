package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDialogueExtractImportIncludesLogBeginWithLabels(t *testing.T) {
	dir := t.TempDir()
	scriptPath := filepath.Join(dir, "seen.txt")
	tsvPath := filepath.Join(dir, "seen.ext.txt")
	outPath := filepath.Join(dir, "seen.out.txt")

	script := strings.Join([]string{
		`MESSAGE_CLEAR ()`,
		`label1: LOG_BEGIN (0x0, 0x0, 0x0, "jp-log", "en-log", "cn-log")`,
		`global2: LOG_BEGIN (0x0, 0x0, 0x0, "jp-global", "en-global", "cn-global")`,
		`global3: label4: LOG_BEGIN (0x0, 0x0, 0x0, "jp-both", "en-both", "cn-both")`,
		`MESSAGE_WAIT ()`,
		`MESSAGE (1, "jp-msg", "en-msg")`,
		`SELECT (0, 0, 0, 0, "jp-choice", "en-choice")`,
	}, "\n")
	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		t.Fatal(err)
	}

	app := &App{}
	count, err := app.extractDialoguesFromFile(scriptPath, tsvPath, []int{1, 2})
	if err != nil {
		t.Fatalf("extractDialoguesFromFile returned error: %v", err)
	}
	if count != 5 {
		t.Fatalf("expected 5 extracted dialogue rows, got %d", count)
	}

	tsvBytes, err := os.ReadFile(tsvPath)
	if err != nil {
		t.Fatal(err)
	}
	tsv := string(tsvBytes)
	if strings.Contains(tsv, "MESSAGE_CLEAR") || strings.Contains(tsv, "MESSAGE_WAIT") {
		t.Fatalf("non-dialogue MESSAGE_* opcode was exported:\n%s", tsv)
	}
	if got := strings.Count(tsv, "\tLOG_BEGIN\t"); got != 3 {
		t.Fatalf("expected 3 LOG_BEGIN rows, got %d:\n%s", got, tsv)
	}

	tsv = strings.ReplaceAll(tsv, "en-", "fr-")
	if err := os.WriteFile(tsvPath, []byte(tsv), 0644); err != nil {
		t.Fatal(err)
	}

	replaced, err := app.importDialoguesToFile(scriptPath, tsvPath, 2, outPath)
	if err != nil {
		t.Fatalf("importDialoguesToFile returned error: %v", err)
	}
	if replaced != 5 {
		t.Fatalf("expected 5 replacements, got %d", replaced)
	}

	outBytes, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	out := string(outBytes)
	for _, want := range []string{"fr-log", "fr-global", "fr-both", "fr-msg", "fr-choice"} {
		if !strings.Contains(out, want) {
			t.Fatalf("output does not contain %q:\n%s", want, out)
		}
	}
}
