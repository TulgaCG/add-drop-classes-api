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

func New(ctx context.Context, path string) (*gendb.Queries, error) {
	d, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	// Create tables if not exist
	if _, err := d.ExecContext(ctx, database.Schema); err != nil {
		return nil, fmt.Errorf("failed to create tables1: %w", err)
	}

	return gendb.New(d), nil
}
