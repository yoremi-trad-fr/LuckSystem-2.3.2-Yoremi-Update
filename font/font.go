package font

import (
	"bytes"
	"errors"
	"github.com/golang/glog"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"lucksystem/czimage"
	"lucksystem/pak"
	"math"
	"os"
	"strconv"
)

type LucaFont struct {
	Size    int
	CzImage czimage.CzImage
	Info    *Info
	Image   *image.NRGBA
}

// LoadLucaFontPak 通过pak加载LucaFont
//
//	Description
//	Param pak *pak.PakFile
//	Param fontName string モダン/明朝/丸ゴシック/ゴシック
//	Param size int 12 14 16 18 20 24 28 30 32 36 72
//	Return *LucaFont
func LoadLucaFontPak(pak *pak.Pak, fontName string, size int) *LucaFont {
	infoFile, err := pak.Get("info" + strconv.Itoa(size))
	if err != nil {
		glog.Fatalln(err)
	}
	imageFile, err := pak.Get(fontName + strconv.Itoa(size))
	if err != nil {
		glog.Fatalln(err)
	}
	return LoadLucaFont(infoFile.Data, imageFile.Data)
}

// LoadLucaFontFile 通过文件名加载LucaFont
//
//	Description
//	Param infoFilename string
//	Param imageFilename string
//	Return *LucaFont
func LoadLucaFontFile(infoFilename, imageFilename string) *LucaFont {
	infoFile, err := os.ReadFile(infoFilename)
	if err != nil {
		glog.Fatalln(err)
	}
	imageFile, err := os.ReadFile(imageFilename)
	if err != nil {
		glog.Fatalln(err)
	}
	return LoadLucaFont(infoFile, imageFile)
}

// LoadLucaFont 通过字节数据加载LucaFont
//
//	Description
//	Param infoFile []byte
//	Param imageFile []byte
//	Return *LucaFont
func LoadLucaFont(infoFile, imageFile []byte) *LucaFont {
	font := &LucaFont{}
	font.Info = LoadFontInfo(infoFile)
	font.Size = int(font.Info.FontSize)
	font.CzImage = czimage.LoadCzImage(imageFile)

	font.Image = font.CzImage.GetImage().(*image.NRGBA)
	return font
}

// GetCharImage 获取单个字符图像和偏移信息
//
//	Description
//	Receiver f *LucaFont
//	Param unicode rune
//	Return image.Image
//	Return DrawSize
func (f *LucaFont) GetCharImage(unicode rune) (image.Image, DrawSize) {

	index, draw, _ := f.Info.Get(unicode)
	size := int(f.Info.BlockSize)
	y := index / 100
	x := index % 100
	return f.Image.SubImage(image.Rect(x*size, y*size, (x+1)*size, (y+1)*size)), draw
}

// GetStringImageList 获取字符串每个字符的图像和偏移信息
//
//	Description
//	Receiver f *LucaFont
//	Param str string
//	Return []image.Image
//	Return []DrawSize
//	Return []CharSize
func (f *LucaFont) GetStringImageList(str string) ([]image.Image, []DrawSize, []CharSize) {
	imgs := make([]image.Image, 0, len(str))
	draws := make([]DrawSize, 0, len(str))
	sizes := make([]CharSize, 0, len(str))
	for _, r := range str {
		index, drawSize, charSize := f.Info.Get(r)
		size := int(f.Info.BlockSize)
		y := index / 100
		x := index % 100
		img := f.Image.SubImage(image.Rect(x*size, y*size, (x+1)*size, (y+1)*size))
		imgs = append(imgs, img)
		draws = append(draws, drawSize)
		sizes = append(sizes, charSize)
	}
	return imgs, draws, sizes
}

// GetStringImage 将字符串转化为图像
//
//	Description
//	Receiver f *LucaFont
//	Param str string
//	Return image.Image
func (f *LucaFont) GetStringImage(str string) image.Image {
	imgW := int(f.Info.BlockSize)
	imgs, draws, sizes := f.GetStringImageList(str)
	pic := image.NewNRGBA(image.Rect(0, 0, len(imgs)*imgW, imgW*2))
	X := 0
	for i, img := range imgs {

		draw.Draw(pic, pic.Bounds().Add(image.Pt(X+signedMetricInt(draws[i].X), signedMetricInt(draws[i].Y))), img, img.Bounds().Min, draw.Src)
		advance := draws[i].W
		if sizes[i].W != 0 {
			advance = sizes[i].W
		}
		X += int(advance)
	}
	return pic
}

// BleedGlyphEdges extends row edge pixels inside selected glyph cells.
// It is intended for bitmap renderers that leave visible gaps between
// pre-shaped connecting glyphs even when metric advances are tightened.
func (f *LucaFont) BleedGlyphEdges(chars []rune, pixels int) {
	if f == nil || f.Info == nil || f.Image == nil || pixels <= 0 {
		return
	}
	if pixels > 8 {
		pixels = 8
	}
	size := int(f.Info.BlockSize)
	if size <= 0 {
		return
	}
	seen := make(map[uint16]bool, len(chars))
	for _, char := range chars {
		if char < 0 || int(char) >= len(f.Info.UnicodeIndex) {
			continue
		}
		index := f.Info.UnicodeIndex[int(char)]
		if index == 0 && char != ' ' {
			continue
		}
		if int(index) >= len(f.Info.DrawSize) || seen[index] {
			continue
		}
		seen[index] = true
		f.bleedGlyphCell(int(index), size, pixels)
	}
}

func (f *LucaFont) bleedGlyphCell(index, size, pixels int) {
	cellX := (index % 100) * size
	cellY := (index / 100) * size
	cellRect := image.Rect(cellX, cellY, cellX+size, cellY+size).Intersect(f.Image.Bounds())
	if cellRect.Empty() || cellRect.Dx() != size || cellRect.Dy() != size {
		return
	}

	src := image.NewNRGBA(image.Rect(0, 0, size, size))
	draw.Draw(src, src.Bounds(), f.Image, cellRect.Min, draw.Src)
	maxTouchedX := -1
	for y := 0; y < size; y++ {
		left := -1
		right := -1
		for x := 0; x < size; x++ {
			if src.NRGBAAt(x, y).A <= 16 {
				continue
			}
			if left < 0 {
				left = x
			}
			right = x
		}
		if left < 0 {
			continue
		}
		if right-left+1 < 5 {
			continue
		}
		leftPixel := src.NRGBAAt(left, y)
		rightPixel := src.NRGBAAt(right, y)
		for step := 1; step <= pixels; step++ {
			if x := left - step; x >= 0 {
				setMaxAlphaNRGBA(f.Image, cellX+x, cellY+y, leftPixel)
			}
			if x := right + step; x < size {
				setMaxAlphaNRGBA(f.Image, cellX+x, cellY+y, rightPixel)
				if x > maxTouchedX {
					maxTouchedX = x
				}
			}
		}
	}
	if maxTouchedX >= 0 && index < len(f.Info.DrawSize) {
		drawSize := &f.Info.DrawSize[index]
		width := maxTouchedX + 1
		if width > size {
			width = size
		}
		if width > 255 {
			width = 255
		}
		if width > int(drawSize.W) {
			drawSize.W = uint8(width)
		}
	}
}

func setMaxAlphaNRGBA(img *image.NRGBA, x, y int, c color.NRGBA) {
	if c.A == 0 || !image.Pt(x, y).In(img.Bounds()) {
		return
	}
	current := img.NRGBAAt(x, y)
	if current.A >= c.A {
		return
	}
	img.SetNRGBA(x, y, c)
}

// CreateLucaFont 创建全新的字体
//
//	Description
//	Param fontSize int 字体大小
//	Param fontFile io.Reader 字体文件
//	Param allChar string 所有字符
//	Return *LucaFont
func CreateLucaFont(fontSize int, fontFile io.Reader, allChar string) *LucaFont {
	font := &LucaFont{
		Size: fontSize,
	}
	font.Info = CreateFontInfo(fontSize, fontSize+1)
	//font.Info.SetChars(, 20)
	font.ReplaceChars(fontFile, allChar, 0, true)

	return font
}

// ReplaceChars 替换字体中的字符
//
//	Description 替换字体中的字符信息以及图像, 如果startIndex=0且allChar为空，则为修改原字体
//	Receiver f *LucaFont
//	Param fontFile io.Reader 字体文件
//	Param allChar string 所替换的字符
//	Param startIndex int 开始序号（图像从上到下，从左到右计算）
//	Param reDraw bool 是否用新字体重绘startIndex之前的字符
func (f *LucaFont) ReplaceChars(fontFile io.Reader, allChar string, startIndex int, reDraw bool) {

	if f.Info == nil {
		glog.Fatalln("需要先载入或创建LucaFont")
		return
	}
	if len(allChar) == 0 && startIndex == 0 && !reDraw {
		// 什么都不做
		return
	}
	f.Info.SetChars(fontFile, allChar, startIndex, reDraw)
	size := int(f.Info.BlockSize)
	imageW := size*100 + 4                                         // 100个字符宽度+4
	imageH := size * int(math.Ceil(float64(f.Info.CharNum)/100.0)) // 对应行数高度
	if f.Image != nil {
		oldSize := f.Image.Bounds().Size()
		if oldSize.X > imageW {
			imageW = oldSize.X
		}
		if oldSize.Y > imageH {
			imageH = oldSize.Y
		}
	}

	pic := image.NewNRGBA(image.Rect(0, 0, imageW, imageH))
	if !reDraw && f.Image != nil {
		oldBounds := f.Image.Bounds()
		copyRect := image.Rect(0, 0, oldBounds.Dx(), oldBounds.Dy()).Intersect(pic.Bounds())
		draw.Draw(pic, copyRect, f.Image, oldBounds.Min, draw.Src)
	}

	alphaMask := image.NewAlpha(image.Rect(0, 0, size, size))
	if reDraw {
		startIndex = 0
	}
	for i := startIndex; i < int(f.Info.CharNum); i++ {
		y := i / 100
		x := i % 100
		point := fixed.Point26_6{
			X: fixed.Int26_6(x * 64),
			Y: fixed.Int26_6(y * 64),
		}
		_, img, _, _, _ := f.Info.FontFace.Glyph(point, f.Info.IndexUnicode[i])
		// yOffset := dr.Min.Y + fontSize
		// fmt.Println(string(font.Info.IndexFont[i]), " ", dr.Min.Y+fontSize)
		draw.Draw(pic, pic.Bounds().Add(image.Pt(x*size, y*size)), alphaMask, alphaMask.Bounds().Min, draw.Src)
		draw.Draw(pic, pic.Bounds().Add(image.Pt(x*size, y*size)), img, img.Bounds().Min, draw.Src)
	}
	f.Image = pic
}

// Export
//
//	Description
//	Receiver f *LucaFont
//	Param w io.Writer
//	Param allCharFile string 导出的全字符文件名
//	Return error
func (f *LucaFont) Export(w io.Writer, allCharFile string) error {
	err := png.Encode(w, f.Image)
	if err != nil {
		return err
	}
	if len(allCharFile) > 0 {
		fs, _ := os.Create(allCharFile)
		err = f.Info.Export(fs)
		fs.Close()
	}
	return err
}

// Import
//
//	Description 若startIndex=0, redraw=true, allChar="", 则仅使用字体重绘原字符集
//	Receiver f *LucaFont
//	Param r io.Reader 字体文件
//	Param startIndex int 开始位置。前面跳过字符数量，-1为添加到最后
//	Param redraw bool 是否用新字体重绘startIndex之前的字符
//	Param allChar string 增加的全字符，若startIndex==0，且第一个字符不是空格，会自动补充为空格
//	Return error
func (f *LucaFont) Import(r io.Reader, startIndex int, redraw bool, allCharFile string) error {

	if len(allCharFile) == 0 {
		if redraw {
			// 仅重绘
			f.ReplaceChars(r, "", 0, true)
		}
		return nil
	}
	if startIndex == -1 {
		startIndex = int(f.Info.CharNum)
	}
	data, err := os.ReadFile(allCharFile)
	if err != nil {
		return err
	}
	f.ReplaceChars(r, string(data), startIndex, redraw)
	return nil
}

// Write
//
//	Description
//	Receiver f *LucaFont
//	Param w io.Writer
//	Param infoW io.Writer 导出新的info文件
//	Return error
func (f *LucaFont) Write(w io.Writer, infoW io.Writer) error {
	var err error
	if f.CzImage != nil {
		// load

		img := bytes.NewBuffer(nil)
		err = png.Encode(img, f.Image)
		if err != nil {
			return err
		}
		czImg := bytes.NewBuffer(nil)
		// Yoremi Patch 3: sync CzHeader dimensions with the (possibly resized) image
		// before calling Import, so CZ2.Import sees the correct target dimensions.
		if setter, ok := f.CzImage.(interface {
			SetDimensions(w, h uint16)
		}); ok {
			setter.SetDimensions(
				uint16(f.Image.Bounds().Size().X),
				uint16(f.Image.Bounds().Size().Y),
			)
		}
		err = f.CzImage.Import(img, true)
		if err != nil {
			return err
		}
		err = f.CzImage.Write(czImg)
		if err != nil {
			return err
		}
		_, err = w.Write(czImg.Bytes())
		if err != nil {
			return err
		}
	} else {
		// create
		return errors.New("LucaFont.Write 目前不支持创建的字体")
	}
	if infoW != nil {
		err = f.Info.Write(infoW)
	}
	return err
}
