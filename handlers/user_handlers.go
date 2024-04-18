package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/jbaikge/sparky/models/session"
	"github.com/jbaikge/sparky/models/user"
	"github.com/jbaikge/sparky/modules/database"
	"github.com/jbaikge/sparky/modules/middleware"
	"github.com/jbaikge/sparky/modules/password"
)

func adminLogin(w http.ResponseWriter, r *http.Request) {
	tpl := r.Context().Value(middleware.ContextTemplate).(*template.Template)
	tpl.ExecuteTemplate(w, "user/login", nil)
}

func adminLoginAuth(w http.ResponseWriter, r *http.Request) {
	type Repost struct {
		Error string
		Email string
	}

	repost := Repost{
		Email: r.PostFormValue("email"),
	}

	defer func(repost *Repost) {
		if repost.Error != "" {
			slog.Warn("invalid login", "email", repost.Email, "error", repost.Error)
		}
		tpl := r.Context().Value(middleware.ContextTemplate).(*template.Template)
		tpl.ExecuteTemplate(w, "user/login-form", repost)
	}(&repost)

	db := r.Context().Value(middleware.ContextDatabase).(database.Database)
	userRepo := user.NewUserRepository(db)

	user, err := userRepo.GetUserByEmail(r.Context(), repost.Email)
	if errors.Is(err, sql.ErrNoRows) {
		repost.Error = "Invalid email address"
		return
	} else if err != nil {
		repost.Error = fmt.Sprintf("Error from the database: %v", err)
		return
	}

	pw := r.PostFormValue("password")
	if !password.Validate(user.Password, pw) {
		repost.Error = "Invalid password"
		return
	}

	sessionRepo := session.NewSessionRepository(db)
	sessionId, err := sessionRepo.NewSession(r.Context(), user.UserId)
	if err != nil {
		repost.Error = fmt.Sprintf("Unable to create new session: %v", err)
		return
	}

	cookie := &http.Cookie{
		Name:  middleware.AdminSessionCookieName,
		Value: sessionId,
	}
	http.SetCookie(w, cookie)
	w.Header().Set("HX-Redirect", "/admin/dashboard")
}
