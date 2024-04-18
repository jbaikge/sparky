package handlers

import "net/http"

type Mux interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler http.HandlerFunc)
}

func Apply(mux Mux) {
	mux.HandleFunc("GET /admin/login", adminLogin)
	mux.HandleFunc("POST /admin/login", adminLoginAuth)
	mux.HandleFunc("GET /admin/dashboard", dashboard)
}
