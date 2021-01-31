package httpx

import (
	"encoding/json"
	"io/ioutil"
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

func ReadJson(r *http.Request, i interface{}) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(reqBody, &i)
	if err != nil {
		return err
	}

	return nil
}
