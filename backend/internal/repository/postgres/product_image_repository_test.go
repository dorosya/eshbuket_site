package postgres

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestProductImageRepository_ProductExists_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	repo := NewProductImageRepository(sqlDB)
	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM products WHERE id=$1)")).
		WithArgs("10").
		WillReturnRows(rows)

	exists, err := repo.ProductExists(context.Background(), "10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Fatal("expected exists=true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestProductImageRepository_ProductExists_Error(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	repo := NewProductImageRepository(sqlDB)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM products WHERE id=$1)")).
		WithArgs("10").
		WillReturnError(errors.New("db error"))

	_, err = repo.ProductExists(context.Background(), "10")
	if err == nil {
		t.Fatal("expected error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestProductImageRepository_GetImagePath_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	repo := NewProductImageRepository(sqlDB)
	rows := sqlmock.NewRows([]string{"image_path"}).AddRow("products/10/a.png")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT image_path FROM products WHERE id=$1")).
		WithArgs("10").
		WillReturnRows(rows)

	path, ok, err := repo.GetImagePath(context.Background(), "10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok || path != "products/10/a.png" {
		t.Fatalf("unexpected result: path=%q ok=%v", path, ok)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestProductImageRepository_GetImagePath_NullOrEmpty(t *testing.T) {
	cases := []struct {
		name string
		row  interface{}
	}{
		{name: "null", row: nil},
		{name: "empty", row: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer func() { _ = sqlDB.Close() }()

			repo := NewProductImageRepository(sqlDB)
			rows := sqlmock.NewRows([]string{"image_path"}).AddRow(tc.row)
			mock.ExpectQuery(regexp.QuoteMeta("SELECT image_path FROM products WHERE id=$1")).
				WithArgs("10").
				WillReturnRows(rows)

			path, ok, err := repo.GetImagePath(context.Background(), "10")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ok || path != "" {
				t.Fatalf("unexpected result: path=%q ok=%v", path, ok)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet sqlmock expectations: %v", err)
			}
		})
	}
}

func TestProductImageRepository_GetImagePath_Error(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	repo := NewProductImageRepository(sqlDB)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT image_path FROM products WHERE id=$1")).
		WithArgs("10").
		WillReturnError(sql.ErrNoRows)

	_, _, err = repo.GetImagePath(context.Background(), "10")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestProductImageRepository_SetImagePath(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create sqlmock: %v", err)
		}
		defer func() { _ = sqlDB.Close() }()

		repo := NewProductImageRepository(sqlDB)
		mock.ExpectExec(regexp.QuoteMeta("UPDATE products SET image_path=$1 WHERE id=$2")).
			WithArgs("products/10/a.png", "10").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err = repo.SetImagePath(context.Background(), "10", "products/10/a.png")
		if err != nil {
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

		repo := NewProductImageRepository(sqlDB)
		mock.ExpectExec(regexp.QuoteMeta("UPDATE products SET image_path=$1 WHERE id=$2")).
			WithArgs("products/10/a.png", "10").
			WillReturnError(errors.New("update failed"))

		err = repo.SetImagePath(context.Background(), "10", "products/10/a.png")
		if err == nil {
			t.Fatal("expected error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet sqlmock expectations: %v", err)
		}
	})
}

