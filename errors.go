package decaf377

import "errors"

var (
	ErrInvalidEncoding = errors.New("decaf377: invalid encoding")
	ErrInvalidScalar   = errors.New("decaf377: invalid scalar")
	ErrInvalidPoint    = errors.New("decaf377: invalid point")
)
