package migrations

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jbaikge/sparky/modules/database"
	"github.com/jbaikge/sparky/modules/password"
	"github.com/jbaikge/sparky/repositories"
)

func init() {
	AddMigration(CreateAdminUser, "Create super administrator", migrateCreateAdministrator)
}

func migrateCreateAdministrator(ctx context.Context, db database.Database) error {
	repo := repositories.User(db)
	params := repositories.CreateUserParams{
		FirstName: "Super",
		LastName:  "Administrator",
		Email:     "admin@sparky.lan",
		Password:  password.TemporaryPassword(16),
		Active:    true,
	}
	if _, err := repo.CreateUser(ctx, params); err != nil {
		return fmt.Errorf("failed to create administrator: %w", err)
	}
	slog.Warn("created Super Administrator", "email", params.Email, "password", params.Password)
	return nil
}
