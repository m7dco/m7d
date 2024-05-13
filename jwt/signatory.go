package jwt

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

type HMACSignatory struct {
	Key    []byte
	Method jwt.SigningMethod
}

func (s *HMACSignatory) Sign(c *Claims) (string, error) {
	raw := jwt.NewWithClaims(s.Method, c)
	return raw.SignedString(s.Key)
}
