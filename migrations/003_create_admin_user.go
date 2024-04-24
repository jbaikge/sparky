package migrations

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jbaikge/sparky/models/user"
	"github.com/jbaikge/sparky/modules/database"
	"github.com/jbaikge/sparky/modules/password"
)

func init() {
	AddMigration(CreateAdminUser, "Create super administrator", migrateCreateAdministrator)
}

func migrateCreateAdministrator(ctx context.Context, db database.Database) error {
	repo := user.NewUserRepository(db)

	u := &user.User{
		FirstName: "Super",
		LastName:  "Administrator",
		Email:     "admin@sparky.lan",
		Active:    true,
	}
	if err := repo.CreateUser(ctx, u); err != nil {
		return fmt.Errorf("failed to create administrator: %w", err)
	}

	pw := password.TemporaryPassword(16)
	if err := repo.SetPassword(ctx, u.UserId, pw); err != nil {
		return fmt.Errorf("failed to set administrator password: %w", err)
	}

	slog.Warn("created Super Administrator", "email", u.Email, "password", pw)
	return nil
}
