package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"eshbuket/internal/service/order"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type orderRepoStub struct {
	beginTxErr error
}

func (s *orderRepoStub) BeginTx(ctx context.Context) (*sql.Tx, error) {
	if s.beginTxErr != nil {
		return nil, s.beginTxErr
	}
	return nil, nil
}

func (s *orderRepoStub) GetProductPrice(ctx context.Context, tx *sql.Tx, productID int) (int64, error) {
	return 0, nil
}

func (s *orderRepoStub) CreateOrder(ctx context.Context, tx *sql.Tx, contactData string, comment string, totalPriceCents int64) (int, error) {
	return 0, nil
}

func (s *orderRepoStub) AddOrderProduct(ctx context.Context, tx *sql.Tx, orderID int, productID int, quantity int) error {
	return nil
}

func TestOrdersHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewOrderHandler(*order.NewOrderService(&orderRepoStub{}))

	router := gin.New()
	router.POST("/api/orders", h.OrdersHandler)

	req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewBufferString(`{"contact_data":"x"`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestOrdersHandler_EmptyProducts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewOrderHandler(*order.NewOrderService(&orderRepoStub{}))

	router := gin.New()
	router.POST("/api/orders", h.OrdersHandler)

	req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewBufferString(`{"contact_data":"x","products":[]}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestOrdersHandler_ServiceError_ReturnsSingle500Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewOrderHandler(*order.NewOrderService(&orderRepoStub{
		beginTxErr: errors.New("begin tx failed"),
	}))

	router := gin.New()
	router.POST("/api/orders", h.OrdersHandler)

	req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewBufferString(`{"contact_data":"x","products":[{"product_id":1,"quantity":1}]}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
	body := rec.Body.String()
	if strings.Contains(body, "\"order_id\"") || strings.Contains(body, "\"total_price\"") {
		t.Fatalf("unexpected success payload in error response: %s", body)
	}
}
