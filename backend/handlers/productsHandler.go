package handlers

import (
	"database/sql"
	db "eshbuket/database"
	handlers "eshbuket/handlers/structures"

	"github.com/gin-gonic/gin"
)

// /api/products - endpoint со списком всех товаров. Парсится из bd и загружает туда же
func ProductsHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		// Get /api/products - select всех товаров из БД по категории
		category := c.Query("category")

		var rows *sql.Rows
		var err error

		if category != "" {
			//Задана определенная категория в Query запроса
			rows, err = db.DB.Query("SELECT id, name, price, category FROM products WHERE category = $1", category)
		} else {
			//Нет категории
			rows, err = db.DB.Query("SELECT id, name, price, category FROM products")
		}

		if err != nil {
			c.JSON(500, gin.H{"error": "db error"})
			return
		}
		defer rows.Close()

		products := []handlers.Product{}
		for rows.Next() {
			var p handlers.Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category); err != nil {
				c.JSON(500, gin.H{"error": "scan error"})
				return
			}
			products = append(products, p)
		}

		c.JSON(200, products)
	} else if c.Request.Method == "POST" {
		// Post /api/products - загрузка в БД товаров (только для админки, пропускает через authMiddleware перед обработкой )*/
		var req handlers.Product
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
