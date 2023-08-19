package user

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

func createUser(ctx context.Context, db *gendb.Queries, username, password string) (gendb.CreateUserRow, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return gendb.CreateUserRow{}, fmt.Errorf("failed to generate hashed password: %w", err)
	}

	row, err := db.CreateUser(ctx, gendb.CreateUserParams{
		Username: username,
		Password: string(hashedPassword),
	})
	if err != nil {
		return gendb.CreateUserRow{}, fmt.Errorf("failed to create user: %w", err)
	}

	return row, nil
}

func listUsers(ctx context.Context, db *gendb.Queries) ([]gendb.ListUsersRow, error) {
	rows, err := db.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return rows, nil
}

func getUser(ctx context.Context, db *gendb.Queries, id types.UserID) (gendb.GetUserRow, error) {
	row, err := db.GetUser(ctx, id)
	if err != nil {
		return gendb.GetUserRow{}, fmt.Errorf("failed to get user: %w", err)
	}

	return row, nil
}

func getUserByUsername(ctx context.Context, db *gendb.Queries, username string) (gendb.GetUserByUsernameRow, error) {
	row, err := db.GetUserByUsername(ctx, username)
	if err != nil {
		return gendb.GetUserByUsernameRow{}, fmt.Errorf("failed to get user by username: %w", err)
	}

	return row, nil
}

func updateUser(ctx context.Context, db *gendb.Queries, params gendb.UpdateUserParams) (gendb.UpdateUserRow, error) {
	u, err := db.UpdateUser(ctx, params)
	if err != nil {
		return gendb.UpdateUserRow{}, fmt.Errorf("failed to update the user: %w", err)
	}

	return u, nil
}
