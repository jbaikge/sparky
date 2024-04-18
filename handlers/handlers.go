package handlers

import "net/http"

func ApplyHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /admin/login", adminLogin)
}
