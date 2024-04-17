package migrations

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/jbaikge/sparky/modules/database"
)

const (
	CreateSchemaVersions = iota
	CreatePermissions
	CreateUsers
	CreateAdminUser
	CreateUserSessions
	CreateRoles
	CreateRolePermissions
	CreateUserRoles
	NumMigrations
)

var state State = State{
	migrations: make([]Migration, 0, NumMigrations),
}

type migrater interface {
	CurrentVersion(ctx context.Context, db database.Database) (version int, err error)
	Start(ctx context.Context, db database.Database, migration Migration) (err error)
	Success(ctx context.Context, db database.Database, migration Migration) (err error)
}

type MigrationFunc func(ctx context.Context, db database.Database) (err error)

type Migration struct {
	Version int
	Note    string
	Handler MigrationFunc
}

type State struct {
	migrations []Migration
}

func AddMigration(version int, note string, handler MigrationFunc) {
	for _, migration := range state.migrations {
		if migration.Version == version {
			panic(fmt.Sprintf("migration version already exists: %d", version))
		}
	}

	state.migrations = append(state.migrations, Migration{
		Version: version,
		Note:    note,
		Handler: handler,
	})
}

func Migrate(db database.Database) (err error) {
	maxVersion := NumMigrations - 1
	ctx := context.Background()

	var m migrater
	switch db.Engine() {
	case database.EngineMySQL:
		m = new(mysqlMigrater)
	case database.EngineSQLite:
		m = new(sqliteMigrater)
	default:
		return fmt.Errorf("unknown database engine: %s", db.Engine())
	}

	current, err := m.CurrentVersion(ctx, db)
	if err != nil {
		return
	}

	// something has gone terribly wrong!
	if maxVersion < current {
		return fmt.Errorf(
			"max version less than current version: %d < %d",
			maxVersion,
			current,
		)
	}

	// up to date, nothing to do
	if maxVersion == current {
		slog.Info("database schema up to date", "version", current)
		return
	}

	slices.SortFunc(state.migrations, func(a, b Migration) int {
		return cmp.Compare(a.Version, b.Version)
	})

	for _, migration := range state.migrations {
		if migration.Version <= current {
			continue
		}
		if err = m.Start(ctx, db, migration); err != nil {
			return
		}
		slog.Info("applying migration", "version", migration.Version)
		if err = migration.Handler(ctx, db); err != nil {
			return
		}
		if err = m.Success(ctx, db, migration); err != nil {
			return
		}
	}
	return
}

func choose(db database.Database, mysqlQuery, sqliteQuery string) string {
	switch db.Engine() {
	case database.EngineMySQL:
		return mysqlQuery
	case database.EngineSQLite:
		return sqliteQuery
	}
	panic(fmt.Sprintf("unknown database engine: %s", db.Engine()))
}
