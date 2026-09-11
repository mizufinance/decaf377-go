package decaf377

import "filippo.io/bigmod"

// Scalar is a fixed-width little-endian integer. ScalarMul consumes all 256 bits.
// Canonical signature scalars must be parsed with ScalarFromCanonicalBytes.
type Scalar [ScalarSize]byte

var nativeScalarModulus = func() *bigmod.Modulus {
	m, err := bigmod.NewModulus(scalarOrder.Bytes())
	if err != nil {
		panic(err)
	}
	return m
}()

func ScalarFromCanonicalBytes(bytes []byte) (Scalar, error) {
	if len(bytes) != ScalarSize {
		return Scalar{}, ErrInvalidScalar
	}
	var be [32]byte
	for i := range be {
		be[i] = bytes[31-i]
	}
	if _, err := bigmod.NewNat().SetBytes(be[:], nativeScalarModulus); err != nil {
		return Scalar{}, ErrInvalidScalar
	}
	return Scalar(bytes), nil
}

// ScalarFromUniformBytes reduces a little-endian integer with a schedule depending
// only on its public byte length. Use 64 bytes for hash-to-scalar material.
func ScalarFromUniformBytes(bytes []byte) Scalar {
	m := nativeScalarModulus
	n := bigmod.NewNat().SetUint(0).ExpandFor(m)
	radix := bigmod.NewNat().SetUint(256).ExpandFor(m)
	digit := bigmod.NewNat()
	for i := len(bytes) - 1; i >= 0; i-- {
		digit.SetUint(uint(bytes[i])).ExpandFor(m)
		n.Mul(radix, m).Add(digit, m)
	}
	be := n.Bytes(m)
	var out Scalar
	for i := range out {
		out[i] = be[31-i]
	}
	return out
}
