package role

import (
	"context"
	"fmt"

	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
)

func createRole(ctx context.Context, db *gendb.Queries, role string) (gendb.Role, error) {
	row, err := db.CreateRole(ctx, role)
	if err != nil {
		return gendb.Role{}, fmt.Errorf("failed to create role: %w", err)
	}

	return row, nil
}

func deleteRole(ctx context.Context, db *gendb.Queries, role string) error {
	rows, err := db.DeleteRole(ctx, role)
	if rows <= 0 {
		return fmt.Errorf("role to delete not found")
	}
	return err
}
