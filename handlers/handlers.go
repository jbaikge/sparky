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
	mux.HandleFunc("GET /admin/{path...}", base)
	mux.HandleFunc("GET /admin/htmx/dashboard", dashboard)
	mux.HandleFunc("GET /admin/htmx/users/list", userList)
	mux.HandleFunc("GET /admin/htmx/users/list/{page}", userList)
	mux.HandleFunc("GET /admin/htmx/users/add", userForm)
	mux.HandleFunc("POST /admin/htmx/users/add", userForm)
	mux.HandleFunc("GET /admin/htmx/users/edit/{id}", userForm)
	mux.HandleFunc("POST /admin/htmx/users/edit/{id}", userForm)
}

func Assets(mux Mux, path string) {
	mux.Handle("GET /admin/assets/", http.StripPrefix("/admin/assets/", http.FileServer(http.Dir(path))))
}

func AssetsFS(mux Mux, fs embed.FS) {
	mux.Handle("GET /admin/assets/", http.StripPrefix("/admin/", http.FileServerFS(fs)))
}
