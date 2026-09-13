package decaf377

import (
	"math/big"
	"testing"
)

// Independent integer equation: never trust the cached Valid() in this oracle.
func referenceMembership(p Point) bool {
	x, y := new(big.Int).SetBytes(p.x[:]), new(big.Int).SetBytes(p.y[:])
	x2 := new(big.Int).Mul(x, x)
	y2 := new(big.Int).Mul(y, y)
	left := new(big.Int).Sub(y2, x2)
	left.Mod(left, fieldModulus)
	right := new(big.Int).Mul(x2, y2)
	right.Mul(right, big.NewInt(3021))
	right.Add(right, big.NewInt(1))
	right.Mod(right, fieldModulus)
	return left.Cmp(right) == 0
}

func TestPointMembershipInvariant(t *testing.T) {
	base, err := Generator()
	if err != nil {
		t.Fatal(err)
	}
	minusOne := new(big.Int).Sub(fieldModulus, big.NewInt(1))
	i := new(big.Int).ModSqrt(minusOne, fieldModulus)
	if i == nil {
		t.Fatal("missing order-four point")
	}
	points := []Point{
		{}, Identity(), base, Neg(base), NewPoint(big.NewInt(0), minusOne),
		NewPoint(i, big.NewInt(0)), NewPoint(big.NewInt(1), big.NewInt(1)),
		NewPoint(big.NewInt(2), big.NewInt(3)),
		NewPoint(new(big.Int).Add(fieldModulus, big.NewInt(1)), big.NewInt(-1)),
	}
	check := func(p Point) {
		t.Helper()
		if p.Valid() != referenceMembership(p) {
			t.Fatal("membership flag disagrees with coordinates")
		}
	}
	var zero, one, maximal Scalar
	one[0] = 1
	for index := range maximal {
		maximal[index] = 255
	}
	for _, point := range points {
		check(point)
		check(Neg(point))
		copy := point
		check(copy)
		x, y := point.AffineBytes()
		decoded, err := PointFromAffineBytes(x, y)
		if point.Valid() {
			if err != nil || !Equal(point, decoded) {
				t.Fatal("valid coordinate roundtrip failed")
			}
			check(decoded)
			sum, err := Add(point, Neg(point))
			if err != nil || !Equal(sum, Identity()) {
				t.Fatal("inverse sum failed")
			}
			check(sum)
			for _, scalar := range []Scalar{zero, one, maximal} {
				result, err := ScalarMul(point, scalar)
				if err != nil {
					t.Fatal(err)
				}
				check(result)
				if scalar == zero && !Equal(result, Identity()) {
					t.Fatal("zero scalar")
				}
				if scalar == one && !Equal(result, point) {
					t.Fatal("unit scalar")
				}
			}
		} else {
			if err == nil {
				t.Fatal("invalid coordinates accepted")
			}
			if _, err := Add(point, Identity()); err == nil {
				t.Fatal("invalid addition accepted")
			}
			if _, err := ScalarMul(point, one); err == nil {
				t.Fatal("invalid multiplication accepted")
			}
			if _, err := CompressToFieldBytes(point); err == nil {
				t.Fatal("invalid compression accepted")
			}
		}
	}
}
