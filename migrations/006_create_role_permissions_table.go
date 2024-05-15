package migrations

import (
	"context"

	"github.com/jbaikge/sparky/modules/database"
)

func init() {
	AddMigration(CreateRolePermissions, "Create role_permissions table", migrateCreateRolePermissionsTable)
}

const mysqlCreateRolePermissionsTable = `
CREATE TABLE role_permissions (
    role_id INT(10) UNSIGNED NOT NULL,
    permission_id VARCHAR(64) NOT NULL,
    UNIQUE(role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles (role_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions (permission_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB CHARSET=utf8mb4
`

const sqliteCreateRolePermissionsTable = `
CREATE TABLE role_permissions (
    role_id INTEGER NOT NULL,
    permission_id VARCHAR(64) NOT NULL,
    UNIQUE(role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles (role_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions (permission_id) ON DELETE CASCADE ON UPDATE CASCADE
)
`

func migrateCreateRolePermissionsTable(ctx context.Context, db database.Database) error {
	query := choose(db, mysqlCreateRolePermissionsTable, sqliteCreateRolePermissionsTable)
	_, err := db.ExecContext(ctx, query)
	return err
}
