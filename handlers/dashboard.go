package handlers

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/jbaikge/sparky/models/user"
	"github.com/jbaikge/sparky/modules/middleware"
)

func dashboard(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.ContextAdminUser).(*user.User)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Welcome to the dashboard, %s", user.Name())
	w.Write(buf.Bytes())
}
