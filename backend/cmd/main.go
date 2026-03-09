package main

import (
	"eshbuket/internal/repository/postgres"
	"eshbuket/internal/transport/http/dto"
	"eshbuket/internal/transport/http/handlers"
	"log"
	"os"

	"github.com/gin-contrib/cors"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db := postgres.Connect(dbHost, dbPort, dbUser, dbPassword, dbName)

	if err := postgres.RunMigrations(db, os.Getenv("MIGRATIONS_DIR")); err != nil {
		log.Fatalf("failed to run DB migrations: %v", err)
	}

	env := os.Getenv("APP_ENV")
	router := dto.NewRouter(dto.Config{Env: env})
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))
	handlers.RegisterRoutes(router, db)

	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}

