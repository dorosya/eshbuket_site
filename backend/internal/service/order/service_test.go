package order

import (
	"context"
	"database/sql"
	"errors"
	"eshbuket/internal/transport/http/dto"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type orderRepoMock struct {
	beginTxFn      func(ctx context.Context) (*sql.Tx, error)
	getPriceFn     func(ctx context.Context, tx *sql.Tx, productID int) (int64, error)
	createOrderFn  func(ctx context.Context, tx *sql.Tx, contactData string, comment string, totalPriceCents int64) (int, error)
	addOrderProdFn func(ctx context.Context, tx *sql.Tx, orderID int, productID int, quantity int) error
}

func (m *orderRepoMock) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return m.beginTxFn(ctx)
}

func (m *orderRepoMock) GetProductPrice(ctx context.Context, tx *sql.Tx, productID int) (int64, error) {
	return m.getPriceFn(ctx, tx, productID)
}

func (m *orderRepoMock) CreateOrder(ctx context.Context, tx *sql.Tx, contactData string, comment string, totalPriceCents int64) (int, error) {
	return m.createOrderFn(ctx, tx, contactData, comment, totalPriceCents)
}

func (m *orderRepoMock) AddOrderProduct(ctx context.Context, tx *sql.Tx, orderID int, productID int, quantity int) error {
	return m.addOrderProdFn(ctx, tx, orderID, productID, quantity)
}

func TestOrderService_OrderServiceFunc_Success(t *testing.T) {
	sqlDB, sm, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	sm.ExpectBegin()
	tx, err := sqlDB.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	sm.ExpectCommit()

	added := 0
	svc := NewOrderService(&orderRepoMock{
		beginTxFn: func(ctx context.Context) (*sql.Tx, error) { return tx, nil },
		getPriceFn: func(ctx context.Context, tx *sql.Tx, productID int) (int64, error) {
			if productID == 1 {
				return 2500, nil
			}
			return 1000, nil
		},
		createOrderFn: func(ctx context.Context, tx *sql.Tx, contactData string, comment string, totalPriceCents int64) (int, error) {
			if totalPriceCents != 6000 {
				t.Fatalf("unexpected total: got %d want 6000", totalPriceCents)
			}
			return 42, nil
		},
		addOrderProdFn: func(ctx context.Context, tx *sql.Tx, orderID int, productID int, quantity int) error {
			added++
			return nil
		},
	})

	orderID, total, err := svc.OrderServiceFunc(context.Background(), dto.OrderRequest{
		ContactData: "phone",
		Comment:     "fast",
		Products: []dto.Product{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 1},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if orderID != 42 {
		t.Fatalf("unexpected order id: %d", orderID)
	}
	if total != 6000 {
		t.Fatalf("unexpected total: %d", total)
	}
	if added != 2 {
		t.Fatalf("expected 2 add-order-product calls, got %d", added)
	}
	if err := sm.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestOrderService_OrderServiceFunc_BeginTxError(t *testing.T) {
	expected := errors.New("begin failed")
	svc := NewOrderService(&orderRepoMock{
		beginTxFn: func(ctx context.Context) (*sql.Tx, error) { return nil, expected },
	})

	_, _, err := svc.OrderServiceFunc(context.Background(), dto.OrderRequest{})
	if !errors.Is(err, expected) {
		t.Fatalf("expected begin error, got %v", err)
	}
}

func TestOrderService_OrderServiceFunc_RollbackOnPriceError(t *testing.T) {
	sqlDB, sm, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	sm.ExpectBegin()
	tx, err := sqlDB.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	sm.ExpectRollback()

	expected := errors.New("price error")
	svc := NewOrderService(&orderRepoMock{
		beginTxFn:  func(ctx context.Context) (*sql.Tx, error) { return tx, nil },
		getPriceFn: func(ctx context.Context, tx *sql.Tx, productID int) (int64, error) { return 0, expected },
	})

	_, _, err = svc.OrderServiceFunc(context.Background(), dto.OrderRequest{
		Products: []dto.Product{{ProductID: 1, Quantity: 1}},
	})
	if !errors.Is(err, expected) {
		t.Fatalf("expected price error, got %v", err)
	}
	if err := sm.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestOrderService_OrderServiceFunc_RollbackOnCreateOrderError(t *testing.T) {
	sqlDB, sm, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	sm.ExpectBegin()
	tx, err := sqlDB.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	sm.ExpectRollback()

	expected := errors.New("create error")
	svc := NewOrderService(&orderRepoMock{
		beginTxFn:  func(ctx context.Context) (*sql.Tx, error) { return tx, nil },
		getPriceFn: func(ctx context.Context, tx *sql.Tx, productID int) (int64, error) { return 1000, nil },
		createOrderFn: func(ctx context.Context, tx *sql.Tx, contactData string, comment string, totalPriceCents int64) (int, error) {
			return 0, expected
		},
	})

	_, _, err = svc.OrderServiceFunc(context.Background(), dto.OrderRequest{
		Products: []dto.Product{{ProductID: 1, Quantity: 1}},
	})
	if !errors.Is(err, expected) {
		t.Fatalf("expected create error, got %v", err)
	}
	if err := sm.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestOrderService_OrderServiceFunc_RollbackOnAddOrderProductError(t *testing.T) {
	sqlDB, sm, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	sm.ExpectBegin()
	tx, err := sqlDB.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	sm.ExpectRollback()

	expected := errors.New("add error")
	svc := NewOrderService(&orderRepoMock{
		beginTxFn:  func(ctx context.Context) (*sql.Tx, error) { return tx, nil },
		getPriceFn: func(ctx context.Context, tx *sql.Tx, productID int) (int64, error) { return 1000, nil },
		createOrderFn: func(ctx context.Context, tx *sql.Tx, contactData string, comment string, totalPriceCents int64) (int, error) {
			return 77, nil
		},
		addOrderProdFn: func(ctx context.Context, tx *sql.Tx, orderID int, productID int, quantity int) error {
			return expected
		},
	})

	_, _, err = svc.OrderServiceFunc(context.Background(), dto.OrderRequest{
		Products: []dto.Product{{ProductID: 1, Quantity: 1}},
	})
	if !errors.Is(err, expected) {
		t.Fatalf("expected add-order-product error, got %v", err)
	}
	if err := sm.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestOrderService_OrderServiceFunc_CommitErrorTriggersRollback(t *testing.T) {
	sqlDB, sm, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	sm.ExpectBegin()
	tx, err := sqlDB.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	sm.ExpectCommit().WillReturnError(errors.New("commit failed"))

	svc := NewOrderService(&orderRepoMock{
		beginTxFn:  func(ctx context.Context) (*sql.Tx, error) { return tx, nil },
		getPriceFn: func(ctx context.Context, tx *sql.Tx, productID int) (int64, error) { return 1000, nil },
		createOrderFn: func(ctx context.Context, tx *sql.Tx, contactData string, comment string, totalPriceCents int64) (int, error) {
			return 1, nil
		},
		addOrderProdFn: func(ctx context.Context, tx *sql.Tx, orderID int, productID int, quantity int) error { return nil },
	})

	_, _, err = svc.OrderServiceFunc(context.Background(), dto.OrderRequest{
		Products: []dto.Product{{ProductID: 1, Quantity: 1}},
	})
	if err == nil {
		t.Fatal("expected commit error")
	}
	if err := sm.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
