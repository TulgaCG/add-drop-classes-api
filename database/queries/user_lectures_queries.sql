-- name: GetUserLectures :many
SELECT l.name, l.code, l.credit, l.type
FROM users u
JOIN user_lectures ul ON u.id = ul.user_id
JOIN lectures l ON ul.lecture_id = l.id
WHERE u.id = $1;

-- name: AddLectureToUser :one
INSERT INTO user_lectures ( user_id, lecture_id ) VALUES ( $1, $2 ) RETURNING user_id, lecture_id;

-- name: RemoveLectureFromUser :execrows
DELETE FROM user_lectures WHERE user_id = $1 AND lecture_id = $2;
