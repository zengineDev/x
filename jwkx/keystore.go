package jwkx

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type KeyStore struct {
	Keys []JWK `json:"keys"`
}

func NewFromUrl(url string) (*KeyStore, error) {
	var store *KeyStore
	client := resty.New()

	resp, err := client.R().
		EnableTrace().
		Get(url)

	if err != nil {
		return store, err
	}

	err = json.Unmarshal(resp.Body(), &store)

	if err != nil {
		return store, err
	}

	return store, nil

}

func (s *KeyStore) FindKey(keyID string) (JWK, error) {
	for _, k := range s.Keys {
		if k.Kid == keyID {
			return k, nil
		}
	}
	return JWK{}, errors.New("key not found")
}
