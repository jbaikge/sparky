package migrations

import (
	"context"

	"github.com/jbaikge/sparky/modules/database"
)

func init() {
	AddMigration(CreateRoles, "Create roles table", migrateCreateRolesTable)
}

const mysqlCreateRolesTable = `
CREATE TABLE roles (
    role_id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY
    name VARCHAR(128) NOT NULL DEFAULT '',
    created_at BIGINT(20) UNSIGNED NOT NULL,
    expires_at BIGINT(20) UNSIGNED NOT NULL
) ENGINE=InnoDB CHARSET=utf8mb4
`

const sqliteCreateRolesTable = `
CREATE TABLE roles (
    role_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(128) NOT NULL DEFAULT '',
    created_at INTEGER NOT NULL,
    expires_at INTEGER NOT NULL
)
`

func migrateCreateRolesTable(ctx context.Context, db database.Database) error {
	query := choose(db, mysqlCreateRolesTable, sqliteCreateRolesTable)
	_, err := db.ExecContext(ctx, query)
	return err
}
