package jwk

import (
	"crypto/rsa"
	"fmt"
	"io"
	"log/slog"
	"sync"
)

type KeySet struct {
	rwm *sync.RWMutex
	// keys are mapped from issuer -> kid -> key
	keys map[string]map[string]rsa.PublicKey
}

func NewKeySet(issuers ...Issuer) *KeySet {
	keys := map[string]map[string]rsa.PublicKey{}
	for _, i := range issuers {
		keys[i.String()] = map[string]rsa.PublicKey{}
	}
	return &KeySet{&sync.RWMutex{}, keys}
}

func (k *KeySet) upsert(issuer Issuer, r io.Reader) error {
	entries, err := issuer.parseJWK(r)
	if err != nil {
		slog.Error("failed to parse", "issuer", issuer)
		return err
	}

	slog.Debug("keys loaded", "count", entries)
	k.rwm.Lock()
	defer k.rwm.Unlock()

	iid := issuer.String()
	for _, entry := range entries {
		if entry.Entry.Issuer != iid {
			return fmt.Errorf("wrong Issuer; got:%s want:%s", entry.Entry.Issuer, iid)
		}

		//...don't add new issuer entries here, they must always be added by Load
		k.keys[iid][entry.Entry.Kid] = entry.Key
		slog.Debug("adding key", "issuer", iid, "kid", entry.Entry.Kid)
	}

	return nil
}

func (k *KeySet) findKey(issuer Issuer, kid string) (rsa.PublicKey, bool) {
	k.rwm.RLock()
	defer k.rwm.RUnlock()

	pub, ok := k.keys[issuer.String()][kid]
	return pub, ok
}

func (k *KeySet) Get(issuer Issuer, kid string) (rsa.PublicKey, bool) {
	if k, ok := k.findKey(issuer, kid); ok {
		return k, ok
	}

	latest, err := issuer.latest()
	if err != nil {
		slog.Error("failed to fetch latest key", "issuer", issuer.String(), "err", err)
		return rsa.PublicKey{}, false
	}

	defer latest.Close()
	err = k.upsert(issuer, latest)
	if err != nil {
		return rsa.PublicKey{}, false
	}

	pub, ok := k.findKey(issuer, kid)
	return pub, ok
}
