package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phayes/freeport"

	"github.com/TulgaCG/add-drop-classes-api/database/schema"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
)

const (
	localhost = "127.0.0.1"
	pg        = "postgres"
)

type DB struct {
	*gendb.Queries
	Pool    *pgxpool.Pool
	closeFn func() error
}

func (db *DB) WithTx(ctx context.Context, fn func(dbTx *gendb.Queries) error) error {
	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	if err := fn(db.Queries.WithTx(tx)); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return errors.Join(err, fmt.Errorf("failed to rollback transaction after an error in operation: %w", rbErr))
		}
		return fmt.Errorf("failed to execute operation with transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return errors.Join(err, fmt.Errorf("failed to rollback transaction after an error in commit: %w", rbErr))
		}
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (db *DB) Close() error {
	return db.closeFn()
}

func New(ctx context.Context, user, password, host, db string, port uint32) (*DB, error) {
	p, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, db))
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &DB{
		closeFn: func() error {
			p.Close()
			return nil
		},
		Pool:    p,
		Queries: gendb.New(p),
	}, nil
}

func NewTestDB(ctx context.Context) (*DB, func(t *testing.T), error) {
	rtPath, err := os.MkdirTemp("", "add-drop-classes-*")
	if err != nil {
		return nil, nil, err
	}

	port, err := freeport.GetFreePort()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get free port: %w", err)
	}

	pgDB := embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Database(pg).
			Username(pg).
			Password(pg).
			Port(uint32(port)).
			RuntimePath(rtPath),
	)
	if err := pgDB.Start(); err != nil {
		return nil, nil, fmt.Errorf("failed to start database: %w", err)
	}

	if err := ExecuteSQL(ctx, schema.Schema, pg, pg, localhost, pg, uint32(port)); err != nil {
		return nil, nil, fmt.Errorf("failed to execute schema SQL: %w", err)
	}

	db, err := New(ctx, pg, pg, localhost, pg, uint32(port))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create database: %w", err)
	}

	return db, func(t *testing.T) {
		t.Helper()

		if err := db.Close(); err != nil {
			t.Fatal(err.Error())
		}

		if err := pgDB.Stop(); err != nil {
			t.Fatal(err.Error())
		}

		if err := os.RemoveAll(rtPath); err != nil {
			t.Fatal(err.Error())
		}
	}, nil
}

func ExecuteSQL(ctx context.Context, sql, user, password, host, db string, port uint32) error {
	execConn, err := pgx.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, db))
	if err != nil {
		return fmt.Errorf("failed to create postgresql connection: %w", err)
	}
	defer func(execConn *pgx.Conn, ctx context.Context) {
		err := execConn.Close(ctx)
		if err != nil {
			log.Fatal(err.Error())
		}
	}(execConn, ctx)

	if _, err := execConn.Exec(ctx, sql); err != nil {
		return fmt.Errorf("failed to execute sql: %w", err)
	}

	return nil
}
