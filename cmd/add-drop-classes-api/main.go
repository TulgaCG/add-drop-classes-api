package main

import (
	"log"

	"github.com/TulgaCG/add-drop-classes-api/database"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server"
)

const dbPath = "test.sqlite"

func main() {
	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("failed to create db connection: %s", err.Error())
	}

	_ = database.AddMockData(dbPath)

	s := server.New(db)
	err = s.Run()
	if err != nil {
		log.Fatalf("failed to run server instance: %s", err.Error())
	}
}
