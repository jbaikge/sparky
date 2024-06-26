package migrations

import (
	"context"

	"github.com/jbaikge/sparky/modules/database"
)

func init() {
	AddMigration(CreateUserSessions, "Create user_sessions table", migrateCreateUserSessionsTable)
}

const mysqlCreateUserSessionsTable = `
CREATE TABLE user_sessions (
    session_id CHAR(36) NOT NULL PRIMARY KEY,
    user_id INT(10) UNSIGNED NOT NULL,
    created_at BIGINT(20) UNSIGNED NOT NULL,
    expires_at BIGINT(20) UNSIGNED NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB
`

const sqliteCreateUserSessionsTable = `
CREATE TABLE user_sessions (
    session_id CHAR(36) NOT NULL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    expires_at INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE ON UPDATE CASCADE
)
`

func migrateCreateUserSessionsTable(ctx context.Context, db database.Database) error {
	query := choose(db, mysqlCreateUserSessionsTable, sqliteCreateUserSessionsTable)
	_, err := db.ExecContext(ctx, query)
	return err
}
