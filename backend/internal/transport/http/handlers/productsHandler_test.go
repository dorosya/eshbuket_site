package handlers

import (
	"bytes"
	"eshbuket/internal/service/product"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestProductsHandler_InvalidJSON_ReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewProductHandler(product.ProductService{})

	router := gin.New()
	router.POST("/api/products", h.ProductsHandler)

	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBufferString(`{"name":"Rose"`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
