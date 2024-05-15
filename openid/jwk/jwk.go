package jwk

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"math/big"
)

type JWKEntry struct {
	Alg    string `json:"alg"`
	Crv    string `json:"crv"`
	Kid    string `json:"kid"`
	Kty    string `json:"kty"`
	Use    string `json:"use"`
	E      string `json:"e"`
	N      string `json:"n"`
	X      string `json:"x"`
	Y      string `json:"y"`
	Issuer string `json:"issuer"`
}

var nokey = rsa.PublicKey{}

func (j JWKEntry) ToRsaPublicKey() (rsa.PublicKey, error) {
	dn, err := base64.RawURLEncoding.DecodeString(j.N)
	if err != nil {
		return nokey, err
	}

	de, err := base64.RawURLEncoding.DecodeString(j.E)
	if err != nil {
		return nokey, err
	}

	pk := rsa.PublicKey{
		N: new(big.Int).SetBytes(dn),
		E: int(new(big.Int).SetBytes(de).Int64()),
	}

	return pk, nil
}

type jwkContainer struct {
	Entries []JWKEntry `json:"keys"`
}

type ParsedKey struct {
	Entry JWKEntry
	Key   rsa.PublicKey
}

// Parses a JWK into their corresponding keys.
func ParseJWK(r io.Reader) ([]ParsedKey, error) {
	jwk, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var container jwkContainer
	err = json.Unmarshal(jwk, &container)
	if err != nil {
		return nil, err
	}

	res := []ParsedKey{}
	for _, e := range container.Entries {
		if e.Kid == "" {
			return nil, errors.New("kid must not be empty")
		}

		if e.Kty != "RSA" {
			return nil, errors.New("unsupported algo:" + e.Alg)
		}

		k, err := e.ToRsaPublicKey()
		if err != nil {
			return nil, err
		}

		res = append(res, ParsedKey{e, k})
	}
	return res, nil
}
