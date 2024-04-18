package handlers

import (
	"bytes"
	"fmt"
	"net/http"
)

func dashboard(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Welcome to the dashboard")
	w.Write(buf.Bytes())
}
