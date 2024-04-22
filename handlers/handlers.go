package handlers

import (
	"embed"
	"net/http"
)

type Mux interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler http.HandlerFunc)
}

func Apply(mux Mux) {
	mux.HandleFunc("GET /admin/login", adminLogin)
	mux.HandleFunc("POST /admin/login", adminLoginAuth)
	mux.HandleFunc("GET /admin/dashboard", dashboard)
	mux.HandleFunc("GET /admin/users", userList)
}

func Assets(mux Mux, path string) {
	mux.Handle("/admin/assets/", http.StripPrefix("/admin/assets/", http.FileServer(http.Dir(path))))
}

func AssetsFS(mux Mux, fs embed.FS) {
	mux.Handle("/admin/assets/", http.StripPrefix("/admin/", http.FileServerFS(fs)))
}
