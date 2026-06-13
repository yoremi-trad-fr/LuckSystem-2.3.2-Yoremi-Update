package operator

import (
	"testing"

	"lucksystem/charset"
)

func TestGetParamEmptyUTF8LStringConsumesTerminator(t *testing.T) {
	data := []byte{0x00, 0x00, 0x00, 0x12, 0x00}
	var value lstring

	next := GetParam(data, &value, 0, 0, charset.UTF_8)

	if value != "" {
		t.Fatalf("value = %q, want empty string", value)
	}
	if next != 3 {
		t.Fatalf("next = %d, want 3", next)
	}
}
