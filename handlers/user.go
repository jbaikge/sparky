package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

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

func userForm(w http.ResponseWriter, r *http.Request) {
	tpl := "user/form"

	p := page.New(r.Context())
	p.Data["PageActiveUsers"] = true
	p.Data["FormAction"] = r.URL.Path

	urlId := r.PathValue("id")
	isEditing := urlId != ""
	p.Data["IsEditing"] = isEditing

	userRepo := user.NewUserRepository(p.Database())

	u := &user.User{Active: true}

	if isEditing {
		id, err := strconv.Atoi(urlId)
		if err != nil {
			slog.Error("invalid user id", "value", urlId)
			http.Redirect(w, r, "/admin/users/list", http.StatusSeeOther)
			return
		}

		u, err = userRepo.GetUserById(r.Context(), id)
		if err != nil {
			slog.Error("user does not exist", "id", id)
			http.Redirect(w, r, "/admin/users/list", http.StatusSeeOther)
			return
		}
	}
	p.Data["User"] = u

	if r.Method == http.MethodGet {
		p.Render(w, tpl)
		return
	}

	// Process POST data
	u.FirstName = r.PostFormValue("firstName")
	u.LastName = r.PostFormValue("lastName")
	u.Email = r.PostFormValue("email")
	u.Active = r.PostFormValue("active") == "1"

	oldPassword := u.Password
	if newPassword := r.PostFormValue("password"); newPassword != "" {
		u.Password = newPassword
	}
	for key, err := range u.Validate() {
		p.AddError(key, err)
	}
	if p.HasErrors() {
		p.Render(w, tpl)
		return
	}

	if err := userRepo.UpsertUser(r.Context(), u); err != nil {
		slog.Error("failed to upsert user", "id", u.UserId, "error", err)
		p.AddError("Database", err.Error())
		p.Render(w, tpl)
		return
	}

	if oldPassword != u.Password {
		err := userRepo.SetPassword(r.Context(), u.UserId, u.Password)
		if err != nil {
			slog.Error("failed to update password", "error", err)
			p.AddError("Database", err.Error())
			p.Render(w, tpl)
			return
		}
	}

	http.Redirect(w, r, "/admin/users/list", http.StatusSeeOther)
}
