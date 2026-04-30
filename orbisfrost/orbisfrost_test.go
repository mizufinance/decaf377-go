package orbisfrost

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mizufinance/decaf377-go"
)

func TestOrbisRustVectors(t *testing.T) {
	path := filepath.Join("..", "testdata", "orbis_vectors.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read vectors: %v", err)
	}

	var vectors decaf377.OrbisFrostVectorFile
	if err := json.Unmarshal(data, &vectors); err != nil {
		t.Fatalf("decode vectors: %v", err)
	}
	if vectors.Schema != "orbis-frost-decaf377-v1" {
		t.Fatalf("unexpected schema %q", vectors.Schema)
	}
	if len(vectors.Vectors) == 0 {
		t.Fatalf("missing vectors")
	}

	for _, vector := range vectors.Vectors {
		vector := vector
		t.Run(vector.Name, func(t *testing.T) {
			pubkey := mustDecodeHex(t, vector.PublicKeyHex)
			msg := mustDecodeHex(t, vector.MessageHex)
			sig := mustDecodeHex(t, vector.SignatureHex)

			ok, err := Verify(pubkey, msg, sig)
			if vector.ExpectValid && err != nil {
				t.Fatalf("expected valid vector, got error: %v", err)
			}
			if ok != vector.ExpectValid {
				t.Fatalf("verification mismatch: got %v want %v err=%v", ok, vector.ExpectValid, err)
			}
		})
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
