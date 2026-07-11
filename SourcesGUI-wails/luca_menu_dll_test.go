package main

import (
	"debug/pe"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSupportsLucaMenuDLLMatchesPlatform(t *testing.T) {
	app := NewApp()
	if got, want := app.SupportsLucaMenuDLL(), runtime.GOOS == "windows"; got != want {
		t.Fatalf("SupportsLucaMenuDLL() = %v, want %v", got, want)
	}
}

func TestFindVisualStudioDevCmd(t *testing.T) {
	if path := findVisualStudioDevCmd(); path != "" {
		if info, err := os.Stat(path); err != nil || info.IsDir() {
			t.Fatalf("Visual Studio developer command path is invalid: %q", path)
		}
	}
}

func testLucaKitDir(t *testing.T) string {
	t.Helper()
	for _, name := range lucaKitDirNames {
		kit := filepath.Join("..", name)
		if info, err := os.Stat(kit); err == nil && info.IsDir() {
			return kit
		}
	}
	t.Fatal("proxy DLL kit folder not found")
	return ""
}

func TestParseCoreLucaPatchScripts(t *testing.T) {
	kit := testLucaKitDir(t)
	profile, err := parseLucaPatchScript(kit, filepath.Join(kit, "AIR", "patches.py"))
	if err != nil {
		t.Fatalf("parse AIR patches: %v", err)
	}
	if profile.RvaDelta != "0x1800" {
		t.Fatalf("AIR RVA delta = %q, want 0x1800", profile.RvaDelta)
	}
	if len(profile.Entries) < 100 {
		t.Fatalf("AIR entries = %d, want at least 100", len(profile.Entries))
	}
	if profile.Entries[0].Source != "Close" || profile.Entries[0].Slot != "en" {
		t.Fatalf("first AIR entry = %#v", profile.Entries[0])
	}
}

func TestParseLBEEMixedEncodingProfile(t *testing.T) {
	kit := testLucaKitDir(t)
	profile, err := parseLucaPatchScript(kit, filepath.Join(kit, "LBEE", "patches.py"))
	if err != nil {
		t.Fatalf("parse LBEE patches: %v", err)
	}
	if profile.ProxyDLL != "winmm" || profile.Architecture != "x86" || profile.RvaMode != "pe" {
		t.Fatalf("LBEE build profile = proxy %q arch %q RVA %q", profile.ProxyDLL, profile.Architecture, profile.RvaMode)
	}
	if len(profile.Entries) < 40 {
		t.Fatalf("LBEE entries = %d, want at least 40", len(profile.Entries))
	}
	wide := 0
	for _, entry := range profile.Entries {
		if entry.Encoding == "utf-16-le" {
			wide++
			if entry.SourceBytes != len([]rune(entry.Source))*2 {
				t.Fatalf("LBEE UTF-16 byte length for %q = %d", entry.Source, entry.SourceBytes)
			}
		}
	}
	if wide < 7 {
		t.Fatalf("LBEE UTF-16 entries = %d, want at least 7", wide)
	}
	definition, err := os.ReadFile(filepath.Join(kit, "winmm.def"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(definition), "timeGetDevCaps") {
		t.Fatal("winmm.def must forward timeGetDevCaps for dynamically loaded graphics drivers")
	}
}

func TestLucaCommonFrenchPreset(t *testing.T) {
	kit := testLucaKitDir(t)
	inv := LucaMenuInventory{KitDir: kit}
	for _, game := range []string{"AIR", "Kanon", "HarmoniaHD", "Loopers"} {
		profile, err := parseLucaPatchScript(kit, filepath.Join(kit, game, "patches.py"))
		if err != nil {
			t.Fatalf("parse %s patches: %v", game, err)
		}
		inv.Profiles = append(inv.Profiles, profile)
	}
	catalog, err := loadLucaMenuCatalog(kit)
	if err != nil {
		t.Fatalf("load Luca menu catalog: %v", err)
	}
	annotateLucaInventory(&inv, catalog)

	var found bool
	for _, profile := range inv.Profiles {
		if profile.ID != "Kanon" {
			continue
		}
		for _, entry := range profile.Entries {
			if entry.Source == "Close" {
				found = true
				if entry.SuggestedFr != "Fermer" || !entry.SafeAuto {
					t.Fatalf("Kanon Close suggestion = %q safe=%v", entry.SuggestedFr, entry.SafeAuto)
				}
			}
		}
	}
	if !found {
		t.Fatal("Kanon Close entry not found")
	}
}

func TestParseKanonJapaneseSlotProfile(t *testing.T) {
	kit := testLucaKitDir(t)
	patch := filepath.Join(kit, "Kanon", "Arabic-B", "patches.py")
	profile, err := parseLucaPatchScript(kit, patch)
	if err != nil {
		t.Fatalf("parse Kanon Japanese slot patches: %v", err)
	}

	var japanese int
	for _, entry := range profile.Entries {
		if entry.Slot == "jp" {
			japanese++
		}
	}
	if japanese < 100 {
		t.Fatalf("Kanon Japanese slot entries = %d, want at least 100", japanese)
	}
}

func TestAllGameSlotCatalogProfiles(t *testing.T) {
	kit := testLucaKitDir(t)
	catalog, err := loadLucaMenuCatalog(kit)
	if err != nil {
		t.Fatalf("load Luca menu catalog: %v", err)
	}
	for _, game := range []string{"AIR", "HarmoniaHD", "Kanon", "Loopers"} {
		t.Run(game, func(t *testing.T) {
			inv := LucaMenuInventory{KitDir: kit}
			for _, relative := range []string{
				filepath.Join(game, "patches.py"),
				filepath.Join(game, "Slots-JP", "patches.py"),
				filepath.Join(game, "Slots-CN", "patches.py"),
			} {
				profile, err := parseLucaPatchScript(kit, filepath.Join(kit, relative))
				if err != nil {
					t.Fatalf("parse %s: %v", relative, err)
				}
				inv.Profiles = append(inv.Profiles, profile)
			}
			annotateLucaInventory(&inv, catalog)

			counts := map[string]int{}
			safeFrench := map[string]int{}
			for _, profile := range inv.Profiles {
				for _, entry := range profile.Entries {
					counts[entry.Slot]++
					if entry.SafeAuto && entry.SuggestedFr != "" {
						safeFrench[entry.Slot]++
					}
				}
			}
			if counts["jp"] < 80 || safeFrench["jp"] < 45 {
				t.Fatalf("%s JP inventory=%d safe FR=%d", game, counts["jp"], safeFrench["jp"])
			}
			if counts["cn"] < 70 || safeFrench["cn"] < 45 {
				t.Fatalf("%s CN inventory=%d safe FR=%d", game, counts["cn"], safeFrench["cn"])
			}
		})
	}
}

func TestAIRJapaneseFrenchDLLIntegration(t *testing.T) {
	if os.Getenv("LUCA_DLL_INTEGRATION") != "1" {
		t.Skip("set LUCA_DLL_INTEGRATION=1 to run against the installed AIR.exe")
	}
	gameExe := `C:\Program Files (x86)\Steam\steamapps\common\AIR\AIR.exe`
	if _, err := os.Stat(gameExe); err != nil {
		t.Skipf("AIR.exe is not installed: %v", err)
	}
	kit := testLucaKitDir(t)
	catalog, err := loadLucaMenuCatalog(kit)
	if err != nil {
		t.Fatal(err)
	}
	profile, err := parseLucaPatchScript(kit, filepath.Join(kit, "AIR", "Slots-JP", "patches.py"))
	if err != nil {
		t.Fatal(err)
	}
	inv := LucaMenuInventory{KitDir: kit, Profiles: []LucaMenuProfile{profile}}
	annotateLucaInventory(&inv, catalog)

	edits := make([]LucaMenuPatchEdit, 0, len(inv.Profiles[0].Entries))
	for _, entry := range inv.Profiles[0].Entries {
		if !entry.SafeAuto || entry.SuggestedFr == "" {
			continue
		}
		edits = append(edits, LucaMenuPatchEdit{
			RawOffset: entry.RawOffset,
			Source:    entry.Source,
			Target:    entry.SuggestedFr,
			Context:   entry.Context,
			Note:      entry.Note,
			Include:   true,
			Budget:    entry.Budget,
		})
	}
	if len(edits) < 55 {
		t.Fatalf("safe AIR JP -> FR edits = %d", len(edits))
	}

	outputDir := t.TempDir()
	req := LucaMenuGenerateRequest{
		ProfileID:     "AIR",
		GameExe:       gameExe,
		OutputDir:     outputDir,
		PatchGameName: "AIR FR JP slot integration",
		PatchVersion:  "test",
		Slot:          "jp",
		BuildDLL:      true,
		Entries:       edits,
	}
	baseProfile, err := parseLucaPatchScript(kit, filepath.Join(kit, "AIR", "patches.py"))
	if err != nil {
		t.Fatal(err)
	}
	script := buildGeneratedLucaPatchesPy(baseProfile, req, edits, req.PatchGameName, req.PatchVersion)
	if err := os.WriteFile(filepath.Join(outputDir, "patches.py"), []byte(script), 0644); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"version.c", "version.def"} {
		if err := copyFile(filepath.Join(kit, name), filepath.Join(outputDir, name)); err != nil {
			t.Fatal(err)
		}
	}
	python, args, err := findPythonCommand()
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(python, append(args, "patches.py")...)
	cmd.Dir = outputDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("patches.py failed: %v\n%s", err, output)
	}

	devCmd := findVisualStudioDevCmd()
	if devCmd == "" {
		t.Skip("Visual Studio Build Tools are not installed")
	}
	comspec := os.Getenv("ComSpec")
	if comspec == "" {
		comspec = "cmd.exe"
	}
	compile := fmt.Sprintf(
		"@echo off\r\ncall \"%s\" -no_logo -arch=amd64 -host_arch=amd64\r\nif errorlevel 1 exit /b %%errorlevel%%\r\ncl.exe /nologo /O2 /W3 /LD /D_CRT_SECURE_NO_WARNINGS /I . /Fe:version.dll version.c /link /DEF:version.def /SUBSYSTEM:WINDOWS /NOLOGO\r\n",
		devCmd,
	)
	buildScript := filepath.Join(outputDir, ".luca-build.cmd")
	if err := os.WriteFile(buildScript, []byte(compile), 0600); err != nil {
		t.Fatal(err)
	}
	cmd = exec.Command(comspec, "/d", "/c", filepath.Base(buildScript))
	cmd.Dir = outputDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("version.dll build failed: %v\n%s", err, output)
	}
	if info, err := os.Stat(filepath.Join(outputDir, "version.dll")); err != nil || info.Size() == 0 {
		t.Fatalf("version.dll was not generated: %v", err)
	}
}

func TestLBEEMixedWinMMDLLIntegration(t *testing.T) {
	if os.Getenv("LUCA_LBEE_DLL_INTEGRATION") != "1" {
		t.Skip("set LUCA_LBEE_DLL_INTEGRATION=1 and LBEE_EXE to run the x86 winmm test")
	}
	gameExe := os.Getenv("LBEE_EXE")
	if gameExe == "" {
		t.Fatal("LBEE_EXE is required")
	}
	if _, err := os.Stat(gameExe); err != nil {
		t.Fatalf("LBEE executable not found: %v", err)
	}

	kit := testLucaKitDir(t)
	profile, err := parseLucaPatchScript(kit, filepath.Join(kit, "LBEE", "patches.py"))
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := loadLucaMenuCatalog(kit)
	if err != nil {
		t.Fatal(err)
	}
	inv := LucaMenuInventory{KitDir: kit, Profiles: []LucaMenuProfile{profile}}
	annotateLucaInventory(&inv, catalog)
	profile = inv.Profiles[0]

	edits := make([]LucaMenuPatchEdit, 0, len(profile.Entries))
	for _, entry := range profile.Entries {
		if !entry.SafeAuto || entry.SuggestedFr == "" {
			continue
		}
		edits = append(edits, LucaMenuPatchEdit{
			RawOffset: entry.RawOffset,
			Source:    entry.Source,
			Target:    entry.SuggestedFr,
			Context:   entry.Context,
			Note:      entry.Note,
			Encoding:  entry.Encoding,
			Include:   true,
			Budget:    entry.Budget,
		})
	}
	if len(edits) < 25 {
		t.Fatalf("safe LBEE EN -> FR edits = %d, want at least 25", len(edits))
	}

	outputDir := t.TempDir()
	req := LucaMenuGenerateRequest{
		ProfileID:     "LBEE",
		GameExe:       gameExe,
		OutputDir:     outputDir,
		PatchGameName: "LBEE FR GUI integration",
		PatchVersion:  "test",
		Slot:          "en",
		BuildDLL:      true,
		Entries:       edits,
	}
	script := buildGeneratedLucaPatchesPy(profile, req, edits, req.PatchGameName, req.PatchVersion)
	if err := os.WriteFile(filepath.Join(outputDir, "patches.py"), []byte(script), 0644); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"version.c", "winmm.def"} {
		if err := copyFile(filepath.Join(kit, name), filepath.Join(outputDir, name)); err != nil {
			t.Fatal(err)
		}
	}
	python, args, err := findPythonCommand()
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(python, append(args, "patches.py")...)
	cmd.Dir = outputDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("LBEE patches.py failed: %v\n%s", err, output)
	}

	devCmd := findVisualStudioDevCmd()
	if devCmd == "" {
		t.Skip("Visual Studio Build Tools are not installed")
	}
	comspec := os.Getenv("ComSpec")
	if comspec == "" {
		comspec = "cmd.exe"
	}
	compile := fmt.Sprintf(
		"@echo off\r\ncall \"%s\" -no_logo -arch=x86 -host_arch=amd64\r\nif errorlevel 1 exit /b %%errorlevel%%\r\ncl.exe /nologo /O2 /W3 /LD /D_CRT_SECURE_NO_WARNINGS /DLUCKPROXY_WINMM /I . /Fe:winmm.dll version.c /link /DEF:winmm.def /SUBSYSTEM:WINDOWS /MACHINE:X86 /NOLOGO\r\n",
		devCmd,
	)
	buildScript := filepath.Join(outputDir, ".luca-build.cmd")
	if err := os.WriteFile(buildScript, []byte(compile), 0600); err != nil {
		t.Fatal(err)
	}
	cmd = exec.Command(comspec, "/d", "/c", filepath.Base(buildScript))
	cmd.Dir = outputDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("LBEE winmm.dll build failed: %v\n%s", err, output)
	}
	dllPath := filepath.Join(outputDir, "winmm.dll")
	dll, err := pe.Open(dllPath)
	if err != nil {
		t.Fatalf("open generated winmm.dll: %v", err)
	}
	defer dll.Close()
	if dll.FileHeader.Machine != pe.IMAGE_FILE_MACHINE_I386 {
		t.Fatalf("winmm.dll machine = 0x%X, want x86", dll.FileHeader.Machine)
	}
}
