package decaf377

import "math/big"

func Decode(bytes []byte) (Point, error) {
	if len(bytes) != ElementSize {
		return Point{}, ErrInvalidEncoding
	}
	if bytes[31]>>5 != 0 {
		return Point{}, ErrInvalidEncoding
	}

	s := littleEndianToBigInt(bytes)
	if s.Cmp(fieldModulus) >= 0 || isNegative(s) {
		return Point{}, ErrInvalidEncoding
	}

	ss := square(s)
	u1 := sub(big.NewInt(1), ss)
	u2 := sub(square(u1), mul(big.NewInt(4), curveD, ss))

	wasSquare, v, err := sqrtRatioZetaDen(mul(u2, square(u1)))
	if err != nil {
		return Point{}, err
	}
	if !wasSquare {
		return Point{}, ErrInvalidEncoding
	}

	twoSU1 := mul(big.NewInt(2), s, u1)
	if isNegative(mul(twoSU1, v)) {
		v = neg(v)
	}

	x := mul(twoSU1, square(v), u2)
	y := mul(add(big.NewInt(1), ss), v, u1)
	point := NewPoint(x, y)
	if !point.Valid() {
		return Point{}, ErrInvalidPoint
	}
	return point, nil
}

func CompressToField(point Point) (*big.Int, error) {
	if !point.Valid() {
		return nil, ErrInvalidPoint
	}
	x, y := mod(point.X), mod(point.Y)
	t := mul(x, y)

	u1 := mul(add(x, t), sub(x, t))
	_, v, err := sqrtRatioZetaDen(mul(u1, aMinusD, square(x)))
	if err != nil {
		return nil, err
	}

	u2 := abs(mul(v, u1))
	u3 := sub(u2, t)
	return abs(mul(aMinusD, v, u3, x)), nil
}

func Encode(point Point) ([]byte, error) {
	s, err := CompressToField(point)
	if err != nil {
		return nil, err
	}
	out := bigIntToLittleEndian32(s)
	out[31] &= 0b00011111
	return out[:], nil
}
