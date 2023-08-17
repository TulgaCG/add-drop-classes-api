-- name: GetUserRole :many
SELECT r.role
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
WHERE u.id = ?;

-- name: CreateUserRole :exec
INSERT INTO user_roles (user_id, role_id)
SELECT ?, roles.id
FROM roles
WHERE roles.role = ?;

-- name: DeleteUserRole :exec
DELETE FROM user_roles
WHERE user_id = ?
AND role_id = (
    SELECT id
    FROM roles
    WHERE role = ?
);