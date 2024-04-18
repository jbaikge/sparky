package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/jbaikge/sparky/handlers"
	"github.com/jbaikge/sparky/migrations"
	"github.com/jbaikge/sparky/modules/database"
	"github.com/jbaikge/sparky/modules/middleware"
	"github.com/jbaikge/sparky/modules/web"
)

func main() {
	var development bool
	var dbEngine, dbName, dbHost, dbUser, dbPass string
	var dbPort int
	var address string

	flag.BoolVar(&development, "dev", false, "Operate in development mode")
	flag.StringVar(&dbEngine, "db.engine", "sqlite3", "Database engine: mysql, sqlite3")
	flag.StringVar(&dbName, "db.database", "sparky.db", "Database name (mysql) or file (sqlite)")
	flag.StringVar(&dbHost, "db.host", "", "Database host")
	flag.StringVar(&dbUser, "db.user", "", "Database username")
	flag.StringVar(&dbPass, "db.pass", "", "Database password")
	flag.IntVar(&dbPort, "db.port", 3306, "Database port")
	flag.StringVar(&address, "server.address", "0.0.0.0:3003", "Address and port to listen on")
	flag.Parse()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	db, err := database.Connect(dbEngine, dbName, dbHost, dbPort, dbUser, dbPass)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	if err := migrations.Migrate(db); err != nil {
		slog.Error("failed to migrate", "error", err)
		os.Exit(1)
	}

	app := web.NewApp(address)
	handlers.Apply(app)
	app.AddMiddleware(middleware.NewContentType())
	app.AddMiddleware(middleware.NewLogger(slog.Default()))
	app.ListenAndServe()
}
