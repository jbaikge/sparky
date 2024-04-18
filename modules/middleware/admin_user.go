package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jbaikge/sparky/models/session"
	"github.com/jbaikge/sparky/models/user"
	"github.com/jbaikge/sparky/modules/database"
)

const AdminSessionCookieName = "sparky_admin"
const ContextAdminUser = ContextKey("admin-user")

type AdminUser struct {
	handler     http.Handler
	userRepo    *user.UserRepository
	sessionRepo *session.SessionRepository
}

func NewAdminHandler(db database.Database) *AdminUser {
	return &AdminUser{
		userRepo:    user.NewUserRepository(db),
		sessionRepo: session.NewSessionRepository(db),
	}
}

func (m *AdminUser) SetHandler(handler http.Handler) {
	m.handler = handler
}

func (m *AdminUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Skip user extraction for non-admin pages and the login pages
	path := r.URL.Path
	if !strings.HasPrefix(path, "/admin") || path == "/admin/login" || path == "/admin/login-form" {
		m.handler.ServeHTTP(w, r)
		return
	}

	// Check for cookie
	var user *user.User
	cookie, err := r.Cookie(AdminSessionCookieName)
	if err != nil {
		http.Redirect(w, r, "/admin/login#no-cookie", http.StatusTemporaryRedirect)
		return
	}

	// Check for valid session
	sessionId := cookie.Value
	valid, err := m.sessionRepo.IsValid(r.Context(), sessionId)
	if err != nil {
		http.Redirect(w, r, "/admin/login#"+err.Error(), http.StatusTemporaryRedirect)
		return
	}
	if !valid.Valid {
		http.Redirect(w, r, "/admin/login#session-expired", http.StatusTemporaryRedirect)
		return
	}

	// Find the user associated with the session
	user, err = m.userRepo.GetUserById(r.Context(), valid.UserId)
	if err != nil {
		http.Redirect(w, r, "/admin/login#"+err.Error(), http.StatusTemporaryRedirect)
		return
	}
	// Ensure the user is still active
	if !user.Active {
		http.Redirect(w, r, "/admin/login#inactive-user", http.StatusTemporaryRedirect)
		return
	}

	if err = m.sessionRepo.Extend(r.Context(), sessionId); err != nil {
		slog.Error("failed to extend session", "error", err)
	}

	ctx := context.WithValue(r.Context(), ContextAdminUser, user)
	m.handler.ServeHTTP(w, r.WithContext(ctx))
}
