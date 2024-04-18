package middleware

import (
	"context"
	"net/http"

	"github.com/jbaikge/sparky/modules/database"
)

const ContextDatabase = ContextKey("database")

type Database struct {
	db      database.Database
	handler http.Handler
}

func NewDatabase(db database.Database) *Database {
	return &Database{
		db: db,
	}
}

func (m *Database) SetHandler(handler http.Handler) {
	m.handler = handler
}

func (m *Database) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), ContextDatabase, m.db)
	m.handler.ServeHTTP(w, r.WithContext(ctx))
}
