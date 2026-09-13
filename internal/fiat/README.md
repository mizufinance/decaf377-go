# Generated native arithmetic

`fq.go` and `fr.go` are unedited Go64 outputs of Fiat-Crypto revision
`e0a0a97d201ec1709d9ac11f1dce47c19468bf9e`, using the array-index printer patch
recorded by shieldd-formal. Exact invocation options are in each generated header.

Generator binary SHA-256:
`f197dbebfefb1ff983aee468201b55c0395450fb97e7d08e33c7eb2cba6f94ef`.
Printer patch SHA-256:
`4fff2e3baaea03ea181cad9d825928493be5d67f6b08058879925c0e9faf0870`.

Regenerate through `decaf_generate.py` in shieldd-formal; do not edit generated
arithmetic. The owning library wrappers preserve canonical big-endian Fq bytes
and little-endian scalar bytes at their existing API boundaries.

Generation and regression tests do not establish the native refinement or
compiled constant-time certificate. Those remain separate proof obligations.
