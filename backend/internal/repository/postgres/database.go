package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func Connect(dbHost string, dbPort string, dbUser string, dbPassword string, dbName string) *sql.DB {
	dsn := BuildDSN(dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}

	log.Println("connected to PostgreSQL")
	return db
}

func BuildDSN(dbHost string, dbPort string, dbUser string, dbPassword string, dbName string) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)
}

func RunMigrations(db *sql.DB, migrationsDir string) error {
	if strings.TrimSpace(migrationsDir) == "" {
		migrationsDir = "./migrations"
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}
	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
