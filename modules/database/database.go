package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	EngineMySQL  = "mysql"
	EngineSQLite = "sqlite3"
)

type Database interface {
	Engine() string
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type DB struct {
	*sql.DB
	engine string
}

func Connect(engine string, database string, host string, port int, user string, pass string) (db *DB, err error) {
	var dsn string
	db = &DB{
		engine: engine,
	}

	switch engine {
	case EngineMySQL:
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?multiStatements=true&parseTime=true",
			user,
			pass,
			host,
			port,
			database,
		)
	case EngineSQLite:
		dsn = fmt.Sprintf("file:%s?_foreign_keys=true", database)
	default:
		return nil, fmt.Errorf("invalid engine: %s", engine)
	}

	db.DB, err = sql.Open(engine, dsn)
	return
}

func (db DB) Engine() string {
	return db.engine
}
