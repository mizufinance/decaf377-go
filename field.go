package decaf377

import (
	"errors"
	"math/big"

	fr "github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
)

func mod(v *big.Int) *big.Int {
	out := new(big.Int).Mod(new(big.Int).Set(v), fieldModulus)
	if out.Sign() < 0 {
		out.Add(out, fieldModulus)
	}
	return out
}

func add(a, b *big.Int) *big.Int {
	return mod(new(big.Int).Add(a, b))
}

func sub(a, b *big.Int) *big.Int {
	return mod(new(big.Int).Sub(a, b))
}

func mul(values ...*big.Int) *big.Int {
	out := big.NewInt(1)
	for _, value := range values {
		out.Mul(out, value)
		out.Mod(out, fieldModulus)
	}
	return out
}

func square(v *big.Int) *big.Int {
	return mul(v, v)
}

func neg(v *big.Int) *big.Int {
	return mod(new(big.Int).Neg(v))
}

func inv(v *big.Int) (*big.Int, error) {
	out := new(big.Int).ModInverse(mod(v), fieldModulus)
	if out == nil {
		return nil, errors.New("decaf377: inverse does not exist")
	}
	return out, nil
}

func isNegative(v *big.Int) bool {
	return mod(v).Bit(0) == 1
}

func abs(v *big.Int) *big.Int {
	v = mod(v)
	if !isNegative(v) {
		return v
	}
	return neg(v)
}

func sqrtRatioZetaDen(den *big.Int) (bool, *big.Int, error) {
	var denEl, invDen, y, zetaEl, zetaInvDen, zero fr.Element
	denEl.SetBigInt(mod(den))
	if denEl.Equal(&zero) {
		return false, big.NewInt(0), nil
	}
	invDen.Inverse(&denEl)
	if y.Sqrt(&invDen) != nil {
		out := new(big.Int)
		y.BigInt(out)
		return true, out, nil
	}
	zetaEl.SetBigInt(zeta)
	zetaInvDen.Mul(&zetaEl, &invDen)
	if y.Sqrt(&zetaInvDen) == nil {
		return false, nil, errors.New("decaf377: sqrt_ratio_zeta failed")
	}
	out := new(big.Int)
	y.BigInt(out)
	return false, out, nil
}

func littleEndianToBigInt(le []byte) *big.Int {
	be := append([]byte(nil), le...)
	for i, j := 0, len(be)-1; i < j; i, j = i+1, j-1 {
		be[i], be[j] = be[j], be[i]
	}
	return new(big.Int).SetBytes(be)
}

func bigIntToLittleEndian32(v *big.Int) [32]byte {
	be := mod(v).Bytes()
	var out [32]byte
	for i := 0; i < len(be) && i < 32; i++ {
		out[i] = be[len(be)-1-i]
	}
	return out
}
