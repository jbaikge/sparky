package middleware

import (
	"context"
	"embed"
	"html/template"
	"net/http"
	"path/filepath"
)

type EmbeddedTemplate struct {
	handler  http.Handler
	template *template.Template
}

func NewEmbeddedTemplate(fs embed.FS, rootFolder string) *EmbeddedTemplate {
	tpl := template.New("")
	tpl = template.Must(tpl.ParseFS(
		fs,
		filepath.Join(rootFolder, "*.html"),
		filepath.Join(rootFolder, "*", "*.html"),
	))

	return &EmbeddedTemplate{
		template: tpl,
	}
}

func (m *EmbeddedTemplate) SetHandler(handler http.Handler) {
	m.handler = handler
}

func (m *EmbeddedTemplate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), ContextTemplate, m.template)
	m.handler.ServeHTTP(w, r.WithContext(ctx))
}
