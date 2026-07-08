package czimage

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/go-restruct/restruct"
)

func init() {
	restruct.EnableExprBeta()
}

func TestCz3LoadPreservesExtraHeader(t *testing.T) {
	extra := []byte{0xba, 0x00, 0x21, 0x00, 0x00, 0x00, 0x00, 0x00}
	raw := makeCz3TestData(t, extra)

	img := LoadCzImage(raw)
	cz, ok := img.(*Cz3Image)
	if !ok {
		t.Fatalf("LoadCzImage returned %T, want *Cz3Image", img)
	}
	if !bytes.Equal(cz.ExtraHeader, extra) {
		t.Fatalf("extra header = % x, want % x", cz.ExtraHeader, extra)
	}
}

func TestCz3WriteKeepsOutputInfoAtHeaderLength(t *testing.T) {
	extra := []byte{0xba, 0x00, 0x21, 0x00, 0x00, 0x00, 0x00, 0x00}
	cz := &Cz3Image{
		CzHeader: CzHeader{
			Magic:        []byte{'C', 'Z', '3', 0},
			HeaderLength: uint32(fixedCzSubHeaderLength + len(extra)),
			Width:        376,
			Heigth:       66,
			Colorbits:    32,
			Colorblock:   0,
		},
		Cz3Header: Cz3Header{
			Flag:    3,
			Width1:  373,
			Heigth1: 66,
			Width2:  373,
			Heigth2: 66,
		},
		ExtraHeader: extra,
		CzData: CzData{
			Raw: []byte{0xde, 0xad, 0xbe, 0xef},
			OutputInfo: &CzOutputInfo{
				FileCount: 1,
				BlockInfo: []CzBlockInfo{
					{CompressedSize: 2, RawSize: 4},
				},
			},
		},
	}

	var out bytes.Buffer
	if err := cz.Write(&out); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	data := out.Bytes()
	headerLength := int(cz.HeaderLength)

	if !bytes.Equal(data[fixedCzSubHeaderLength:headerLength], extra) {
		t.Fatalf("extra header bytes = % x, want % x", data[fixedCzSubHeaderLength:headerLength], extra)
	}
	if got := binary.LittleEndian.Uint32(data[headerLength:]); got != 1 {
		t.Fatalf("file count at HeaderLength = %d, want 1", got)
	}
	if got := binary.LittleEndian.Uint32(data[fixedCzSubHeaderLength:]); got == 1 {
		t.Fatalf("output info was written before the preserved extra header")
	}
}

func makeCz3TestData(t *testing.T, extra []byte) []byte {
	t.Helper()

	var buf bytes.Buffer
	header := CzHeader{
		Magic:        []byte{'C', 'Z', '3', 0},
		HeaderLength: uint32(fixedCzSubHeaderLength + len(extra)),
		Width:        376,
		Heigth:       66,
		Colorbits:    32,
		Colorblock:   0,
	}
	cz3Header := Cz3Header{
		Flag:    3,
		Width1:  373,
		Heigth1: 66,
		Width2:  373,
		Heigth2: 66,
	}
	outputInfo := &CzOutputInfo{
		FileCount: 1,
		BlockInfo: []CzBlockInfo{
			{CompressedSize: 2, RawSize: 4},
		},
	}
	if err := WriteStruct(&buf, &header, &cz3Header); err != nil {
		t.Fatalf("write header: %v", err)
	}
	if _, err := buf.Write(extra); err != nil {
		t.Fatalf("write extra header: %v", err)
	}
	if err := WriteStruct(&buf, outputInfo); err != nil {
		t.Fatalf("write output info: %v", err)
	}
	if _, err := buf.Write([]byte{0xde, 0xad, 0xbe, 0xef}); err != nil {
		t.Fatalf("write payload: %v", err)
	}
	return buf.Bytes()
}
