package role

import (
	"context"
	"fmt"

	"github.com/TulgaCG/add-drop-classes-api/pkg/database"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

func addRoleToUser(ctx context.Context, db *database.DB, uid types.UserID, rid types.RoleID) (gendb.UserRole, error) {
	row, err := db.AddRoleToUser(ctx, gendb.AddRoleToUserParams{
		UserID: uid,
		RoleID: rid,
	})
	if err != nil {
		return gendb.UserRole{}, fmt.Errorf("failed to add role to the user: %w", err)
	}

	return row, nil
}

func removeRoleFromUser(ctx context.Context, db *database.DB, uid types.UserID, rid types.RoleID) error {
	row, err := db.RemoveRoleFromUser(ctx, gendb.RemoveRoleFromUserParams{
		UserID: uid,
		RoleID: rid,
	})
	if err != nil {
		return fmt.Errorf("failed to remove user from db")
	}

	if row <= 0 {
		return fmt.Errorf("role not found for given user")
	}

	return nil
}
