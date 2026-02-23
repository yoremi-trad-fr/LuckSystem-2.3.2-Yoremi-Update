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
