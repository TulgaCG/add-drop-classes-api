package database

import (
	"context"
	"database/sql"
	"fmt"

	// Imported to use sqlite3 with sqlc.
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

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
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	// Create roles if not exist
	if _, err := d.ExecContext(ctx, database.Roles); err != nil {
		return nil, fmt.Errorf("failed to create roles: %w", err)
	}

	return gendb.New(d), nil
}

func NewTestDb(ctx context.Context) (*gendb.Queries, error) {
	d, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	if _, err := d.ExecContext(ctx, database.Schema); err != nil {
		return nil, err
	}

	db := gendb.New(d)

	for i := 1; i <= 5; i++ {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("testpassword%d", i)), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to generate hashed password: %w", err)
		}

		if _, err := db.CreateRole(ctx, fmt.Sprintf("testrole%d", i)); err != nil {
			return nil, err
		}

		if _, err := db.CreateUser(ctx, gendb.CreateUserParams{
			Username: fmt.Sprintf("testuser%d", i),
			Password: string(hashedPassword),
		}); err != nil {
			return nil, err
		}
	}

	return db, nil
}
