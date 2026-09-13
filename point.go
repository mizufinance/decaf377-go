package decaf377

import (
	"crypto/subtle"
	"math/big"
)

// Point stores canonical affine coordinates at a fixed width.
// The zero value is invalid; use Identity for the group identity.
type Point struct{ x, y fieldElement }

// NewPoint imports public coordinates. big.Int conversion is variable-time.
func NewPoint(x, y *big.Int) Point {
	var p Point
	mod(x).FillBytes(p.x[:])
	mod(y).FillBytes(p.y[:])
	return p
}

// AffineBytes returns fixed-width little-endian coordinates, including for secret points.
func (p Point) AffineBytes() (x, y [32]byte) {
	for i := range x {
		x[i], y[i] = p.x[31-i], p.y[31-i]
	}
	return
}

// PointFromAffineBytes validates canonical little-endian coordinates.
func PointFromAffineBytes(x, y [32]byte) (Point, error) {
	var p Point
	for i := range x {
		p.x[i], p.y[i] = x[31-i], y[31-i]
	}
	if !fieldCanonical(p.x) {
		return Point{}, ErrInvalidPoint
	}
	if !fieldCanonical(p.y) {
		return Point{}, ErrInvalidPoint
	}
	if !p.Valid() {
		return Point{}, ErrInvalidPoint
	}
	return p, nil
}

func Identity() Point           { return Point{y: fieldUint(1)} }
func Generator() (Point, error) { var enc [32]byte; enc[0] = 8; return Decode(enc[:]) }

func (p Point) Valid() bool {
	x2, y2 := p.x.mul(p.x), p.y.mul(p.y)
	left := y2.sub(x2)
	right := fieldUint(1).add(fieldUint(3021).mul(x2).mul(y2))
	return subtle.ConstantTimeCompare(left[:], right[:]) == 1
}

// extendedPoint uses complete a=-1 twisted Edwards formulas without inversions.
type extendedPoint struct{ x, y, z, t fieldElement }

func (p Point) extended() extendedPoint { return extendedPoint{p.x, p.y, fieldUint(1), p.x.mul(p.y)} }
func (p extendedPoint) affine() Point {
	inverse := p.z.inverse()
	return Point{p.x.mul(inverse), p.y.mul(inverse)}
}
func (p extendedPoint) add(q extendedPoint) extendedPoint {
	a, b := p.x.mul(q.x), p.y.mul(q.y)
	c, d := fieldUint(3021).mul(p.t).mul(q.t), p.z.mul(q.z)
	e := p.x.add(p.y).mul(q.x.add(q.y)).sub(a).sub(b)
	f, g, h := d.sub(c), d.add(c), b.add(a)
	return extendedPoint{e.mul(f), g.mul(h), f.mul(g), e.mul(h)}
}
func selectPoint(p, q extendedPoint, bit byte) extendedPoint {
	return extendedPoint{fieldSelect(p.x, q.x, bit), fieldSelect(p.y, q.y, bit), fieldSelect(p.z, q.z, bit), fieldSelect(p.t, q.t, bit)}
}
func Add(left, right Point) (Point, error) {
	if !left.Valid() || !right.Valid() {
		return Point{}, ErrInvalidPoint
	}
	return left.extended().add(right.extended()).affine(), nil
}
func Neg(p Point) Point                    { return Point{fieldElement{}.sub(p.x), p.y} }
func Sub(left, right Point) (Point, error) { return Add(left, Neg(right)) }

// ScalarMul multiplies by the full 256-bit little-endian integer without reducing
// or truncating it. Arithmetic and selection use a fixed schedule for valid points.
func ScalarMul(base Point, scalar Scalar) (Point, error) {
	if !base.Valid() {
		return Point{}, ErrInvalidPoint
	}
	result, current := Identity().extended(), base.extended()
	for i := 0; i < 256; i++ {
		sum := result.add(current)
		result = selectPoint(result, sum, (scalar[i/8]>>uint(i%8))&1)
		current = current.add(current)
	}
	return result.affine(), nil
}
func Equivalent(left, right Point) bool {
	a, b := left.x.mul(right.y), right.x.mul(left.y)
	return subtle.ConstantTimeCompare(a[:], b[:]) == 1
}
func Equal(left, right Point) bool {
	return (subtle.ConstantTimeCompare(left.x[:], right.x[:]) & subtle.ConstantTimeCompare(left.y[:], right.y[:])) == 1
}
