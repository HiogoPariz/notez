package main

import (
	"log"

	"github.com/HiogoPariz/notez/apps/files/internal/api"
	"github.com/HiogoPariz/notez/apps/files/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	// .env is optional: the root Makefile exports the variables when running
	// the whole stack, and the defaults below cover local development.
	_ = godotenv.Load()

	router := api.SetupRouter()

	if err := router.Run(":" + config.Get("PORT", "3001")); err != nil {
		log.Fatal(err)
	}
}
