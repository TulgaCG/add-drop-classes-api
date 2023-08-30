package user

import (
	"context"
	"fmt"

	"github.com/TulgaCG/add-drop-classes-api/pkg/database"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

func ListUsers(ctx context.Context, db *database.DB) ([]gendb.ListUsersRow, error) {
	rows, err := db.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return rows, nil
}

func GetUser(ctx context.Context, db *database.DB, id types.UserID) (gendb.GetUserRow, error) {
	row, err := db.GetUser(ctx, id)
	if err != nil {
		return gendb.GetUserRow{}, fmt.Errorf("failed to get user: %w", err)
	}

	return row, nil
}

func GetUserByUsername(ctx context.Context, db *database.DB, username string) (gendb.GetUserByUsernameRow, error) {
	row, err := db.GetUserByUsername(ctx, username)
	if err != nil {
		return gendb.GetUserByUsernameRow{}, fmt.Errorf("failed to get user by username: %w", err)
	}

	return row, nil
}

func UpdateUser(ctx context.Context, db *database.DB, params gendb.UpdateUserParams) (gendb.UpdateUserRow, error) {
	u, err := db.UpdateUser(ctx, params)
	if err != nil {
		return gendb.UpdateUserRow{}, fmt.Errorf("failed to update the user: %w", err)
	}

	return u, nil
}
