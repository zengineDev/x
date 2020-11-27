package jwtx

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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
	HasuraClaims HasuraClaims `json:"https://hasura.io/jwt/claims"`
}

type Token struct {
	RawToken string
	TokenHeader
	Claims
	Signature string
}

func New(rawTokenString string) (Token, error) {
	var header TokenHeader
	var claims Claims
	tokenParts := strings.Split(rawTokenString, ".")

	if len(tokenParts) != 3 {
		return Token{}, errors.New("a valid jwt need 3 parts")
	}

	dataH, err := base64.RawURLEncoding.DecodeString(tokenParts[0])
	dataP, err := base64.RawURLEncoding.DecodeString(tokenParts[1])

	fmt.Println(string(dataH))

	if err != nil {
		return Token{}, err
	}

	err = json.Unmarshal(dataH, &header) // This is kind of strange
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

	f := crypto.SHA256.New()
	f.Write([]byte(strings.Join([]string{tokenParts[0], tokenParts[1]}, ".")))

	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, f.Sum(nil), s)
	if err != nil {
		fmt.Print(err)
		return err
	}

	return nil
}

func (t *Token) IsExpired() bool {
	tm := time.Unix(t.ExpiresAt, 0)
	return tm.Before(time.Now())
}

func DecodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}
