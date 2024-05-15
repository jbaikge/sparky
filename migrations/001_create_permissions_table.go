package migrations

import (
	"context"

	"github.com/jbaikge/sparky/modules/database"
)

func init() {
	AddMigration(CreatePermissions, "Create permissions table", createPermissionsTable)
}

const mysqlCreatePermissionsTable = `
CREATE TABLE permissions (
    permission_id VARCHAR(64) NOT NULL PRIMARY KEY,
    category VARCHAR(64) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT ''
) ENGINE=InnoDB CHARSET=utf8mb4
`

const sqliteCreatePermissionsTable = `
CREATE TABLE permissions (
    permission_id VARCHAR(64) NOT NULL PRIMARY KEY,
    category VARCHAR(64) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT ''
)
`

func createPermissionsTable(ctx context.Context, db database.Database) error {
	query := choose(db, mysqlCreatePermissionsTable, sqliteCreatePermissionsTable)
	_, err := db.ExecContext(ctx, query)
	return err
}
