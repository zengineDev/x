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
		req := Request{
			Vars: mux.Vars(request),
		}
		f(res, req)
	}).Methods(http.MethodGet)
}

func (r *Router) Resource(path string, ctr ControllerContract) {
	// Index Route
	r.Mux.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		log.WithField("path", path).Info("http")
		res := Response{
			writer: writer,
		}
		req := Request{
			Vars: mux.Vars(request),
			Req:  request,
		}
		ctr.Index(res, req)
	}).Methods(http.MethodGet)

	// Store
	r.Mux.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		log.WithField("path", path).Info("http")
		res := Response{
			writer: writer,
		}
		req := Request{
			Vars: mux.Vars(request),
			Req:  request,
		}
		ctr.Store(res, req)
	}).Methods(http.MethodPost)

	// Get Route
	r.Mux.HandleFunc(fmt.Sprintf("/%s/{id}", path), func(writer http.ResponseWriter, request *http.Request) {
		res := Response{
			writer: writer,
		}
		req := Request{
			Vars: mux.Vars(request),
			Req:  request,
		}
		ctr.Get(res, req)
	}).Methods(http.MethodGet)

	r.Mux.HandleFunc(fmt.Sprintf("/%s/{id}", path), func(writer http.ResponseWriter, request *http.Request) {
		res := Response{
			writer: writer,
		}
		req := Request{
			Vars: mux.Vars(request),
			Req:  request,
		}
		ctr.Update(res, req)
	}).Methods(http.MethodPut)

	r.Mux.HandleFunc(fmt.Sprintf("/%s/{id}", path), func(writer http.ResponseWriter, request *http.Request) {
		res := Response{
			writer: writer,
		}
		req := Request{
			Vars: mux.Vars(request),
			Req:  request,
		}
		ctr.Delete(res, req)
	}).Methods(http.MethodDelete)

}

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
