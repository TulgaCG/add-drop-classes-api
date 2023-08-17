-- name: GetRoleName :one
SELECT role FROM roles
WHERE id = ?;

-- name: CreateRole :one
INSERT INTO roles(role)
VALUES (?);

-- name: DeleteRole :exec