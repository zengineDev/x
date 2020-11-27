package httpx

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/zengineDev/x/configx"
)

type Router struct {
	Mux *mux.Router
}

func (r Router) Get(path string, f func(response Response, request Request)) {
	r.Mux.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		log.WithField("path", path).Info("http")
		res := Response{
			writer: writer,
		}
		req := Request{}
		f(res, req)
	}).Methods(http.MethodGet)
}

// How dose i use this package after all

// r.get("", hanlder)

func (r Router) ListenAndServe() {
	c := configx.GetConfig()
	srv := &http.Server{
		Handler: r.Mux,
		Addr:    fmt.Sprintf(":%v", c.App.Port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {

	}
}
