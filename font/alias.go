package font

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"strconv"
	"strings"

	"lucksystem/czimage"
)

const fontAtlasColumns = 100

// BuildFontAliasEntry copies the content of one font-size entry into another
// while retaining the structural geometry required by the target size.
//
// A byte-for-byte copy is unsafe for LucaSystem CZ2 fonts: an entry named
// "...32" must keep the size-32 cell layout and texture width. The returned
// data therefore keeps target info geometry, or regrids source CZ2 cells into
// the target CZ2 container.
func BuildFontAliasEntry(sourceName, targetName string, sourceData, targetData []byte) ([]byte, error) {
	sourceSize, err := fontEntrySize(sourceName)
	if err != nil {
		return nil, err
	}
	targetSize, err := fontEntrySize(targetName)
	if err != nil {
		return nil, err
	}
	if sourceSize == targetSize {
		return nil, errors.New("font alias source and target sizes must differ")
	}

	sourceIsInfo := strings.HasPrefix(strings.ToLower(sourceName), "info")
	targetIsInfo := strings.HasPrefix(strings.ToLower(targetName), "info")
	if sourceIsInfo || targetIsInfo {
		if !sourceIsInfo || !targetIsInfo {
			return nil, errors.New("font alias cannot mix info and image entries")
		}
		alias, err := buildFontInfoAlias(sourceData, targetData)
		return padFontAliasToTargetLength(alias, targetData), err
	}

	alias, err := buildFontImageAlias(sourceData, targetData, sourceSize+1, targetSize+1)
	return padFontAliasToTargetLength(alias, targetData), err
}

func padFontAliasToTargetLength(alias, target []byte) []byte {
	if len(alias) >= len(target) {
		return alias
	}
	padded := make([]byte, len(target))
	copy(padded, alias)
	return padded
}

func fontEntrySize(name string) (int, error) {
	end := len(name)
	start := end
	for start > 0 {
		c := name[start-1]
		if c < '0' || c > '9' {
			break
		}
		start--
	}
	if start == end {
		return 0, fmt.Errorf("font entry %q has no trailing size", name)
	}
	size, err := strconv.Atoi(name[start:end])
	if err != nil || size <= 0 {
		return 0, fmt.Errorf("invalid font size in entry %q", name)
	}
	return size, nil
}

func buildFontInfoAlias(sourceData, targetData []byte) ([]byte, error) {
	if len(sourceData) < 6 || len(targetData) < 6 {
		return nil, errors.New("font info entry is too small")
	}

	source := LoadFontInfo(sourceData)
	target := LoadFontInfo(targetData)
	if source == nil || target == nil {
		return nil, errors.New("unable to load font info entries")
	}

	// Keep the source glyph map and metrics, but retain the nominal target
	// size/cell geometry so the engine still validates this as the target size.
	source.FontSize = target.FontSize
	source.BlockSize = target.BlockSize

	var out bytes.Buffer
	if err := source.Write(&out); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func buildFontImageAlias(sourceData, targetData []byte, sourceBlock, targetBlock int) ([]byte, error) {
	if len(sourceData) < 15 || len(targetData) < 15 {
		return nil, errors.New("font image entry is too small")
	}
	if string(sourceData[:3]) != "CZ2" || string(targetData[:3]) != "CZ2" {
		return nil, errors.New("target-compatible font image alias currently requires CZ2 source and target entries")
	}

	source := czimage.LoadCzImage(sourceData)
	target := czimage.LoadCzImage(targetData)
	if source == nil || target == nil {
		return nil, errors.New("unable to load CZ2 font entries")
	}

	targetWidth := target.GetImage().Bounds().Dx()
	aliasImage, err := regridFontAtlas(source.GetImage(), sourceBlock, targetBlock, targetWidth)
	if err != nil {
		return nil, err
	}

	var pngData bytes.Buffer
	if err := png.Encode(&pngData, aliasImage); err != nil {
		return nil, err
	}
	// Do not fill to the old target height. The source may contain appended
	// glyph rows; Import(false) updates only the height while preserving the
	// target CZ2 header, palette, nominal width, and encoding path.
	if err := target.Import(&pngData, false); err != nil {
		return nil, err
	}

	var out bytes.Buffer
	if err := target.Write(&out); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func regridFontAtlas(source image.Image, sourceBlock, targetBlock, targetWidth int) (*image.NRGBA, error) {
	if source == nil {
		return nil, errors.New("source font atlas is nil")
	}
	if sourceBlock <= 0 || targetBlock <= 0 {
		return nil, errors.New("font cell sizes must be positive")
	}
	if targetBlock < sourceBlock {
		return nil, errors.New("font alias target size must be greater than the source size")
	}

	bounds := source.Bounds()
	if bounds.Dx() < fontAtlasColumns*sourceBlock {
		return nil, fmt.Errorf("source atlas width %d is smaller than %d font cells", bounds.Dx(), fontAtlasColumns)
	}
	if bounds.Dy()%sourceBlock != 0 {
		return nil, fmt.Errorf("source atlas height %d is not divisible by cell size %d", bounds.Dy(), sourceBlock)
	}
	if targetWidth < fontAtlasColumns*targetBlock {
		return nil, fmt.Errorf("target atlas width %d is smaller than %d font cells", targetWidth, fontAtlasColumns)
	}

	rows := bounds.Dy() / sourceBlock
	target := image.NewNRGBA(image.Rect(0, 0, targetWidth, rows*targetBlock))
	for row := 0; row < rows; row++ {
		for column := 0; column < fontAtlasColumns; column++ {
			sourceMin := image.Pt(bounds.Min.X+column*sourceBlock, bounds.Min.Y+row*sourceBlock)
			targetRect := image.Rect(
				column*targetBlock,
				row*targetBlock,
				column*targetBlock+sourceBlock,
				row*targetBlock+sourceBlock,
			)
			draw.Draw(target, targetRect, source, sourceMin, draw.Src)
		}
	}
	return target, nil
}
