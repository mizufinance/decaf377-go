package decaf377

type OrbisFrostVector struct {
	Name          string `json:"name"`
	PublicKeyHex  string `json:"public_key_hex"`
	MessageHex    string `json:"message_hex"`
	SignatureHex  string `json:"signature_hex"`
	ExpectValid   bool   `json:"expect_valid"`
	DerivationHex string `json:"derivation_hex,omitempty"`
	MetadataHex   string `json:"metadata_hex,omitempty"`
	Description   string `json:"description,omitempty"`
}

type OrbisFrostVectorFile struct {
	Schema  string             `json:"schema"`
	Vectors []OrbisFrostVector `json:"vectors"`
}
