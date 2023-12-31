// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package gendb

import (
	"context"

	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

type Querier interface {
	AddLectureToUser(ctx context.Context, arg AddLectureToUserParams) (AddLectureToUserRow, error)
	AddRoleToUser(ctx context.Context, arg AddRoleToUserParams) (UserRole, error)
	CreateLecture(ctx context.Context, arg CreateLectureParams) (Lecture, error)
	CreateRole(ctx context.Context, role types.Role) (Role, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error)
	DeleteRole(ctx context.Context, role types.Role) (int64, error)
	GetLecture(ctx context.Context, id types.LectureID) (Lecture, error)
	GetLectureByCode(ctx context.Context, code string) (Lecture, error)
	GetRoleByName(ctx context.Context, role types.Role) (Role, error)
	GetUser(ctx context.Context, id types.UserID) (GetUserRow, error)
	GetUserByUsername(ctx context.Context, username string) (GetUserByUsernameRow, error)
	GetUserCredentialsWithUsername(ctx context.Context, username string) (User, error)
	GetUserLectures(ctx context.Context, id types.UserID) ([]GetUserLecturesRow, error)
	GetUserRoles(ctx context.Context, id types.UserID) ([]types.Role, error)
	ListLectures(ctx context.Context) ([]Lecture, error)
	ListUsers(ctx context.Context) ([]ListUsersRow, error)
	RemoveLectureFromUser(ctx context.Context, arg RemoveLectureFromUserParams) (int64, error)
	RemoveRoleFromUser(ctx context.Context, arg RemoveRoleFromUserParams) (int64, error)
	TestDeleteUser(ctx context.Context, username string) error
	UpdateToken(ctx context.Context, arg UpdateTokenParams) (UpdateTokenRow, error)
	UpdateTokenExpirationDate(ctx context.Context, arg UpdateTokenExpirationDateParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) (UpdateUserRow, error)
}

var _ Querier = (*Queries)(nil)
