package decaf377

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
)

const (
	ElementSize = 32
	ScalarSize  = 32
	FieldBits   = 253
)

const fieldBits = FieldBits

var (
	fieldModulus = new(big.Int).Set(ecc.BLS12_377.ScalarField())

	scalarOrder = mustDecimal("2111115437357092606062206234695386632838870926408408195193685246394721360383")
	curveD      = mustDecimal("3021")
	curveA      = new(big.Int).Sub(fieldModulus, big.NewInt(1))
	aMinusD     = mustDecimal("8444461749428370424248824938781546531375899335154063827935233455917409236019")
	zeta        = mustDecimal("2841681278031794617739547238867782961338435681360110683443920362658525667816")
)

func mustDecimal(s string) *big.Int {
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("invalid decimal constant")
	}
	return v
}

func FieldModulus() *big.Int {
	return new(big.Int).Set(fieldModulus)
}

func ScalarOrder() *big.Int {
	return new(big.Int).Set(scalarOrder)
}

func CurveA() *big.Int {
	return new(big.Int).Set(curveA)
}

func CurveD() *big.Int {
	return new(big.Int).Set(curveD)
}

func CurveAMinusD() *big.Int {
	return new(big.Int).Set(aMinusD)
}

func CurveZeta() *big.Int {
	return new(big.Int).Set(zeta)
}
