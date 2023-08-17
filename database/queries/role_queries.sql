-- name: GetRoleName :one
SELECT role FROM roles
WHERE id = ?;

-- name: CreateRole :one
INSERT INTO roles(role)
VALUES (?)
RETURNING role;

-- name: DeleteRoleByName :exec
DELETE FROM roles WHERE role = ?;

-- name: DeleteRoleByID :exec
DELETE FROM roles WHERE id = ?;