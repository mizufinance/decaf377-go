package decaf377

import "math/big"

func ScalarFromCanonicalBytes(bytes []byte) (*big.Int, error) {
	if len(bytes) != ScalarSize {
		return nil, ErrInvalidScalar
	}
	v := littleEndianToBigInt(bytes)
	if v.Cmp(scalarOrder) >= 0 {
		return nil, ErrInvalidScalar
	}
	return v, nil
}

func ScalarFromUniformBytes(bytes []byte) *big.Int {
	return new(big.Int).Mod(littleEndianToBigInt(bytes), scalarOrder)
}
