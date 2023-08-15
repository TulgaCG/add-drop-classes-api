package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	// Imported to use sqlite3 with sqlc.
	_ "github.com/mattn/go-sqlite3"

	"github.com/TulgaCG/add-drop-classes-api/pkg/gen/db"
)

//go:embed schema.sql
var ddl string

//go:embed mockdata.sql
var mockdata string

func New(path string) (*db.Queries, error) {
	d, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	// Create tables if not exist
	if _, err := d.ExecContext(context.Background(), ddl); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	query := db.New(d)

	return query, nil
}

func AddMockData(path string) error {
	d, err := sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("failed to add mock data: %w", err)
	}

	if _, err := d.ExecContext(context.Background(), mockdata); err != nil {
		return fmt.Errorf("failed to exec query: %w", err)
	}

	return nil
}
