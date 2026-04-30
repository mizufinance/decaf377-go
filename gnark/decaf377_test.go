package gnark

import (
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	gnarkte "github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/test"
	decaf377 "github.com/mizufinance/decaf377-go"
)

type compressToFieldCircuit struct {
	X frontend.Variable
	Y frontend.Variable

	Expected frontend.Variable `gnark:",public"`
}

func (c *compressToFieldCircuit) Define(api frontend.API) error {
	result, err := CompressToField(api, gnarkte.Point{X: c.X, Y: c.Y})
	if err != nil {
		return err
	}
	api.AssertIsEqual(result, c.Expected)
	return nil
}

type encodeToCurveCircuit struct {
	Input frontend.Variable

	ExpectedX        frontend.Variable `gnark:",public"`
	ExpectedY        frontend.Variable `gnark:",public"`
	ExpectedCompress frontend.Variable `gnark:",public"`
}

func (c *encodeToCurveCircuit) Define(api frontend.API) error {
	point, err := EncodeToCurve(api, c.Input)
	if err != nil {
		return err
	}
	api.AssertIsEqual(point.X, c.ExpectedX)
	api.AssertIsEqual(point.Y, c.ExpectedY)

	compressed, err := CompressToField(api, point)
	if err != nil {
		return err
	}
	api.AssertIsEqual(compressed, c.ExpectedCompress)
	return nil
}

type isqrtZeroCircuit struct {
	ExpectedWasSquare frontend.Variable `gnark:",public"`
}

func (c *isqrtZeroCircuit) Define(api frontend.API) error {
	wasSquare, _, err := decaf377Isqrt(api, 0)
	if err != nil {
		return err
	}
	api.AssertIsEqual(wasSquare, c.ExpectedWasSquare)
	return nil
}

func TestCompressToFieldMatchesNative(t *testing.T) {
	generator, err := decaf377.Generator()
	if err != nil {
		t.Fatalf("generator: %v", err)
	}
	expected, err := decaf377.CompressToField(generator)
	if err != nil {
		t.Fatalf("compress generator: %v", err)
	}

	assignment := &compressToFieldCircuit{
		X:        generator.X,
		Y:        generator.Y,
		Expected: expected,
	}
	if err := test.IsSolved(&compressToFieldCircuit{}, assignment, ecc.BLS12_377.ScalarField()); err != nil {
		t.Fatalf("compress_to_field circuit: %v", err)
	}
}

func TestEncodeToCurveMatchesNative(t *testing.T) {
	input := big.NewInt(123456789)
	point, err := EncodeToCurveNative(input)
	if err != nil {
		t.Fatalf("encode_to_curve native: %v", err)
	}
	expected, err := CompressToFieldNative(point)
	if err != nil {
		t.Fatalf("compress encoded point: %v", err)
	}

	assignment := &encodeToCurveCircuit{
		Input:            input,
		ExpectedX:        point.X,
		ExpectedY:        point.Y,
		ExpectedCompress: expected,
	}
	if err := test.IsSolved(&encodeToCurveCircuit{}, assignment, ecc.BLS12_377.ScalarField()); err != nil {
		t.Fatalf("encode_to_curve circuit: %v", err)
	}
}

func TestEncodeToCurveCompiles(t *testing.T) {
	if _, err := frontend.Compile(ecc.BLS12_377.ScalarField(), r1cs.NewBuilder, &encodeToCurveCircuit{}); err != nil {
		t.Fatalf("compile encode_to_curve circuit: %v", err)
	}
}

func TestIsqrtZeroAllowsNonSquareBranch(t *testing.T) {
	assignment := &isqrtZeroCircuit{ExpectedWasSquare: 0}
	if err := test.IsSolved(&isqrtZeroCircuit{}, assignment, ecc.BLS12_377.ScalarField()); err != nil {
		t.Fatalf("expected den=0 isqrt to solve with wasSquare=0: %v", err)
	}
}

func TestIsqrtZeroRejectsSquareBranch(t *testing.T) {
	assignment := &isqrtZeroCircuit{ExpectedWasSquare: 1}
	if err := test.IsSolved(&isqrtZeroCircuit{}, assignment, ecc.BLS12_377.ScalarField()); err == nil {
		t.Fatalf("expected den=0 isqrt to reject wasSquare=1")
	}
}
