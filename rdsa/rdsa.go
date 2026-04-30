package rdsa

import "github.com/mizufinance/decaf377-go"

const spendAuthPersonal = "decaf377-rdsa---"

func VerifySpendAuth(publicKeyBytes, msg, signatureBytes []byte) (bool, error) {
	if len(publicKeyBytes) != decaf377.ElementSize {
		return false, decaf377.ErrInvalidEncoding
	}
	if len(signatureBytes) != decaf377.ElementSize+decaf377.ScalarSize {
		return false, decaf377.ErrInvalidScalar
	}

	rBytes := signatureBytes[:decaf377.ElementSize]
	sBytes := signatureBytes[decaf377.ElementSize:]

	publicKey, err := decaf377.Decode(publicKeyBytes)
	if err != nil {
		return false, err
	}
	rPoint, err := decaf377.Decode(rBytes)
	if err != nil {
		return false, err
	}
	s, err := decaf377.ScalarFromCanonicalBytes(sBytes)
	if err != nil {
		return false, err
	}

	h := newPersonalizedBlake2b64(spendAuthPersonal)
	h.Write(rBytes)
	h.Write(publicKeyBytes)
	h.Write(msg)
	c := decaf377.ScalarFromUniformBytes(h.Sum())

	generator, err := decaf377.Generator()
	if err != nil {
		return false, err
	}
	lhs, err := decaf377.ScalarMul(generator, s)
	if err != nil {
		return false, err
	}
	cA, err := decaf377.ScalarMul(publicKey, c)
	if err != nil {
		return false, err
	}
	rhs, err := decaf377.Add(rPoint, cA)
	if err != nil {
		return false, err
	}
	return decaf377.Equivalent(lhs, rhs), nil
}
