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

	mux.HandleFunc("GET /admin/roles/list", roleList)
	mux.HandleFunc("GET /admin/roles/add", roleForm)
	mux.HandleFunc("POST /admin/roles/add", roleForm)
	mux.HandleFunc("GET /admin/roles/edit/{id}", roleForm)
	mux.HandleFunc("POST /admin/roles/edit/{id}", roleForm)

	mux.HandleFunc("GET /admin/users/list", userList)
	mux.HandleFunc("GET /admin/users/add", userForm)
	mux.HandleFunc("POST /admin/users/add", userForm)
	mux.HandleFunc("GET /admin/users/edit/{id}", userForm)
	mux.HandleFunc("POST /admin/users/edit/{id}", userForm)
}

func Assets(mux Mux, path string) {
	mux.Handle("GET /admin/assets/", http.StripPrefix("/admin/assets/", http.FileServer(http.Dir(path))))
}

func AssetsFS(mux Mux, fs embed.FS) {
	mux.Handle("GET /admin/assets/", http.StripPrefix("/admin/", http.FileServerFS(fs)))
}
