package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(dbHost string, dbPort string, dbUser string, dbPassword string, dbName string) {

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Не удалось пинговать БД: %v", err)
	}

	log.Println("Успешное подключение к PostgreSQL!")
}

func InitSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			price INT NOT NULL,
			category TEXT
		);

		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			contact_phone TEXT NOT NULL,
			comment TEXT,
			product_id INT REFERENCES products(id)
		);
	`)
	return err
}
