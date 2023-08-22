-- name: GetUserLectures :many
SELECT l.name, l.code, l.credit, l.type
FROM users u
JOIN user_lectures ul ON u.id = ul.user_id
JOIN lectures l ON ul.lecture_id = l.id
WHERE u.id = ?;

-- name: AddLectureToUser :one
INSERT INTO user_lectures ( user_id, lecture_id ) VALUES ( ?, ? ) RETURNING user_id, lecture_id;

-- name: RemoveLectureFromUser :execrows
DELETE FROM user_lectures WHERE user_id = ? AND lecture_id = ?;