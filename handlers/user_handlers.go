package handlers

import "net/http"

func adminLogin(w http.ResponseWriter, r *http.Request) {
	//tpl := r.Context().Value(middleware.ContextTemplate).(*template.Template)
	w.Write([]byte(`This is the admin login`))
}
