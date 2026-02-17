package czimage

import (
	"github.com/golang/glog"
	"image"
	"image/color"
)

func PanelImage(header *CzHeader, colorPanel [][]byte, data []byte) image.Image {
	width := int(header.Width)
	height := int(header.Heigth)
	pic := image.NewNRGBA(image.Rect(0, 0, width, height))
	// B,G,R,A
	// 0,1,2,3
	i := 0
	for y := 0; y < int(header.Heigth); y++ {
		for x := 0; x < int(header.Width); x++ {
			pic.SetNRGBA(x, y, color.NRGBA{
				R: colorPanel[data[i]][2],
				G: colorPanel[data[i]][1],
				B: colorPanel[data[i]][0],
				A: colorPanel[data[i]][3],
			})
			i++
		}
	}
	return pic
}

// DiffLine 图像拆分
//  Description 图像拆分，cz3用 png->data
//  Param header CzHeader
//  Param img image.Image
//  Return data
//
func DiffLine(header CzHeader, pic *image.NRGBA) (data []byte) {
	width := int(header.Width)
	height := int(header.Heigth)

	if width != pic.Rect.Size().X || height != pic.Rect.Size().Y {
		glog.V(2).Infof("图片大小不匹配，应该为 w%d h%d\n", width, height)
		return nil
	}
	data = make([]byte, len(pic.Pix))
	
	// PATCH YOREMI: blockHeight calculation (GARbro algorithm)
	blockHeight := (height + 2) / 3
	
	glog.V(0).Infof("DiffLine: height=%d, colorblock=%d, blockHeight=%d\n", height, header.Colorblock, blockHeight)
	
	pixelByteCount := int(header.Colorbits >> 3)
	lineByteCount := width * pixelByteCount
	
	preLine := make([]byte, lineByteCount)
	currLine := make([]byte, lineByteCount)  // Buffer pour la ligne courante
	
	i := 0
	for y := 0; y < height; y++ {
		// PATCH YOREMI: Copier pic.Pix dans currLine au lieu de créer un slice alias
		// Bug original: currLine = pic.Pix[i:...] créait un alias qui modifiait pic.Pix
		// Solution: Copier dans un buffer séparé
		copy(currLine, pic.Pix[i:i+lineByteCount])
		
		// Algorithme EXACT du code original
		if y%blockHeight != 0 {
			for x := 0; x < lineByteCount; x++ {
				currLine[x] -= preLine[x]
				// 因为是每一行较上一行的变化，故逆向执行时需要累加差异
				preLine[x] += currLine[x]
			}
		} else {
			copy(preLine, currLine)
		}

		copy(data[i:i+lineByteCount], currLine)
		i += lineByteCount
	}
	
	return data
}

// LineDiff 拆分图像还原
//  Description 拆分图像还原，cz3用 data->png
//  Param header *CzHeader
//  Param data []byte
//  Return image.Image
//
func LineDiff(header *CzHeader, data []byte) image.Image {
	//os.WriteFile("../data/LB_EN/IMAGE/ld.data", data, 0666)
	width := int(header.Width)
	height := int(header.Heigth)
	pic := image.NewNRGBA(image.Rect(0, 0, width, height))
	
	// PATCH YOREMI: blockHeight calculation (must match DiffLine)
	blockHeight := (height + 2) / 3
	
	glog.V(0).Infof("LineDiff: height=%d, colorblock=%d, blockHeight=%d\n", height, header.Colorblock, blockHeight)
	
	pixelByteCount := int(header.Colorbits >> 3)
	lineByteCount := width * pixelByteCount
	var currLine []byte
	preLine := make([]byte, lineByteCount)
	i := 0
	for y := 0; y < height; y++ {
		currLine = data[i : i+lineByteCount]
		if y%blockHeight != 0 {
			for x := 0; x < lineByteCount; x++ {
				currLine[x] += preLine[x]
			}
		}
		// PATCH YOREMI: Copier les données au lieu d'aliaser
		// Bug original: preLine = currLine créait un alias, pas une copie
		copy(preLine, currLine)
		
		if pixelByteCount == 4 {
			// y*pic.Stride : (y+1)*pic.Stride
			copy(pic.Pix[i:i+lineByteCount], currLine)
		} else if pixelByteCount == 3 {
			for x := 0; x < lineByteCount; x += 3 {
				pic.SetNRGBA(x/3, y, color.NRGBA{R: currLine[x], G: currLine[x+1], B: currLine[x+2], A: 0xFF})
			}
		}
		i += lineByteCount
	}
	
	//os.WriteFile("../data/LB_EN/IMAGE/ld.data.pix", pic.Pix, 0666)
	return pic
}
