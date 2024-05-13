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
			"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQWxwaGEgQnJhdm8iLCJpc3MiOiJNN0MuY28ifQ.WcJOo0qGbHZB5yVj-O7K_mZ8lQhJ_FSJNYe2BjpNDx5kMsxv_mHh5bE8F3EjxjKDl8Kv8rzvUC6JfbntpC-QDQ",
		},
	}

	for _, tc := range tests {
		t.Log(tc)

		c := &Claims{
			"Alpha Bravo",
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
