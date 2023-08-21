package role

import (
	"context"
	"fmt"

	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

// func GetUserRoles(ctx context.Context, db *gendb.Queries, id types.UserID) ([]string, error) {
// 	roles, err := db.GetUserRoles(ctx, id)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get roles of the user: %w", err)
// 	}
// 	return roles, nil
// }

func addRoleToUser(ctx context.Context, db *gendb.Queries, uid types.UserID, rid types.RoleID) (gendb.AddRoleToUserRow, error) {
	row, err := db.AddRoleToUser(ctx, gendb.AddRoleToUserParams{
		UserID: uid,
		RoleID: rid,
	})
	if err != nil {
		return gendb.AddRoleToUserRow{}, fmt.Errorf("failed to add role to the user: %w", err)
	}

	return row, nil
}

func removeRoleFromUser(ctx context.Context, db *gendb.Queries, uid types.UserID, rid types.RoleID) error {
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
