package main

import (
	"github.com/TulgaCG/add-drop-classes-api/pkg/app"
	"log"
	"log/slog"
	"os"
)

const dbPath = "test.sqlite"

func main() {
	slogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	if err := app.Run(&app.Conf{
		DbPath: dbPath,
		Log:    slogger,
	}); err != nil {
		log.Fatal("failed to run application: %w", err)
	}
}
