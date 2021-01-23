package errorsx

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

type HandlerE = func(w http.ResponseWriter, r *http.Request) error

type ErrorResponder interface {
	// RespondError writes an error message to w. If it doesn't know what to
	// respond, it returns false.
	RespondError(w http.ResponseWriter, r *http.Request) bool
}

func WithError(h HandlerE) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if er, ok := err.(ErrorResponder); ok {
				if er.RespondError(w, r) {
					return
				}
			}

			log.Printf("Something went wrong: %v", err)

			http.Error(w, "Internal server error", 500)
		}
	}
}
