package handlers

import (
	"html/template"
	"net/http"

	"github.com/jbaikge/sparky/models/user"
	"github.com/jbaikge/sparky/modules/middleware"
)

func dashboard(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
		User  *user.User
	}{
		Title: "Dashboard",
	}

	tpl := r.Context().Value(middleware.ContextTemplate).(*template.Template)
	data.User = r.Context().Value(middleware.ContextAdminUser).(*user.User)

	tpl.ExecuteTemplate(w, "admin/dashboard", &data)
}
