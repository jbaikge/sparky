package handlers

import (
	"net/http"
	"strings"
)

func redirect(w http.ResponseWriter, r *http.Request, path string) {
	w.Header().Add("HX-Push-Url", path)
	hxUrl := strings.Replace(path, "/admin/", "/admin/htmx/", 1)
	http.Redirect(w, r, hxUrl, http.StatusSeeOther)
}
