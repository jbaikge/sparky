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

type SessionRepository struct {
	db database.Database
}

func NewSessionRepository(db database.Database) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository) NewSession(ctx context.Context, userId int) (sessionId string, err error) {
	query := `
    INSERT INTO user_sessions (session_id, user_id, created_at, expires_at) VALUES (?, ?, ?, ?)
    `

	createdTime := time.Now()
	expiresTime := createdTime.Add(SessionDuration)
	created := createdTime.UnixMicro()
	expires := expiresTime.UnixMicro()

	sessionId = uuid.NewString()
	_, err = r.db.ExecContext(ctx, query, sessionId, userId, created, expires)
	return
}

// Prevent beating the database up by only refreshing the session if less
// than half of it remains
func (r *SessionRepository) Extend(ctx context.Context, sessionId string) error {
	selectQuery := `
    SELECT expires_at FROM user_sessions WHERE session_id = ?
    `

	var expires int64
	err := r.db.QueryRowContext(ctx, selectQuery, sessionId).Scan(&expires)
	if err != nil {
		return err
	}

	if time.UnixMicro(expires).Sub(time.Now()) > SessionDuration/2 {
		return nil
	}

	updateQuery := `
    UPDATE user_sessions SET expires_at = ? WHERE session_id = ?
    `

	expiresAt := time.Now().Add(SessionDuration)
	_, err = r.db.ExecContext(ctx, updateQuery, expiresAt, sessionId)
	return err
}

type ValidSession struct {
	Valid  bool
	UserId int
}

func (r *SessionRepository) IsValid(ctx context.Context, sessionId string) (valid ValidSession, err error) {
	query := `
    SELECT user_id, expires_at FROM user_sessions WHERE session_id = ?
    `

	var expires int64
	err = r.db.QueryRowContext(ctx, query, sessionId).Scan(&valid.UserId, &expires)
	if errors.Is(err, sql.ErrNoRows) {
		return valid, nil
	}
	if err != nil {
		return
	}

	valid.Valid = time.UnixMicro(expires).After(time.Now())
	return
}
