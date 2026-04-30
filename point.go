package decaf377

import "math/big"

type Point struct {
	X *big.Int
	Y *big.Int
}

func NewPoint(x, y *big.Int) Point {
	return Point{X: mod(x), Y: mod(y)}
}

func Identity() Point {
	return Point{X: big.NewInt(0), Y: big.NewInt(1)}
}

func Generator() (Point, error) {
	var enc [32]byte
	enc[0] = 8
	return Decode(enc[:])
}

func (p Point) Valid() bool {
	x2 := square(p.X)
	y2 := square(p.Y)
	left := add(mul(curveA, x2), y2)
	right := add(big.NewInt(1), mul(curveD, x2, y2))
	return left.Cmp(right) == 0
}

func Add(left, right Point) (Point, error) {
	x1, y1 := mod(left.X), mod(left.Y)
	x2, y2 := mod(right.X), mod(right.Y)

	x1x2 := mul(x1, x2)
	y1y2 := mul(y1, y2)
	dxxyy := mul(curveD, x1x2, y1y2)

	xNum := add(mul(x1, y2), mul(y1, x2))
	xDen, err := inv(add(big.NewInt(1), dxxyy))
	if err != nil {
		return Point{}, err
	}
	x := mul(xNum, xDen)

	yNum := sub(y1y2, mul(curveA, x1x2))
	yDen, err := inv(sub(big.NewInt(1), dxxyy))
	if err != nil {
		return Point{}, err
	}
	y := mul(yNum, yDen)

	return NewPoint(x, y), nil
}

func Neg(p Point) Point {
	return NewPoint(neg(p.X), p.Y)
}

func Sub(left, right Point) (Point, error) {
	return Add(left, Neg(right))
}

func ScalarMul(base Point, scalar *big.Int) (Point, error) {
	result := Identity()
	current := NewPoint(base.X, base.Y)
	for i := 0; i < fieldBits; i++ {
		if scalar.Bit(i) == 1 {
			var err error
			result, err = Add(result, current)
			if err != nil {
				return Point{}, err
			}
		}
		var err error
		current, err = Add(current, current)
		if err != nil {
			return Point{}, err
		}
	}
	return result, nil
}

func Equivalent(left, right Point) bool {
	return mul(left.X, right.Y).Cmp(mul(right.X, left.Y)) == 0
}

func Equal(left, right Point) bool {
	return mod(left.X).Cmp(mod(right.X)) == 0 && mod(left.Y).Cmp(mod(right.Y)) == 0
}
