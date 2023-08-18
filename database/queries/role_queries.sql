-- name: GetRole :one
SELECT * FROM roles
WHERE id = ?;

-- name: GetRoleByName :one
SELECT * FROM roles
WHERe role = ?;

-- name: CreateRole :one
INSERT INTO roles(role)
VALUES (?)
RETURNING *;

-- name: DeleteRoleByName :exec
DELETE FROM roles WHERE role = ?;

-- name: DeleteRoleByID :exec
DELETE FROM roles WHERE id = ?;