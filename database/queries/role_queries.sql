-- name: GetRoleByName :one
SELECT * FROM roles WHERE role = $1;

-- name: CreateRole :one
INSERT INTO roles (role) VALUES ($1)
RETURNING *;

-- name: DeleteRole :execrows
DELETE FROM roles WHERE role = $1;
