package movie

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

var ebmlMagic = []byte{0x1A, 0x45, 0xDF, 0xA3}

// ExtractWebMFromMVT strips the LucaSystem MVT wrapper and writes the embedded WebM.
func ExtractWebMFromMVT(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}
	payload, err := WebMPayload(data)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0777); err != nil {
		return err
	}
	return os.WriteFile(outputPath, payload, 0666)
}

func WebMPayload(data []byte) ([]byte, error) {
	if len(data) < 4 || !bytes.Equal(data[:4], []byte{'M', 'V', 'T', 0}) {
		return nil, fmt.Errorf("not an MVT movie")
	}

	offset := bytes.Index(data, ebmlMagic)
	if offset < 0 {
		return nil, fmt.Errorf("MVT movie does not contain a WebM/EBML payload")
	}

	return data[offset:], nil
}
