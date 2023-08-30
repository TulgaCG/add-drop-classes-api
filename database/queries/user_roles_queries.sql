-- name: GetUserRoles :many
SELECT r.role
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
WHERE u.id = $1;

-- name: AddRoleToUser :one
INSERT INTO user_roles (user_id, role_id) VALUES (
    $1, $2
) RETURNING *;

-- name: RemoveRoleFromUser :execrows
DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2;
