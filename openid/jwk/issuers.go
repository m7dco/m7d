package jwk

import (
	"io"
	"log/slog"
	"net/http"
)

var (
	ISSUER_GOOGLE Issuer = &issuerGoogle{
		issuerBase{
			"https://accounts.google.com",
			"https://www.googleapis.com/oauth2/v3/certs",
		},
	}

	ISSUER_MS Issuer = &issuerMS{
		issuerBase{
			"https://login.microsoftonline.com/9188040d-6c67-4c5b-b112-36a304b66dad/v2.0",
			"https://login.microsoftonline.com/common/discovery/v2.0/keys",
		},
	}
)

type Issuer interface {
	String() string
	parseJWK(r io.Reader) ([]ParsedKey, error)
	latest() (io.ReadCloser, error)
}

type issuerBase struct {
	id     string
	jwkURL string
}

func (i issuerBase) String() string {
	return i.id
}

func (i issuerBase) latest() (io.ReadCloser, error) {
	resp, err := http.Get(i.jwkURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		slog.Error("failed to get JWK",
			"issuer", i.id,
			"url", i.jwkURL,
			"error", err)

		if err != nil {
			slog.Error("http status", "status", resp.Status)
		}
	}

	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

type issuerGoogle struct {
	issuerBase
}

func (i *issuerGoogle) parseJWK(r io.Reader) ([]ParsedKey, error) {
	entries, err := ParseJWK(r)
	if err != nil {
		return nil, err
	}

	// Google JWK do not contain any issuer attribute and therefore
	// we patch up the parsed result to match the expectation from KeySet.
	for j, e := range entries {
		e.Entry.Issuer = i.id
		entries[j] = e
	}

	return entries, nil
}

type issuerMS struct {
	issuerBase
}

func (i *issuerMS) parseJWK(r io.Reader) ([]ParsedKey, error) {
	entries, err := ParseJWK(r)
	if err != nil {
		return nil, err
	}

	// MS JWK contains tenantid specific keys, we don't use them or care
	// about them and therefore we just drop them from the result.
	res := []ParsedKey{}
	for _, e := range entries {
		if e.Entry.Issuer != i.id {
			slog.Debug("dropping", "issuer", e.Entry.Issuer, "kid", e.Entry.Kid)
			continue
		}

		res = append(res, e)
	}

	return res, nil
}
