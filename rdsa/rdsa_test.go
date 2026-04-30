package rdsa

import (
	"encoding/hex"
	"testing"
)

func TestVerifySpendAuthRustVector(t *testing.T) {
	pubkey := mustDecodeHex(t, "48b01e513dd37d94c3b48940dc133b92ccba7f546e99d3fc2e602d284f609f00")
	msg := mustDecodeHex(t, "70656e756d627261207264736120696e7465726f70206d657373616765")
	sig := mustDecodeHex(t, "506c8f3e38bd6d47088a7de148b1b6361e2abd9b692b2fc060a4e6312ba97904e9a98b503d8ad1c848abd5e863c18c6ffecf4dae7fd930c7f8c7bc480c581c03")

	ok, err := VerifySpendAuth(pubkey, msg, sig)
	if err != nil {
		t.Fatalf("verify valid vector: %v", err)
	}
	if !ok {
		t.Fatalf("expected valid vector to verify")
	}

	tampered := append([]byte(nil), sig...)
	tampered[0] ^= 0x01
	ok, err = VerifySpendAuth(pubkey, msg, tampered)
	if err == nil && ok {
		t.Fatalf("expected tampered signature to fail")
	}

	ok, err = VerifySpendAuth(pubkey, []byte("wrong message"), sig)
	if err == nil && ok {
		t.Fatalf("expected wrong message to fail")
	}
}

func mustDecodeHex(t *testing.T, value string) []byte {
	t.Helper()
	out, err := hex.DecodeString(value)
	if err != nil {
		t.Fatalf("decode hex %q: %v", value, err)
	}
	return out
}
