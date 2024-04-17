package middlewares

import "net/http"

type Middleware interface {
	SetHandler(handler http.Handler)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
