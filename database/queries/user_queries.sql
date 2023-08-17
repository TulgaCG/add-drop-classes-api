-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username;

-- name: CreateUser :one
INSERT INTO users (
    username, password
) VALUES (
    ?, ?
)
RETURNING *;

-- name: DeleteUser :execrows
DELETE FROM users
WHERE id = ?;

-- name: UpdateUser :one
UPDATE users
set username = ?,
password = ?
WHERE id = ?
RETURNING *;