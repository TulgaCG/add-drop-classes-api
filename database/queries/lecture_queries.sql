-- name: GetLecture :one
SELECT * FROM lectures WHERE id = $1 LIMIT 1;

-- name: GetLectureByCode :one
SELECT * FROM lectures WHERE code = $1;

-- name: ListLectures :many
SELECT * FROM lectures ORDER BY code;

-- name: CreateLecture :one
INSERT INTO lectures (name, code, credit, type) VALUES (
    $1, $2, $3, $4
) RETURNING *;
