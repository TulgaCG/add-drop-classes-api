package server

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"

	"github.com/TulgaCG/add-drop-classes-api/cmd/add-drop-classes-api/server/internal/logger"
	"github.com/TulgaCG/add-drop-classes-api/pkg/database"
	"github.com/TulgaCG/add-drop-classes-api/pkg/env"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server"
)

//nolint:govet
type Cmd struct {
	Database struct {
		User     string `help:"PostgreSQL database username" default:"postgres" env:"ADC_DB_USER"`
		Password string `help:"PostgreSQL database password" default:"postgres" env:"ADC_DB_PASSWORD"`
		Host     string `help:"PostgreSQL database host"     default:"127.0.0.1" env:"ADC_DB_HOST"`
		Port     uint32 `help:"PostgreSQL database port"     default:"5432" env:"ADC_DB_PORT"`
		Database string `help:"PostgreSQL database name"     default:"postgres" env:"ADC_DB_NAME"`
	} `embed:"" prefix:"db."`
	Port     uint32       `help:"Port the server listens" default:"8080"`
	EnvLevel env.Level    `help:"Environment level (${enum})" default:"dev" enum:"dev," env:"ADC_ENV"`
	LogLevel logger.Level `help:"Log level (${enum})" default:"info" enum:"debug,info,warn,error" env:"ADC_LOG_LEVEL"`
}

func (c *Cmd) Run(_ *kong.Kong) error {
	ctx := context.Background()
	log := logger.New(c.LogLevel)

	db, err := database.New(ctx, c.Database.User, c.Database.Password, c.Database.Host, c.Database.Database, c.Database.Port)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	s := server.New(
		server.WithPort(c.Port),
		server.WithDB(db),
		server.WithEnv(c.EnvLevel),
		server.WithLogger(log),
	)
	if err = s.Run(); err != nil {
		return fmt.Errorf("failed to run server: %w", err)
	}

	return nil
}
