package middleware

import "net/http"

type ContextKey string

type Middleware interface {
	SetHandler(handler http.Handler)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
