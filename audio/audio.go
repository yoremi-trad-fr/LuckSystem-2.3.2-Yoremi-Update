package audio

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const (
	DefaultMP3Quality = "2"
	DefaultOggQuality = "5"
)

type Entry struct {
	Index  int
	ID     int
	Offset int64
	Length int64
	Name   string
}

type PakInfo struct {
	HeaderLength uint32
	FileCount    uint32
	IDStart      uint32
	BlockSize    uint32
	Flags        uint32
	OffsetPos    int64
	Entries      []Entry
}

type ExtractOptions struct {
	PakFile    string
	OutputDir  string
	Kind       string
	ConvertMP3 bool
	FFmpegPath string
	MP3Quality string
}

type ConvertOptions struct {
	InputPath  string
	OutputDir  string
	Direction  string
	FFmpegPath string
	MP3Quality string
	OggQuality string
	Overwrite  bool
}

type Summary struct {
	Files     int
	Converted int
	Skipped   int
	Errors    int
	ListFile  string
	NativeDir string
	MP3Dir    string
}

func ExtractPak(opts ExtractOptions, logf func(string)) (Summary, error) {
	if opts.PakFile == "" || opts.OutputDir == "" {
		return Summary{}, errors.New("PAK file and output directory are required")
	}
	if opts.MP3Quality == "" {
		opts.MP3Quality = DefaultMP3Quality
	}

	info, err := ReadPakInfo(opts.PakFile)
	if err != nil {
		return Summary{}, err
	}

	nativeDir := filepath.Join(opts.OutputDir, "ogg")
	if err := os.MkdirAll(nativeDir, 0755); err != nil {
		return Summary{}, err
	}

	pakBase := strings.TrimSuffix(filepath.Base(opts.PakFile), filepath.Ext(opts.PakFile))
	listFile := filepath.Join(opts.OutputDir, pakBase+"_audio_list.txt")
	list, err := os.Create(listFile)
	if err != nil {
		return Summary{}, err
	}
	defer list.Close()
	listWriter := bufio.NewWriter(list)
	defer listWriter.Flush()

	in, err := os.Open(opts.PakFile)
	if err != nil {
		return Summary{}, err
	}
	defer in.Close()

	summary := Summary{
		ListFile:  listFile,
		NativeDir: nativeDir,
	}

	var ffmpeg string
	mp3Dir := ""
	if opts.ConvertMP3 {
		ffmpeg, err = ResolveFFmpeg(opts.FFmpegPath)
		if err != nil {
			return summary, err
		}
		mp3Dir = filepath.Join(opts.OutputDir, "mp3")
		if err := os.MkdirAll(mp3Dir, 0755); err != nil {
			return summary, err
		}
		summary.MP3Dir = mp3Dir
	}

	for _, entry := range info.Entries {
		name := audioEntryName(entry)
		oggPath := filepath.Join(nativeDir, name+".ogg")
		if err := copyEntry(in, entry, oggPath); err != nil {
			summary.Errors++
			if logf != nil {
				logf(fmt.Sprintf("  [ERR] %s: %v", filepath.Base(oggPath), err))
			}
			continue
		}
		summary.Files++
		_, _ = fmt.Fprintf(listWriter, "id:%d,%s\n", entry.ID, oggPath)

		if opts.ConvertMP3 {
			mp3Path := filepath.Join(mp3Dir, name+".mp3")
			if err := convertFile(ffmpeg, oggPath, mp3Path, "mp3", opts.MP3Quality, true); err != nil {
				summary.Errors++
				if logf != nil {
					logf(fmt.Sprintf("  [ERR] %s -> mp3: %v", filepath.Base(oggPath), err))
				}
			} else {
				summary.Converted++
			}
		}

		if logf != nil && (summary.Files <= 5 || summary.Files%250 == 0 || summary.Files == len(info.Entries)) {
			if opts.ConvertMP3 {
				logf(fmt.Sprintf("  [%d/%d] %s + mp3", summary.Files, len(info.Entries), filepath.Base(oggPath)))
			} else {
				logf(fmt.Sprintf("  [%d/%d] %s", summary.Files, len(info.Entries), filepath.Base(oggPath)))
			}
		}
	}

	if summary.Files == 0 {
		return summary, errors.New("no audio entries extracted")
	}
	return summary, nil
}

func ConvertPath(opts ConvertOptions, logf func(string)) (Summary, error) {
	if opts.InputPath == "" || opts.OutputDir == "" {
		return Summary{}, errors.New("input and output directory are required")
	}
	if opts.MP3Quality == "" {
		opts.MP3Quality = DefaultMP3Quality
	}
	if opts.OggQuality == "" {
		opts.OggQuality = DefaultOggQuality
	}
	direction := strings.ToLower(strings.TrimSpace(opts.Direction))
	if direction != "mp3" && direction != "native" && direction != "ogg" {
		return Summary{}, errors.New("direction must be mp3 or native")
	}
	if direction == "ogg" {
		direction = "native"
	}

	ffmpeg, err := ResolveFFmpeg(opts.FFmpegPath)
	if err != nil {
		return Summary{}, err
	}
	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return Summary{}, err
	}

	inInfo, err := os.Stat(opts.InputPath)
	if err != nil {
		return Summary{}, err
	}

	var inputs []string
	if inInfo.IsDir() {
		inputs, err = collectAudioInputs(opts.InputPath, direction)
		if err != nil {
			return Summary{}, err
		}
	} else if matchesDirection(opts.InputPath, direction) {
		inputs = []string{opts.InputPath}
	}
	sort.Strings(inputs)
	if len(inputs) == 0 {
		return Summary{}, errors.New("no matching audio files found")
	}

	summary := Summary{}
	for i, inPath := range inputs {
		rel := filepath.Base(inPath)
		if inInfo.IsDir() {
			if r, err := filepath.Rel(opts.InputPath, inPath); err == nil {
				rel = r
			}
		}
		outRel := replaceExt(rel, outputExt(direction))
		outPath := filepath.Join(opts.OutputDir, outRel)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			summary.Errors++
			if logf != nil {
				logf(fmt.Sprintf("  [ERR] %s: %v", rel, err))
			}
			continue
		}
		if !opts.Overwrite {
			if _, err := os.Stat(outPath); err == nil {
				summary.Skipped++
				continue
			}
		}
		quality := opts.MP3Quality
		if direction == "native" {
			quality = opts.OggQuality
		}
		if err := convertFile(ffmpeg, inPath, outPath, direction, quality, opts.Overwrite); err != nil {
			summary.Errors++
			if logf != nil {
				logf(fmt.Sprintf("  [ERR] %s: %v", rel, err))
			}
			continue
		}
		summary.Converted++
		if logf != nil && (summary.Converted <= 5 || summary.Converted%250 == 0 || i == len(inputs)-1) {
			logf(fmt.Sprintf("  [%d/%d] %s -> %s", i+1, len(inputs), rel, outRel))
		}
	}
	if summary.Converted == 0 && summary.Errors > 0 {
		return summary, errors.New("audio conversion failed")
	}
	return summary, nil
}

func ReadPakInfo(filename string) (PakInfo, error) {
	f, err := os.Open(filename)
	if err != nil {
		return PakInfo{}, err
	}
	defer f.Close()

	var headerLenBytes [4]byte
	if _, err := f.ReadAt(headerLenBytes[:], 0); err != nil {
		return PakInfo{}, err
	}
	headerLen := binary.LittleEndian.Uint32(headerLenBytes[:])
	if headerLen < 36 {
		return PakInfo{}, fmt.Errorf("invalid PAK header length: %d", headerLen)
	}

	header := make([]byte, headerLen)
	if _, err := f.ReadAt(header, 0); err != nil {
		return PakInfo{}, err
	}
	info := PakInfo{
		HeaderLength: binary.LittleEndian.Uint32(header[0:4]),
		FileCount:    binary.LittleEndian.Uint32(header[4:8]),
		IDStart:      binary.LittleEndian.Uint32(header[8:12]),
		BlockSize:    binary.LittleEndian.Uint32(header[12:16]),
		Flags:        binary.LittleEndian.Uint32(header[32:36]),
	}
	if info.BlockSize == 0 {
		return PakInfo{}, errors.New("invalid PAK block size: 0")
	}

	needle := info.HeaderLength / info.BlockSize
	offsetPos := int64(32)
	for offsetPos+4 <= int64(len(header)) && binary.LittleEndian.Uint32(header[offsetPos:offsetPos+4]) != needle {
		offsetPos += 4
	}
	if offsetPos+int64(8*info.FileCount) > int64(len(header)) {
		return PakInfo{}, errors.New("cannot locate PAK entry table")
	}
	info.OffsetPos = offsetPos

	fileInfo, err := f.Stat()
	if err != nil {
		return PakInfo{}, err
	}
	info.Entries = make([]Entry, 0, info.FileCount)
	for i := 0; i < int(info.FileCount); i++ {
		pos := int(offsetPos) + i*8
		offsetBlocks := binary.LittleEndian.Uint32(header[pos : pos+4])
		length := binary.LittleEndian.Uint32(header[pos+4 : pos+8])
		offset := int64(offsetBlocks) * int64(info.BlockSize)
		if offset < 0 || int64(length) < 0 || offset+int64(length) > fileInfo.Size() {
			return PakInfo{}, fmt.Errorf("entry %d points outside file", i)
		}
		info.Entries = append(info.Entries, Entry{
			Index:  i,
			ID:     int(info.IDStart) + i,
			Offset: offset,
			Length: int64(length),
		})
	}
	return info, nil
}

func ResolveFFmpeg(preferred string) (string, error) {
	if preferred != "" {
		if _, err := os.Stat(preferred); err == nil {
			return preferred, nil
		}
	}
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		return path, nil
	}
	if path, err := exec.LookPath("ffmpeg.exe"); err == nil {
		return path, nil
	}
	return "", errors.New("ffmpeg not found in PATH")
}

func copyEntry(in *os.File, entry Entry, outPath string) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}
	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()
	reader := io.NewSectionReader(in, entry.Offset, entry.Length)
	written, err := io.Copy(out, reader)
	if err != nil {
		return err
	}
	if written != entry.Length {
		return fmt.Errorf("short write: %d/%d", written, entry.Length)
	}
	return nil
}

func convertFile(ffmpeg, inPath, outPath, direction, quality string, overwrite bool) error {
	args := []string{"-hide_banner", "-loglevel", "error"}
	if overwrite {
		args = append(args, "-y")
	} else {
		args = append(args, "-n")
	}
	args = append(args, "-i", inPath)
	if direction == "mp3" {
		args = append(args, "-vn", "-codec:a", "libmp3lame", "-q:a", quality, outPath)
	} else {
		args = append(args, "-vn", "-codec:a", "libvorbis", "-q:a", quality, outPath)
	}
	cmd := exec.Command(ffmpeg, args...)
	hideWindow(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return errors.New(msg)
	}
	return nil
}

func collectAudioInputs(root, direction string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if matchesDirection(path, direction) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func matchesDirection(path, direction string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	if direction == "mp3" {
		return ext == ".ogg"
	}
	return ext == ".mp3"
}

func outputExt(direction string) string {
	if direction == "mp3" {
		return ".mp3"
	}
	return ".ogg"
}

func replaceExt(path, ext string) string {
	return strings.TrimSuffix(path, filepath.Ext(path)) + ext
}

func audioEntryName(entry Entry) string {
	if entry.Name != "" {
		return sanitizeBaseName(strings.TrimSuffix(entry.Name, filepath.Ext(entry.Name)))
	}
	return fmt.Sprintf("%d", entry.ID)
}

func sanitizeBaseName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "audio"
	}
	replacer := strings.NewReplacer(
		"<", "_", ">", "_", ":", "_", "\"", "_", "/", "_", "\\", "_",
		"|", "_", "?", "_", "*", "_",
	)
	name = replacer.Replace(name)
	return strings.Trim(name, ". ")
}
