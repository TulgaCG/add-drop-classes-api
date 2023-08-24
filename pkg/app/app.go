package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/TulgaCG/add-drop-classes-api/pkg/database"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server"
)

type Conf struct {
	Log    *slog.Logger
	DbPath string
}

func Run(conf *Conf) error {
	// db, err := database.New(context.Background(), conf.DbPath)
	// if err != nil {
	// 	return fmt.Errorf("failed to create db connection: %w", err)
	// }

	db, err := database.NewTestDb(context.Background())
	if err != nil {
		return fmt.Errorf("failed to create db connection: %w", err)
	}

	s := server.New(db, conf.Log)
	if err = s.Run(); err != nil {
		return fmt.Errorf("failed to run server instance: %w", err)
	}

	return nil
}
