package jwt

import (
	"testing"

	jwt "github.com/golang-jwt/jwt/v5"
)

func TestSignAndParse(t *testing.T) {
	tests := []struct {
		signatory interface {
			Sign(*Claims) (string, error)
			Parse(*Claims, string) (*jwt.Token, error)
		}
		want string
	}{
		{
			&HMACSignatory{[]byte("the-key"), jwt.SigningMethodHS512},
			"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQWxwaGEgQnJhdm8iLCJlbWFpbCI6ImFscGhhQGJyYXZvLmNvbSIsImlzcyI6Ik03Qy5jbyJ9.lK9QUtOzR-I4LUlX83xfc_c5zfsSq1P-XHGybztBvBP2TkAfdyySP5oWaHKRtVQHB_FcN3vX4pwxV9ebfM5fsg",
		},
	}

	for _, tc := range tests {
		t.Log(tc)

		c := &Claims{
			"Alpha Bravo",
			"alpha@bravo.com",
			jwt.RegisteredClaims{
				Issuer: "M7C.co",
			},
		}

		g, err := tc.signatory.Sign(c)
		if err != nil || g != tc.want {
			t.Fatalf("wrong result; got:%q want:%q err:%+v", g, tc.want, err)
		}

		c2 := Claims{}
		tok, err := tc.signatory.Parse(&c2, g)
		if err != nil {
			t.Fatalf("failed to parse token; err:%+v", err)
		}
		t.Log(tok)
	}
}
