package jwkx

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/pkg/errors"
	"math/big"
)

type JWK struct {
	E   string `json:"e"`
	Kty string `json:"kty"`
	N   string `json:"n"`
	Kid string `json:"kid"`
}

func Encode(pub *rsa.PublicKey) (string, error) {
	// https://tools.ietf.org/html/rfc7518#section-6.3.1
	n := pub.N
	e := big.NewInt(int64(pub.E))
	// Field order is important.
	// See https://tools.ietf.org/html/rfc7638#section-3.3 for details.
	return fmt.Sprintf(`{"e":"%s","kty":"RSA","n":"%s"}`,
		base64.RawURLEncoding.EncodeToString(e.Bytes()),
		base64.RawURLEncoding.EncodeToString(n.Bytes()),
	), nil
}

func (j *JWK) ToPublicKey() (*rsa.PublicKey, error) {
	var publicKey *rsa.PublicKey
	// decode the base64 bytes for n
	nb, err := base64.RawURLEncoding.DecodeString(j.N)
	if err != nil {
		return publicKey, err
	}

	e := 0
	// The default exponent is usually 65537, so just compare the
	// base64 for [1,0,1] or [0,1,0,1]
	if j.E == "AQAB" || j.E == "AAEAAQ" {
		e = 65537
	} else {
		// need to decode "e" as a big-endian int
		return publicKey, errors.New("need to deocde e:")
	}

	pk := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nb),
		E: e,
	}

	der, err := x509.MarshalPKIXPublicKey(pk)
	if err != nil {
		return publicKey, err
	}

	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: der,
	}

	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return publicKey, err
	}
	var pubKey *rsa.PublicKey
	pubKey = parsedKey.(*rsa.PublicKey)
	return pubKey, nil
}

func (j *JWK) ToPEM() (string, error) {
	// decode the base64 bytes for n
	nb, err := base64.RawURLEncoding.DecodeString(j.N)
	if err != nil {
		panic(err)
	}

	e := 0
	// The default exponent is usually 65537, so just compare the
	// base64 for [1,0,1] or [0,1,0,1]
	if j.E == "AQAB" || j.E == "AAEAAQ" {
		e = 65537
	} else {
		// need to decode "e" as a big-endian int
		return "", errors.New("need to deocde e:")
	}

	pk := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nb),
		E: e,
	}

	der, err := x509.MarshalPKIXPublicKey(pk)
	if err != nil {
		return "", err
	}

	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: der,
	}

	var out bytes.Buffer
	err = pem.Encode(&out, block)

	return out.String(), nil

}
