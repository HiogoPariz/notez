package main

import (
	"log"

	"github.com/HiogoPariz/notez/apps/api/internal/api"
	"github.com/HiogoPariz/notez/apps/api/internal/db"
	"github.com/HiogoPariz/notez/apps/api/internal/migration"
)

func main() {
	// Init db connection
	repo, err := db.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	// Run migrations
	if err := migration.Run(repo); err != nil {
		log.Fatal(err)
	}

	// Start api server
	api.Init(repo)
}
