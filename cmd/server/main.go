package main

import (
	"log"
	"portfolio-server/internal/config"
	"portfolio-server/internal/database"
	"portfolio-server/internal/middleware"
	"portfolio-server/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	if err := database.InitDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	middleware.InitJWT(&cfg.JWT)

	if cfg.Server.ENV == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(middleware.CORS())
	router.Use(middleware.RecoveryHandler())
	router.Use(middleware.ErrorHandler())

	routes.SetupRoutes(router)

	addr := ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
