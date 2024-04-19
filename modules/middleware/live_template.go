package middleware

import (
	"context"
	"errors"
	"fmt"
	"html/template"
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
	patterns := []string{
		filepath.Join(m.path, "*.html"),
		filepath.Join(m.path, "*", "*.html"),
	}
	errs := make([]error, 0, len(patterns))

	var err error
	tpl := template.New("")
	for _, pattern := range patterns {
		if tpl, err = tpl.ParseGlob(pattern); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		allErrs := errors.Join(errs...)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Errors while processing templates:\n%v\n", allErrs)
		return
	}

	ctx := context.WithValue(r.Context(), ContextTemplate, tpl)
	m.handler.ServeHTTP(w, r.WithContext(ctx))
}
