package user

import (
	"context"
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

// Creates a new user in the repository, side effects:
// - Sets the UserID after insertion
// - Sets CreatedAt to the current time
// - Sets UpdatedAt to the same value as CreatedAt
// - Does NOT set the password, use SetPassword for that
func (r *UserRepository) CreateUser(ctx context.Context, u *User) (err error) {
	query := `
    INSERT INTO users (
        first_name,
        last_name,
        email,
        active,
        created_at,
        updated_at
    ) VALUES (?, ?, ?, ?, ?, ?)
    `

	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	result, err := r.db.ExecContext(
		ctx,
		query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Active,
		u.CreatedAt.UnixMicro(),
		u.UpdatedAt.UnixMicro(),
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get ID of created user: %w", err)
	}

	u.UserId = int(id)
	return
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
        active,
        created_at,
        updated_at
    FROM users
    WHERE user_id = ?
    `

	user = new(User)
	var created, updated int64
	err = r.db.QueryRowContext(ctx, query, id).Scan(
		&user.UserId,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Active,
		&created,
		&updated,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	user.CreatedAt = time.UnixMicro(created)
	user.UpdatedAt = time.UnixMicro(updated)

	return
}

func (r *UserRepository) SetLastLogin(ctx context.Context, u *User) (err error) {
	query := `UPDATE users SET last_login = ? WHERE user_id = ?`
	u.LastLogin = time.Now()
	_, err = r.db.ExecContext(ctx, query, u.LastLogin.UnixMicro(), u.UserId)
	return
}

func (r *UserRepository) SetPassword(ctx context.Context, id int, pw string) (err error) {
	query := `UPDATE users SET password = ? WHERE user_id = ?`
	hashedPassword, err := password.Hash(pw)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	_, err = r.db.ExecContext(ctx, query, hashedPassword, id)
	return
}

func (r *UserRepository) UpdateUser(ctx context.Context, u *User) (err error) {
	query := `
    UPDATE users SET
        first_name = ?,
        last_name = ?,
        email = ?,
        active = ?,
        updated_at = ?
    WHERE user_id = ?
    `

	u.UpdatedAt = time.Now()
	_, err = r.db.ExecContext(
		ctx,
		query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Active,
		u.UpdatedAt.UnixMicro(),
		u.UserId,
	)

	return
}

func (r *UserRepository) UpsertUser(ctx context.Context, u *User) (err error) {
	if u.UserId == 0 {
		return r.CreateUser(ctx, u)
	}
	return r.UpdateUser(ctx, u)
}

func (r *UserRepository) UserList(ctx context.Context, params *UserListParams) (items []*User, err error) {
	query := `
    SELECT
        user_id,
        first_name,
        last_name,
        email,
        active,
        last_login,
        created_at,
        updated_at
    FROM
        users
    ORDER BY
        last_name ASC,
        first_name ASC
    LIMIT ? OFFSET ?
    `

	rows, err := r.db.QueryContext(
		ctx,
		query,
		params.Pagination.PerPage(),
		params.Pagination.Offset(),
	)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	var lastLogin, created, updated int64
	items = make([]*User, 0, params.Pagination.PerPage())
	for rows.Next() {
		i := new(User)
		err = rows.Scan(
			&i.UserId,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Active,
			&lastLogin,
			&created,
			&updated,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		i.LastLogin = time.UnixMicro(lastLogin)
		i.CreatedAt = time.UnixMicro(created)
		i.UpdatedAt = time.UnixMicro(updated)
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rowset: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row traversal: %w", err)
	}

	countQuery := `SELECT COUNT(1) FROM users`
	var total int
	if err = r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count: %w", err)
	}
	params.Pagination.SetTotal(total)

	return
}
