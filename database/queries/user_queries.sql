-- name: GetUser :one
SELECT id, username FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT id, username FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserCredentialsWithUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: ListUsers :many
SELECT id, username FROM users
ORDER BY username;

-- name: CreateUser :one
INSERT INTO users (
    username, password
) VALUES (
    $1, $2
)
RETURNING id, username;

-- name: UpdateUser :one
UPDATE users
set username = $1,
password = $2
WHERE id = $3
RETURNING id, username;

-- name: UpdateToken :one
UPDATE users
set token = $1,
    token_expire_at = $2
WHERE id = $3
RETURNING username, token;

-- name: UpdateTokenExpirationDate :exec
UPDATE users
set token_expire_at = $1
WHERE id = $2;

-- name: TestDeleteUser :exec
DELETE FROM users
WHERE username = $1;
