package main

// ═══════════════════════════════════════════════════════════════════
// dialogue.go — Extract / Import MESSAGE dialogues from script .txt
// Supports two MESSAGE formats:
//   AIR-type:  MESSAGE ("JAP", "ENG", "CN")       → 3 strings
//   LB_EN-type: MESSAGE (voiceId, "JAP", "ENG")   → number + 2 strings
// Auto-detects format from the first MESSAGE line in each file.
// ═══════════════════════════════════════════════════════════════════

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// ───────────────────────────────────────
// MESSAGE arg parser (handles nested quotes, escapes)
// ───────────────────────────────────────

func splitMessageArgs(s string) []string {
	var args []string
	var buf []byte
	inStr := false
	esc := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if esc {
			buf = append(buf, ch)
			esc = false
			continue
		}
		if ch == '\\' {
			buf = append(buf, ch)
			esc = true
			continue
		}
		if ch == '"' {
			buf = append(buf, ch)
			inStr = !inStr
			continue
		}
		if ch == ',' && !inStr {
			part := strings.TrimSpace(string(buf))
			if part != "" {
				args = append(args, part)
			}
			buf = buf[:0]
			continue
		}
		buf = append(buf, ch)
	}
	tail := strings.TrimSpace(string(buf))
	if tail != "" {
		args = append(args, tail)
	}
	return args
}

// unquoteArg strips surrounding quotes and processes escape sequences
func unquoteArg(token string) string {
	token = strings.TrimSpace(token)
	if len(token) >= 2 && token[0] == '"' && token[len(token)-1] == '"' {
		body := token[1 : len(token)-1]
		body = strings.ReplaceAll(body, `\\`, "\x00BSLASH\x00")
		body = strings.ReplaceAll(body, `\"`, `"`)
		body = strings.ReplaceAll(body, `\n`, "\n")
		body = strings.ReplaceAll(body, `\t`, "\t")
		body = strings.ReplaceAll(body, `\r`, "\r")
		body = strings.ReplaceAll(body, "\x00BSLASH\x00", `\`)
		return body
	}
	return token
}

// quoteArg wraps a string in quotes with proper escaping
func quoteArg(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	return `"` + s + `"`
}

// ───────────────────────────────────────
// Format detection
// ───────────────────────────────────────

type messageFormat int

const (
	formatUnknown messageFormat = iota
	formatAIR                   // MESSAGE ("JAP", "ENG", "CN") — 3 strings
	formatLBEN                  // MESSAGE (voiceId, "JAP", "ENG") — number + 2 strings
)

func (f messageFormat) String() string {
	switch f {
	case formatAIR:
		return "AIR-type (JAP, ENG, CN)"
	case formatLBEN:
		return "LB_EN-type (voiceId, JAP, ENG)"
	default:
		return "Unknown"
	}
}

// detectFormat checks the first arg of a MESSAGE line
func detectFormat(firstArg string) messageFormat {
	firstArg = strings.TrimSpace(firstArg)
	if len(firstArg) == 0 {
		return formatUnknown
	}
	if firstArg[0] == '"' {
		return formatAIR
	}
	// Check if it's a number (voiceId)
	for _, r := range firstArg {
		if !unicode.IsDigit(r) && r != '-' {
			return formatUnknown
		}
	}
	return formatLBEN
}

// Language field indices by format
// Returns: map[languageLabel] -> argIndex
func langIndices(fmt messageFormat) map[string]int {
	switch fmt {
	case formatAIR:
		return map[string]int{"JAP": 0, "ENG": 1, "CN": 2}
	case formatLBEN:
		return map[string]int{"JAP": 1, "ENG": 2}
	default:
		return nil
	}
}

// availableLanguages returns ordered language list for a format
func availableLanguages(fmt messageFormat) []string {
	switch fmt {
	case formatAIR:
		return []string{"JAP", "ENG", "CN"}
	case formatLBEN:
		return []string{"JAP", "ENG"}
	default:
		return nil
	}
}

// ───────────────────────────────────────
// Parse a MESSAGE line → returns args (raw tokens)
// ───────────────────────────────────────

type messageLine struct {
	lineNo  int      // 1-based line number in file
	prefix  string   // everything before the opening '(' (incl. "MESSAGE ")
	args    []string // raw token strings (quoted or number)
	suffix  string   // everything after the closing ')'
	format  messageFormat
}

func parseMessageLine(line string, lineNo int) *messageLine {
	idx := strings.Index(line, "MESSAGE")
	if idx < 0 {
		return nil
	}
	lp := strings.Index(line[idx:], "(")
	if lp < 0 {
		return nil
	}
	lp += idx
	rp := strings.LastIndex(line, ")")
	if rp < 0 || rp <= lp {
		return nil
	}
	inside := strings.TrimSpace(line[lp+1 : rp])
	if inside == "" {
		return nil
	}
	args := splitMessageArgs(inside)
	if len(args) < 2 {
		return nil
	}
	f := detectFormat(args[0])
	if f == formatUnknown {
		return nil
	}
	return &messageLine{
		lineNo: lineNo,
		prefix: line[:lp+1],
		args:   args,
		suffix: line[rp:],
		format: f,
	}
}

// ───────────────────────────────────────
// EXTRACT — single file
// ───────────────────────────────────────

type extractRow struct {
	file   string
	lineNo int
	idx    int
	texts  map[string]string // lang -> text
}

func extractFile(filePath string, languages []string) ([]extractRow, messageFormat, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, formatUnknown, err
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")

	var rows []extractRow
	var detectedFmt messageFormat
	msgIdx := 0
	fileName := filepath.Base(filePath)

	for lineNo, line := range lines {
		if !strings.Contains(line, "MESSAGE") {
			continue
		}
		ml := parseMessageLine(line, lineNo+1) // 1-based
		if ml == nil {
			continue
		}

		// Detect format from first MESSAGE
		if detectedFmt == formatUnknown {
			detectedFmt = ml.format
		}

		indices := langIndices(ml.format)
		if indices == nil {
			continue
		}

		msgIdx++
		row := extractRow{
			file:   fileName,
			lineNo: lineNo + 1,
			idx:    msgIdx,
			texts:  make(map[string]string),
		}
		for _, lang := range languages {
			argIdx, ok := indices[lang]
			if ok && argIdx < len(ml.args) {
				row.texts[lang] = unquoteArg(ml.args[argIdx])
			}
		}
		rows = append(rows, row)
	}
	return rows, detectedFmt, nil
}

// writeExtractTSV writes extraction results to a TSV file
func writeExtractTSV(outPath string, rows []extractRow, languages []string) error {
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = '\t'

	// Header
	header := []string{"file", "line", "idx"}
	header = append(header, languages...)
	if err := w.Write(header); err != nil {
		return err
	}

	// Rows
	for _, row := range rows {
		record := []string{row.file, strconv.Itoa(row.lineNo), strconv.Itoa(row.idx)}
		for _, lang := range languages {
			record = append(record, row.texts[lang])
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

// ───────────────────────────────────────
// IMPORT — single file
// ───────────────────────────────────────

// loadImportMap reads a TSV and builds: lineNo -> map[lang]text
func loadImportMap(tsvPath string, targetFile string) (map[int]map[string]string, []string, error) {
	f, err := os.Open(tsvPath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = '\t'
	r.LazyQuotes = true

	header, err := r.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot read TSV header: %v", err)
	}

	// Find column indices
	colIdx := make(map[string]int)
	for i, h := range header {
		colIdx[h] = i
	}
	lineCol, ok := colIdx["line"]
	if !ok {
		return nil, nil, fmt.Errorf("TSV missing 'line' column")
	}
	fileCol, hasFile := colIdx["file"]

	// Determine which language columns exist
	var langs []string
	for _, lang := range []string{"JAP", "ENG", "CN"} {
		if _, ok := colIdx[lang]; ok {
			langs = append(langs, lang)
		}
	}

	mp := make(map[int]map[string]string)
	for {
		record, err := r.Read()
		if err != nil {
			break
		}
		// Filter by filename if present
		if hasFile && fileCol < len(record) {
			if record[fileCol] != targetFile {
				continue
			}
		}
		if lineCol >= len(record) {
			continue
		}
		lineNo, err := strconv.Atoi(record[lineCol])
		if err != nil {
			continue
		}
		texts := make(map[string]string)
		for _, lang := range langs {
			idx := colIdx[lang]
			if idx < len(record) {
				texts[lang] = record[idx]
			}
		}
		mp[lineNo] = texts
	}
	return mp, langs, nil
}

// importFile applies TSV translations back into a script file
func importFile(scriptPath string, tsvMap map[int]map[string]string, importLang string, outputPath string) (int, error) {
	data, err := os.ReadFile(scriptPath)
	if err != nil {
		return 0, err
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")

	changed := 0
	for lineNo := range lines {
		lineNum := lineNo + 1 // 1-based
		texts, ok := tsvMap[lineNum]
		if !ok {
			continue
		}
		newText, ok := texts[importLang]
		if !ok || strings.TrimSpace(newText) == "" {
			continue
		}

		ml := parseMessageLine(lines[lineNo], lineNum)
		if ml == nil {
			continue
		}

		indices := langIndices(ml.format)
		if indices == nil {
			continue
		}
		argIdx, ok := indices[importLang]
		if !ok || argIdx >= len(ml.args) {
			continue
		}

		// Check if text actually changed
		oldText := unquoteArg(ml.args[argIdx])
		if oldText == newText {
			continue
		}

		// Rebuild the line with the modified arg
		ml.args[argIdx] = quoteArg(newText)
		// Rebuild: non-string args (voiceId) stay as-is, string args get quoted
		var rebuiltArgs []string
		for _, a := range ml.args {
			// Already quoted or a number — keep as-is
			rebuiltArgs = append(rebuiltArgs, a)
		}
		lines[lineNo] = ml.prefix + strings.Join(rebuiltArgs, ", ") + ml.suffix
		changed++
	}

	err = os.WriteFile(outputPath, []byte(strings.Join(lines, "\n")), 0644)
	return changed, err
}

// ═══════════════════════════════════════
// Wails-exposed methods
// ═══════════════════════════════════════

// DialogueDetectFormat scans a .txt file and returns the detected format + available languages
func (a *App) DialogueDetectFormat(filePath string) map[string]interface{} {
	result := map[string]interface{}{
		"format":    "Unknown",
		"languages": []string{},
	}
	if filePath == "" {
		return result
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return result
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if !strings.Contains(line, "MESSAGE") {
			continue
		}
		ml := parseMessageLine(line, 0)
		if ml != nil {
			result["format"] = ml.format.String()
			result["languages"] = availableLanguages(ml.format)
			return result
		}
	}
	return result
}

// DialogueExtractFile extracts MESSAGE lines from a single .txt file
func (a *App) DialogueExtractFile(inputFile string, outputFile string, languages []string) string {
	if inputFile == "" || outputFile == "" || len(languages) == 0 {
		a.logError("Input file, output file, and at least one language are required")
		return "ERROR"
	}

	a.log("════════════════════════════════════════")
	a.log("  DIALOGUE EXTRACT (single file)")
	a.log("════════════════════════════════════════")
	a.log(fmt.Sprintf("Input:     %s", inputFile))
	a.log(fmt.Sprintf("Output:    %s", outputFile))
	a.log(fmt.Sprintf("Languages: %s", strings.Join(languages, ", ")))
	a.log("────────────────────────────────────────")

	rows, detectedFmt, err := extractFile(inputFile, languages)
	if err != nil {
		a.logError(fmt.Sprintf("Read error: %v", err))
		return "ERROR"
	}
	a.log(fmt.Sprintf("Format detected: %s", detectedFmt))

	if len(rows) == 0 {
		a.logError("No MESSAGE lines found in file")
		return "ERROR"
	}

	if err := writeExtractTSV(outputFile, rows, languages); err != nil {
		a.logError(fmt.Sprintf("Write error: %v", err))
		return "ERROR"
	}

	a.logOK(fmt.Sprintf("Extracted %d MESSAGE lines → %s", len(rows), filepath.Base(outputFile)))
	a.log("════════════════════════════════════════")
	return "OK"
}

// DialogueExtractBatch extracts MESSAGE lines from all .txt files in a directory
func (a *App) DialogueExtractBatch(inputDir string, outputDir string, languages []string) string {
	if inputDir == "" || outputDir == "" || len(languages) == 0 {
		a.logError("Input directory, output directory, and at least one language are required")
		return "ERROR"
	}

	a.log("════════════════════════════════════════")
	a.log("  DIALOGUE EXTRACT (batch)")
	a.log("════════════════════════════════════════")
	a.log(fmt.Sprintf("Input:     %s", inputDir))
	a.log(fmt.Sprintf("Output:    %s", outputDir))
	a.log(fmt.Sprintf("Languages: %s", strings.Join(languages, ", ")))
	a.log("────────────────────────────────────────")

	os.MkdirAll(outputDir, os.ModePerm)

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		a.logError(fmt.Sprintf("Cannot read directory: %v", err))
		return "ERROR"
	}

	totalFiles := 0
	totalMessages := 0
	errors := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".txt" {
			continue
		}

		inPath := filepath.Join(inputDir, name)
		outName := strings.TrimSuffix(name, ext) + ".ext.txt"
		outPath := filepath.Join(outputDir, outName)

		rows, detectedFmt, err := extractFile(inPath, languages)
		if err != nil {
			a.logError(fmt.Sprintf("  [ERR] %s: %v", name, err))
			errors++
			continue
		}
		if len(rows) == 0 {
			continue // skip files without MESSAGE
		}

		if totalFiles == 0 {
			a.log(fmt.Sprintf("Format detected: %s", detectedFmt))
		}

		if err := writeExtractTSV(outPath, rows, languages); err != nil {
			a.logError(fmt.Sprintf("  [ERR] %s: %v", name, err))
			errors++
			continue
		}

		a.log(fmt.Sprintf("  [%d] %s → %d messages", totalFiles+1, name, len(rows)))
		totalFiles++
		totalMessages += len(rows)
	}

	result := fmt.Sprintf("%d files processed, %d messages total, %d errors", totalFiles, totalMessages, errors)
	a.logOK(result)
	a.log("════════════════════════════════════════")
	return "OK"
}

// DialogueImportFile reinjects translated TSV into a single script file
func (a *App) DialogueImportFile(scriptFile string, tsvFile string, targetLang string, outputFile string) string {
	if scriptFile == "" || tsvFile == "" || targetLang == "" || outputFile == "" {
		a.logError("Script file, TSV file, target language, and output file are required")
		return "ERROR"
	}

	a.log("════════════════════════════════════════")
	a.log("  DIALOGUE IMPORT (single file)")
	a.log("════════════════════════════════════════")
	a.log(fmt.Sprintf("Script:  %s", scriptFile))
	a.log(fmt.Sprintf("TSV:     %s", tsvFile))
	a.log(fmt.Sprintf("Target:  %s", targetLang))
	a.log(fmt.Sprintf("Output:  %s", outputFile))
	a.log("────────────────────────────────────────")

	fileName := filepath.Base(scriptFile)
	tsvMap, langs, err := loadImportMap(tsvFile, fileName)
	if err != nil {
		a.logError(fmt.Sprintf("TSV read error: %v", err))
		return "ERROR"
	}
	a.log(fmt.Sprintf("TSV columns: %s | %d entries for %s", strings.Join(langs, ", "), len(tsvMap), fileName))

	changed, err := importFile(scriptFile, tsvMap, targetLang, outputFile)
	if err != nil {
		a.logError(fmt.Sprintf("Write error: %v", err))
		return "ERROR"
	}

	a.logOK(fmt.Sprintf("%d MESSAGE lines updated → %s", changed, filepath.Base(outputFile)))
	a.log("════════════════════════════════════════")
	return "OK"
}

// DialogueImportBatch reinjects translated TSV files into all matching scripts
func (a *App) DialogueImportBatch(scriptDir string, tsvDir string, targetLang string, outputDir string) string {
	if scriptDir == "" || tsvDir == "" || targetLang == "" || outputDir == "" {
		a.logError("Script directory, TSV directory, target language, and output directory are required")
		return "ERROR"
	}

	a.log("════════════════════════════════════════")
	a.log("  DIALOGUE IMPORT (batch)")
	a.log("════════════════════════════════════════")
	a.log(fmt.Sprintf("Scripts: %s", scriptDir))
	a.log(fmt.Sprintf("TSV:     %s", tsvDir))
	a.log(fmt.Sprintf("Target:  %s", targetLang))
	a.log(fmt.Sprintf("Output:  %s", outputDir))
	a.log("────────────────────────────────────────")

	os.MkdirAll(outputDir, os.ModePerm)

	entries, err := os.ReadDir(tsvDir)
	if err != nil {
		a.logError(fmt.Sprintf("Cannot read TSV directory: %v", err))
		return "ERROR"
	}

	totalFiles := 0
	totalChanged := 0
	errors := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".ext.txt") {
			continue
		}

		tsvPath := filepath.Join(tsvDir, name)

		// Derive script filename: SEEN8736.ext.txt → SEEN8736.txt
		scriptName := strings.TrimSuffix(name, ".ext.txt") + ".txt"
		scriptPath := filepath.Join(scriptDir, scriptName)
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			a.log(fmt.Sprintf("  [SKIP] %s (no matching script: %s)", name, scriptName))
			continue
		}

		outputPath := filepath.Join(outputDir, scriptName)

		// Load TSV for this file
		tsvMap, _, err := loadImportMap(tsvPath, scriptName)
		if err != nil {
			a.logError(fmt.Sprintf("  [ERR] %s: %v", name, err))
			errors++
			continue
		}
		if len(tsvMap) == 0 {
			continue
		}

		changed, err := importFile(scriptPath, tsvMap, targetLang, outputPath)
		if err != nil {
			a.logError(fmt.Sprintf("  [ERR] %s: %v", name, err))
			errors++
			continue
		}

		a.log(fmt.Sprintf("  [%d] %s → %d lines updated", totalFiles+1, scriptName, changed))
		totalFiles++
		totalChanged += changed
	}

	result := fmt.Sprintf("%d files processed, %d lines updated, %d errors", totalFiles, totalChanged, errors)
	a.logOK(result)
	a.log("════════════════════════════════════════")
	return "OK"
}

// ───────────────────────────────────────
// File dialogs for dialogue operations
// ───────────────────────────────────────

func (a *App) SelectScriptTxtFile() string {
	file, _ := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select script .txt file",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "Text Files (*.txt)", Pattern: "*.txt"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
	return file
}

func (a *App) SelectTsvFile() string {
	file, _ := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select dialogue TSV file",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "TSV/Text Files (*.txt)", Pattern: "*.txt"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
	return file
}

func (a *App) SelectSaveTsvFile(defaultName string) string {
	file, _ := wailsRuntime.SaveFileDialog(a.ctx, wailsRuntime.SaveDialogOptions{
		Title:           "Save dialogue TSV",
		DefaultFilename: defaultName,
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "TSV/Text Files (*.txt)", Pattern: "*.txt"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
	return file
}

func (a *App) SelectSaveScriptFile(defaultName string) string {
	file, _ := wailsRuntime.SaveFileDialog(a.ctx, wailsRuntime.SaveDialogOptions{
		Title:           "Save modified script",
		DefaultFilename: defaultName,
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "Text Files (*.txt)", Pattern: "*.txt"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
	return file
}
