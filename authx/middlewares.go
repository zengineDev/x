package authx

import (
	"context"
	"encoding/json"
	"github.com/zengineDev/x/jwtx"
	"net/http"
)

type FailedResponse struct {
	Message string `json:"message"`
}

func BearerAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := StripBearerPrefixFromTokenString(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resBody := FailedResponse{Message: "bearer token missing"}
			data, _ := json.Marshal(resBody)
			_, err = w.Write(data)
			return
		}
		token, err := jwt.New(tokenString)
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

		// verify // The key store .. who is the key store
		//verifyErr := token.Verify()

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

		} else {
			ctx := context.WithValue(r.Context(), "auth", token)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

	})
}
