package jwt

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

type Signatory struct {
	Key    string
	Method jwt.SigningMethod
}

func (s *Signatory) Sign(c *Claims) (string, error) {
	raw := jwt.NewWithClaims(s.Method, c)
	return raw.SignedString(s.Key)
}
