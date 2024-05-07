package jwt

import jwt "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}
