package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

var _ Middleware = new(Logger)

type Logger struct {
	logger  *slog.Logger
	handler http.Handler
}

func NewLogger(logger *slog.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

func (m *Logger) SetHandler(handler http.Handler) {
	m.handler = handler
}

func (m *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	m.handler.ServeHTTP(w, r)
	m.logger.Info("request", "ip", r.RemoteAddr, "method", r.Method, "path", r.URL.Path, "time", time.Since(start))
}
