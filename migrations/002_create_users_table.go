package migrations

import (
	"context"

	"github.com/jbaikge/sparky/modules/database"
)

func init() {
	AddMigration(CreateUsers, "Create users table", migrateCreateUsersTable)
}

const mysqlCreateUsersTable = `
CREATE TABLE users (
    user_id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(64) NOT NULL DEFAULT '',
    last_name VARCHAR(64) NOT NULL DEFAULT '',
    email VARCHAR(128) NOT NULL DEFAULT '',
    password CHAR(60) NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login BIGINT(20) UNSIGNED NOT NULL DEFAULT 0,
    created_at BIGINT(20) UNSIGNED NOT NULL,
    updated_at BIGINT(20) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`

const sqliteCreateUsersTable = `
CREATE TABLE users (
    user_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    first_name VARCHAR(64) NOT NULL DEFAULT '',
    last_name VARCHAR(64) NOT NULL DEFAULT '',
    email VARCHAR(128) NOT NULL DEFAULT '',
    password CHAR(60) NULL,
    active INTEGER NOT NULL DEFAULT TRUE,
    last_login INTEGER NOT NULL DEFAULT 0,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
)
`

func migrateCreateUsersTable(ctx context.Context, db database.Database) error {
	query := choose(db, mysqlCreateUsersTable, sqliteCreateUsersTable)
	_, err := db.ExecContext(ctx, query)
	return err
}
