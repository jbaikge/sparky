package middleware

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
)

const ContextTemplate = ContextKey("template")

type LiveTemplate struct {
	handler http.Handler
	path    string
}

func NewLiveTemplate(path string) *LiveTemplate {
	return &LiveTemplate{
		path: path,
	}
}

func (m *LiveTemplate) SetHandler(handler http.Handler) {
	m.handler = handler
}

func (m *LiveTemplate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tpl := template.New("")
	tpl = template.Must(tpl.ParseGlob(filepath.Join(m.path, "*.html")))
	tpl = template.Must(tpl.ParseGlob(filepath.Join(m.path, "*", "*.html")))
	ctx := context.WithValue(r.Context(), ContextTemplate, tpl)
	slog.Debug("live templates" + tpl.DefinedTemplates())
	m.handler.ServeHTTP(w, r.WithContext(ctx))
}
