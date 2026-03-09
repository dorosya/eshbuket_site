package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestProductRepository_FindProducts_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	repo := NewProductRepository(sqlDB)
	query := regexp.QuoteMeta(`
			SELECT id, name, price, category, image_path
			FROM products
			WHERE ($1 = '' OR category = $1)
		`)

	rows := sqlmock.NewRows([]string{"id", "name", "price", "category", "image_path"}).
		AddRow(1, "Rose", 100, "flowers", nil)
	mock.ExpectQuery(query).WithArgs("flowers").WillReturnRows(rows)

	gotRows, err := repo.FindProducts(context.Background(), "flowers")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = gotRows.Close() }()

	if !gotRows.Next() {
		t.Fatal("expected at least one row")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestProductRepository_FindProducts_Error(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	repo := NewProductRepository(sqlDB)
	query := regexp.QuoteMeta(`
			SELECT id, name, price, category, image_path
			FROM products
			WHERE ($1 = '' OR category = $1)
		`)

	mock.ExpectQuery(query).WithArgs("").WillReturnError(errors.New("db down"))

	_, err = repo.FindProducts(context.Background(), "")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "repository error: cannot find an item during scanning products rows" {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestProductRepository_InsertProduct_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	repo := NewProductRepository(sqlDB)
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO products (name, price, category) VALUES ($1, $2, $3)")).
		WithArgs("Rose", "100", "flowers").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.InsertProduct(context.Background(), "Rose", "100", "flowers")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestProductRepository_InsertProduct_Error(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	repo := NewProductRepository(sqlDB)
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO products (name, price, category) VALUES ($1, $2, $3)")).
		WithArgs("Rose", "100", "flowers").
		WillReturnError(errors.New("insert failed"))

	err = repo.InsertProduct(context.Background(), "Rose", "100", "flowers")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "repository error: insert failed" {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

