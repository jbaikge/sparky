package handlers

import (
	"net/http"
)

func init() {
	mux.HandleFunc("GET /admin/login", handleAdminLogin)
}

func handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	//tpl := r.Context().Value(middleware.ContextTemplate).(*template.Template)
}
