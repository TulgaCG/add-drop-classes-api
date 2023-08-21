-- name: CreateRole :one
INSERT INTO roles (role) VALUES (?)
RETURNING *;

-- name: DeleteRole :execrows
DELETE FROM roles WHERE role = ?;