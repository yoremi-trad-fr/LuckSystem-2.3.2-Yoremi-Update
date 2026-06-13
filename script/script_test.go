package script

import (
	"bytes"
	"testing"

	"lucksystem/charset"
)

func TestCodeStringEmptyUTF8LStringWritesTerminator(t *testing.T) {
	var buf bytes.Buffer

	size := CodeString(&buf, "", true, charset.UTF_8)

	if size != 3 {
		t.Fatalf("size = %d, want 3", size)
	}
	if got, want := buf.Bytes(), []byte{0x00, 0x00, 0x00}; !bytes.Equal(got, want) {
		t.Fatalf("bytes = % X, want % X", got, want)
	}
}
