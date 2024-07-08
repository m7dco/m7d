package jwt

import (
	"testing"

	jwt "github.com/golang-jwt/jwt/v5"
)

type AppClaims struct {
	Claims
	SubscriptionExpiresAt int64 `json:"subexp"`
}

func TestSignAndParse(t *testing.T) {
	tests := []struct {
		signatory interface {
			Sign(jwt.Claims) (string, error)
			Parse(jwt.Claims, string) (*jwt.Token, error)
		}
		want string
	}{
		{
			&HMACSignatory{[]byte("the-key"), jwt.SigningMethodHS512},
			"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQWxwaGEgQnJhdm8iLCJlbWFpbCI6ImFscGhhQGJyYXZvLmNvbSIsImlzcyI6Ik03Qy5jbyIsInN1YmV4cCI6NDg4NjcxODM0NX0.tHlLYzekbbrOWnR3sl5_h3lBrVjMhGmg4VuR4hgHXchrJZ55L6UpWr6-cdpT-4KsqrQKUFf6Ampp_ih0ucg_YQ",
		},
	}

	for _, tc := range tests {
		t.Log(tc)

		c := AppClaims{
			Claims{
				"Alpha Bravo",
				"alpha@bravo.com",
				jwt.RegisteredClaims{
					Issuer: "M7C.co",
				},
			},
			0x123456789,
		}

		g, err := tc.signatory.Sign(c)
		if err != nil || g != tc.want {
			t.Fatalf("wrong result; got:%q want:%q err:%+v", g, tc.want, err)
		}

		c2 := AppClaims{}
		tok, err := tc.signatory.Parse(&c2, g)
		if err != nil {
			t.Fatalf("failed to parse token; err:%+v", err)
		}
		t.Log(tok)
		t.Log(c)
		t.Log(c2)

		if c2.Name != "Alpha Bravo" || c2.Email != "alpha@bravo.com" || c2.Issuer != "M7C.co" || c2.SubscriptionExpiresAt != 0x123456789 {
			t.Fatal("Parse don't produce the expected results")
		}
	}
}
