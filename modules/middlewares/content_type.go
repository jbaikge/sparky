package middlewares

import "net/http"

var _ Middleware = new(ContentType)

type ContentType struct {
	handler http.Handler
}

func NewContentType() *ContentType {
	return &ContentType{}
}

func (m *ContentType) SetHandler(handler http.Handler) {
	m.handler = handler
}

func (m *ContentType) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf=8")
	m.handler.ServeHTTP(w, r)
}
