package movie

import (
	"bytes"
	"testing"
)

func TestWebMPayload(t *testing.T) {
	data := append([]byte{'M', 'V', 'T', 0, 1, 2, 3}, ebmlMagic...)
	data = append(data, 4, 5, 6)

	payload, err := WebMPayload(data)
	if err != nil {
		t.Fatalf("WebMPayload returned error: %v", err)
	}

	want := append([]byte{}, ebmlMagic...)
	want = append(want, 4, 5, 6)
	if !bytes.Equal(payload, want) {
		t.Fatalf("payload mismatch: got %x, want %x", payload, want)
	}
}

func TestWebMPayloadRejectsNonMVT(t *testing.T) {
	if _, err := WebMPayload([]byte{'C', 'Z', '4', 0}); err == nil {
		t.Fatal("expected non-MVT data to fail")
	}
}

func TestWebMPayloadRejectsMissingWebM(t *testing.T) {
	if _, err := WebMPayload([]byte{'M', 'V', 'T', 0, 1, 2, 3}); err == nil {
		t.Fatal("expected MVT without EBML payload to fail")
	}
}
