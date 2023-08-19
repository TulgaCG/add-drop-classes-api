-- name: GetUser :one
SELECT id, username FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByUsername :one
SELECT id, username FROM users
WHERE username = ? LIMIT 1;

-- name: GetUserCredentialsWithUsername :one
SELECT * FROM users
WHERE username = ? LIMIT 1;

-- name: ListUsers :many
SELECT id, username FROM users
ORDER BY username;

-- name: CreateUser :one
INSERT INTO users (
    username, password
) VALUES (
    ?, ?
)
RETURNING id, username;

-- name: UpdateUser :one
UPDATE users
set username = ?,
password = ?
WHERE id = ?
RETURNING id, username;

-- name: UpdateToken :one
UPDATE users
set token = ?,
    token_expire_at = ?
WHERE id = ?
RETURNING username, token;

-- name: UpdateTokenExpirationDate :exec
UPDATE users
set token_expire_at = ?
WHERE id = ?;
