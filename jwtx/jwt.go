package jwtx

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"log"
	"strings"
	"time"
)

type TokenHeader struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
	Kid string `json:"kid"`
}

type HasuraClaims struct {
	UserID       uuid.UUID `json:"x-hasura-user-id,omitempty"`
	DefaultRole  string    `json:"x-hasura-default-role,omitempty"`
	AllowedRoles []string  `json:"x-hasura-allowed-roles,omitempty"`
}

type Claims struct {
	Subject      uuid.UUID    `json:"sub,omitempty"`
	Id           string       `json:"jti,omitempty"`
	Acr          string       `json:"acr,omitempty"`
	Issuer       string       `json:"iss,omitempty"`
	IssuedAt     int64        `json:"iat,omitempty"`
	ExpiresAt    int64        `json:"exp,omitempty"`
	NotBefore    int64        `json:"nbf,omitempty"`
	Scope        string       `json:"scope"`
	Audience     []string     `json:"aud,omitempty"`
	Nonce        string       `json:"nonce,omitempty"`
	HasuraClaims HasuraClaims `json:"https://hasura.io/jwt/claims"`
}

type Token struct {
	RawToken string
	Header   TokenHeader
	Claims
	Signature string
}

func NewFor() Token {
	return Token{
		Header: TokenHeader{
			Typ: "JWT",
		},
	}
}

func Parse(rawTokenString string) (Token, error) {
	var header TokenHeader
	var claims Claims
	tokenParts := strings.Split(rawTokenString, ".")

	if len(tokenParts) != 3 {
		return Token{}, errors.New("a valid jwt need 3 parts")
	}

	dataH, err := base64.RawURLEncoding.DecodeString(tokenParts[0])
	if err != nil {
		return Token{}, err
	}
	dataP, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return Token{}, err
	}

	err = json.Unmarshal(dataH, &header) // This is kind of strange
	if err != nil {
		return Token{}, err
	}
	err = json.Unmarshal(dataP, &claims)

	if err != nil {
		return Token{}, err
	}

	return Token{
		rawTokenString,
		header,
		claims,
		tokenParts[2],
	}, nil

}

func (t *Token) Verify(pubKey *rsa.PublicKey) error {
	tokenParts := strings.Split(t.RawToken, ".")

	s, err := DecodeSegment(t.Signature)
	if err != nil {
		return err
	}

	f := crypto.SHA256.New()
	_, err = f.Write([]byte(strings.Join([]string{tokenParts[0], tokenParts[1]}, ".")))
	if err != nil {
		return err
	}
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, f.Sum(nil), s)
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) Sign(signingString string, key *rsa.PrivateKey) (string, error) {
	f := crypto.SHA256.New()

	_, err := f.Write([]byte(signingString))
	if err != nil {
		return "", err
	}

	// Sign the string and return the encoded bytes
	if sigBytes, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, f.Sum(nil)); err == nil {
		return EncodeSegment(sigBytes), nil
	} else {
		return "", err
	}
}

func (t *Token) IsExpired() bool {
	tm := time.Unix(t.ExpiresAt, 0)
	return tm.Before(time.Now())
}

func (t *Token) SignedString(key *rsa.PrivateKey) (string, error) {
	var sig, sstr string
	var err error
	if sstr, err = t.SigningString(); err != nil {
		return "", err
	}
	if sig, err = t.Sign(sstr, key); err != nil {
		return "", err
	}
	return strings.Join([]string{sstr, sig}, "."), nil
}

func (t *Token) SigningString() (string, error) {
	var err error
	parts := make([]string, 2)
	for i := range parts {
		var jsonValue []byte
		if i == 0 {
			if jsonValue, err = json.Marshal(t.Header); err != nil {
				return "", err
			}
		} else {
			if jsonValue, err = json.Marshal(t.Claims); err != nil {
				return "", err
			}
		}

		parts[i] = EncodeSegment(jsonValue)
	}
	return strings.Join(parts, "."), nil
}

func EncodeSegment(seg []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(seg), "=")
}

func DecodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}

func keyIDEncode(b []byte) string {
	s := strings.TrimRight(base32.StdEncoding.EncodeToString(b), "=")
	var buf bytes.Buffer
	var i int
	for i = 0; i < len(s)/4-1; i++ {
		start := i * 4
		end := start + 4
		buf.WriteString(s[start:end] + ":")
	}
	buf.WriteString(s[i*4:])
	return buf.String()
}

func KeyIDFromCryptoKey(pubKey *rsa.PublicKey) string {
	// Generate and return a 'libtrust' fingerprint of the public key.
	// For an RSA key this should be:
	//   SHA256(DER encoded ASN1)
	// Then truncated to 240 bits and encoded into 12 base32 groups like so:
	//   ABCD:EFGH:IJKL:MNOP:QRST:UVWX:YZ23:4567:ABCD:EFGH:IJKL:MNOP
	derBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	hasher := crypto.SHA256.New()
	_, err = hasher.Write(derBytes)
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	return keyIDEncode(hasher.Sum(nil)[:30])
}
