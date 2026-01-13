package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	db "eshbuket/database"
	"eshbuket/handlers"
	structures "eshbuket/handlers/structures"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func setup_Router() *gin.Engine {
	r := gin.Default()
	r.POST("/api/orders", handlers.OrdersHandler)
	return r
}

func TestOrdersHandler(t *testing.T) {
	// Подключаемся к тестовой БД
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db.Connect(dbHost, dbPort, dbUser, dbPassword, dbName)

	// Очищаем таблицы перед тестом
	// _, err := db.DB.Exec("DELETE FROM order_products")
	// if err != nil {
	// 	t.Fatalf("Failed to clear order_products table: %v", err)
	// }
	// _, err = db.DB.Exec("DELETE FROM orders")
	// if err != nil {
	// 	t.Fatalf("Failed to clear orders table: %v", err)
	// }

	router := setup_Router()

	// Подготовим тестовые продукты
	_, err := db.DB.Exec(`INSERT INTO products (name, price, category) VALUES
		('Товар A', '100', 'cat1'),
		('Товар B', '200', 'cat2')`)
	if err != nil {
		t.Fatalf("Failed to insert test products: %v", err)
	}

	// Считаем их ID для запроса
	rows, _ := db.DB.Query("SELECT id FROM products ORDER BY id")
	var productIDs []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		productIDs = append(productIDs, id)
	}
	rows.Close()

	t.Run("Create order with multiple products", func(t *testing.T) {
		reqBody := structures.OrderRequest{
			ContactData: "test@example.com",
			Comment:     "Позвонить после 18:00",
			Products: []structures.OrderProductRequest{
				{ProductID: productIDs[0], Quantity: 2},
				{ProductID: productIDs[1], Quantity: 1},
			},
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected 201, got %d, body: %s", w.Code, w.Body.String())
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		orderID, ok := resp["order_id"].(float64)
		if !ok || orderID == 0 {
			t.Fatalf("Invalid order_id in response: %v", resp)
		}

		rows, _ := db.DB.Query("SELECT order_id, product_id FROM order_products")
		for rows.Next() {
			var o, p int
			rows.Scan(&o, &p)
			t.Logf("order_product: order=%d product=%d", o, p)
		}

		// Проверяем, что order_products заполнена
		rows, _ = db.DB.Query("SELECT product_id, quantity FROM order_products WHERE order_id=$1", int(orderID))
		defer rows.Close()

		count := 0
		for rows.Next() {
			count++
		}
		if count != len(reqBody.Products) {
			t.Fatalf("Expected %d products in order_products, got %d", len(reqBody.Products), count)
		}
	})

	t.Run("Fail on empty products", func(t *testing.T) {
		reqBody := structures.OrderRequest{
			ContactData: "test@example.com",
			Comment:     "No products",
			Products:    []structures.OrderProductRequest{},
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected 400, got %d", w.Code)
		}
	})
}
