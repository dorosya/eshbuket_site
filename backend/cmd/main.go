package main

import (
	"eshbuket/internal/repository/postgres"
	"eshbuket/internal/transport/http/dto"
	"eshbuket/internal/transport/http/handlers"
	"log"
	"os"
	"strings"
	"time"

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
	config := cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		MaxAge:       12 * time.Hour,
	}
	origins := parseOrigins(os.Getenv("FRONTEND_ORIGINS"))
	if len(origins) == 0 {
		origins = []string{
			"http://localhost:5500",
			"http://127.0.0.1:5500",
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		}
	}
	config.AllowOrigins = origins
	config.AllowCredentials = true
	router.Use(cors.New(config))
	handlers.RegisterRoutes(router, db)

	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}

func parseOrigins(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			origins = append(origins, v)
		}
	}
	return origins
}
