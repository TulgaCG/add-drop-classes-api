-- name: UpdateToken :one
UPDATE users
set token = ?,
token_expire_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING token;

-- name: UpdateExpirationToken :one
UPDATE users
set token_expire_at = ?
WHERE id = ?
RETURNING *;