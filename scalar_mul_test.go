package decaf377

import (
	"bytes"
	"math/big"
	"math/rand"
	"testing"
)

// The independent oracle deliberately uses affine big.Int arithmetic, never the
// production extended-point engine or fixed-width field operations.
func referenceAdd(p, q [2]*big.Int) [2]*big.Int {
	x1x2, y1y2 := mul(p[0], q[0]), mul(p[1], q[1])
	dxy := mul(curveD, x1x2, y1y2)
	dx, _ := inv(add(big.NewInt(1), dxy))
	dy, _ := inv(sub(big.NewInt(1), dxy))
	return [2]*big.Int{mul(add(mul(p[0], q[1]), mul(p[1], q[0])), dx), mul(add(y1y2, x1x2), dy)}
}
func referenceMul(base Point, scalar Scalar) Point {
	acc := [2]*big.Int{big.NewInt(0), big.NewInt(1)}
	cur := [2]*big.Int{new(big.Int).SetBytes(base.x[:]), new(big.Int).SetBytes(base.y[:])}
	for i := 0; i < 256; i++ {
		if scalar[i/8]&(1<<uint(i%8)) != 0 {
			acc = referenceAdd(acc, cur)
		}
		cur = referenceAdd(cur, cur)
	}
	return NewPoint(acc[0], acc[1])
}
func TestScalarMulMatchesIndependentAffineOracle(t *testing.T) {
	g, err := Generator()
	if err != nil {
		t.Fatal(err)
	}
	rng := rand.New(rand.NewSource(377))
	cases := []Scalar{{}, {1}, {0: 255}, {31: 128}, {31: 32}}
	minusOne := bigIntToLittleEndian32(new(big.Int).Sub(scalarOrder, big.NewInt(1)))
	cases = append(cases, Scalar(minusOne))
	for i := 0; i < 16; i++ {
		var s Scalar
		rng.Read(s[:])
		cases = append(cases, s)
	}
	for _, base := range []Point{Identity(), g, Neg(g)} {
		for _, scalar := range cases {
			got, err := ScalarMul(base, scalar)
			if err != nil {
				t.Fatal(err)
			}
			want := referenceMul(base, scalar)
			if !Equal(got, want) {
				t.Fatalf("scalar %x mismatch", scalar)
			}
			encoded, err := Encode(got)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := Decode(encoded)
			if err != nil {
				t.Fatal(err)
			}
			if !Equivalent(decoded, got) {
				t.Fatalf("scalar %x encoding mismatch", scalar)
			}
		}
	}
}
func TestFieldArithmeticAgainstBigInt(t *testing.T) {
	rng := rand.New(rand.NewSource(378))
	for i := 0; i < 64; i++ {
		var a, b fieldElement
		rng.Read(a[:])
		rng.Read(b[:])
		av, bv := mod(new(big.Int).SetBytes(a[:])), mod(new(big.Int).SetBytes(b[:]))
		av.FillBytes(a[:])
		bv.FillBytes(b[:])
		for _, pair := range []struct {
			got  fieldElement
			want *big.Int
		}{{a.add(b), add(av, bv)}, {a.sub(b), sub(av, bv)}, {a.mul(b), mul(av, bv)}} {
			if new(big.Int).SetBytes(pair.got[:]).Cmp(pair.want) != 0 {
				t.Fatal("field mismatch")
			}
		}
		inverse := a.inverse()
		if av.Sign() != 0 && new(big.Int).SetBytes(inverse[:]).Cmp(new(big.Int).ModInverse(av, fieldModulus)) != 0 {
			t.Fatal("inverse mismatch")
		}
		square := a.mul(a)
		root := square.sqrt()
		if root.mul(root) != square {
			t.Fatal("sqrt mismatch")
		}
	}
}
func TestScalarCanonicalAndUniform(t *testing.T) {
	if _, err := ScalarFromCanonicalBytes(make([]byte, 31)); err == nil {
		t.Fatal("short scalar accepted")
	}
	order := bigIntToLittleEndian32(scalarOrder)
	if _, err := ScalarFromCanonicalBytes(order[:]); err == nil {
		t.Fatal("order accepted")
	}
	for _, width := range []int{0, 1, 32, 48, 64, 96} {
		input := bytes.Repeat([]byte{255}, width)
		got := ScalarFromUniformBytes(input)
		want := new(big.Int).Mod(littleEndianToBigInt(input), scalarOrder)
		if littleEndianToBigInt(got[:]).Cmp(want) != 0 {
			t.Fatal("uniform reduction mismatch")
		}
	}
}
func TestPointValidationAndIdentity(t *testing.T) {
	if _, err := ScalarMul(Point{}, Scalar{1}); err == nil {
		t.Fatal("invalid point accepted")
	}
	g, _ := Generator()
	x, y := g.AffineBytes()
	p, err := PointFromAffineBytes(x, y)
	if err != nil || !Equal(p, g) {
		t.Fatal("affine import mismatch")
	}
	bad := bigIntToLittleEndian32(new(big.Int).Sub(fieldModulus, big.NewInt(1)))
	for i := range bad {
		bad[i] = 255
	}
	if _, err := PointFromAffineBytes(bad, y); err == nil {
		t.Fatal("noncanonical coordinate accepted")
	}
	sum, err := Add(g, Neg(g))
	if err != nil || !Equivalent(sum, Identity()) {
		t.Fatal("inverse addition")
	}
	encoded, err := Encode(Identity())
	if err != nil || !bytes.Equal(encoded, make([]byte, 32)) {
		t.Fatal("identity encoding")
	}
}
