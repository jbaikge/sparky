package migrations

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jbaikge/sparky/modules/database"
)

var _ migrater = new(sqliteMigrater)

type sqliteMigrater struct{}

const sqliteCreateTable = `
CREATE TABLE schema_versions (
    version INTEGER NOT NULL PRIMARY KEY,
    success INTEGER NOT NULL DEFAULT FALSE,
    note TEXT,
    applied INTEGER
)
`

func (m *sqliteMigrater) CreateTable(ctx context.Context, db database.Database) (err error) {
	check := `
        SELECT 1
        FROM sqlite_master
        WHERE
            type = 'table'
            AND name = 'schema_versions'
    `
	var found int
	row := db.QueryRowContext(ctx, check)
	if err = row.Scan(&found); err == nil {
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return
	}

	if _, err = db.ExecContext(ctx, sqliteCreateTable); err != nil {
		return
	}

	init := `
        INSERT INTO schema_versions
            (version, success, note, applied)
        VALUES
            (?, 1, 'Create schema_versions table', unixepoch())
    `
	_, err = db.ExecContext(ctx, init, CreateSchemaVersions)
	return
}

// CurrentVersion implements migrator.
func (m *sqliteMigrater) CurrentVersion(ctx context.Context, db database.Database) (version int, err error) {
	if err = m.CreateTable(ctx, db); err != nil {
		return
	}

	query := `SELECT MAX(version) FROM schema_versions WHERE success = 1`
	row := db.QueryRowContext(ctx, query)
	err = row.Scan(&version)
	return
}

// Start implements migrator.
func (m *sqliteMigrater) Start(ctx context.Context, db database.Database, migration Migration) (err error) {
	query := `INSERT OR IGNORE INTO schema_versions (version, note, applied) VALUES (?, ?, ?)`
	_, err = db.ExecContext(ctx, query, migration.Version, migration.Note, time.Now())
	return
}

// Success implements migrator.
func (m *sqliteMigrater) Success(ctx context.Context, db database.Database, migration Migration) (err error) {
	query := `UPDATE schema_versions SET success = TRUE WHERE version = ?`
	_, err = db.ExecContext(ctx, query, migration.Version)
	return
}
