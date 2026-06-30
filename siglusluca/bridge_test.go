package siglusluca

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunKeepsMergedSiglusLinesForReviewAndReportsHD(t *testing.T) {
	dir := t.TempDir()
	lucaDir := filepath.Join(dir, "luca")
	siglusDir := filepath.Join(dir, "siglus")
	outDir := filepath.Join(dir, "out")
	if err := os.MkdirAll(lucaDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(siglusDir, 0755); err != nil {
		t.Fatal(err)
	}

	lucaScript := strings.Join([]string{
		`MESSAGE (0, "jp", "Alpha machine battery road.")`,
		`MESSAGE (0, "jp", "The person must solve the wall problem through personal effort.")`,
		`MESSAGE (0, "jp", "Help them only after they request aid and support.")`,
		`MESSAGE (0, "jp", "Final sunset melody.")`,
	}, "\n")
	if err := os.WriteFile(filepath.Join(lucaDir, "scene.txt"), []byte(lucaScript), 0644); err != nil {
		t.Fatal(err)
	}

	siglusScript := strings.Join([]string{
		`○0000000001○Alpha machine battery road.`,
		`●0000000001●Route française alpha.`,
		`○0000000002○The person must solve the wall problem through personal effort. Help them only after they request aid and support.`,
		`●0000000002●Bloc français fusionné à découper manuellement.`,
		`○0000000003○Final sunset melody.`,
		`●0000000003●Mélodie finale du coucher de soleil.`,
	}, "\n\n")
	if err := os.WriteFile(filepath.Join(siglusDir, "scene.ss.txt"), []byte(siglusScript), 0644); err != nil {
		t.Fatal(err)
	}

	hdLucaScript := strings.Join([]string{
		`MESSAGE (0, "jp", "Opening shared line.")`,
		`MESSAGE (0, "jp", "This HD exclusive epilogue sentence has no Siglus counterpart.")`,
		`MESSAGE (0, "jp", "Closing shared line.")`,
	}, "\n")
	if err := os.WriteFile(filepath.Join(lucaDir, "hdscene.txt"), []byte(hdLucaScript), 0644); err != nil {
		t.Fatal(err)
	}
	hdSiglusScript := strings.Join([]string{
		`○0000000001○Opening shared line.`,
		`●0000000001●Ligne commune d’ouverture.`,
		`○0000000002○Closing shared line.`,
		`●0000000002●Ligne commune de fermeture.`,
	}, "\n\n")
	if err := os.WriteFile(filepath.Join(siglusDir, "hdscene.ss.txt"), []byte(hdSiglusScript), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Run(Options{
		LucaDir:   lucaDir,
		SiglusDir: siglusDir,
		OutputDir: outDir,
		TargetCol: 2,
	})
	if err != nil {
		t.Fatal(err)
	}

	outBytes, err := os.ReadFile(filepath.Join(outDir, "scene.txt"))
	if err != nil {
		t.Fatal(err)
	}
	out := string(outBytes)
	if !strings.Contains(out, "Route française alpha.") || !strings.Contains(out, "Mélodie finale du coucher de soleil.") {
		t.Fatalf("expected normal aligned lines to be imported:\n%s", out)
	}
	if strings.Contains(out, "Bloc français fusionné") {
		t.Fatalf("merged Siglus text was imported into a single Luca line:\n%s", out)
	}
	if !strings.Contains(out, "The person must solve the wall problem through personal effort.") ||
		!strings.Contains(out, "Help them only after they request aid and support.") {
		t.Fatalf("merged Luca lines should have been left in English for review:\n%s", out)
	}

	reviewBytes, err := os.ReadFile(filepath.Join(outDir, "review.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(reviewBytes), "siglus_merged_or_luca_split") {
		t.Fatalf("review report does not contain merged/split rows:\n%s", string(reviewBytes))
	}

	hdBytes, err := os.ReadFile(filepath.Join(outDir, "hd_candidates.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(hdBytes), "HD exclusive epilogue") {
		t.Fatalf("HD candidate report does not contain Luca-only line:\n%s", string(hdBytes))
	}
}

func TestRunUsesSoftEnglishNormalizationForRewordedLines(t *testing.T) {
	dir := t.TempDir()
	lucaDir := filepath.Join(dir, "luca")
	siglusDir := filepath.Join(dir, "siglus")
	outDir := filepath.Join(dir, "out")
	if err := os.MkdirAll(lucaDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(siglusDir, 0755); err != nil {
		t.Fatal(err)
	}

	lucaScript := strings.Join([]string{
		`MESSAGE (0, "jp", "Cleaning.")`,
		`MESSAGE (0, "jp", "That sounded like a nice, productive thing to do.")`,
		`MESSAGE (0, "jp", "I was sure I would be helping Shiona out if I cleaned.")`,
		`MESSAGE (0, "jp", "Shiona pointed to the door with her ladle.")`,
	}, "\n")
	if err := os.WriteFile(filepath.Join(lucaDir, "clean.txt"), []byte(lucaScript), 0644); err != nil {
		t.Fatal(err)
	}

	siglusScript := strings.Join([]string{
		`○0000000001○Cleaning.`,
		`●0000000001●Menage.`,
		`○0000000002○That had a good sound to it.`,
		`●0000000002●Ca sonnait bien.`,
		`○0000000003○It sounded really helpful.`,
		`●0000000003●Ca semblait vraiment utile.`,
		`○0000000004○Shiona pointed to the door with her ladle.`,
		`●0000000004●Shiona pointa vers la porte avec sa louche.`,
	}, "\n\n")
	if err := os.WriteFile(filepath.Join(siglusDir, "clean.ss.txt"), []byte(siglusScript), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Run(Options{
		LucaDir:   lucaDir,
		SiglusDir: siglusDir,
		OutputDir: outDir,
		TargetCol: 2,
	})
	if err != nil {
		t.Fatal(err)
	}

	outBytes, err := os.ReadFile(filepath.Join(outDir, "clean.txt"))
	if err != nil {
		t.Fatal(err)
	}
	out := string(outBytes)
	if !strings.Contains(out, "Ca sonnait bien.") || !strings.Contains(out, "Ca semblait vraiment utile.") {
		t.Fatalf("expected reworded lines to be imported:\n%s", out)
	}
	if strings.Contains(out, "I was sure I would be helping Shiona out if I cleaned.") {
		t.Fatalf("reworded Luca line should not be reported as HD-only:\n%s", out)
	}

	hdBytes, err := os.ReadFile(filepath.Join(outDir, "hd_candidates.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(hdBytes), "helping Shiona") {
		t.Fatalf("reworded line was exported as HD-only:\n%s", string(hdBytes))
	}
}

func TestReplaceNthQuotedStringPreservesOriginalLineBreakSuffix(t *testing.T) {
	line := `MESSAGE (0, "jp\n", "English line.\n", "cn\n", 1, 2, 0x0)`

	got := replaceNthQuotedString(line, 1, "Ligne française.")

	if !strings.Contains(got, `"Ligne française.\n"`) {
		t.Fatalf("expected replacement to preserve \\n suffix:\n%s", got)
	}
}

func TestReplaceNthQuotedStringDoesNotAddLineBreakWithoutOriginalSuffix(t *testing.T) {
	line := `MESSAGE (1, "jp", "English line.", "cn", 1, 2, 0x0)`

	got := replaceNthQuotedString(line, 1, "Ligne française.")

	if strings.Contains(got, `"Ligne française.\n"`) {
		t.Fatalf("replacement unexpectedly gained \\n suffix:\n%s", got)
	}
}
