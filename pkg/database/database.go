package database

import (
	"context"
	"database/sql"
	"fmt"

	// Imported to use sqlite3 with sqlc.
	_ "github.com/mattn/go-sqlite3"

	"github.com/TulgaCG/add-drop-classes-api/database"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
)

// go embed's doesnt work when we move this file to pkg/database.

func New(path string) (*gendb.Queries, error) {
	d, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	// Create tables if not exist
	if _, err := d.ExecContext(context.Background(), database.Schema); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	db := gendb.New(d)

	return db, nil
}

func AddMockData(path string) error {
	d, err := sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("failed to add mock data: %w", err)
	}

	if _, err := d.ExecContext(context.Background(), database.Mockdata); err != nil {
		return fmt.Errorf("failed to exec query: %w", err)
	}

	return nil
}
