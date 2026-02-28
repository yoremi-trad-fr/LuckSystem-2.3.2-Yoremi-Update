package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct - GUI backend, calls lucksystem.exe via subprocess
type App struct {
	ctx        context.Context
	lucksystem string // path to lucksystem.exe
	mu         sync.Mutex
	cancelFunc context.CancelFunc // cancels the running subprocess
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.findLuckSystem()
}

// ───────────────────────────────────────
// Find lucksystem executable
// ───────────────────────────────────────

func (a *App) findLuckSystem() {
	// 1. Same directory as this GUI executable
	exePath, err := os.Executable()
	if err == nil {
		candidate := filepath.Join(filepath.Dir(exePath), "lucksystem.exe")
		if _, err := os.Stat(candidate); err == nil {
			a.lucksystem = candidate
			return
		}
		// Also try without .exe (Linux/Mac)
		candidate = filepath.Join(filepath.Dir(exePath), "lucksystem")
		if _, err := os.Stat(candidate); err == nil {
			a.lucksystem = candidate
			return
		}
	}

	// 2. Current working directory
	cwd, err := os.Getwd()
	if err == nil {
		candidate := filepath.Join(cwd, "lucksystem.exe")
		if _, err := os.Stat(candidate); err == nil {
			a.lucksystem = candidate
			return
		}
		candidate = filepath.Join(cwd, "lucksystem")
		if _, err := os.Stat(candidate); err == nil {
			a.lucksystem = candidate
			return
		}
	}

	// 3. In PATH
	path, err := exec.LookPath("lucksystem")
	if err == nil {
		a.lucksystem = path
		return
	}
	path, err = exec.LookPath("lucksystem.exe")
	if err == nil {
		a.lucksystem = path
		return
	}

	a.lucksystem = "" // not found
}

// GetLuckSystemPath returns the detected path (for UI display)
func (a *App) GetLuckSystemPath() string {
	return a.lucksystem
}

// SetLuckSystemPath allows the user to manually set the path
func (a *App) SetLuckSystemPath() string {
	file, err := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Locate lucksystem executable",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "Executable", Pattern: "*.exe;lucksystem"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
	if err != nil || file == "" {
		return a.lucksystem
	}
	a.lucksystem = file
	return a.lucksystem
}

// ───────────────────────────────────────
// Console Logging
// ───────────────────────────────────────

func (a *App) log(msg string) {
	wailsRuntime.EventsEmit(a.ctx, "log", msg)
}

func (a *App) logError(msg string) {
	wailsRuntime.EventsEmit(a.ctx, "log", "[ERROR] "+msg)
}

func (a *App) logOK(msg string) {
	wailsRuntime.EventsEmit(a.ctx, "log", "[OK] "+msg)
}

// ───────────────────────────────────────
// Run lucksystem subprocess
// ───────────────────────────────────────

// runLuckSystem executes lucksystem with given arguments, streaming output to console
func (a *App) runLuckSystem(args ...string) error {
	if a.lucksystem == "" {
		a.logError("lucksystem.exe not found! Place it next to the GUI or use Settings to locate it.")
		return fmt.Errorf("lucksystem not found")
	}

	// Log the command being executed
	a.log(fmt.Sprintf("> %s %s", filepath.Base(a.lucksystem), strings.Join(args, " ")))

	// Create a cancellable context for this subprocess
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

	cmd := exec.CommandContext(ctx, a.lucksystem, args...)

	// Hide the CMD window on Windows (no console popup during batch operations)
	hideWindow(cmd)

	// Capture stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		a.logError(fmt.Sprintf("stdout pipe: %v", err))
		return err
	}

	// Capture stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		a.logError(fmt.Sprintf("stderr pipe: %v", err))
		return err
	}

	// Start the process
	if err := cmd.Start(); err != nil {
		a.logError(fmt.Sprintf("Failed to start: %v", err))
		return err
	}

	// Stream stdout/stderr with batched logging for performance
	done := make(chan struct{}, 2)

	streamLines := func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 1024*1024)
		for scanner.Scan() {
			a.log(scanner.Text())
		}
		done <- struct{}{}
	}

	go streamLines(stdout)
	go streamLines(stderr)

	// Wait for both goroutines
	<-done
	<-done

	// Wait for process to finish
	if err := cmd.Wait(); err != nil {
		if ctx.Err() != nil {
			// Cancelled by user — not an error
			a.log("[STOPPED] Process cancelled by user.")
			return fmt.Errorf("cancelled")
		}
		a.logError(fmt.Sprintf("Process exited with error: %v", err))
		return err
	}

	return nil
}

// StopProcess cancels the currently running subprocess (called from frontend)
func (a *App) StopProcess() {
	a.mu.Lock()
	cancel := a.cancelFunc
	a.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

// ───────────────────────────────────────
// File Dialogs (generic)
// ───────────────────────────────────────

func (a *App) SelectPakFile() string {
	file, _ := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select .PAK file",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "PAK Files (*.PAK)", Pattern: "*.PAK;*.pak"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
	return file
}

func (a *App) SelectFile(title string, pattern string, desc string) string {
	file, _ := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: title,
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: desc, Pattern: pattern},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
	return file
}

func (a *App) SelectDirectory(title string) string {
	dir, _ := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: title,
	})
	return dir
}

func (a *App) SelectSaveFile(title string, defaultName string, pattern string, desc string) string {
	file, _ := wailsRuntime.SaveFileDialog(a.ctx, wailsRuntime.SaveDialogOptions{
		Title:           title,
		DefaultFilename: defaultName,
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: desc, Pattern: pattern},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
	return file
}

// SelectOutputPath opens a save dialog that accepts any filename (no extension enforcement).
// Returns the full path as typed by the user. Uses OpenFile dialog in directory mode
// so Windows cannot reject the filename for missing extension.
func (a *App) SelectOutputDir(title string) string {
	dir, _ := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: title,
	})
	return dir
}

// stripScriptPakSuffix removes trailing "SCRIPT.PAK" or "SCRIPT.pak" from a directory path.
// game.go automatically appends /SCRIPT.PAK/ to the import/export path,
// so if the user already selected a SCRIPT.PAK folder, it would be doubled.
func stripScriptPakSuffix(dir string) string {
	base := strings.ToUpper(filepath.Base(dir))
	if base == "SCRIPT.PAK" {
		return filepath.Dir(dir)
	}
	return dir
}

// ═══════════════════════════════════════
// SCRIPT DECOMPILE
// ═══════════════════════════════════════
// lucksystem script decompile -s PAK -c charset -O opcode -p plugin -o outputdir

func (a *App) ScriptDecompile(pakFile, opcodeFile, pluginFile, charsetStr, outputDir string) string {
	if pakFile == "" || outputDir == "" {
		a.logError("SCRIPT.PAK and output directory are required")
		return "ERROR"
	}
	if charsetStr == "" {
		charsetStr = "UTF-8"
	}

	// Strip trailing SCRIPT.PAK from output dir (game.go adds it automatically)
	outputDir = stripScriptPakSuffix(outputDir)

	a.log("════════════════════════════════════════")
	a.log("  SCRIPT DECOMPILE")
	a.log("════════════════════════════════════════")

	args := []string{"script", "decompile", "-s", pakFile, "-c", charsetStr, "-o", outputDir}
	if opcodeFile != "" {
		args = append(args, "-O", opcodeFile)
	}
	if pluginFile != "" {
		args = append(args, "-p", pluginFile)
	}

	err := a.runLuckSystem(args...)
	if err != nil {
		return "ERROR"
	}

	a.logOK("Script decompile completed")
	a.log("════════════════════════════════════════")
	return "OK"
}

// ═══════════════════════════════════════
// SCRIPT COMPILE (IMPORT)
// ═══════════════════════════════════════
// lucksystem script import -s PAK -c charset -O opcode -p plugin -i importdir -o output.PAK

func (a *App) ScriptCompile(pakFile, opcodeFile, pluginFile, charsetStr, importDir, outputPak string) string {
	if pakFile == "" || importDir == "" || outputPak == "" {
		a.logError("SCRIPT.PAK, translated folder, and output PAK are required")
		return "ERROR"
	}
	if charsetStr == "" {
		charsetStr = "UTF-8"
	}

	// Strip trailing SCRIPT.PAK from import dir (game.go adds it automatically)
	importDir = stripScriptPakSuffix(importDir)

	a.log("════════════════════════════════════════")
	a.log("  SCRIPT COMPILE (IMPORT)")
	a.log("════════════════════════════════════════")

	args := []string{"script", "import", "-s", pakFile, "-c", charsetStr, "-i", importDir, "-o", outputPak}
	if opcodeFile != "" {
		args = append(args, "-O", opcodeFile)
	}
	if pluginFile != "" {
		args = append(args, "-p", pluginFile)
	}

	err := a.runLuckSystem(args...)
	if err != nil {
		return "ERROR"
	}

	a.logOK(fmt.Sprintf("Script compile completed -> %s", outputPak))
	a.log("════════════════════════════════════════")
	return "OK"
}

// ═══════════════════════════════════════
// PAK EXTRACT
// ═══════════════════════════════════════
// lucksystem pak extract -i PAK -o listfile --all outputdir -c charset

func (a *App) PakExtract(pakFile, outputDir string) string {
	if pakFile == "" || outputDir == "" {
		a.logError("PAK file and output directory are required")
		return "ERROR"
	}

	a.log("════════════════════════════════════════")
	a.log("  PAK EXTRACT")
	a.log("════════════════════════════════════════")

	os.MkdirAll(outputDir, os.ModePerm)

	// Nom du fichier liste = <NomDuPak>_list.txt  (ex: SYSCG.PAK → SYSCG_list.txt)
	pakBase := strings.TrimSuffix(filepath.Base(pakFile), filepath.Ext(pakFile))
	listFile := filepath.Join(outputDir, pakBase+"_list.txt")
	a.log(fmt.Sprintf("List file: %s", listFile))

	args := []string{"pak", "extract", "-i", pakFile, "-o", listFile, "--all", outputDir}

	err := a.runLuckSystem(args...)
	if err != nil {
		return "ERROR"
	}

	a.logOK(fmt.Sprintf("PAK extracted to %s", outputDir))
	a.log("════════════════════════════════════════")
	return "OK"
}

// ═══════════════════════════════════════
// PAK REPLACE
// ═══════════════════════════════════════
// Mode dossier : lucksystem pak replace -s PAK -i inputdir -o output.PAK
// Mode liste   : lucksystem pak replace -s PAK -i listfile -l -o output.PAK

func (a *App) PakReplace(pakSource, inputDir, listFile, outputPak string) string {
	if pakSource == "" || outputPak == "" {
		a.logError("Original PAK and output PAK are required")
		return "ERROR"
	}

	// Exactly one of inputDir or listFile must be set
	useList := listFile != ""
	useDir := inputDir != ""
	if !useList && !useDir {
		a.logError("Provide either a folder or a list file as input")
		return "ERROR"
	}

	a.log("════════════════════════════════════════")
	a.log("  PAK REPLACE")
	a.log("════════════════════════════════════════")

	var args []string
	if useList {
		a.log(fmt.Sprintf("Mode: list file → %s", listFile))
		args = []string{"pak", "replace", "-s", pakSource, "-i", listFile, "-l", "-o", outputPak}
	} else {
		a.log(fmt.Sprintf("Mode: directory → %s", inputDir))
		args = []string{"pak", "replace", "-s", pakSource, "-i", inputDir, "-o", outputPak}
	}

	err := a.runLuckSystem(args...)
	if err != nil {
		return "ERROR"
	}

	a.logOK(fmt.Sprintf("PAK written -> %s", outputPak))
	a.log("════════════════════════════════════════")
	return "OK"
}

// ═══════════════════════════════════════
// PAK FONT EXTRACT
// ═══════════════════════════════════════
// lucksystem pak extract -i PAK -o listfile --all outputdir -c charset

func (a *App) PakFontExtract(pakFile, charsetStr, outputDir string) string {
	if pakFile == "" || outputDir == "" {
		a.logError("PAK file and output directory are required")
		return "ERROR"
	}
	if charsetStr == "" {
		charsetStr = "UTF-8"
	}

	a.log("════════════════════════════════════════")
	a.log("  PAK (FONT) EXTRACT")
	a.log("════════════════════════════════════════")

	os.MkdirAll(outputDir, os.ModePerm)

	pakBase := strings.TrimSuffix(filepath.Base(pakFile), filepath.Ext(pakFile))
	listFile := filepath.Join(outputDir, pakBase+"_list.txt")
	a.log(fmt.Sprintf("List file: %s", listFile))

	args := []string{"pak", "extract", "-i", pakFile, "-o", listFile, "--all", outputDir, "-c", charsetStr}

	err := a.runLuckSystem(args...)
	if err != nil {
		return "ERROR"
	}

	a.logOK(fmt.Sprintf("PAK (Font) extracted to %s", outputDir))
	a.log("════════════════════════════════════════")
	return "OK"
}

// ═══════════════════════════════════════
// PAK FONT REPLACE
// ═══════════════════════════════════════
// Mode liste   : lucksystem pak replace -s PAK -i listfile --list -o output.PAK -c charset
// Mode dossier : lucksystem pak replace -s PAK -i inputdir  -o output.PAK -c charset

func (a *App) PakFontReplace(pakSource, charsetStr, inputDir, listFile, outputPak string) string {
	if pakSource == "" || outputPak == "" {
		a.logError("Original PAK and output PAK are required")
		return "ERROR"
	}
	useList := listFile != ""
	useDir := inputDir != ""
	if !useList && !useDir {
		a.logError("Provide either a list file or a folder as input")
		return "ERROR"
	}
	if charsetStr == "" {
		charsetStr = "UTF-8"
	}

	a.log("════════════════════════════════════════")
	a.log("  PAK (FONT) REPLACE")
	a.log("════════════════════════════════════════")

	var args []string
	if useList {
		a.log(fmt.Sprintf("Mode: list file → %s", listFile))
		args = []string{"pak", "replace", "-s", pakSource, "-i", listFile, "-l", "-o", outputPak, "-c", charsetStr}
	} else {
		a.log(fmt.Sprintf("Mode: directory → %s", inputDir))
		args = []string{"pak", "replace", "-s", pakSource, "-i", inputDir, "-o", outputPak, "-c", charsetStr}
	}

	err := a.runLuckSystem(args...)
	if err != nil {
		return "ERROR"
	}

	a.logOK(fmt.Sprintf("PAK (Font) written -> %s", outputPak))
	a.log("════════════════════════════════════════")
	return "OK"
}

// ═══════════════════════════════════════
// FONT EXTRACT
// ═══════════════════════════════════════
// lucksystem font extract -s czfile -S infofile -o output.png -O charset.txt

func (a *App) FontExtract(czFile, infoFile, outputPng, outputCharset string) string {
	if czFile == "" || infoFile == "" || outputPng == "" {
		a.logError("Font CZ file, info file, and output PNG are required")
		return "ERROR"
	}

	a.log("════════════════════════════════════════")
	a.log("  FONT EXTRACT")
	a.log("════════════════════════════════════════")

	args := []string{"font", "extract", "-s", czFile, "-S", infoFile, "-o", outputPng}
	if outputCharset != "" {
		args = append(args, "-O", outputCharset)
	}

	err := a.runLuckSystem(args...)
	if err != nil {
		return "ERROR"
	}

	a.logOK(fmt.Sprintf("Font extracted -> %s", outputPng))
	a.log("════════════════════════════════════════")
	return "OK"
}

// ═══════════════════════════════════════
// FONT EDIT
// ═══════════════════════════════════════
// lucksystem font edit -s cz -S info -f ttf -o outcz -O outinfo [-r] [-a] [-i idx] [-c charset]

func (a *App) FontEdit(czFile, infoFile, ttfFile, outputCz, outputInfo, charsetFile string, redraw, appendMode bool, startIndex int) string {
	if czFile == "" || infoFile == "" || ttfFile == "" || outputCz == "" {
		a.logError("Font CZ, info, TTF, and output CZ are required")
		return "ERROR"
	}

	a.log("════════════════════════════════════════")
	a.log("  FONT EDIT")
	a.log("════════════════════════════════════════")

	args := []string{"font", "edit", "-s", czFile, "-S", infoFile, "-f", ttfFile, "-o", outputCz}
	if outputInfo != "" {
		args = append(args, "-O", outputInfo)
	}
	if charsetFile != "" {
		args = append(args, "-c", charsetFile)
	}
	if redraw {
		args = append(args, "-r")
	}
	if appendMode {
		args = append(args, "-a")
	} else if !redraw && startIndex > 0 {
		args = append(args, "-i", fmt.Sprintf("%d", startIndex))
	}

	err := a.runLuckSystem(args...)
	if err != nil {
		return "ERROR"
	}

	a.logOK(fmt.Sprintf("Font edited -> %s", outputCz))
	a.log("════════════════════════════════════════")
	return "OK"
}

// ═══════════════════════════════════════
// IMAGE EXPORT (single)
// ═══════════════════════════════════════
// lucksystem image export -i czfile -o output.png

func (a *App) ImageExport(czFile, outputPng string) string {
	if czFile == "" || outputPng == "" {
		a.logError("CZ input and PNG output are required")
		return "ERROR"
	}

	a.log("════════════════════════════════════════")
	a.log("  IMAGE EXPORT")
	a.log("════════════════════════════════════════")

	args := []string{"image", "export", "-i", czFile, "-o", outputPng}

	err := a.runLuckSystem(args...)
	if err != nil {
		return "ERROR"
	}

	a.logOK(fmt.Sprintf("Image exported -> %s", outputPng))
	a.log("════════════════════════════════════════")
	return "OK"
}

// ═══════════════════════════════════════
// IMAGE IMPORT (single)
// ═══════════════════════════════════════
// lucksystem image import -s source.cz -i input.png -o output.cz

func (a *App) ImageImport(sourceCz, inputPng, outputCz string, fill bool) string {
	if sourceCz == "" || inputPng == "" || outputCz == "" {
		a.logError("Source CZ, input PNG, and output CZ are required")
		return "ERROR"
	}

	a.log("════════════════════════════════════════")
	a.log("  IMAGE IMPORT")
	a.log("════════════════════════════════════════")

	args := []string{"image", "import", "-s", sourceCz, "-i", inputPng, "-o", outputCz}
	if fill {
		args = append(args, "-f")
	}

	err := a.runLuckSystem(args...)
	if err != nil {
		return "ERROR"
	}

	a.logOK(fmt.Sprintf("Image imported -> %s", outputCz))
	a.log("════════════════════════════════════════")
	return "OK"
}

// ═══════════════════════════════════════
// IMAGE BATCH EXPORT (directory)
// ═══════════════════════════════════════
// Iterates over all CZ files in inputDir, converts each to PNG in outputDir

func (a *App) ImageBatchExport(inputDir, outputDir string) string {
	if inputDir == "" || outputDir == "" {
		a.logError("Input and output directories are required")
		return "ERROR"
	}

	a.log("════════════════════════════════════════")
	a.log("  IMAGE BATCH EXPORT (CZ → PNG)")
	a.log("════════════════════════════════════════")
	a.log(fmt.Sprintf("Input:  %s", inputDir))
	a.log(fmt.Sprintf("Output: %s", outputDir))
	a.log("────────────────────────────────────────")

	os.MkdirAll(outputDir, os.ModePerm)

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		a.logError(fmt.Sprintf("Cannot read directory: %v", err))
		return "ERROR"
	}

	count := 0
	errors := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Skip files that already have a known extension (not CZ files)
		ext := strings.ToLower(filepath.Ext(name))
		if ext == ".png" || ext == ".txt" || ext == ".json" || ext == ".xml" {
			continue
		}

		inFile := filepath.Join(inputDir, name)
		outFile := filepath.Join(outputDir, name+".png")

		a.log(fmt.Sprintf("  [%d] %s ...", count+1, name))
		args := []string{"image", "export", "-i", inFile, "-o", outFile}
		if err := a.runLuckSystem(args...); err != nil {
			errors++
		} else {
			count++
		}
	}

	result := fmt.Sprintf("%d images exported, %d errors", count, errors)
	a.logOK(result)
	a.log("════════════════════════════════════════")
	return "OK: " + result
}

// ═══════════════════════════════════════
// IMAGE BATCH IMPORT (directory)
// ═══════════════════════════════════════
// For each PNG in inputDir, finds matching CZ in sourceDir, imports, saves to outputDir

func (a *App) ImageBatchImport(sourceDir, inputDir, outputDir string, fill bool) string {
	if sourceDir == "" || inputDir == "" || outputDir == "" {
		a.logError("Source CZ dir, input PNG dir, and output dir are required")
		return "ERROR"
	}

	a.log("════════════════════════════════════════")
	a.log("  IMAGE BATCH IMPORT (PNG → CZ)")
	a.log("════════════════════════════════════════")
	a.log(fmt.Sprintf("Source CZ: %s", sourceDir))
	a.log(fmt.Sprintf("Input PNG: %s", inputDir))
	a.log(fmt.Sprintf("Output:    %s", outputDir))
	a.log("────────────────────────────────────────")

	os.MkdirAll(outputDir, os.ModePerm)

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		a.logError(fmt.Sprintf("Cannot read directory: %v", err))
		return "ERROR"
	}

	count := 0
	errors := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.ToLower(filepath.Ext(name)) != ".png" {
			continue
		}

		// Derive original CZ name: "filename.png" -> "filename"
		czName := strings.TrimSuffix(name, filepath.Ext(name))
		sourceCz := filepath.Join(sourceDir, czName)
		inputPng := filepath.Join(inputDir, name)
		outputCz := filepath.Join(outputDir, czName)

		// Check source CZ exists
		if _, err := os.Stat(sourceCz); os.IsNotExist(err) {
			a.log(fmt.Sprintf("  [SKIP] %s (no matching CZ: %s)", name, czName))
			continue
		}

		a.log(fmt.Sprintf("  [%d] %s ...", count+1, name))
		args := []string{"image", "import", "-s", sourceCz, "-i", inputPng, "-o", outputCz}
		if fill {
			args = append(args, "-f")
		}
		if err := a.runLuckSystem(args...); err != nil {
			errors++
		} else {
			count++
		}
	}

	result := fmt.Sprintf("%d images imported, %d errors", count, errors)
	a.logOK(result)
	a.log("════════════════════════════════════════")
	return "OK: " + result
}

// ═══════════════════════════════════════
// FILE SELECTION HELPERS (Dialogue)
// ═══════════════════════════════════════

func (a *App) SelectScriptTxtFile() string {
	file, _ := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select decompiled script (.txt)",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "Text Files (*.txt)", Pattern: "*.txt"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
	return file
}

func (a *App) SelectTsvFile() string {
	file, _ := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select TSV dialogue file",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "TSV/Text Files (*.txt;*.tsv)", Pattern: "*.txt;*.tsv"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
	return file
}

func (a *App) SelectSaveTsvFile(defaultName string) string {
	file, _ := wailsRuntime.SaveFileDialog(a.ctx, wailsRuntime.SaveDialogOptions{
		Title:           "Save TSV dialogue file",
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
		Title:           "Save patched script file",
		DefaultFilename: defaultName,
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "Text Files (*.txt)", Pattern: "*.txt"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
	return file
}

// ═══════════════════════════════════════
// DIALOGUE EXTRACT / IMPORT
// ═══════════════════════════════════════
// Internal Go functions — NO lucksystem subprocess.
// Parses decompiled script .txt files from LuckSystem
// to extract/inject translatable entries as TSV.
//
// Supported line types:
//   MESSAGE(...)  — dialogue lines (all LuckEngine games)
//   LOG_BEGIN(...) — log/title entries (e.g. AIR, CLANNAD)
//
// Each line contains N quoted strings = N language columns.
// The user picks columns by number (Lang 1, Lang 2, ...).
// Column assignment varies by game — user must verify.

// DialogueFormatInfo is returned by DialogueDetectFormat
type DialogueFormatInfo struct {
	Format   string `json:"format"`
	MaxCols  int    `json:"maxCols"`
}

// isDialogueLine returns true if the line starts with MESSAGE or LOG_BEGIN.
func isDialogueLine(trimmed string) bool {
	return strings.HasPrefix(trimmed, "MESSAGE") || strings.HasPrefix(trimmed, "LOG_BEGIN")
}

// lineTag returns "MESSAGE" or "LOG_BEGIN" for tagging in the TSV ID column.
func lineTag(trimmed string) string {
	if strings.HasPrefix(trimmed, "LOG_BEGIN") {
		return "LOG_BEGIN"
	}
	return "MESSAGE"
}

// DialogueDetectFormat reads a decompiled script and detects the format.
// Scans MESSAGE and LOG_BEGIN lines, counts max quoted strings.
func (a *App) DialogueDetectFormat(scriptFile string) DialogueFormatInfo {
	result := DialogueFormatInfo{Format: "Unknown", MaxCols: 0}

	if scriptFile == "" {
		return result
	}

	data, err := os.ReadFile(scriptFile)
	if err != nil {
		a.logError(fmt.Sprintf("Cannot read file: %v", err))
		return result
	}

	lines := strings.Split(string(data), "\n")

	maxQuotes := 0
	msgCount := 0
	logCount := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !isDialogueLine(trimmed) {
			continue
		}
		if strings.HasPrefix(trimmed, "LOG_BEGIN") {
			logCount++
		} else {
			msgCount++
		}
		quotes := len(extractQuotedStrings(trimmed))
		if quotes > maxQuotes {
			maxQuotes = quotes
		}
		if msgCount+logCount >= 50 {
			break
		}
	}

	if msgCount+logCount == 0 {
		result.Format = "No MESSAGE / LOG_BEGIN found"
		return result
	}

	// Cap at 4 columns
	if maxQuotes > 4 {
		maxQuotes = 4
	}
	result.MaxCols = maxQuotes

	parts := []string{}
	if msgCount > 0 {
		parts = append(parts, fmt.Sprintf("%d MESSAGE", msgCount))
	}
	if logCount > 0 {
		parts = append(parts, fmt.Sprintf("%d LOG_BEGIN", logCount))
	}
	result.Format = fmt.Sprintf("%d columns detected (%s sampled)", maxQuotes, strings.Join(parts, " + "))

	return result
}

// extractQuotedStrings extracts all quoted strings from a line.
// Handles escaped quotes inside strings.
func extractQuotedStrings(line string) []string {
	var results []string
	inQuote := false
	var current strings.Builder
	runes := []rune(line)

	for i := 0; i < len(runes); i++ {
		ch := runes[i]
		if !inQuote {
			if ch == '"' {
				inQuote = true
				current.Reset()
			}
		} else {
			if ch == '\\' && i+1 < len(runes) && runes[i+1] == '"' {
				current.WriteRune('"')
				i++ // skip escaped quote
			} else if ch == '"' {
				results = append(results, current.String())
				inQuote = false
			} else {
				current.WriteRune(ch)
			}
		}
	}
	return results
}

// DialogueExtractFile extracts MESSAGE + LOG_BEGIN entries from a single script file to TSV.
func (a *App) DialogueExtractFile(inputFile, outputFile string, cols []int) string {
	if inputFile == "" || outputFile == "" {
		a.logError("Input script and output TSV file are required")
		return "ERROR"
	}
	if len(cols) == 0 {
		a.logError("At least one column must be selected")
		return "ERROR"
	}

	a.log("════════════════════════════════════════")
	a.log("  DIALOGUE EXTRACT (single file)")
	a.log("════════════════════════════════════════")
	a.log(fmt.Sprintf("Input:  %s", inputFile))
	a.log(fmt.Sprintf("Output: %s", outputFile))
	colNames := make([]string, len(cols))
	for i, c := range cols {
		colNames[i] = fmt.Sprintf("Lang %d", c)
	}
	a.log(fmt.Sprintf("Columns: %s", strings.Join(colNames, ", ")))

	count, err := a.extractDialoguesFromFile(inputFile, outputFile, cols)
	if err != nil {
		a.logError(fmt.Sprintf("Error: %v", err))
		return "ERROR"
	}

	result := fmt.Sprintf("%d entries extracted", count)
	a.logOK(result)
	a.log("════════════════════════════════════════")
	return "OK: " + result
}

// DialogueExtractBatch extracts MESSAGE + LOG_BEGIN entries from all .txt scripts in a folder.
func (a *App) DialogueExtractBatch(inputDir, outputDir string, cols []int) string {
	if inputDir == "" || outputDir == "" {
		a.logError("Input folder and output folder are required")
		return "ERROR"
	}
	if len(cols) == 0 {
		a.logError("At least one column must be selected")
		return "ERROR"
	}

	a.log("════════════════════════════════════════")
	a.log("  DIALOGUE EXTRACT (batch)")
	a.log("════════════════════════════════════════")
	a.log(fmt.Sprintf("Input:  %s", inputDir))
	a.log(fmt.Sprintf("Output: %s", outputDir))
	colNames := make([]string, len(cols))
	for i, c := range cols {
		colNames[i] = fmt.Sprintf("Lang %d", c)
	}
	a.log(fmt.Sprintf("Columns: %s", strings.Join(colNames, ", ")))

	os.MkdirAll(outputDir, 0755)

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		a.logError(fmt.Sprintf("Cannot read directory: %v", err))
		return "ERROR"
	}

	totalEntries := 0
	fileCount := 0
	errors := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".txt") {
			continue
		}
		// Skip files that are already extracted TSV (*.ext.txt)
		if strings.HasSuffix(strings.ToLower(e.Name()), ".ext.txt") {
			continue
		}

		inPath := filepath.Join(inputDir, e.Name())
		outName := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name())) + ".ext.txt"
		outPath := filepath.Join(outputDir, outName)

		count, err := a.extractDialoguesFromFile(inPath, outPath, cols)
		if err != nil {
			a.log(fmt.Sprintf("  [SKIP] %s: %v", e.Name(), err))
			errors++
			continue
		}
		if count > 0 {
			a.log(fmt.Sprintf("  [%d] %s → %s (%d entries)", fileCount+1, e.Name(), outName, count))
			totalEntries += count
			fileCount++
		}
	}

	result := fmt.Sprintf("%d files processed, %d entries total, %d errors", fileCount, totalEntries, errors)
	a.logOK(result)
	a.log("════════════════════════════════════════")
	return "OK: " + result
}

// extractDialoguesFromFile does the actual extraction work.
// cols contains 1-based column indices (e.g. [1, 2] for Lang 1 and Lang 2).
func (a *App) extractDialoguesFromFile(inputFile, outputFile string, cols []int) (int, error) {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return 0, fmt.Errorf("cannot read %s: %v", inputFile, err)
	}

	lines := strings.Split(string(data), "\n")

	// Build TSV header: ID | TAG | Lang N | Lang M | ...
	var sb strings.Builder
	sb.WriteString("ID\tTAG")
	for _, col := range cols {
		sb.WriteString(fmt.Sprintf("\tLang %d", col))
	}
	sb.WriteString("\n")

	count := 0
	seqID := 0 // sequential ID for stable matching
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !isDialogueLine(trimmed) {
			continue
		}

		seqID++
		tag := lineTag(trimmed)
		quoted := extractQuotedStrings(trimmed)

		sb.WriteString(fmt.Sprintf("%d\t%s", seqID, tag))
		for _, col := range cols {
			sb.WriteString("\t")
			idx := col - 1 // convert 1-based to 0-based
			if idx >= 0 && idx < len(quoted) {
				text := strings.ReplaceAll(quoted[idx], "\t", "\\t")
				text = strings.ReplaceAll(text, "\n", "\\n")
				text = strings.ReplaceAll(text, "\r", "")
				sb.WriteString(text)
			}
		}
		sb.WriteString("\n")
		count++
	}

	if count == 0 {
		return 0, nil
	}

	if err := os.WriteFile(outputFile, []byte(sb.String()), 0644); err != nil {
		return 0, fmt.Errorf("cannot write %s: %v", outputFile, err)
	}

	return count, nil
}

// DialogueImportFile re-injects a TSV column back into a single script file.
// targetCol is 1-based (Lang 1, Lang 2, etc.).
func (a *App) DialogueImportFile(scriptFile, tsvFile string, targetCol int, outputFile string) string {
	if scriptFile == "" || tsvFile == "" || outputFile == "" {
		a.logError("Script file, TSV file, and output file are required")
		return "ERROR"
	}

	a.log("════════════════════════════════════════")
	a.log("  DIALOGUE IMPORT (single file)")
	a.log("════════════════════════════════════════")
	a.log(fmt.Sprintf("Script: %s", scriptFile))
	a.log(fmt.Sprintf("TSV:    %s", tsvFile))
	a.log(fmt.Sprintf("Target: Lang %d (quoted string #%d)", targetCol, targetCol))
	a.log(fmt.Sprintf("Output: %s", outputFile))

	count, err := a.importDialoguesToFile(scriptFile, tsvFile, targetCol, outputFile)
	if err != nil {
		a.logError(fmt.Sprintf("Error: %v", err))
		return "ERROR"
	}

	result := fmt.Sprintf("%d entries injected", count)
	a.logOK(result)
	a.log("════════════════════════════════════════")
	return "OK: " + result
}

// DialogueImportBatch re-injects TSV columns into all matching scripts in a folder.
func (a *App) DialogueImportBatch(scriptsDir, tsvDir string, targetCol int, outputDir string) string {
	if scriptsDir == "" || tsvDir == "" || outputDir == "" {
		a.logError("Scripts folder, TSV folder, and output folder are required")
		return "ERROR"
	}

	a.log("════════════════════════════════════════")
	a.log("  DIALOGUE IMPORT (batch)")
	a.log("════════════════════════════════════════")
	a.log(fmt.Sprintf("Scripts: %s", scriptsDir))
	a.log(fmt.Sprintf("TSV:     %s", tsvDir))
	a.log(fmt.Sprintf("Target:  Lang %d", targetCol))
	a.log(fmt.Sprintf("Output:  %s", outputDir))

	os.MkdirAll(outputDir, 0755)

	entries, err := os.ReadDir(tsvDir)
	if err != nil {
		a.logError(fmt.Sprintf("Cannot read TSV directory: %v", err))
		return "ERROR"
	}

	totalEntries := 0
	fileCount := 0
	errors := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".ext.txt") {
			continue
		}

		// Derive script name: SEEN0001.ext.txt → SEEN0001.txt
		scriptName := strings.TrimSuffix(name, ".ext.txt") + ".txt"
		if strings.HasSuffix(strings.ToLower(name), ".EXT.txt") {
			scriptName = strings.TrimSuffix(name, ".EXT.txt") + ".txt"
		}

		scriptPath := filepath.Join(scriptsDir, scriptName)
		tsvPath := filepath.Join(tsvDir, name)
		outPath := filepath.Join(outputDir, scriptName)

		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			a.log(fmt.Sprintf("  [SKIP] %s (no matching script: %s)", name, scriptName))
			continue
		}

		count, err := a.importDialoguesToFile(scriptPath, tsvPath, targetCol, outPath)
		if err != nil {
			a.log(fmt.Sprintf("  [WARN] %s: %v", name, err))
			errors++
			continue
		}

		a.log(fmt.Sprintf("  [%d] %s + %s → %s (%d replaced)", fileCount+1, scriptName, name, scriptName, count))
		totalEntries += count
		fileCount++
	}

	result := fmt.Sprintf("%d files processed, %d entries injected, %d errors", fileCount, totalEntries, errors)
	a.logOK(result)
	a.log("════════════════════════════════════════")
	return "OK: " + result
}

// importDialoguesToFile does the actual import work.
// It reads the TSV, builds a map of seqID→translated text, then replaces
// the corresponding quoted string in each MESSAGE/LOG_BEGIN line of the script.
// targetCol is 1-based: Lang 1 = replace quoted string #0, Lang 2 = #1, etc.
func (a *App) importDialoguesToFile(scriptFile, tsvFile string, targetCol int, outputFile string) (int, error) {
	// --- Read TSV ---
	tsvData, err := os.ReadFile(tsvFile)
	if err != nil {
		return 0, fmt.Errorf("cannot read TSV: %v", err)
	}

	tsvLines := strings.Split(string(tsvData), "\n")
	if len(tsvLines) < 2 {
		return 0, fmt.Errorf("TSV file is empty or has no data rows")
	}

	// Parse header to find the target column.
	// New format: ID | TAG | Lang 1 | Lang 2 | ...
	// We look for "Lang N" where N == targetCol, OR fall back to column index.
	header := strings.Split(strings.TrimSpace(tsvLines[0]), "\t")
	targetTsvCol := -1
	targetHeader := fmt.Sprintf("Lang %d", targetCol)
	for i, col := range header {
		if strings.EqualFold(strings.TrimSpace(col), targetHeader) {
			targetTsvCol = i
			break
		}
	}
	// Fallback: also accept old-style named headers (JAP=col1→0, ENG=col1→1, CN→2)
	if targetTsvCol < 0 {
		oldNames := map[int][]string{
			1: {"JAP"},
			2: {"ENG"},
			3: {"CN"},
		}
		if names, ok := oldNames[targetCol]; ok {
			for i, col := range header {
				for _, name := range names {
					if strings.EqualFold(strings.TrimSpace(col), name) {
						targetTsvCol = i
						break
					}
				}
				if targetTsvCol >= 0 {
					break
				}
			}
		}
	}
	if targetTsvCol < 0 {
		return 0, fmt.Errorf("column '%s' not found in TSV header: %v", targetHeader, header)
	}

	// Build map: sequential_ID → translated text
	translations := make(map[int]string)
	for _, tsvLine := range tsvLines[1:] {
		tsvLine = strings.TrimSpace(tsvLine)
		if tsvLine == "" {
			continue
		}
		tsvCols := strings.Split(tsvLine, "\t")
		if len(tsvCols) <= targetTsvCol {
			continue
		}
		// First column is the sequential ID
		idStr := strings.TrimSpace(tsvCols[0])
		id := 0
		fmt.Sscanf(idStr, "%d", &id)
		if id <= 0 {
			continue
		}
		text := tsvCols[targetTsvCol]
		text = strings.ReplaceAll(text, "\\t", "\t")
		text = strings.ReplaceAll(text, "\\n", "\n")
		if text != "" {
			translations[id] = text
		}
	}

	// --- Read script and replace ---
	scriptData, err := os.ReadFile(scriptFile)
	if err != nil {
		return 0, fmt.Errorf("cannot read script: %v", err)
	}

	replaceIdx := targetCol - 1 // 1-based to 0-based

	scriptLines := strings.Split(string(scriptData), "\n")
	count := 0
	seqID := 0

	for i, line := range scriptLines {
		trimmed := strings.TrimSpace(line)
		if !isDialogueLine(trimmed) {
			continue
		}

		seqID++
		newText, ok := translations[seqID]
		if !ok || newText == "" {
			continue
		}

		replaced := replaceNthQuotedString(line, replaceIdx, newText)
		if replaced != line {
			scriptLines[i] = replaced
			count++
		}
	}

	output := strings.Join(scriptLines, "\n")
	if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
		return 0, fmt.Errorf("cannot write output: %v", err)
	}

	return count, nil
}

// replaceNthQuotedString replaces the Nth (0-based) quoted string in a line.
func replaceNthQuotedString(line string, n int, newText string) string {
	// Escape special chars in newText for reinsertion
	escaped := strings.ReplaceAll(newText, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
	escaped = strings.ReplaceAll(escaped, "\n", "\\n")
	escaped = strings.ReplaceAll(escaped, "\t", "\\t")

	runes := []rune(line)
	quoteCount := 0
	result := make([]rune, 0, len(runes)+len(escaped))
	inQuote := false
	skipUntilClose := false

	for i := 0; i < len(runes); i++ {
		ch := runes[i]

		if !inQuote {
			if ch == '"' {
				if quoteCount == n {
					// This is the opening quote of the target string
					result = append(result, '"')
					result = append(result, []rune(escaped)...)
					// Skip until closing quote
					skipUntilClose = true
					inQuote = true
					continue
				}
				inQuote = true
				result = append(result, ch)
			} else {
				result = append(result, ch)
			}
		} else {
			if skipUntilClose {
				// Skip original content until we find the unescaped closing quote
				if ch == '\\' && i+1 < len(runes) {
					i++ // skip escaped char
					continue
				}
				if ch == '"' {
					result = append(result, '"')
					inQuote = false
					skipUntilClose = false
					quoteCount++
				}
				continue
			}
			// Normal pass-through of non-target quoted strings
			if ch == '\\' && i+1 < len(runes) && runes[i+1] == '"' {
				result = append(result, ch, runes[i+1])
				i++
				continue
			}
			if ch == '"' {
				result = append(result, ch)
				inQuote = false
				quoteCount++
			} else {
				result = append(result, ch)
			}
		}
	}

	return string(result)
}
