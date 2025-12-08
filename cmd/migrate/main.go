package main

import (
	"log"
	"portfolio-server/internal/config"
	"portfolio-server/internal/database"
)

func main() {
	cfg := config.LoadConfig()

	if err := database.InitDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migration completed successfully!")
}
