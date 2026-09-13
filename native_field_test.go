package decaf377

import (
	"math/big"
	"math/rand"
	"testing"
)

var nativeFieldSink fieldElement
var nativeScalarSink Scalar
var nativePointSink Point
var nativeEncodingSink [32]byte

func TestNativeFieldArithmetic(t *testing.T) {
	rng := rand.New(rand.NewSource(377))
	values := []*big.Int{new(big.Int), big.NewInt(1), new(big.Int).Sub(fieldModulus, big.NewInt(1))}
	var maximum fieldElement
	values[2].FillBytes(maximum[:])
	if fieldUint(0).sub(fieldUint(1)) != maximum || maximum.add(fieldUint(1)) != (fieldElement{}) {
		t.Fatal("field borrow/carry boundary")
	}
	for i := 0; i < 32; i++ {
		var x fieldElement
		x[i] = 1
		if fieldFromWords(fieldWords(x)) != x {
			t.Fatal("field endian round trip")
		}
	}
	for i := 0; i < 40; i++ {
		values = append(values, new(big.Int).Rand(rng, fieldModulus))
	}
	for i, a := range values {
		b := values[(i+2)%len(values)]
		var x, y fieldElement
		a.FillBytes(x[:])
		b.FillBytes(y[:])
		check := func(name string, got fieldElement, want *big.Int) {
			t.Helper()
			want.Mod(want, fieldModulus)
			if new(big.Int).SetBytes(got[:]).Cmp(want) != 0 || !fieldCanonical(got) {
				t.Fatalf("%s case %d: got %x want %x", name, i, got, want)
			}
		}
		check("add", x.add(y), new(big.Int).Add(a, b))
		check("sub", x.sub(y), new(big.Int).Sub(a, b))
		check("mul", x.mul(y), new(big.Int).Mul(a, b))
		check("inverse", x.inverse(), new(big.Int).Exp(a, new(big.Int).Sub(fieldModulus, big.NewInt(2)), fieldModulus))
		check("pow", x.pow([]byte{0x81, 0, 0x7f}), new(big.Int).Exp(a, big.NewInt(0x81007f), fieldModulus))
		check("empty exponent", x.pow(nil), big.NewInt(1))
	}
	for _, a := range []*big.Int{new(big.Int).Set(fieldModulus), new(big.Int).Add(fieldModulus, big.NewInt(1)), new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))} {
		var x fieldElement
		a.FillBytes(x[:])
		if fieldCanonical(x) {
			t.Fatal("accepted noncanonical field element")
		}
	}
}

func TestNativeScalarReduction(t *testing.T) {
	rng := rand.New(rand.NewSource(378))
	for _, size := range []int{0, 1, 31, 32, 33, 64, 65, 256} {
		for iteration := 0; iteration < 12; iteration++ {
			input := make([]byte, size)
			rng.Read(input)
			if iteration == 0 {
				for i := range input {
					input[i] = 0xff
				}
			}
			be := make([]byte, size)
			for i := range input {
				be[size-1-i] = input[i]
			}
			want := new(big.Int).Mod(new(big.Int).SetBytes(be), scalarOrder)
			got := ScalarFromUniformBytes(input)
			var gotBE [32]byte
			for i := range got {
				gotBE[31-i] = got[i]
			}
			if new(big.Int).SetBytes(gotBE[:]).Cmp(want) != 0 {
				t.Fatalf("reduction size %d", size)
			}
			if parsed, err := ScalarFromCanonicalBytes(got[:]); err != nil || parsed != got {
				t.Fatal("canonical round trip")
			}
		}
	}
	for _, value := range []*big.Int{new(big.Int).Sub(scalarOrder, big.NewInt(1)), new(big.Int).Set(scalarOrder), new(big.Int).Add(scalarOrder, big.NewInt(1))} {
		var be, le [32]byte
		value.FillBytes(be[:])
		for i := range be {
			le[31-i] = be[i]
		}
		_, err := ScalarFromCanonicalBytes(le[:])
		if (err == nil) != (value.Cmp(scalarOrder) < 0) {
			t.Fatal("scalar canonical boundary")
		}
	}
}

func TestNativeFieldAllocations(t *testing.T) {
	x, y := fieldUint(12), fieldUint(35)
	var uniform [64]byte
	uniform[63] = 0xff
	for name, operation := range map[string]func(){
		"add":              func() { nativeFieldSink = x.add(y) },
		"sub":              func() { nativeFieldSink = x.sub(y) },
		"mul":              func() { nativeFieldSink = x.mul(y) },
		"inverse":          func() { nativeFieldSink = x.inverse() },
		"sqrt":             func() { nativeFieldSink = x.sqrt() },
		"select":           func() { nativeFieldSink = fieldSelect(x, y, 1) },
		"scalar reduction": func() { nativeScalarSink = ScalarFromUniformBytes(uniform[:]) },
	} {
		t.Run(name, func(t *testing.T) {
			if count := testing.AllocsPerRun(20, operation); count != 0 {
				t.Fatalf("allocations = %g", count)
			}
		})
	}
}

func TestNativeGroupAllocations(t *testing.T) {
	base, err := Generator()
	if err != nil {
		t.Fatal(err)
	}
	scalar := Scalar{1, 2, 3, 4, 5}
	for name, operation := range map[string]func(){
		"scalar multiplication": func() {
			var err error
			nativePointSink, err = ScalarMul(base, scalar)
			if err != nil {
				panic(err)
			}
		},
		"compression": func() {
			var err error
			nativeEncodingSink, err = CompressToFieldBytes(base)
			if err != nil {
				panic(err)
			}
		},
	} {
		t.Run(name, func(t *testing.T) {
			if count := testing.AllocsPerRun(5, operation); count != 0 {
				t.Fatalf("allocations = %g", count)
			}
		})
	}
}
