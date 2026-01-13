package handlers

import (
	db "eshbuket/database"
	handlers "eshbuket/handlers/structures"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OrdersHandler обрабатывает POST /api/orders
func OrdersHandler(c *gin.Context) {
	var req handlers.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if len(req.Products) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no products in order"})
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create order"})
		return
	}

	var totalPrice int
	for _, p := range req.Products {
		// Получаем цену товара из БД
		var price int
		err := tx.QueryRow("SELECT price FROM products WHERE id=$1", p.ProductID).Scan(&price)
		if err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "product not found"})
			return
		}
		totalPrice += price * p.Quantity
	}

	// Вставляем заказ с total_price
	var orderID int
	err = tx.QueryRow(
		"INSERT INTO orders (contact_data, comment, total_price) VALUES ($1, $2, $3) RETURNING id",
		req.ContactData, req.Comment, totalPrice,
	).Scan(&orderID)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "failed to create order"})
		return
	}

	// Вставляем записи в order_products
	for _, p := range req.Products {
		_, err = tx.Exec(
			"INSERT INTO order_products (order_id, product_id, quantity) VALUES ($1, $2, $3)",
			orderID, p.ProductID, p.Quantity,
		)
		if err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "failed to insert order products"})
			return
		}

	}
	tx.Commit()
	c.JSON(201, gin.H{"order_id": orderID, "total_price": totalPrice})

}
