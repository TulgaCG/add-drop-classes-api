// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: auth_queries.sql

package gendb

import (
	"context"
	"database/sql"

	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

const updateExpirationToken = `-- name: UpdateExpirationToken :one
UPDATE users
set token_expire_at = ?
WHERE id = ?
RETURNING id, username, password, token, token_expire_at
`

type UpdateExpirationTokenParams struct {
	TokenExpireAt sql.NullTime `db:"token_expire_at" json:"tokenExpireAt"`
	ID            types.UserID `db:"id" json:"id"`
}

func (q *Queries) UpdateExpirationToken(ctx context.Context, arg UpdateExpirationTokenParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateExpirationToken, arg.TokenExpireAt, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Token,
		&i.TokenExpireAt,
	)
	return i, err
}

const updateToken = `-- name: UpdateToken :one
UPDATE users
set token = ?,
token_expire_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING token
`

type UpdateTokenParams struct {
	Token sql.NullString `db:"token" json:"token"`
	ID    types.UserID   `db:"id" json:"id"`
}

func (q *Queries) UpdateToken(ctx context.Context, arg UpdateTokenParams) (sql.NullString, error) {
	row := q.db.QueryRowContext(ctx, updateToken, arg.Token, arg.ID)
	var token sql.NullString
	err := row.Scan(&token)
	return token, err
}