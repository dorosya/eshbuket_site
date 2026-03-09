package postgres

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInitSchema(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create sqlmock: %v", err)
		}
		defer func() { _ = sqlDB.Close() }()

		db = sqlDB
		mock.ExpectExec(regexp.QuoteMeta(`
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
	`)).WillReturnResult(sqlmock.NewResult(0, 0))

		if err := InitSchema(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet sqlmock expectations: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create sqlmock: %v", err)
		}
		defer func() { _ = sqlDB.Close() }()

		db = sqlDB
		mock.ExpectExec(regexp.QuoteMeta(`
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
	`)).WillReturnError(errors.New("exec failed"))

		err = InitSchema()
		if err == nil {
			t.Fatal("expected error")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet sqlmock expectations: %v", err)
		}
	})
}

