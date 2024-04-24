package handlers

import (
	"net/http"

	"github.com/jbaikge/sparky/modules/page"
)

func base(w http.ResponseWriter, r *http.Request) {
	p := page.New(r.Context())
	p.Data["Path"] = "/admin/htmx/" + r.PathValue("path")
	p.Render(w, "base")
}
