package role

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/jbaikge/sparky/modules/database"
)

type RoleRepository struct {
	db database.Database
}

func NewRoleRepository(db database.Database) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) CreateRole(ctx context.Context, role *Role) (err error) {
	query := `INSERT INTO roles (name, created_at, updated_at) VALUES (?, ?, ?)`

	now := time.Now()
	role.CreatedAt = now
	role.UpdatedAt = now
	nowMicro := now.UnixMicro()
	result, err := r.db.ExecContext(ctx, query, role.Name, nowMicro, nowMicro)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get ID of created role: %w", err)
	}
	role.RoleId = int(id)

	return r.syncPermissions(ctx, role)
}

func (r *RoleRepository) GetRoleById(ctx context.Context, id int) (role *Role, err error) {
	return
}

func (r *RoleRepository) ListRoles(ctx context.Context, params *RoleListParams) (roles []*Role, err error) {
	roleQuery := `
    SELECT role_id, name, created_at, updated_at
    FROM roles
    ORDER BY name ASC
    LIMIT ? OFFSET ?
    `
	roleRows, err := r.db.QueryContext(
		ctx,
		roleQuery,
		params.Pagination.PerPage(),
		params.Pagination.Offset(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %w", err)
	}

	permQuery := `SELECT permission_id FROM role_permissions WHERE role_id = ?`
	var created, updated int64
	var id string
	for roleRows.Next() {
		role := new(Role)
		err = roleRows.Scan(&role.RoleId, &role.Name, &created, &updated)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		role.CreatedAt = time.UnixMicro(created)
		role.UpdatedAt = time.UnixMicro(updated)

		permRows, err := r.db.QueryContext(ctx, permQuery, role.RoleId)
		if err != nil {
			return nil, fmt.Errorf("failed to query permissions: %w", err)
		}

		for permRows.Next() {
			if err = permRows.Scan(&id); err != nil {
				return nil, fmt.Errorf("failed to scan permission: %w", err)
			}
			role.Permissions = append(role.Permissions, id)
		}

		roles = append(roles, role)
	}

	countQuery := `SELECT COUNT(1) FROM roles`
	var total int
	if err = r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count roles: %w", err)
	}
	params.Pagination.SetTotal(total)

	return
}

func (r *RoleRepository) UpdateRole(ctx context.Context, role *Role) (err error) {
	query := `UPDATE roles SET name = ?, updated_at = ? WHERE role_id = ?`
	role.UpdatedAt = time.Now()
	_, err = r.db.ExecContext(
		ctx,
		query,
		role.Name,
		role.RoleId,
		role.UpdatedAt.UnixMicro(),
	)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	return r.syncPermissions(ctx, role)
}

func (r *RoleRepository) UpsertRole(ctx context.Context, role *Role) (err error) {
	if role.RoleId == 0 {
		return r.CreateRole(ctx, role)
	}
	return r.UpdateRole(ctx, role)
}

func (r *RoleRepository) syncPermissions(ctx context.Context, role *Role) (err error) {
	// Pull in existing permissions associated with this role
	existingQuery := `
    SELECT permission_id
    FROM role_permissions
    WHERE role_id = ?
    ORDER BY permission_id
    `
	var existing []string
	existingRows, err := r.db.QueryContext(ctx, existingQuery, role.RoleId)
	if err != nil {
		return fmt.Errorf("failed to query existing permissions: %w", err)
	}
	for existingRows.Next() {
		var id string
		if err = existingRows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan permission ID: %w", err)
		}
		existing = append(existing, id)
	}
	if err = existingRows.Err(); err != nil {
		return fmt.Errorf("failed to iterate rows: %w", err)
	}

	// Sort incoming permissions to make searches easier
	slices.Sort(role.Permissions)

	// Items in existing and not role need to be deleted
	deleteQuery := `
    DELETE FROM role_permissions
    WHERE role_id = ? AND permission_id = ?
    `
	for _, item := range existing {
		if slices.Contains(role.Permissions, item) {
			continue
		}
		if _, err = r.db.ExecContext(ctx, deleteQuery, role.RoleId, item); err != nil {
			return fmt.Errorf("failed to delete permission: %w", err)
		}
	}

	// Items in role and not existing need to be added
	insertQuery := `
    INSERT INTO role_permissions (role_id, permission_id)
    VALUES (?, ?)
    `
	for _, item := range role.Permissions {
		if slices.Contains(existing, item) {
			continue
		}
		if _, err = r.db.ExecContext(ctx, insertQuery, role.RoleId, item); err != nil {
			return fmt.Errorf("failed to insert permission: %w", err)
		}
	}
	return
}
