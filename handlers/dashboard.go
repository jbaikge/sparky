package handlers

import (
	"net/http"

	"github.com/jbaikge/sparky/modules/page"
)

func dashboard(w http.ResponseWriter, r *http.Request) {
	p := page.New(r.Context())
	p.Render(w, "dashboard")
}
