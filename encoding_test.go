package decaf377

import (
	"encoding/hex"
	"testing"
)

func TestGeneratorEncodingRoundTrip(t *testing.T) {
	generator, err := Generator()
	if err != nil {
		t.Fatalf("generator: %v", err)
	}
	encoding, err := Encode(generator)
	if err != nil {
		t.Fatalf("encode generator: %v", err)
	}
	if got, want := hex.EncodeToString(encoding), "0800000000000000000000000000000000000000000000000000000000000000"; got != want {
		t.Fatalf("generator encoding mismatch: got %s want %s", got, want)
	}
	decoded, err := Decode(encoding)
	if err != nil {
		t.Fatalf("decode generator: %v", err)
	}
	if !Equal(generator, decoded) {
		t.Fatalf("generator roundtrip mismatch")
	}
}

func TestDecodeRejectsNegativeEncoding(t *testing.T) {
	bytes := make([]byte, ElementSize)
	bytes[0] = 1
	if _, err := Decode(bytes); err == nil {
		t.Fatalf("expected negative s encoding to fail")
	}
}

func TestDecodeRejectsHighBits(t *testing.T) {
	bytes := make([]byte, ElementSize)
	bytes[31] = 0xe0
	if _, err := Decode(bytes); err == nil {
		t.Fatalf("expected high-bit encoding to fail")
	}
}
