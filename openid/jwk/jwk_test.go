package jwk

import (
	"bytes"
	"os"
	"testing"

	. "github.com/m7dco/m7d"
)

var googleKey = `
{
  "keys": [
    {
      "use": "sig",
      "e": "AQAB",
      "kty": "RSA",
      "kid": "e847d9948e8545948fa8157b73e915c567302d4e",
      "alg": "RS256",
      "n": "21mw5OBuQDYONRNRyek-5Mwe2anpgn-1Ny_RGKU9eNO6_wWg-emzTpwKt4c7dDXgfyJEJ63L0zD_CS-FSyzksHKoGGySsDVX-6nD6n36MGxVCz5Z60wgM5FaSKpf7G3iOJi0IiutLcoYv5jl72g6k6nqrRTe5BSm7JfNedjpRzOeBm3IPQChW9OSW_fufV8q7Ty09ZbS0fU6KRnsMyCi80EYYg0ondJDd56iVUKR4f_OivS-EAZSUzjcu4uWYDzc9lOw8sCbb9oJE4HWLE1bgbQ05jxIqzD-6oztB1Mi-0fT5A8BV26MXnSLVPiTCgbSmQSiTq-I__uqxAfsg2v6OQ"
    },
    {
      "alg": "RS256",
      "kid": "caabf6908191616a908a1389220a975b3c0fbca1",
      "kty": "RSA",
      "e": "AQAB",
      "use": "sig",
      "n": "s0WZ_ZkzbW3vUuUiWy3u2D1RNLWDM02VeuizCFj16xX2Swd2WyS4-m0kLeOBgxU6zhenpzrU4aQypv4YFJMaB2QOvPXrLtcF4re3quSbjxjWqDKc4fJkOYMVV6X6GpaUV6FYdiYNDMiIBctPMoWVSpYhdHulKCXz366BmrAYFqUO1sUYEvkA8RKpnMyiOq85EZFXhTsoBc7OUjXrRGJQh_pUOF49_fZpCCW_y3BA9xDmxfw4AUn9ehwhmZ0J6ZicVNY10Axt7mTpilPtP--rNMYCDRXGODttMSmEJO8bCb5h7hvCM4y9cpsYBR5oq953Ik5Hm24Mub99MwsJGzyp5Q"
    }
  ]
}
`

func TestParseJWK(t *testing.T) {
	keys, err := ParseJWK(bytes.NewBufferString(googleKey))
	t.Log(keys, err)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseGolden(t *testing.T) {
	tests := []struct {
		file string
	}{
		{"golden/google_rsa_pub.jwk"},
		{"golden/ms_rsa_pub.jwk"},
	}

	for _, tc := range tests {
		f := Check(os.Open(tc.file))
		defer f.Close()

		keys, err := ParseJWK(f)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(keys)
	}
}
