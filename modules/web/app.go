package web

import (
	"net/http"

	"github.com/jbaikge/sparky/modules/middleware"
)

type App struct {
	mux    *http.ServeMux
	server *http.Server
}

func NewApp(address string) *App {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}
	return &App{
		mux:    mux,
		server: server,
	}
}

func (a *App) AddMiddleware(m middleware.Middleware) {
	m.SetHandler(a.server.Handler)
	a.server.Handler = m
}

func (a *App) Handle(pattern string, handler http.Handler) {
	a.mux.Handle(pattern, handler)
}

func (a *App) HandleFunc(pattern string, handler http.HandlerFunc) {
	a.mux.HandleFunc(pattern, handler)
}

func (a *App) ListenAndServe() error {
	return a.server.ListenAndServe()
}
