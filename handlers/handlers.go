package handlers

import (
	"net/http"

	"github.com/jbaikge/sparky/modules/middleware"
)

var mux *http.ServeMux = &http.ServeMux{}
var server *http.Server = &http.Server{
	Handler: mux,
}

func AddHandler(path string, handler http.Handler) {
	mux.Handle(path, handler)
}

func AddMiddleware(m middleware.Middleware) {
	m.SetHandler(server.Handler)
	server.Handler = m
}
