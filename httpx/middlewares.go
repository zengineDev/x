package httpx

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/zengineDev/x/jwkx"
	"github.com/zengineDev/x/jwtx"
)

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src https:")
		w.Header().Set("X-Content-Type-Options", "'nosniff' always;")
		next.ServeHTTP(w, r)
	})
}

func CorsHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type FailedResponse struct {
	Message string `json:"message"`
}

// Keys for authentication related keys in the context
type AuthenticationContextKey string

var (
	// Key for an jwtx.jwt in the context
	JwtContextKey AuthenticationContextKey = "auth"
)

func AuthenticationMiddleware(jwkUrl string) func(next http.Handler) http.Handler {
	keystore, err := jwkx.NewFromUrl(jwkUrl)
	if err != nil {
		log.WithField("src", "authentication middleware").Error(err)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := StripBearerPrefixFromTokenString(r.Header.Get("Authorization"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				resBody := FailedResponse{Message: "bearer token missing"}
				data, _ := json.Marshal(resBody)
				_, err = w.Write(data)
				return
			}

			token, err := jwtx.Parse(tokenString)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				resBody := FailedResponse{Message: "invalid token string"}
				data, _ := json.Marshal(resBody)
				_, err = w.Write(data)
				return
			}

			// expired
			if token.IsExpired() {
				w.WriteHeader(http.StatusUnauthorized)
				resBody := FailedResponse{Message: "the token is expired"}
				data, _ := json.Marshal(resBody)
				_, err = w.Write(data)
				return
			}

			// find the key in the store
			jwk, err := keystore.FindKey(token.Header.Kid)
			if err != nil {
				log.WithField("src", "authentication middleware").Error(err)
				w.WriteHeader(http.StatusUnauthorized)
				resBody := FailedResponse{Message: "not authorized"}
				data, _ := json.Marshal(resBody)
				_, err = w.Write(data)
				return
			}

			// build a public key from the jwk
			pubKey, err := jwk.ToPublicKey()
			if err != nil {
				log.WithField("src", "authentication middleware").Error(err)
				w.WriteHeader(http.StatusUnauthorized)
				resBody := FailedResponse{Message: "not authorized"}
				data, _ := json.Marshal(resBody)
				_, err = w.Write(data)
				return
			}

			// Verify the tokens signature
			err = token.Verify(pubKey)
			if err != nil {
				log.WithField("src", "authentication middleware").Error(err)
				w.WriteHeader(http.StatusUnauthorized)
				resBody := FailedResponse{Message: "not authorized"}
				data, _ := json.Marshal(resBody)
				_, err = w.Write(data)
				return
			}

			ctx := context.WithValue(r.Context(), JwtContextKey, token)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
