package permission

import (
	"context"
	"fmt"

	"github.com/jbaikge/sparky/modules/database"
)

type PermissionRepository struct {
	db database.Database
}

func NewPermissionRepository(db database.Database) *PermissionRepository {
	return &PermissionRepository{
		db: db,
	}
}

type CreatePermissionParams struct {
	PermissionId string
	Description  string
}

func (r *PermissionRepository) CreatePermission(ctx context.Context, arg CreatePermissionParams) error {
	query := `
    INSERT INTO permissions (permission_id, description) VALUES (?, ?)
    `
	_, err := r.db.ExecContext(ctx, query, arg.PermissionId, arg.Description)
	return fmt.Errorf("failed to create permission: %w", err)
}

func (r *PermissionRepository) DeletePermission(ctx context.Context, permissionId string) error {
	query := `
    DELETE FROM permissions WHERE permission_id = ?
    `
	_, err := r.db.ExecContext(ctx, query, permissionId)
	return fmt.Errorf("failed to delete permission (%s): %w", permissionId, err)
}

func (r *PermissionRepository) GetPermissions(ctx context.Context) ([]Permission, error) {
	query := `
    SELECT permission_id, description FROM permissions ORDER BY permission_id
    `
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	defer rows.Close()
	var items []Permission
	for rows.Next() {
		var i Permission
		if err := rows.Scan(&i.PermissionId, &i.Description); err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rowset: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row traversal: %w", err)
	}
	return items, nil
}
