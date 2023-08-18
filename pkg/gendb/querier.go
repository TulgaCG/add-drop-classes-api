// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package gendb

import (
	"context"
	"database/sql"

	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

type Querier interface {
	CreateRole(ctx context.Context, role string) (Role, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateUserRole(ctx context.Context, arg CreateUserRoleParams) error
	DeleteRoleByID(ctx context.Context, id types.RoleID) error
	DeleteRoleByName(ctx context.Context, role string) error
	DeleteUser(ctx context.Context, id types.UserID) (int64, error)
	DeleteUserRole(ctx context.Context, arg DeleteUserRoleParams) (int64, error)
	GetRole(ctx context.Context, id types.RoleID) (Role, error)
	GetRoleByName(ctx context.Context, role string) (Role, error)
	GetUser(ctx context.Context, id types.UserID) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetUserRole(ctx context.Context, id types.UserID) ([]string, error)
	ListUsers(ctx context.Context) ([]User, error)
	UpdateExpirationToken(ctx context.Context, arg UpdateExpirationTokenParams) (User, error)
	UpdateToken(ctx context.Context, arg UpdateTokenParams) (sql.NullString, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
