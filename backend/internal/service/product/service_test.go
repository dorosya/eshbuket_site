package product

import (
	"context"
	"database/sql"
	"errors"
	"eshbuket/internal/transport/http/dto"
	"testing"
)

type productRepoStub struct {
	insertFn func(ctx context.Context, name string, priceRub int64, category string) error
}

func (s *productRepoStub) FindProducts(ctx context.Context, category string) (*sql.Rows, error) {
	return nil, errors.New("not used in this test")
}

func (s *productRepoStub) InsertProduct(ctx context.Context, name string, priceRub int64, category string) error {
	if s.insertFn != nil {
		return s.insertFn(ctx, name, priceRub, category)
	}
	return nil
}

func TestProductPostService_InvalidPriceFormat(t *testing.T) {
	svc := NewProductService(&productRepoStub{})

	err := svc.ProductPostService(context.Background(), dto.ProductRequest{
		Name:     "Rose",
		Price:    "12.5",
		Category: "flowers",
	})

	if err == nil || err.Error() != "invalid price format" {
		t.Fatalf("expected invalid price format, got: %v", err)
	}
}

func TestProductPostService_NonPositivePrice(t *testing.T) {
	svc := NewProductService(&productRepoStub{})

	err := svc.ProductPostService(context.Background(), dto.ProductRequest{
		Name:     "Rose",
		Price:    "0",
		Category: "flowers",
	})

	if err == nil || err.Error() != "price must be positive" {
		t.Fatalf("expected price must be positive, got: %v", err)
	}
}

func TestProductPostService_PropagatesRepositoryError(t *testing.T) {
	expectedErr := errors.New("repository error")
	svc := NewProductService(&productRepoStub{
		insertFn: func(ctx context.Context, name string, priceRub int64, category string) error {
			return expectedErr
		},
	})

	err := svc.ProductPostService(context.Background(), dto.ProductRequest{
		Name:     "Rose",
		Price:    "100",
		Category: "flowers",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected repo error to be returned, got: %v", err)
	}
}

func TestProductPostService_SuccessCallsRepository(t *testing.T) {
	called := false
	svc := NewProductService(&productRepoStub{
		insertFn: func(ctx context.Context, name string, priceRub int64, category string) error {
			called = true
			if name != "Rose" || priceRub != 100 || category != "flowers" {
				t.Fatalf("unexpected insert args: %q %d %q", name, priceRub, category)
			}
			return nil
		},
	})

	err := svc.ProductPostService(context.Background(), dto.ProductRequest{
		Name:     "Rose",
		Price:    "100",
		Category: "flowers",
	})

	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
	if !called {
		t.Fatal("expected repository InsertProduct to be called")
	}
}

