package middleware

import (
	"context"
	"net/http"
)

const ContextHTMX = ContextKey("htmx")

var _ Middleware = new(HTMX)

type HTMX struct {
	handler http.Handler
}

func NewHTMX() *HTMX {
	return &HTMX{}
}

func (m *HTMX) SetHandler(handler http.Handler) {
	m.handler = handler
}

func (m *HTMX) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, exists := r.Header["Hx-Request"]
	ctx := context.WithValue(r.Context(), ContextHTMX, exists)
	m.handler.ServeHTTP(w, r.WithContext(ctx))
}
