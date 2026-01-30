package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	db "eshbuket/internal/database"
	"eshbuket/internal/handlers"
	models "eshbuket/internal/models"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/api/products", handlers.ProductsHandler)
	return r
}

func TestProductsHandler(t *testing.T) {
	// Подключаемся к тестовой БД
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db.Connect(dbHost, dbPort, dbUser, dbPassword, dbName)

	// Очищаем таблицу перед тестом
	_, err := db.DB.Exec("DELETE FROM products")
	if err != nil {
		t.Fatalf("Failed to clear products table: %v", err)
	}

	// Добавим тестовые данные
	_, err = db.DB.Exec(`INSERT INTO products (name, price, category) VALUES
		('Шоколадный букет', '1500', 'chocolate'),
		('Фруктовый букет', '1200', 'fruits')`)
	if err != nil {
		t.Fatalf("Failed to insert test products: %v", err)
	}

	router := setupRouter()

	t.Run("Get all products", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/products", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatalf("Expected 200, got %d", w.Code)
		}

		var products []models.Product
		if err := json.Unmarshal(w.Body.Bytes(), &products); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(products) != 2 {
			t.Fatalf("Expected 2 products, got %d", len(products))
		}
	})

	t.Run("Get products by category", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/products?category=chocolate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatalf("Expected 200, got %d", w.Code)
		}

		var products []models.Product
		if err := json.Unmarshal(w.Body.Bytes(), &products); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(products) != 1 || products[0].Category != "chocolate" {
			t.Fatalf("Expected 1 chocolate product, got %+v", products)
		}
	})

	t.Run("Get products by nonexistent category", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/products?category=nonexistent", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatalf("Expected 200, got %d", w.Code)
		}

		var products []models.Product
		if err := json.Unmarshal(w.Body.Bytes(), &products); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(products) != 0 {
			t.Fatalf("Expected 0 products, got %d", len(products))
		}
	})
}
