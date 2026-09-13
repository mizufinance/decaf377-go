package decaf377

import "testing"

var benchmarkField fieldElement
var benchmarkScalar Scalar
var benchmarkPoint Point

func BenchmarkNativeFieldMul(b *testing.B) {
	x, y := fieldUint(123), fieldUint(456)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkField = x.mul(y)
	}
}

func BenchmarkNativeScalarReduction(b *testing.B) {
	var input [64]byte
	for i := range input {
		input[i] = byte(i*17 + 3)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkScalar = ScalarFromUniformBytes(input[:])
	}
}

func BenchmarkNativeScalarMul(b *testing.B) {
	base, err := Generator()
	if err != nil {
		b.Fatal(err)
	}
	var scalar Scalar
	for i := range scalar {
		scalar[i] = byte(i*13 + 7)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		point, err := ScalarMul(base, scalar)
		if err != nil {
			b.Fatal(err)
		}
		benchmarkPoint = point
	}
}
