package decaf377

import (
	"crypto/subtle"
	"filippo.io/bigmod"
	"math/big"
)

// fieldElement holds a canonical, fixed-width big-endian residue.
type fieldElement [32]byte

var nativeModulus = func() *bigmod.Modulus {
	m, err := bigmod.NewModulus(fieldModulus.Bytes())
	if err != nil {
		panic(err)
	}
	return m
}()

func fieldNat(x fieldElement) *bigmod.Nat {
	n, err := bigmod.NewNat().SetBytes(x[:], nativeModulus)
	// Internal elements are canonical. This is an invariant check, not reduction.
	if err != nil {
		panic(err)
	}
	return n
}

func fieldFromNat(n *bigmod.Nat) (out fieldElement) {
	copy(out[:], n.Bytes(nativeModulus))
	return
}

func fieldUint(x uint) fieldElement {
	return fieldFromNat(bigmod.NewNat().SetUint(x).ExpandFor(nativeModulus))
}
func (x fieldElement) add(y fieldElement) fieldElement {
	return fieldFromNat(fieldNat(x).Add(fieldNat(y), nativeModulus))
}
func (x fieldElement) sub(y fieldElement) fieldElement {
	return fieldFromNat(fieldNat(x).Sub(fieldNat(y), nativeModulus))
}
func (x fieldElement) mul(y fieldElement) fieldElement {
	return fieldFromNat(fieldNat(x).Mul(fieldNat(y), nativeModulus))
}
func (x fieldElement) inverse() fieldElement {
	// Public exponent q - 2; Exp also has a fixed schedule for this byte length.
	exponent := inverseExponent
	return fieldFromNat(bigmod.NewNat().Exp(fieldNat(x), exponent[:], nativeModulus))
}
func fieldSelect(x, y fieldElement, choice byte) (out fieldElement) {
	for i := range out {
		out[i] = byte(subtle.ConstantTimeSelect(int(choice), int(y[i]), int(x[i])))
	}
	return
}

var inverseExponent = func() (out [32]byte) { new(big.Int).Sub(fieldModulus, big.NewInt(2)).FillBytes(out[:]); return }()

func bigmodFieldCheck(x fieldElement) (*bigmod.Nat, error) {
	return bigmod.NewNat().SetBytes(x[:], nativeModulus)
}

func (x fieldElement) pow(exponent []byte) fieldElement {
	return fieldFromNat(bigmod.NewNat().Exp(fieldNat(x), exponent, nativeModulus))
}
func (x fieldElement) abs() fieldElement {
	return fieldSelect(x, fieldElement{}.sub(x), x[31]&1)
}

var sqrtExponent = func() []byte {
	// (odd part of q - 1 minus 1) / 2; q has 2-adicity 47.
	q := new(big.Int).Sub(fieldModulus, big.NewInt(1))
	q.Rsh(q, 47).Sub(q, big.NewInt(1)).Rsh(q, 1)
	return q.Bytes()
}()
var sqrtRoot = func() (out fieldElement) {
	odd := new(big.Int).Sub(fieldModulus, big.NewInt(1))
	odd.Rsh(odd, 47)
	new(big.Int).Exp(zeta, odd, fieldModulus).FillBytes(out[:])
	return
}()
var nativeZeta = func() (out fieldElement) { zeta.FillBytes(out[:]); return }()

// Constant-time Tonelli-Shanks, RFC 9380 Appendix I.4. Exponents and loop
// bounds are public; nonsquares return a candidate checked by sqrtRatioZeta.
func (x fieldElement) sqrt() fieldElement {
	z := x.pow(sqrtExponent)
	t := z.mul(z).mul(x)
	z = z.mul(x)
	b, c := t, sqrtRoot
	one := fieldUint(1)
	for i := 47; i >= 2; i-- {
		for j := 1; j <= i-2; j++ {
			b = b.mul(b)
		}
		choice := byte(1 - subtle.ConstantTimeCompare(b[:], one[:]))
		z = fieldSelect(z, z.mul(c), choice)
		c = c.mul(c)
		t = fieldSelect(t, t.mul(c), choice)
		b = t
	}
	return z
}
func (x fieldElement) sqrtRatioZeta() fieldElement {
	root := x.sqrt()
	alternate := nativeZeta.mul(x).sqrt()
	square := root.mul(root)
	return fieldSelect(alternate, root, byte(subtle.ConstantTimeCompare(square[:], x[:])))
}
