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

// CompressToField converts the encoded result to big.Int. This final conversion is
// variable-time; use CompressToFieldBytes when the result must remain secret.
func CompressToField(point Point) (*big.Int, error) {
	bytes, err := CompressToFieldBytes(point)
	if err != nil {
		return nil, err
	}
	return littleEndianToBigInt(bytes[:]), nil
}

// CompressToFieldBytes keeps intermediate arithmetic and the result at fixed width.
func CompressToFieldBytes(point Point) ([32]byte, error) {
	if !point.Valid() {
		return [32]byte{}, ErrInvalidPoint
	}
	x, y := point.x, point.y
	t := x.mul(y)
	u1 := x.add(t).mul(x.sub(t))
	aMinusD := fieldElement{}.sub(fieldUint(3022))
	v := u1.mul(aMinusD).mul(x.mul(x)).inverse().sqrtRatioZeta()
	u2 := v.mul(u1).abs()
	u3 := u2.sub(t)
	s := aMinusD.mul(v).mul(u3).mul(x).abs()
	var out [32]byte
	for i := range out {
		out[i] = s[31-i]
	}
	return out, nil
}

func Encode(point Point) ([]byte, error) {
	s, err := CompressToFieldBytes(point)
	if err != nil {
		return nil, err
	}
	return s[:], nil
}
