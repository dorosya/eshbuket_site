package handlers

import (
	db "eshbuket/database"
	handlers "eshbuket/handlers/structures"

	"github.com/gin-gonic/gin"
)

// /api/products - endpoint со списком всех товаров. Парсится из bd и загружает туда же
func ProductsHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		// Get /api/products - select всех товаров из БД по категории
	} else if c.Request.Method == "POST" {
		// Post /api/products - загрузка в БД товаров (только для админки, пропускает через authMiddleware перед обработкой )*/
		var req handlers.ProductRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}
		if req.Category != "" {
			_, err := db.DB.Exec("insert into Products (name, price, category) values ($1, $2, $3)",
				req.Name, req.Price, req.Category)
			if err != nil {
				c.JSON(500, gin.H{"error": "database error"})
				return

			}
		} else {
			_, err := db.DB.Exec("insert into Products (name, price) values ($1, $2)",
				req.Name, req.Price)
			if err != nil {
				c.JSON(500, gin.H{"error": "database error"})
				return
			}
		}

	}
}
