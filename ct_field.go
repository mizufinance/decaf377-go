package decaf377

import (
	"crypto/subtle"
	"encoding/binary"
	"math/big"
	"math/bits"

	"github.com/mizufinance/decaf377-go/internal/fiat"
)

// fieldElement holds a canonical, fixed-width big-endian residue.
type fieldElement [32]byte

func fieldWords(x fieldElement) (out [4]uint64) {
	for i := range out {
		out[i] = binary.BigEndian.Uint64(x[24-8*i : 32-8*i])
	}
	return
}

func fieldFromWords(x [4]uint64) (out fieldElement) {
	for i := range x {
		binary.BigEndian.PutUint64(out[24-8*i:32-8*i], x[i])
	}
	return
}

func fieldUint(x uint) fieldElement {
	return fieldFromWords([4]uint64{uint64(x)})
}
func (x fieldElement) add(y fieldElement) fieldElement {
	a, b := fieldWords(x), fieldWords(y)
	var out [4]uint64
	fiat.FqAdd(&out, &a, &b)
	return fieldFromWords(out)
}
func (x fieldElement) sub(y fieldElement) fieldElement {
	a, b := fieldWords(x), fieldWords(y)
	var out [4]uint64
	fiat.FqSub(&out, &a, &b)
	return fieldFromWords(out)
}
func (x fieldElement) mul(y fieldElement) fieldElement {
	a, b := fieldWords(x), fieldWords(y)
	var am, bm, out [4]uint64
	fiat.FqToMontgomery(&am, &a)
	fiat.FqToMontgomery(&bm, &b)
	fiat.FqMul(&out, &am, &bm)
	fiat.FqFromMontgomery(&a, &out)
	return fieldFromWords(a)
}
func (x fieldElement) inverse() fieldElement {
	exponent := inverseExponent
	return x.pow(exponent[:])
}
func fieldSelect(x, y fieldElement, choice byte) (out fieldElement) {
	for i := range out {
		out[i] = byte(subtle.ConstantTimeSelect(int(choice), int(y[i]), int(x[i])))
	}
	return
}

var inverseExponent = func() (out [32]byte) { new(big.Int).Sub(fieldModulus, big.NewInt(2)).FillBytes(out[:]); return }()

func fieldCanonical(x fieldElement) bool {
	a := fieldWords(x)
	var modulus [5]uint64
	fiat.FqMsat(&modulus)
	var borrow uint64
	for i := range a {
		_, borrow = bits.Sub64(a[i], modulus[i], borrow)
	}
	return borrow == 1
}

func (x fieldElement) pow(exponent []byte) fieldElement {
	a := fieldWords(x)
	var base, out, product [4]uint64
	fiat.FqToMontgomery(&base, &a)
	fiat.FqSetOne(&out)
	for _, digit := range exponent {
		for bit := 7; bit >= 0; bit-- {
			fiat.FqSquare(&out, &out)
			fiat.FqMul(&product, &out, &base)
			fiat.FqSelectznz(&out, fiat.FqUint1((digit>>bit)&1), &out, &product)
		}
	}
	fiat.FqFromMontgomery(&a, &out)
	return fieldFromWords(a)
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
