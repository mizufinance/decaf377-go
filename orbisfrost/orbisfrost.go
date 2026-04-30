package orbisfrost

import (
	"crypto/sha256"

	"github.com/mizufinance/decaf377-go"
)

const ChallengeDomain = "FROST-decaf377-challenge"

func Verify(publicKeyBytes, msg, signatureBytes []byte) (bool, error) {
	if len(publicKeyBytes) != decaf377.ElementSize {
		return false, decaf377.ErrInvalidEncoding
	}
	if len(signatureBytes) != decaf377.ElementSize+decaf377.ScalarSize {
		return false, decaf377.ErrInvalidScalar
	}

	rBytes := signatureBytes[:decaf377.ElementSize]
	zBytes := signatureBytes[decaf377.ElementSize:]

	rPoint, err := decaf377.Decode(rBytes)
	if err != nil {
		return false, err
	}
	yPoint, err := decaf377.Decode(publicKeyBytes)
	if err != nil {
		return false, err
	}
	z, err := decaf377.ScalarFromCanonicalBytes(zBytes)
	if err != nil {
		return false, err
	}

	h := sha256.New()
	h.Write([]byte(ChallengeDomain))
	h.Write(rBytes)
	h.Write(publicKeyBytes)
	h.Write(msg)
	c := decaf377.ScalarFromUniformBytes(h.Sum(nil))

	generator, err := decaf377.Generator()
	if err != nil {
		return false, err
	}
	lhs, err := decaf377.ScalarMul(generator, z)
	if err != nil {
		return false, err
	}
	cY, err := decaf377.ScalarMul(yPoint, c)
	if err != nil {
		return false, err
	}
	rhs, err := decaf377.Add(rPoint, cY)
	if err != nil {
		return false, err
	}
	return decaf377.Equivalent(lhs, rhs), nil
}
