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

func (s *HMACSignatory) Parse(c *Claims, txt string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(txt, c, func(t *jwt.Token) (interface{}, error) {
		return s.Key, nil
	})
}
