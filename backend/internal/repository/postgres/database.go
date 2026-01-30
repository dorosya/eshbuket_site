package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Connect(dbHost string, dbPort string, dbUser string, dbPassword string, dbName string) *sql.DB {

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Не удалось пинговать БД: %v", err)
	}

	log.Println("Успешное подключение к PostgreSQL!")
	return db
}

// Сделать нормальную миграцию надо
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
