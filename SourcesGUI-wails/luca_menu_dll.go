package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

type LucaMenuInventory struct {
	KitDir   string              `json:"kitDir"`
	Profiles []LucaMenuProfile   `json:"profiles"`
	Common   []LucaMenuCommonRow `json:"common"`
}

type LucaMenuProfile struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Folder        string          `json:"folder"`
	PatchFile     string          `json:"patchFile"`
	GameExe       string          `json:"gameExe"`
	RvaDelta      string          `json:"rvaDelta"`
	PatchGameName string          `json:"patchGameName"`
	PatchVersion  string          `json:"patchVersion"`
	Entries       []LucaMenuEntry `json:"entries"`
	TotalCount    int             `json:"totalCount"`
	EnglishCount  int             `json:"englishCount"`
	JapaneseCount int             `json:"japaneseCount"`
	ChineseCount  int             `json:"chineseCount"`
	SafeAutoCount int             `json:"safeAutoCount"`
}

type LucaMenuEntry struct {
	RawOffset   string   `json:"rawOffset"`
	Source      string   `json:"source"`
	Target      string   `json:"target"`
	SuggestedFr string   `json:"suggestedFr"`
	SuggestedEn string   `json:"suggestedEn"`
	SuggestedAr string   `json:"suggestedAr"`
	SuggestedJp string   `json:"suggestedJp"`
	SuggestedCn string   `json:"suggestedCn"`
	CatalogID   string   `json:"catalogId"`
	Context     string   `json:"context"`
	Note        string   `json:"note"`
	Slot        string   `json:"slot"`
	Category    string   `json:"category"`
	TextKind    string   `json:"textKind"`
	SourceBytes int      `json:"sourceBytes"`
	TargetBytes int      `json:"targetBytes"`
	Budget      int      `json:"budget"`
	CommonCount int      `json:"commonCount"`
	CommonGames []string `json:"commonGames"`
	SafeAuto    bool     `json:"safeAuto"`
	Risk        string   `json:"risk"`
	Include     bool     `json:"include"`
}

type LucaMenuCommonRow struct {
	Source      string   `json:"source"`
	SuggestedFr string   `json:"suggestedFr"`
	Games       []string `json:"games"`
	Count       int      `json:"count"`
}

type LucaMenuPatchEdit struct {
	RawOffset string `json:"rawOffset"`
	Source    string `json:"source"`
	Target    string `json:"target"`
	Context   string `json:"context"`
	Note      string `json:"note"`
	Include   bool   `json:"include"`
	Budget    int    `json:"budget"`
}

type LucaMenuGenerateRequest struct {
	ProfileID     string              `json:"profileId"`
	GameExe       string              `json:"gameExe"`
	OutputDir     string              `json:"outputDir"`
	PatchGameName string              `json:"patchGameName"`
	PatchVersion  string              `json:"patchVersion"`
	Slot          string              `json:"slot"`
	BuildDLL      bool                `json:"buildDll"`
	Entries       []LucaMenuPatchEdit `json:"entries"`
}

type parsedPatchTuple struct {
	RawOffset int64
	Source    string
	Target    string
	Context   string
	Note      string
}

type lucaMenuCatalogFile struct {
	Version int                    `json:"version"`
	Entries []lucaMenuCatalogEntry `json:"entries"`
}

type lucaMenuCatalogEntry struct {
	ID      string `json:"id"`
	Context string `json:"context"`
	En      string `json:"en"`
	Fr      string `json:"fr"`
	Ar      string `json:"ar"`
	Jp      string `json:"jp"`
	Cn      string `json:"cn"`
	Safe    bool   `json:"safe"`
}

var lucaKitDirNames = []string{"proxy dll", "LuckEngine_proxy_DLL.-KIT"}

// SupportsLucaMenuDLL keeps the Win32 version.dll workflow out of native Linux builds.
func (a *App) SupportsLucaMenuDLL() bool {
	return runtime.GOOS == "windows"
}

// ScanLucaMenuKit inventories every patches.py table in the proxy DLL resource folder.
func (a *App) ScanLucaMenuKit() LucaMenuInventory {
	inv := LucaMenuInventory{}
	if !a.SupportsLucaMenuDLL() {
		return inv
	}
	kitDir := a.findLucaKitDir()
	if kitDir == "" {
		a.logError("proxy dll folder not found next to the executables.")
		return inv
	}
	inv.KitDir = kitDir

	_ = filepath.WalkDir(kitDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || strings.ToLower(d.Name()) != "patches.py" {
			return nil
		}
		if samePath(filepath.Dir(path), kitDir) {
			return nil
		}
		profile, err := parseLucaPatchScript(kitDir, path)
		if err != nil {
			a.log(fmt.Sprintf("[LUCA] skipped %s: %v", path, err))
			return nil
		}
		if len(profile.Entries) > 0 {
			inv.Profiles = append(inv.Profiles, profile)
		}
		return nil
	})

	catalog, err := loadLucaMenuCatalog(kitDir)
	if err != nil {
		a.log(fmt.Sprintf("[LUCA] menu catalog skipped: %v", err))
	}
	annotateLucaInventory(&inv, catalog)
	sort.Slice(inv.Profiles, func(i, j int) bool {
		return lucaProfileSortKey(inv.Profiles[i].ID) < lucaProfileSortKey(inv.Profiles[j].ID)
	})
	return inv
}

func (a *App) LucaMenuGenerate(req LucaMenuGenerateRequest) string {
	if !a.SupportsLucaMenuDLL() {
		a.logError("Luca DLL: the version.dll hook is available only in the Windows GUI.")
		return "ERROR"
	}
	if req.ProfileID == "" || req.GameExe == "" || req.OutputDir == "" {
		a.logError("Luca DLL: game profile, EXE and output folder are required.")
		return "ERROR"
	}
	gameExe := strings.TrimSpace(req.GameExe)
	if !filepath.IsAbs(gameExe) {
		gameExe = filepath.Join(req.OutputDir, gameExe)
	}
	gameExe, _ = filepath.Abs(gameExe)
	if info, err := os.Stat(gameExe); err != nil || info.IsDir() {
		a.logError(fmt.Sprintf("Luca DLL: game EXE not found: %s", gameExe))
		return "ERROR"
	}
	req.GameExe = gameExe
	if len(req.Entries) == 0 {
		a.logError("Luca DLL: no patch entry selected.")
		return "ERROR"
	}

	inv := a.ScanLucaMenuKit()
	var profile *LucaMenuProfile
	for i := range inv.Profiles {
		if inv.Profiles[i].ID == req.ProfileID {
			profile = &inv.Profiles[i]
			break
		}
	}
	if profile == nil {
		a.logError("Luca DLL: selected profile was not found in the kit.")
		return "ERROR"
	}

	selected := make([]LucaMenuPatchEdit, 0, len(req.Entries))
	seen := map[string]bool{}
	for _, e := range req.Entries {
		if !e.Include || strings.TrimSpace(e.Target) == "" {
			continue
		}
		key := strings.ToUpper(strings.TrimSpace(e.RawOffset))
		if seen[key] {
			continue
		}
		seen[key] = true
		if e.Budget >= 0 && len([]byte(e.Target)) > e.Budget {
			a.logError(fmt.Sprintf("Luca DLL: target too long at %s (%d bytes > budget %d): %s", e.RawOffset, len([]byte(e.Target)), e.Budget, e.Target))
			return "ERROR"
		}
		selected = append(selected, e)
	}
	if len(selected) == 0 {
		a.logError("Luca DLL: selected entries are empty.")
		return "ERROR"
	}

	if err := os.MkdirAll(req.OutputDir, 0755); err != nil {
		a.logError(fmt.Sprintf("Luca DLL: cannot create output folder: %v", err))
		return "ERROR"
	}

	patchName := strings.TrimSpace(req.PatchGameName)
	if patchName == "" {
		patchName = profile.PatchGameName
	}
	if patchName == "" {
		patchName = profile.Name
	}
	patchVersion := strings.TrimSpace(req.PatchVersion)
	if patchVersion == "" {
		patchVersion = profile.PatchVersion
	}
	if patchVersion == "" {
		patchVersion = "0.1-gui"
	}

	script := buildGeneratedLucaPatchesPy(*profile, req, selected, patchName, patchVersion)
	patchPath := filepath.Join(req.OutputDir, "patches.py")
	if err := os.WriteFile(patchPath, []byte(script), 0644); err != nil {
		a.logError(fmt.Sprintf("Luca DLL: cannot write patches.py: %v", err))
		return "ERROR"
	}
	for _, name := range []string{"version.c", "version.def", "Makefile"} {
		if err := copyFile(filepath.Join(inv.KitDir, name), filepath.Join(req.OutputDir, name)); err != nil {
			a.logError(fmt.Sprintf("Luca DLL: cannot copy %s from proxy dll folder: %v", name, err))
			return "ERROR"
		}
	}

	a.log("════════════════════════════════════════")
	a.log("  LUCA MENU DLL")
	a.log("════════════════════════════════════════")
	a.log(fmt.Sprintf("Profile: %s", profile.Name))
	a.log(fmt.Sprintf("Slot:    %s", req.Slot))
	a.log(fmt.Sprintf("Entries: %d", len(selected)))
	a.log(fmt.Sprintf("Output:  %s", req.OutputDir))

	python, args, err := findPythonCommand()
	if err != nil {
		a.logError(err.Error())
		return "ERROR"
	}
	args = append(args, "patches.py")
	if err := a.runLucaCommand(req.OutputDir, python, args...); err != nil {
		a.logError(fmt.Sprintf("Luca DLL: patches.py failed: %v", err))
		return "ERROR"
	}

	if req.BuildDLL {
		if err := a.buildLucaVersionDLL(req.OutputDir); err != nil {
			a.logError(err.Error())
			a.log("Generated patches.py, patches.h and patches.csv are ready; install a Windows C compiler to build version.dll.")
			return "ERROR"
		}
		a.logOK(fmt.Sprintf("version.dll generated: %s", filepath.Join(req.OutputDir, "version.dll")))
	} else {
		a.logOK("patches.py, patches.h and patches.csv generated")
	}
	a.log("════════════════════════════════════════")
	return "OK"
}

func (a *App) findLucaKitDir() string {
	var seeds []string
	if a.lucksystem != "" {
		seeds = append(seeds, filepath.Dir(a.lucksystem))
	}
	if exe, err := os.Executable(); err == nil {
		seeds = append(seeds, filepath.Dir(exe))
	}
	if cwd, err := os.Getwd(); err == nil {
		seeds = append(seeds, cwd)
	}

	seen := map[string]bool{}
	for _, seed := range seeds {
		dir, _ := filepath.Abs(seed)
		for {
			if dir == "" {
				break
			}
			for _, name := range lucaKitDirNames {
				candidate := filepath.Join(dir, name)
				key := strings.ToLower(filepath.Clean(candidate))
				if seen[key] {
					continue
				}
				seen[key] = true
				if info, err := os.Stat(candidate); err == nil && info.IsDir() {
					return candidate
				}
			}
			next := filepath.Dir(dir)
			if next == dir {
				break
			}
			dir = next
		}
	}
	return ""
}

func parseLucaPatchScript(kitDir, patchPath string) (LucaMenuProfile, error) {
	data, err := os.ReadFile(patchPath)
	if err != nil {
		return LucaMenuProfile{}, err
	}
	text := string(data)
	rel, _ := filepath.Rel(kitDir, filepath.Dir(patchPath))
	rel = filepath.ToSlash(rel)

	profile := LucaMenuProfile{
		ID:        rel,
		Folder:    filepath.Dir(patchPath),
		PatchFile: patchPath,
		GameExe:   extractPyAssignmentString(text, "GAME_EXE"),
		RvaDelta:  extractPyAssignmentExpr(text, "RVA_DELTA"),
	}
	profile.PatchGameName = extractPyAssignmentString(text, "PATCH_GAME_NAME")
	profile.PatchVersion = extractPyAssignmentString(text, "PATCH_VERSION")
	if profile.GameExe == "" || profile.RvaDelta == "" || profile.PatchGameName == "" || profile.PatchVersion == "" {
		parentPatch := filepath.Join(filepath.Dir(filepath.Dir(patchPath)), "patches.py")
		if parentData, err := os.ReadFile(parentPatch); err == nil && !samePath(parentPatch, patchPath) {
			parentText := string(parentData)
			if profile.GameExe == "" {
				profile.GameExe = extractPyAssignmentString(parentText, "GAME_EXE")
			}
			if profile.RvaDelta == "" {
				profile.RvaDelta = extractPyAssignmentExpr(parentText, "RVA_DELTA")
			}
			if profile.PatchGameName == "" {
				profile.PatchGameName = extractPyAssignmentString(parentText, "PATCH_GAME_NAME")
			}
			if profile.PatchVersion == "" {
				profile.PatchVersion = extractPyAssignmentString(parentText, "PATCH_VERSION")
			}
		}
	}
	profile.Name = profile.PatchGameName
	if profile.Name == "" {
		profile.Name = strings.ReplaceAll(rel, "/", " / ")
	}

	tuples := parsePatchTuples(text)
	if len(tuples) == 0 {
		return profile, fmt.Errorf("no literal patch tuples found")
	}
	for _, t := range tuples {
		entry := LucaMenuEntry{
			RawOffset:   fmt.Sprintf("0x%X", t.RawOffset),
			Source:      t.Source,
			Target:      t.Target,
			Context:     t.Context,
			Note:        t.Note,
			Slot:        detectLucaSlot(t.Source, t.Note),
			Category:    detectLucaCategory(t.Context),
			TextKind:    detectLucaTextKind(t.Context, t.Note, t.Source),
			SourceBytes: len([]byte(t.Source)),
			TargetBytes: len([]byte(t.Target)),
			Budget:      parseBudget(t.Note),
			Include:     false,
		}
		entry.Risk = detectLucaRisk(entry)
		profile.Entries = append(profile.Entries, entry)
	}
	return profile, nil
}

func parsePatchTuples(text string) []parsedPatchTuple {
	assignRe := regexp.MustCompile(`(?m)(PATCHES|[A-Za-z_][A-Za-z0-9_]*PATCHES)\s*=\s*\[`)
	matches := assignRe.FindAllStringSubmatchIndex(text, -1)
	var tuples []parsedPatchTuple
	for _, m := range matches {
		start := m[1] - 1
		if start < 0 || start >= len(text) || text[start] != '[' {
			continue
		}
		end := findPythonClosing(text, start, '[', ']')
		if end < 0 {
			continue
		}
		body := text[start+1 : end]
		for _, expr := range splitTopLevelTuples(body) {
			parts := splitTopLevelCSV(expr)
			if len(parts) < 5 {
				continue
			}
			off, err := parseIntLiteral(parts[0])
			if err != nil {
				continue
			}
			src, ok := parsePythonValue(parts[1])
			if !ok {
				continue
			}
			target, ok := parsePythonValue(parts[2])
			if !ok {
				continue
			}
			context, ok := parsePythonValue(parts[3])
			if !ok {
				continue
			}
			note, ok := parsePythonValue(parts[4])
			if !ok {
				continue
			}
			tuples = append(tuples, parsedPatchTuple{
				RawOffset: off,
				Source:    src,
				Target:    target,
				Context:   context,
				Note:      note,
			})
		}
	}
	return tuples
}

func loadLucaMenuCatalog(kitDir string) (map[string]lucaMenuCatalogEntry, error) {
	data, err := os.ReadFile(filepath.Join(kitDir, "menu_catalog.json"))
	if err != nil {
		return nil, err
	}
	var file lucaMenuCatalogFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, err
	}
	entries := make(map[string]lucaMenuCatalogEntry, len(file.Entries))
	for _, entry := range file.Entries {
		if entry.ID != "" {
			entries[entry.ID] = entry
		}
	}
	return entries, nil
}

func catalogIDFromNote(note string) string {
	for _, part := range strings.Split(note, ";") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "catalog=") {
			return strings.TrimSpace(strings.TrimPrefix(part, "catalog="))
		}
	}
	return ""
}

func annotateLucaInventory(inv *LucaMenuInventory, catalog map[string]lucaMenuCatalogEntry) {
	type sourceInfo struct {
		games map[string]bool
		fr    []string
	}
	coreIDs := map[string]bool{"AIR": true, "HarmoniaHD": true, "Kanon": true, "Loopers": true}
	infos := map[string]*sourceInfo{}
	for _, profile := range inv.Profiles {
		core := strings.Split(profile.ID, "/")[0]
		if !coreIDs[core] {
			continue
		}
		for _, e := range profile.Entries {
			if e.Slot != "en" {
				continue
			}
			key := normalizeLucaSource(e.Source)
			if infos[key] == nil {
				infos[key] = &sourceInfo{games: map[string]bool{}}
			}
			infos[key].games[core] = true
			if strings.TrimSpace(e.Target) != "" {
				infos[key].fr = append(infos[key].fr, e.Target)
			}
		}
	}

	commonRows := make([]LucaMenuCommonRow, 0, len(infos))
	for key, info := range infos {
		games := sortedMapKeys(info.games)
		row := LucaMenuCommonRow{
			Source:      key,
			SuggestedFr: chooseShortestNonEmpty(info.fr),
			Games:       games,
			Count:       len(games),
		}
		commonRows = append(commonRows, row)
	}
	sort.Slice(commonRows, func(i, j int) bool {
		if commonRows[i].Count != commonRows[j].Count {
			return commonRows[i].Count > commonRows[j].Count
		}
		return commonRows[i].Source < commonRows[j].Source
	})
	inv.Common = commonRows

	commonBySource := map[string]LucaMenuCommonRow{}
	for _, row := range commonRows {
		commonBySource[normalizeLucaSource(row.Source)] = row
	}
	catalogByEn := map[string]lucaMenuCatalogEntry{}
	catalogByJp := map[string]lucaMenuCatalogEntry{}
	catalogByCn := map[string]lucaMenuCatalogEntry{}
	for _, row := range catalog {
		if row.En != "" {
			catalogByEn[normalizeLucaSource(row.En)] = row
		}
		if row.Jp != "" {
			catalogByJp[normalizeLucaSource(row.Jp)] = row
		}
		if row.Cn != "" {
			catalogByCn[normalizeLucaSource(row.Cn)] = row
		}
	}

	for pi := range inv.Profiles {
		p := &inv.Profiles[pi]
		p.TotalCount = len(p.Entries)
		for ei := range p.Entries {
			e := &p.Entries[ei]
			row := commonBySource[normalizeLucaSource(e.Source)]
			e.CommonCount = row.Count
			e.CommonGames = row.Games
			e.SuggestedFr = row.SuggestedFr

			catalogID := catalogIDFromNote(e.Note)
			catalogRow, found := catalog[catalogID]
			if !found {
				switch e.Slot {
				case "jp":
					catalogRow, found = catalogByJp[normalizeLucaSource(e.Source)]
				case "cn":
					catalogRow, found = catalogByCn[normalizeLucaSource(e.Source)]
				default:
					catalogRow, found = catalogByEn[normalizeLucaSource(e.Source)]
				}
			}
			if found {
				e.CatalogID = catalogRow.ID
				e.SuggestedFr = catalogRow.Fr
				e.SuggestedEn = catalogRow.En
				e.SuggestedAr = catalogRow.Ar
				e.SuggestedJp = catalogRow.Jp
				e.SuggestedCn = catalogRow.Cn
				e.CommonCount = 4
				e.CommonGames = []string{"AIR", "HarmoniaHD", "Kanon", "Loopers"}
			}
			if e.SuggestedEn == "" && e.Slot == "en" {
				e.SuggestedEn = e.Source
			}
			if e.SuggestedJp == "" && e.Slot == "jp" {
				e.SuggestedJp = e.Source
			}
			if e.SuggestedCn == "" && e.Slot == "cn" {
				e.SuggestedCn = e.Source
			}
			isSafe := e.SuggestedFr != "" && e.SuggestedFr != e.Source && e.CommonCount >= 4
			if found {
				isSafe = isSafe && catalogRow.Safe
			}
			if isSafe {
				e.SafeAuto = e.Budget < 0 || len([]byte(e.SuggestedFr)) <= e.Budget
			}
			if e.SafeAuto {
				p.SafeAutoCount++
			}
			switch e.Slot {
			case "jp":
				p.JapaneseCount++
			case "cn":
				p.ChineseCount++
			default:
				p.EnglishCount++
			}
		}
	}
}

func buildGeneratedLucaPatchesPy(profile LucaMenuProfile, req LucaMenuGenerateRequest, entries []LucaMenuPatchEdit, patchName, patchVersion string) string {
	var b strings.Builder
	b.WriteString("#!/usr/bin/env python3\n")
	b.WriteString("# -*- coding: utf-8 -*-\n")
	b.WriteString("\"\"\"\n")
	b.WriteString("Generated by LuckSystem GUI - Luca Menu DLL.\n")
	b.WriteString("Source profile: " + profile.ID + "\n")
	b.WriteString("Slot: " + req.Slot + "\n")
	b.WriteString("\"\"\"\n\n")
	b.WriteString("from pathlib import Path\n")
	b.WriteString("import sys\n\n")
	b.WriteString("GAME_EXE = " + pyStringLiteral(req.GameExe) + "\n")
	b.WriteString("RVA_DELTA = " + strings.TrimSpace(profile.RvaDelta) + "\n")
	b.WriteString("PATCH_GAME_NAME = " + pyStringLiteral(patchName) + "\n")
	b.WriteString("PATCH_VERSION = " + pyStringLiteral(patchVersion) + "\n\n")
	b.WriteString("PATCHES = [\n")
	for _, e := range entries {
		off, _ := parseIntLiteral(e.RawOffset)
		note := strings.TrimSpace(e.Note)
		guiNote := "GUI slot=" + req.Slot
		if note == "" {
			note = guiNote
		} else if !strings.Contains(note, "GUI slot=") {
			note += "; " + guiNote
		}
		b.WriteString(fmt.Sprintf("    (0x%X, %s, %s, %s, %s),\n",
			off,
			pyBytesLiteral(e.Source),
			pyStringLiteral(e.Target),
			pyStringLiteral(e.Context),
			pyStringLiteral(note),
		))
	}
	b.WriteString("]\n\n")
	b.WriteString(generatedLucaPythonTail())
	return b.String()
}

func generatedLucaPythonTail() string {
	return `def main():
    if hasattr(sys.stdout, 'reconfigure'):
        sys.stdout.reconfigure(encoding='utf-8')
        sys.stderr.reconfigure(encoding='utf-8')

    data = Path(GAME_EXE).read_bytes()

    def slot_size(start):
        i = start
        while data[i] != 0:
            i += 1
        while i < len(data) and data[i] == 0:
            i += 1
        return i - start

    rows = []
    errors = []
    for off, src, target, context, note in PATCHES:
        actual = data[off:off + len(src)]
        if actual != src:
            errors.append(f"0x{off:X}: expected {src!r}, got {actual!r}")
            continue
        target_bytes = target.encode('utf-8')
        slot = slot_size(off)
        budget = slot - 1
        fits = len(target_bytes) <= budget
        rows.append({
            'off': off,
            'src': src.decode('utf-8', errors='replace'),
            'target': target,
            'src_len': len(src),
            'target_len': len(target_bytes),
            'slot': slot,
            'budget': budget,
            'fits': fits,
            'context': context,
            'note': note,
            'src_bytes': src,
            'target_bytes': target_bytes,
        })

    if errors:
        print('=== OFFSET MISMATCH (aborting) ===', file=sys.stderr)
        for e in errors:
            print('  ' + e, file=sys.stderr)
        sys.exit(1)

    print(f"{'off':>8}  {'slot':>4}  {'src':>3}  {'tgt':>3}  {'fit':3}  src -> target")
    print('-' * 100)
    n_ok = n_bad = 0
    for r in rows:
        mark = 'OK' if r['fits'] else 'NO'
        if r['fits']:
            n_ok += 1
        else:
            n_bad += 1
        print(f"0x{r['off']:06X}  {r['slot']:>4}  {r['src_len']:>3}  {r['target_len']:>3}  {mark:3}  {r['src']!r} -> {r['target']!r}")
    print(f"\nTotal: {len(rows)}  OK: {n_ok}  FAIL: {n_bad}")
    if n_bad:
        print('\nFailures (target too long):')
        for r in rows:
            if not r['fits']:
                print(f"  0x{r['off']:X}: target {r['target_len']}B > budget {r['budget']}B: {r['target']!r}")
        sys.exit(2)

    with open('patches.h', 'w', encoding='utf-8') as f:
        f.write('/* Auto-generated from patches.py. Do not edit. */\n')
        f.write('#ifndef LUCKPROXY_PATCHES_H\n#define LUCKPROXY_PATCHES_H\n\n')
        f.write(f'#define PATCH_GAME_NAME "{PATCH_GAME_NAME}"\n')
        f.write(f'#define PATCH_VERSION   "{PATCH_VERSION}"\n\n')
        for i, r in enumerate(rows):
            write_len = max(len(r['src_bytes']), len(r['target_bytes'])) + 1
            src_padded = list(r['src_bytes']) + [0] * (write_len - len(r['src_bytes']))
            target_padded = list(r['target_bytes']) + [0] * (write_len - len(r['target_bytes']))
            src_arr = ','.join(f'0x{x:02X}' for x in src_padded)
            target_arr = ','.join(f'0x{x:02X}' for x in target_padded)
            f.write(f'static const BYTE s_src_{i:03d}[] = {{ {src_arr} }};\n')
            f.write(f'static const BYTE s_tgt_{i:03d}[] = {{ {target_arr} }};\n')
        f.write('\nstatic const LuckPatch g_patches[] = {\n')
        for i, r in enumerate(rows):
            rva = r['off'] + RVA_DELTA
            write_len = max(len(r['src_bytes']), len(r['target_bytes'])) + 1
            ctx = r['context'] + ': ' + r['src'][:30]
            ctx = (ctx.replace('\\', '\\\\')
                      .replace('"', '\\"')
                      .replace('\n', '\\n')
                      .replace('\r', '\\r')
                      .replace('\t', '\\t'))
            f.write(f'    {{ 0x{rva:06X}, {write_len:>4}, s_src_{i:03d}, s_tgt_{i:03d}, "{ctx}" }},\n')
        f.write('};\n\n#define N_PATCHES (sizeof(g_patches)/sizeof(g_patches[0]))\n')
        f.write('\n#endif\n')
    print(f"\nGenerated patches.h with {len(rows)} entries.")

    with open('patches.csv', 'w', encoding='utf-8') as f:
        f.write('raw_offset,rva,slot,budget,src_len,target_len,fits,src,target,context,note\n')
        for r in rows:
            rva = r['off'] + RVA_DELTA
            def esc(s):
                return '"' + s.replace('"', '""') + '"'
            f.write(
                f'0x{r["off"]:X},0x{rva:X},{r["slot"]},{r["budget"]},'
                f'{r["src_len"]},{r["target_len"]},{r["fits"]},'
                f'{esc(r["src"])},{esc(r["target"])},{esc(r["context"])},{esc(r["note"])}\n'
            )
    print(f"Generated patches.csv with {len(rows)} entries.")


if __name__ == '__main__':
    main()
`
}

func (a *App) runLucaCommand(workdir, name string, args ...string) error {
	a.log(fmt.Sprintf("> %s %s", filepath.Base(name), strings.Join(args, " ")))
	ctx, cancel := context.WithCancel(a.ctx)
	a.mu.Lock()
	a.cancelFunc = cancel
	a.mu.Unlock()
	defer func() {
		cancel()
		a.mu.Lock()
		a.cancelFunc = nil
		a.mu.Unlock()
	}()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = workdir
	hideWindow(cmd)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	done := make(chan struct{}, 2)
	stream := func(r io.Reader) {
		scanner := bufio.NewScanner(r)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 1024*1024)
		for scanner.Scan() {
			a.log(scanner.Text())
		}
		done <- struct{}{}
	}
	go stream(stdout)
	go stream(stderr)
	<-done
	<-done

	if err := cmd.Wait(); err != nil {
		if ctx.Err() != nil {
			a.log("[STOPPED] Process cancelled by user.")
			return fmt.Errorf("cancelled")
		}
		return err
	}
	return nil
}

func (a *App) buildLucaVersionDLL(outputDir string) error {
	gcc, err := exec.LookPath("x86_64-w64-mingw32-gcc")
	if err != nil {
		gcc, err = exec.LookPath("gcc")
	}
	if err == nil {
		args := []string{
			"-O2", "-Wall", "-Wextra", "-Wno-unused-parameter",
			"-I.", "-shared", "-s", "-static-libgcc",
			"-Wl,--subsystem,windows,--enable-stdcall-fixup",
			"-o", "version.dll", "version.c", "version.def",
		}
		return a.runLucaCommand(outputDir, gcc, args...)
	}

	cl, err := exec.LookPath("cl")
	if err == nil {
		args := []string{
			"/nologo", "/O2", "/W3", "/LD", "/I", ".",
			"/Fe:version.dll", "version.c",
			"/link", "/DEF:version.def", "/SUBSYSTEM:WINDOWS", "/NOLOGO",
		}
		return a.runLucaCommand(outputDir, cl, args...)
	}

	if devCmd := findVisualStudioDevCmd(); devCmd != "" {
		comspec := os.Getenv("ComSpec")
		if comspec == "" {
			comspec = "cmd.exe"
		}
		buildScript := filepath.Join(outputDir, ".luca-build.cmd")
		command := fmt.Sprintf(
			"@echo off\r\ncall \"%s\" -no_logo -arch=amd64 -host_arch=amd64\r\nif errorlevel 1 exit /b %%errorlevel%%\r\ncl.exe /nologo /O2 /W3 /LD /D_CRT_SECURE_NO_WARNINGS /I . /Fe:version.dll version.c /link /DEF:version.def /SUBSYSTEM:WINDOWS /NOLOGO\r\n",
			devCmd,
		)
		if err := os.WriteFile(buildScript, []byte(command), 0600); err != nil {
			return fmt.Errorf("Luca DLL: cannot prepare Visual Studio build: %w", err)
		}
		defer os.Remove(buildScript)
		if err := a.runLucaCommand(outputDir, comspec, "/d", "/c", filepath.Base(buildScript)); err != nil {
			return err
		}
		for _, name := range []string{"version.obj", "version.lib", "version.exp"} {
			_ = os.Remove(filepath.Join(outputDir, name))
		}
		return nil
	}

	return fmt.Errorf("Luca DLL: no C compiler found (MinGW GCC or Visual Studio Build Tools)")
}

func findVisualStudioDevCmd() string {
	programFilesX86 := os.Getenv("ProgramFiles(x86)")
	if programFilesX86 == "" {
		programFilesX86 = `C:\Program Files (x86)`
	}
	vswhere := filepath.Join(programFilesX86, "Microsoft Visual Studio", "Installer", "vswhere.exe")
	if info, err := os.Stat(vswhere); err == nil && !info.IsDir() {
		cmd := exec.Command(vswhere,
			"-latest", "-products", "*",
			"-requires", "Microsoft.VisualStudio.Component.VC.Tools.x86.x64",
			"-property", "installationPath",
		)
		hideWindow(cmd)
		if output, err := cmd.Output(); err == nil {
			installDir := strings.TrimSpace(string(output))
			candidate := filepath.Join(installDir, "Common7", "Tools", "VsDevCmd.bat")
			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				return candidate
			}
		}
	}

	patterns := []string{
		filepath.Join(programFilesX86, "Microsoft Visual Studio", "*", "*", "Common7", "Tools", "VsDevCmd.bat"),
	}
	for _, pattern := range patterns {
		matches, _ := filepath.Glob(pattern)
		sort.Sort(sort.Reverse(sort.StringSlice(matches)))
		for _, candidate := range matches {
			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				return candidate
			}
		}
	}
	return ""
}

func findPythonCommand() (string, []string, error) {
	for _, spec := range [][]string{
		{"py", "-3"},
		{"python"},
		{"python3"},
	} {
		path, err := exec.LookPath(spec[0])
		if err == nil {
			return path, spec[1:], nil
		}
	}
	return "", nil, fmt.Errorf("Python was not found in PATH.")
}

func extractPyAssignmentString(text, name string) string {
	expr := extractPyAssignmentExpr(text, name)
	value, ok := parsePythonValue(expr)
	if !ok {
		return ""
	}
	return value
}

func extractPyAssignmentExpr(text, name string) string {
	re := regexp.MustCompile(`(?m)^\s*(?:base\.)?` + regexp.QuoteMeta(name) + `\s*=\s*(.+)$`)
	m := re.FindStringSubmatch(text)
	if len(m) < 2 {
		return ""
	}
	return strings.TrimSpace(m[1])
}

func splitTopLevelTuples(body string) []string {
	var out []string
	for i := 0; i < len(body); i++ {
		c := body[i]
		if c != '(' {
			continue
		}
		end := findPythonClosing(body, i, '(', ')')
		if end < 0 {
			continue
		}
		out = append(out, strings.TrimSpace(body[i+1:end]))
		i = end
	}
	return out
}

func splitTopLevelCSV(expr string) []string {
	var parts []string
	start := 0
	depth := 0
	inString := false
	var quote byte
	triple := false
	escape := false
	comment := false
	for i := 0; i < len(expr); i++ {
		c := expr[i]
		if comment {
			if c == '\n' {
				comment = false
			}
			continue
		}
		if inString {
			if escape {
				escape = false
				continue
			}
			if c == '\\' {
				escape = true
				continue
			}
			if c == quote {
				if triple {
					if i+2 < len(expr) && expr[i+1] == quote && expr[i+2] == quote {
						inString = false
						i += 2
					}
				} else {
					inString = false
				}
			}
			continue
		}
		switch c {
		case '#':
			comment = true
		case '\'', '"':
			inString = true
			quote = c
			triple = i+2 < len(expr) && expr[i+1] == c && expr[i+2] == c
			if triple {
				i += 2
			}
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				parts = append(parts, strings.TrimSpace(expr[start:i]))
				start = i + 1
			}
		}
	}
	if start <= len(expr) {
		last := strings.TrimSpace(expr[start:])
		if last != "" {
			parts = append(parts, last)
		}
	}
	return parts
}

func findPythonClosing(text string, start int, open, close byte) int {
	depth := 0
	inString := false
	var quote byte
	triple := false
	escape := false
	comment := false
	for i := start; i < len(text); i++ {
		c := text[i]
		if comment {
			if c == '\n' {
				comment = false
			}
			continue
		}
		if inString {
			if escape {
				escape = false
				continue
			}
			if c == '\\' {
				escape = true
				continue
			}
			if c == quote {
				if triple {
					if i+2 < len(text) && text[i+1] == quote && text[i+2] == quote {
						inString = false
						i += 2
					}
				} else {
					inString = false
				}
			}
			continue
		}
		switch c {
		case '#':
			comment = true
		case '\'', '"':
			inString = true
			quote = c
			triple = i+2 < len(text) && text[i+1] == c && text[i+2] == c
			if triple {
				i += 2
			}
		default:
			if c == open {
				depth++
			} else if c == close {
				depth--
				if depth == 0 {
					return i
				}
			}
		}
	}
	return -1
}

func parsePythonValue(expr string) (string, bool) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return "", false
	}
	if idx := strings.Index(expr, ".encode"); idx > 0 {
		return parsePythonValue(expr[:idx])
	}
	if strings.HasPrefix(expr, "J(") {
		end := findPythonClosing(expr, 1, '(', ')')
		if end > 1 {
			parts := splitTopLevelCSV(expr[2:end])
			if len(parts) > 0 {
				return parsePythonValue(parts[0])
			}
		}
	}
	return parsePythonQuotedLiteral(expr)
}

func parsePythonQuotedLiteral(expr string) (string, bool) {
	expr = strings.TrimSpace(expr)
	raw := false
	i := 0
	for i < len(expr) {
		c := expr[i]
		if c == 'r' || c == 'R' {
			raw = true
			i++
			continue
		}
		if c == 'u' || c == 'U' || c == 'b' || c == 'B' || c == 'f' || c == 'F' {
			i++
			continue
		}
		break
	}
	if i >= len(expr) || (expr[i] != '\'' && expr[i] != '"') {
		return "", false
	}
	quote := expr[i]
	triple := i+2 < len(expr) && expr[i+1] == quote && expr[i+2] == quote
	contentStart := i + 1
	if triple {
		contentStart = i + 3
	}
	contentEnd := -1
	escape := false
	for j := contentStart; j < len(expr); j++ {
		c := expr[j]
		if escape {
			escape = false
			continue
		}
		if c == '\\' && !raw {
			escape = true
			continue
		}
		if c == quote {
			if triple {
				if j+2 < len(expr) && expr[j+1] == quote && expr[j+2] == quote {
					contentEnd = j
					break
				}
			} else {
				contentEnd = j
				break
			}
		}
	}
	if contentEnd < 0 {
		return "", false
	}
	content := expr[contentStart:contentEnd]
	if raw {
		return content, true
	}
	return unescapePythonString(content), true
}

func unescapePythonString(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '\\' || i+1 >= len(s) {
			b.WriteByte(c)
			continue
		}
		i++
		switch s[i] {
		case 'n':
			b.WriteByte('\n')
		case 'r':
			b.WriteByte('\r')
		case 't':
			b.WriteByte('\t')
		case 'b':
			b.WriteByte('\b')
		case 'f':
			b.WriteByte('\f')
		case 'a':
			b.WriteByte('\a')
		case 'v':
			b.WriteByte('\v')
		case '\\', '\'', '"':
			b.WriteByte(s[i])
		case '\n':
			// Line continuation.
		case 'x':
			if i+2 < len(s) {
				if decoded, err := hex.DecodeString(s[i+1 : i+3]); err == nil && len(decoded) == 1 {
					b.WriteByte(decoded[0])
					i += 2
				}
			}
		case 'u':
			if i+4 < len(s) {
				if r, err := strconv.ParseInt(s[i+1:i+5], 16, 32); err == nil {
					b.WriteRune(rune(r))
					i += 4
				}
			}
		case 'U':
			if i+8 < len(s) {
				if r, err := strconv.ParseInt(s[i+1:i+9], 16, 32); err == nil {
					b.WriteRune(rune(r))
					i += 8
				}
			}
		default:
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

func parseIntLiteral(expr string) (int64, error) {
	expr = strings.TrimSpace(expr)
	expr = strings.TrimSuffix(expr, ",")
	expr = strings.ReplaceAll(expr, "_", "")
	return strconv.ParseInt(expr, 0, 64)
}

func detectLucaSlot(source, note string) string {
	lower := strings.ToLower(note)
	if strings.Contains(lower, "japanese slot") || strings.Contains(lower, "jp source slot") || strings.Contains(lower, "slot=jp") || hasJapanese(source) {
		return "jp"
	}
	if strings.Contains(lower, "chinese") || strings.Contains(lower, "cn source slot") || strings.Contains(lower, "slot=cn") || hasChineseOnly(source) {
		return "cn"
	}
	return "en"
}

func hasJapanese(s string) bool {
	for _, r := range s {
		if (r >= 0x3040 && r <= 0x30FF) || (r >= 0x31F0 && r <= 0x31FF) {
			return true
		}
	}
	return false
}

func hasChineseOnly(s string) bool {
	hasHan := false
	for _, r := range s {
		if r >= 0x3040 && r <= 0x30FF {
			return false
		}
		if r >= 0x4E00 && r <= 0x9FFF {
			hasHan = true
		}
	}
	return hasHan
}

func detectLucaCategory(context string) string {
	if idx := strings.Index(context, "/"); idx > 0 {
		return context[:idx]
	}
	lower := strings.ToLower(context)
	for _, cat := range []string{"basic", "text1", "text2", "sound", "keyboard", "mouse", "system", "save", "touch"} {
		if strings.Contains(lower, cat) {
			return strings.Title(cat)
		}
	}
	return "Global"
}

func detectLucaTextKind(context, note, source string) string {
	lower := strings.ToLower(context + " " + note)
	switch {
	case strings.Contains(lower, "tooltip"):
		return "tooltip"
	case strings.Contains(lower, "prompt"):
		return "prompt"
	case strings.Contains(lower, "button"):
		return "button"
	case strings.Contains(lower, "tab"):
		return "tab"
	case strings.Contains(lower, "value"):
		return "value"
	case len([]byte(source)) <= 7:
		return "short"
	default:
		return "label"
	}
}

func detectLucaRisk(e LucaMenuEntry) string {
	var risks []string
	if e.Budget >= 0 && e.Budget <= 3 {
		risks = append(risks, "tiny slot")
	}
	if e.Source == "Yes" || e.Source == "No" {
		risks = append(risks, "global dialog")
	}
	if strings.Contains(e.Source, "%") {
		risks = append(risks, "format tokens")
	}
	if e.Target != "" && !utf8.ValidString(e.Target) {
		risks = append(risks, "invalid utf-8")
	}
	return strings.Join(risks, ", ")
}

func parseBudget(note string) int {
	re := regexp.MustCompile(`(?i)budget\s+(\d+)`)
	m := re.FindStringSubmatch(note)
	if len(m) < 2 {
		return -1
	}
	n, err := strconv.Atoi(m[1])
	if err != nil {
		return -1
	}
	return n
}

func normalizeLucaSource(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, "\r\n", "\n"))
}

func chooseShortestNonEmpty(values []string) string {
	best := ""
	for _, v := range values {
		if strings.TrimSpace(v) == "" {
			continue
		}
		if best == "" || len([]byte(v)) < len([]byte(best)) {
			best = v
		}
	}
	return best
}

func sortedMapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func pyStringLiteral(s string) string {
	repl := strings.NewReplacer(
		"\\", "\\\\",
		"'", "\\'",
		"\n", "\\n",
		"\r", "\\r",
		"\t", "\\t",
	)
	return "'" + repl.Replace(s) + "'"
}

func pyBytesLiteral(s string) string {
	var b strings.Builder
	b.WriteString("b'")
	for _, c := range []byte(s) {
		if c >= 0x20 && c <= 0x7E && c != '\\' && c != '\'' {
			b.WriteByte(c)
		} else {
			b.WriteString(fmt.Sprintf("\\x%02x", c))
		}
	}
	b.WriteString("'")
	return b.String()
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

func samePath(a, b string) bool {
	aa, _ := filepath.Abs(a)
	bb, _ := filepath.Abs(b)
	return strings.EqualFold(filepath.Clean(aa), filepath.Clean(bb))
}

func lucaProfileSortKey(id string) string {
	order := map[string]string{
		"Kanon":      "0",
		"AIR":        "1",
		"HarmoniaHD": "2",
		"Loopers":    "3",
	}
	root := strings.Split(id, "/")[0]
	if prefix, ok := order[root]; ok {
		return prefix + ":" + id
	}
	return "9:" + id
}
