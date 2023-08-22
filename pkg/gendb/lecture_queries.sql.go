// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: lecture_queries.sql

package gendb

import (
	"context"
	"database/sql"

	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

const createLecture = `-- name: CreateLecture :one
INSERT INTO lectures (name, code, credit, type) VALUES (
    ?, ?, ?, ?
) RETURNING id, name, code, credit, type
`

type CreateLectureParams struct {
	Name   string        `db:"name" json:"name"`
	Code   string        `db:"code" json:"code"`
	Credit int64         `db:"credit" json:"credit"`
	Type   sql.NullInt64 `db:"type" json:"type"`
}

func (q *Queries) CreateLecture(ctx context.Context, arg CreateLectureParams) (Lecture, error) {
	row := q.db.QueryRowContext(ctx, createLecture,
		arg.Name,
		arg.Code,
		arg.Credit,
		arg.Type,
	)
	var i Lecture
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Code,
		&i.Credit,
		&i.Type,
	)
	return i, err
}

const getLecture = `-- name: GetLecture :one
SELECT id, name, code, credit, type FROM lectures WHERE id = ? LIMIT 1
`

func (q *Queries) GetLecture(ctx context.Context, id types.LectureID) (Lecture, error) {
	row := q.db.QueryRowContext(ctx, getLecture, id)
	var i Lecture
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Code,
		&i.Credit,
		&i.Type,
	)
	return i, err
}

const getLectureByCode = `-- name: GetLectureByCode :one
SELECT id, name, code, credit, type FROM lectures WHERE code = ?
`

func (q *Queries) GetLectureByCode(ctx context.Context, code string) (Lecture, error) {
	row := q.db.QueryRowContext(ctx, getLectureByCode, code)
	var i Lecture
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Code,
		&i.Credit,
		&i.Type,
	)
	return i, err
}

const listLectures = `-- name: ListLectures :many
SELECT id, name, code, credit, type FROM lectures ORDER BY code
`

func (q *Queries) ListLectures(ctx context.Context) ([]Lecture, error) {
	rows, err := q.db.QueryContext(ctx, listLectures)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Lecture
	for rows.Next() {
		var i Lecture
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Code,
			&i.Credit,
			&i.Type,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}