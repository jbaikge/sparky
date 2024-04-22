package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jbaikge/sparky/modules/database"
	"github.com/jbaikge/sparky/modules/password"
)

type UserRepository struct {
	db database.Database
}

func NewUserRepository(db database.Database) *UserRepository {
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

func (r *UserRepository) CreateUser(ctx context.Context, arg CreateUserParams) (user *User, err error) {
	query := `
    INSERT INTO users (
        first_name,
        last_name,
        email,
        password,
        start_date,
        active,
        created_at,
        updated_at
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `
	hashedPassword, err := password.Hash(arg.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now()
	result, err := r.db.ExecContext(
		ctx,
		query,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		hashedPassword,
		now.Format("2006-01-02"),
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

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (user *User, err error) {
	query := `SELECT user_id FROM users WHERE email = ?`

	var id int
	if err = r.db.QueryRowContext(ctx, query, email).Scan(&id); err != nil {
		return nil, fmt.Errorf("failed to get id by email: %w", err)
	}

	return r.GetUserById(ctx, id)
}

func (r *UserRepository) GetUserById(ctx context.Context, id int) (user *User, err error) {
	query := `
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

	user = new(User)
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

type UserListParams struct {
	Page    int
	PerPage int
}

func (p UserListParams) Limit() int {
	return p.PerPage
}

func (p UserListParams) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func (r *UserRepository) UserList(ctx context.Context, params UserListParams) (items []*User, err error) {
	query := `
    SELECT
        user_id,
        first_name,
        last_name,
        email,
        start_date,
        end_date,
        active,
        created_at,
        updated_at
    FROM
        users
    ORDER BY
        last_name ASC,
        first_name ASC
    LIMIT ? OFFSET ?
    `

	var endDate sql.NullTime
	rows, err := r.db.QueryContext(ctx, query, params.Limit(), params.Offset())
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	items = make([]*User, 0, params.Limit())
	for rows.Next() {
		i := new(User)
		err = rows.Scan(
			&i.UserId,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.StartDate,
			&endDate,
			&i.Active,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		if endDate.Valid {
			i.EndDate = endDate.Time
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rowset: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row traversal: %w", err)
	}

	return
}
