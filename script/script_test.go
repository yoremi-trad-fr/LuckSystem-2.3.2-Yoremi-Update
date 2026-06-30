package script

import (
	"bytes"
	"testing"

	"lucksystem/charset"
	"lucksystem/game/enum"
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

func TestSetOperateParamsImportByteAcceptsShortHex(t *testing.T) {
	s := &Script{
		Info: Info{Name: "test"},
		Codes: []*CodeLine{{
			OpStr:  "MESSAGE",
			Params: []interface{}{"0x0"},
		}},
	}

	err := s.SetOperateParams(0, enum.VMRunImport, byte(0), []bool{true})

	if err != nil {
		t.Fatalf("SetOperateParams returned error: %v", err)
	}
	if got, want := s.Codes[0].RawBytes, []byte{0x00}; !bytes.Equal(got, want) {
		t.Fatalf("raw bytes = % X, want % X", got, want)
	}
}

func TestSetOperateParamsImportByteReportsInvalidValue(t *testing.T) {
	s := &Script{
		Info: Info{Name: "test"},
		Codes: []*CodeLine{{
			OpStr:  "MESSAGE",
			Params: []interface{}{"0x"},
		}},
	}

	err := s.SetOperateParams(0, enum.VMRunImport, byte(0), []bool{true})

	if err == nil {
		t.Fatal("SetOperateParams returned nil error")
	}
}

func TestParseCodeParamsHandlesEscapedQuotes(t *testing.T) {
	code := &CodeLine{}

	ParseCodeParams(code, `MESSAGE (1, "「おかえりなさい」", "\"Bon retour à la maison.\"", "cn", 36363, 8, 0x0)`)

	if got, want := code.OpStr, "MESSAGE"; got != want {
		t.Fatalf("opcode = %q, want %q", got, want)
	}
	if got, want := len(code.Params), 7; got != want {
		t.Fatalf("param count = %d, want %d: %#v", got, want, code.Params)
	}
	if got, want := code.Params[2], `"Bon retour à la maison."`; got != want {
		t.Fatalf("param 2 = %#v, want %#v", got, want)
	}
	if got, want := code.Params[4], "36363"; got != want {
		t.Fatalf("param 4 = %#v, want %#v", got, want)
	}
}

func TestCodeParamsStringRoundTripEscapedQuotes(t *testing.T) {
	src := &CodeLine{
		OpStr:  "MESSAGE",
		Params: []interface{}{`C'était "la colère".`},
	}
	line := ToStringCodeParams(src)
	dst := &CodeLine{}

	ParseCodeParams(dst, line)

	if got, want := dst.Params[0], src.Params[0]; got != want {
		t.Fatalf("round trip = %#v, want %#v (line: %s)", got, want, line)
	}
}
