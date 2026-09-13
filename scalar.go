package decaf377

import (
	"encoding/binary"
	"math/bits"

	"github.com/mizufinance/decaf377-go/internal/fiat"
)

// Scalar is a fixed-width little-endian integer. ScalarMul consumes all 256 bits.
// Canonical signature scalars must be parsed with ScalarFromCanonicalBytes.
type Scalar [ScalarSize]byte

func ScalarFromCanonicalBytes(bytes []byte) (Scalar, error) {
	if len(bytes) != ScalarSize {
		return Scalar{}, ErrInvalidScalar
	}
	var modulus [5]uint64
	fiat.FrMsat(&modulus)
	var borrow uint64
	for i := 0; i < 4; i++ {
		_, borrow = bits.Sub64(binary.LittleEndian.Uint64(bytes[8*i:8*i+8]), modulus[i], borrow)
	}
	if borrow != 1 {
		return Scalar{}, ErrInvalidScalar
	}
	return Scalar(bytes), nil
}

// ScalarFromUniformBytes reduces a little-endian integer with a schedule depending
// only on its public byte length. Use 64 bytes for hash-to-scalar material.
func ScalarFromUniformBytes(bytes []byte) Scalar {
	var accumulator, radix, digit [4]uint64
	radixWords := [4]uint64{256}
	fiat.FrToMontgomery(&radix, &radixWords)
	for i := len(bytes) - 1; i >= 0; i-- {
		words := [4]uint64{uint64(bytes[i])}
		fiat.FrToMontgomery(&digit, &words)
		fiat.FrMul(&accumulator, &accumulator, &radix)
		fiat.FrAdd(&accumulator, &accumulator, &digit)
	}
	var words [4]uint64
	fiat.FrFromMontgomery(&words, &accumulator)
	var out Scalar
	for i := range words {
		binary.LittleEndian.PutUint64(out[8*i:8*i+8], words[i])
	}
	return out
}
