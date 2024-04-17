package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jbaikge/sparky/models"
	"github.com/jbaikge/sparky/modules/database"
	"github.com/jbaikge/sparky/modules/password"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db database.Database
}

func User(db database.Database) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

type CreateUserParams struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Active    bool
}

func (r *UserRepository) CreateUser(ctx context.Context, arg CreateUserParams) (user *models.User, err error) {
	const query = `
    INSERT INTO users (
        first_name,
        last_name,
        email,
        password,
        active,
        created_at,
        updated_at
    ) VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	hashedPassword, err := password.Hash(arg.Password)
	if err != nil {
		return
	}

	now := time.Now()
	result, err := r.db.ExecContext(
		ctx,
		query,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		hashedPassword,
		arg.Active,
		now,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get ID of created user: %w", err)
	}

	return r.GetUserById(ctx, int(id))
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (user *models.User, err error) {
	const query = `SELECT user_id FROM users WHERE email = ?`

	var id int
	if err = r.db.QueryRowContext(ctx, query, email).Scan(&id); err != nil {
		return nil, fmt.Errorf("failed to get id by email: %w", err)
	}

	return r.GetUserById(ctx, id)
}

func (r *UserRepository) GetUserById(ctx context.Context, id int) (user *models.User, err error) {
	const query = `
    SELECT
        user_id,
        first_name,
        last_name,
        email,
        password,
        start_date,
        end_date,
        active,
        created_at,
        updated_at
    FROM users
    WHERE user_id = ?
    `

	user = new(models.User)
	var endDate sql.NullTime
	err = r.db.QueryRowContext(ctx, query, id).Scan(
		&user.UserId,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.StartDate,
		&endDate,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	if endDate.Valid {
		user.EndDate = endDate.Time
	}

	return
}
