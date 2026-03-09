package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestOrderRepository_BeginTx_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	repo := NewOrderRepository(sqlDB)
	mock.ExpectBegin()

	tx, err := repo.BeginTx(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mock.ExpectRollback()
	_ = tx.Rollback()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestOrderRepository_BeginTx_Error(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	repo := NewOrderRepository(sqlDB)
	mock.ExpectBegin().WillReturnError(errors.New("begin failed"))

	_, err = repo.BeginTx(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "repository error: failed to create transaction" {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestOrderRepository_GetProductPrice_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	repo := NewOrderRepository(sqlDB)
	mock.ExpectBegin()
	tx, err := sqlDB.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}

	rows := sqlmock.NewRows([]string{"price"}).AddRow(200)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT price FROM products WHERE id=$1")).
		WithArgs(5).
		WillReturnRows(rows)

	price, err := repo.GetProductPrice(context.Background(), tx, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if price != 200 {
		t.Fatalf("unexpected price: got %d want %d", price, 200)
	}

	mock.ExpectRollback()
	_ = tx.Rollback()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestOrderRepository_GetProductPrice_Error(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	repo := NewOrderRepository(sqlDB)
	mock.ExpectBegin()
	tx, err := sqlDB.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT price FROM products WHERE id=$1")).
		WithArgs(5).
		WillReturnError(errors.New("not found"))

	_, err = repo.GetProductPrice(context.Background(), tx, 5)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "repository error: failed to find a product in order" {
		t.Fatalf("unexpected error: %v", err)
	}

	mock.ExpectRollback()
	_ = tx.Rollback()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestOrderRepository_CreateOrder_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	repo := NewOrderRepository(sqlDB)
	mock.ExpectBegin()
	tx, err := sqlDB.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(42)
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO orders (contact_data, comment, total_price) VALUES ($1, $2, $3) RETURNING id")).
		WithArgs("phone", "comment", 300).
		WillReturnRows(rows)

	orderID, err := repo.CreateOrder(context.Background(), tx, "phone", "comment", 300)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if orderID != 42 {
		t.Fatalf("unexpected order id: got %d want %d", orderID, 42)
	}

	mock.ExpectRollback()
	_ = tx.Rollback()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestOrderRepository_CreateOrder_Error(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	repo := NewOrderRepository(sqlDB)
	mock.ExpectBegin()
	tx, err := sqlDB.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO orders (contact_data, comment, total_price) VALUES ($1, $2, $3) RETURNING id")).
		WithArgs("phone", "comment", 300).
		WillReturnError(errors.New("insert failed"))

	_, err = repo.CreateOrder(context.Background(), tx, "phone", "comment", 300)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "repository error: failed to insert Order" {
		t.Fatalf("unexpected error: %v", err)
	}

	mock.ExpectRollback()
	_ = tx.Rollback()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestOrderRepository_AddOrderProduct_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	repo := NewOrderRepository(sqlDB)
	mock.ExpectBegin()
	tx, err := sqlDB.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO order_products (order_id, product_id, quantity) VALUES ($1, $2, $3)")).
		WithArgs(1, 2, 3).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.AddOrderProduct(context.Background(), tx, 1, 2, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mock.ExpectRollback()
	_ = tx.Rollback()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestOrderRepository_AddOrderProduct_Error(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	repo := NewOrderRepository(sqlDB)
	mock.ExpectBegin()
	tx, err := sqlDB.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO order_products (order_id, product_id, quantity) VALUES ($1, $2, $3)")).
		WithArgs(1, 2, 3).
		WillReturnError(errors.New("insert failed"))

	err = repo.AddOrderProduct(context.Background(), tx, 1, 2, 3)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "failed to insert order products" {
		t.Fatalf("unexpected error: %v", err)
	}

	mock.ExpectRollback()
	_ = tx.Rollback()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

