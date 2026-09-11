package decaf377

import (
	"math/big"
	"testing"
)

func BenchmarkScalarMul(b *testing.B) {
	generator, err := Generator()
	if err != nil {
		b.Fatal(err)
	}
	for name, scalar := range map[string]*big.Int{"zero": big.NewInt(0), "one": big.NewInt(1), "dense": new(big.Int).Sub(ScalarOrder(), big.NewInt(1))} {
		encoded := bigIntToLittleEndian32(scalar)
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if _, err := ScalarMul(generator, Scalar(encoded)); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
