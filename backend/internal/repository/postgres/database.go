package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var db *sql.DB

func Connect(dbHost string, dbPort string, dbUser string, dbPassword string, dbName string) *sql.DB {
	psqlInfo := BuildDSN(dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	err = db.Ping()
	if err != nil {
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

// Используется дл тестов.
func InitSchema() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			price NUMERIC NOT NULL,
			category TEXT,
			image_path TEXT
		);

		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			contact_data TEXT NOT NULL,
			comment TEXT,
			total_price NUMERIC NOT NULL
		);

		CREATE TABLE IF NOT EXISTS order_products (
			order_id INT REFERENCES orders(id) ON DELETE CASCADE,
			product_id INT REFERENCES products(id) ON DELETE CASCADE,
			quantity INT NOT NULL,
			PRIMARY KEY(order_id, product_id)
		);
	`)
	return err
}
