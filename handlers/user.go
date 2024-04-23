package handlers

import (
	"log/slog"
	"net/http"

	"github.com/jbaikge/sparky/models/user"
	"github.com/jbaikge/sparky/modules/page"
)

func userList(w http.ResponseWriter, r *http.Request) {
	p := page.New(r.Context())
	p.Data["PageActiveUsers"] = true

	userRepo := user.NewUserRepository(p.Database())
	params := user.UserListParams{
		Page:    1,
		PerPage: 25,
	}
	users, err := userRepo.UserList(r.Context(), params)
	if err != nil {
		slog.Error("UserList error", "error", err)
	}
	p.Data["Users"] = users
	p.Render(w, "user/list")
}
