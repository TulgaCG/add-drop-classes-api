-- name: GetLecture :one
SELECT * FROM lectures WHERE id = ? LIMIT 1;

-- name: GetLectureByCode :one
SELECT * FROM lectures WHERE code = ?;

-- name: ListLectures :many
SELECT * FROM lectures ORDER BY code;

-- name: CreateLecture :one
INSERT INTO lectures (name, code, credit, type) VALUES (
    ?, ?, ?, ?
) RETURNING *;