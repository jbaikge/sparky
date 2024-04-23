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

func userAddForm(w http.ResponseWriter, r *http.Request) {
	tpl := "user/form"

	p := page.New(r.Context())
	p.Data["PageActiveUsers"] = true
	p.Data["FormAction"] = r.URL.Path

	// Really important or blank records go into the database, whoops
	if r.Method == http.MethodGet {
		p.Data["User"] = user.User{
			Active: true,
		}
		p.Render(w, tpl)
		return
	}

	u := user.User{
		FirstName: r.PostFormValue("firstName"),
		LastName:  r.PostFormValue("lastName"),
		Email:     r.PostFormValue("email"),
		Password:  r.PostFormValue("password"),
		Active:    r.PostFormValue("active") == "1",
	}
	p.Data["User"] = u

	for key, err := range u.Validate() {
		p.AddError(key, err)
	}
	if p.HasErrors() {
		p.Render(w, tpl)
		return
	}

	userRepo := user.NewUserRepository(p.Database())
	params := user.CreateUserParams{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
		Active:    u.Active,
	}
	slog.Debug("creating user", "params", params)
	if _, err := userRepo.CreateUser(r.Context(), params); err != nil {
		p.AddError("CreateUser", err.Error())
	}

	// Great success, go back to the user list
	if !p.HasErrors() {
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
		return
	}

	p.Render(w, tpl)
}
