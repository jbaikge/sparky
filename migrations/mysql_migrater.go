package migrations

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jbaikge/sparky/modules/database"
)

var _ migrater = new(mysqlMigrater)

type mysqlMigrater struct{}

const mysqlCreateTable = `
CREATE TABLE schema_versions (
    version INT UNSIGNED NOT NULL PRIMARY KEY,
    success BOOLEAN NOT NULL DEFAULT FALSE,
    note VARCHAR(255) NOT NULL DEFAULT '',
    created_at BIGINT(20) NOT NULL DEFAULT 0
)
`

func (m *mysqlMigrater) CreateTable(ctx context.Context, db database.Database) (err error) {
	check := `
        SELECT 1
        FROM information_schema.TABLES
        WHERE
            TABLE_SCHEMA = DATABASE()
            AND TABLE_NAME = 'schema_versions'
    `
	var found int
	row := db.QueryRowContext(ctx, check)
	if err = row.Scan(&found); err == nil {
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return
	}

	if _, err = db.ExecContext(ctx, mysqlCreateTable); err != nil {
		return
	}

	init := `
        INSERT INTO schema_versions
            (version, success, note, created_at)
        VALUES
            (?, 1, 'Create schema_versions table', UNIX_TIMESTAMP())
    `
	_, err = db.ExecContext(ctx, init, CreateSchemaVersions)
	return
}

// CurrentVersion implements migrator.
func (m *mysqlMigrater) CurrentVersion(ctx context.Context, db database.Database) (version int, err error) {
	if err = m.CreateTable(ctx, db); err != nil {
		return
	}

	query := `SELECT MAX(version) FROM schema_versions WHERE success = TRUE`
	row := db.QueryRowContext(ctx, query)
	err = row.Scan(&version)
	return
}

// Start implements migrator.
func (m *mysqlMigrater) Start(ctx context.Context, db database.Database, migration Migration) (err error) {
	query := `INSERT IGNORE INTO schema_versions (version, note) VALUES (?, ?)`
	_, err = db.ExecContext(ctx, query, migration.Version, migration.Note)
	return
}

// Success implements migrator.
func (m *mysqlMigrater) Success(ctx context.Context, db database.Database, migration Migration) (err error) {
	query := `UPDATE schema_versions SET success = 1 WHERE version = ?`
	_, err = db.ExecContext(ctx, query, migration.Version)
	return
}
