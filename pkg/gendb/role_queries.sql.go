// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: role_queries.sql

package gendb

import (
	"context"
)

const createRole = `-- name: CreateRole :one
INSERT INTO roles(role)
VALUES (?)
RETURNING id, role
`

func (q *Queries) CreateRole(ctx context.Context, role string) (Role, error) {
	row := q.db.QueryRowContext(ctx, createRole, role)
	var i Role
	err := row.Scan(&i.ID, &i.Role)
	return i, err
}

const deleteRoleByID = `-- name: DeleteRoleByID :exec
DELETE FROM roles WHERE id = ?
`

func (q *Queries) DeleteRoleByID(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteRoleByID, id)
	return err
}

const deleteRoleByName = `-- name: DeleteRoleByName :exec
DELETE FROM roles WHERE role = ?
`

func (q *Queries) DeleteRoleByName(ctx context.Context, role string) error {
	_, err := q.db.ExecContext(ctx, deleteRoleByName, role)
	return err
}

const getRole = `-- name: GetRole :one
SELECT id, role FROM roles
WHERE id = ?
`

func (q *Queries) GetRole(ctx context.Context, id int64) (Role, error) {
	row := q.db.QueryRowContext(ctx, getRole, id)
	var i Role
	err := row.Scan(&i.ID, &i.Role)
	return i, err
}

const getRoleByName = `-- name: GetRoleByName :one
SELECT id, role FROM roles
WHERe role = ?
`

func (q *Queries) GetRoleByName(ctx context.Context, role string) (Role, error) {
	row := q.db.QueryRowContext(ctx, getRoleByName, role)
	var i Role
	err := row.Scan(&i.ID, &i.Role)
	return i, err
}
