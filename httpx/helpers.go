package httpx

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zengineDev/x/jwtx"
)

func StripBearerPrefixFromTokenString(tok string) (string, error) {
	// Should be a bearer token
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}
	return tok, nil
}

func GetJWTFromRequest(r *http.Request) jwtx.Token {
	user := r.Context().Value("auth")
	return user.(jwtx.Token)
}

func WriteJson(w http.ResponseWriter, d interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(d)
}
