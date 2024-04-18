package session

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jbaikge/sparky/modules/database"
)

const SessionDuration = 8 * time.Hour

type UserSessionRepository struct {
	db database.Database
}

func UserSession(db database.Database) *UserSessionRepository {
	return &UserSessionRepository{
		db: db,
	}
}

func (r *UserSessionRepository) NewSession(ctx context.Context, userId int) (sessionId string, err error) {
	const query = `
    INSERT INTO user_sessions (session_id, user_id, created_at, expires_at) VALUES (?, ?, ?, ?)
    `

	sessionId = uuid.NewString()
	created := time.Now()
	expires := created.Add(SessionDuration)
	_, err = r.db.ExecContext(ctx, query, sessionId, userId, created, expires)
	return
}

func (r *UserSessionRepository) Extend(ctx context.Context, sessionId string) error {
	const query = `
    UPDATE user_sessions SET expires_at = ? WHERE session_id = ?
    `

	expires := time.Now().Add(SessionDuration)
	_, err := r.db.ExecContext(ctx, query, expires, sessionId)
	return err
}

func (r *UserSessionRepository) IsValid(ctx context.Context, sessionId string) (valid bool, err error) {
	const query = `
    SELECT expires_at FROM user_sessions WHERE session_id = ?
    `

	var expires time.Time
	err = r.db.QueryRowContext(ctx, query, sessionId).Scan(&expires)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return
	}

	return expires.After(time.Now()), nil
}
