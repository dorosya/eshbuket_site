package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"eshbuket/internal/transport/http/dto"
	"eshbuket/internal/transport/http/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOrdersHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		request        dto.OrderRequest
		serviceFunc    func(ctx context.Context, req dto.OrderRequest) (int, int, error)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "valid request",
			request: dto.OrderRequest{
				ContactData: "Alice",
				Comment:     "Happy birthday",
				Products: []dto.ProductRequest{
					{ProductID: 1, Quantity: 2},
				},
			},
			serviceFunc: func(ctx context.Context, req dto.OrderRequest) (int, int, error) {
				return 123, 5000, nil
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"order_id":    float64(123),
				"total_price": float64(5000),
			},
		},
		{
			name: "service error",
			request: dto.OrderRequest{
				ContactData: "Alice",
				Comment:     "Hello",
				Products: []dto.ProductRequest{
					{ProductID: 1, Quantity: 2},
				},
			},
			serviceFunc: func(ctx context.Context, req dto.OrderRequest) (int, int, error) {
				return 0, 0, errors.New("service failure")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "service failure",
			},
		},
		{
			name: "empty products",
			request: dto.OrderRequest{
				ContactData: "Alice",
				Comment:     "Hello",
				Products:    []dto.ProductRequest{},
			},
			serviceFunc:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "no products in order",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Подготовка мока сервиса
			mockService := &MockOrderService{
				CreateOrderFunc: tt.serviceFunc,
			}

			handler := handlers.NewOrderHandler(mockService)

			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			handler.OrdersHandler(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Equal(t, tt.expectedBody, resp)
		})
	}
}
