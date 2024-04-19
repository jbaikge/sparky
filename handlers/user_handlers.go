package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jbaikge/sparky/models/session"
	"github.com/jbaikge/sparky/models/user"
	"github.com/jbaikge/sparky/modules/database"
	"github.com/jbaikge/sparky/modules/middleware"
	"github.com/jbaikge/sparky/modules/page"
	"github.com/jbaikge/sparky/modules/password"
)

func adminLogin(w http.ResponseWriter, r *http.Request) {
	p := page.New(r.Context())
	p.Data["Title"] = "Login"
	p.Render(w, "user/login")
}

func adminLoginAuth(w http.ResponseWriter, r *http.Request) {
	var err error
	p := page.New(r.Context())

	email := r.PostFormValue("email")

	defer func(p *page.Page) {
		if err != nil {
			slog.Warn("invalid login", "email", email, "error", err)
			p.Data["Error"] = err.Error()
		}
		p.Data["Email"] = email
		p.Render(w, "user/login-form")
	}(p)

	db := r.Context().Value(middleware.ContextDatabase).(database.Database)
	userRepo := user.NewUserRepository(db)

	user, err := userRepo.GetUserByEmail(r.Context(), email)
	if errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("Invalid email address")
		return
	} else if err != nil {
		err = fmt.Errorf("Error from the database: %w", err)
		return
	}

	pw := r.PostFormValue("password")
	if !password.Validate(user.Password, pw) {
		err = fmt.Errorf("Invalid password")
		return
	}

	sessionRepo := session.NewSessionRepository(db)
	sessionId, err := sessionRepo.NewSession(r.Context(), user.UserId)
	if err != nil {
		err = fmt.Errorf("Unable to create new session: %w", err)
		return
	}

	cookie := &http.Cookie{
		Name:  middleware.AdminSessionCookieName,
		Value: sessionId,
	}
	http.SetCookie(w, cookie)
	w.Header().Set("HX-Redirect", "/admin/dashboard")
}
