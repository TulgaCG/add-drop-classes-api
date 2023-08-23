package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"

	"github.com/TulgaCG/add-drop-classes-api/pkg/database"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server"
)

type Conf struct {
	Log    *slog.Logger
	DbPath string
}

func Run(conf *Conf) error {
	db, err := database.New(context.Background(), conf.DbPath)
	if err != nil {
		return fmt.Errorf("failed to create db connection: %w", err)
	}

	v := validator.New()

	s := server.New(db, conf.Log, v)
	if err = s.Run(); err != nil {
		return fmt.Errorf("failed to run server instance: %w", err)
	}

	return nil
}
