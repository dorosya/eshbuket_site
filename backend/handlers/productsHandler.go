package handlers

import (
	"database/sql"
	db "eshbuket/database"
	"eshbuket/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// /api/products - endpoint со списком всех товаров. Парсится из bd и загружает туда же
func ProductsHandler(c *gin.Context) {
	switch c.Request.Method {
	// Get /api/products - select всех товаров из БД по категории
	case "GET":
		category := c.Query("category")

		var rows *sql.Rows
		var err error
		if category != "" {
			rows, err = db.DB.Query("SELECT id, name, price, category, image_path FROM products WHERE category = $1", category)
		} else {
			rows, err = db.DB.Query("SELECT id, name, price, category, image_path FROM products")
		}
		if err != nil {
			c.JSON(500, gin.H{"error": "db error"})
			return
		}
		defer rows.Close()

		products := []models.Product{}
		for rows.Next() {
			var p models.Product
			var imagePath sql.NullString
			if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category, &imagePath); err != nil {
				c.JSON(500, gin.H{"error": "scan error"})
				return
			}
			if imagePath.Valid && imagePath.String != "" {
				p.ImageURL = "/api/products/" + fmt.Sprint(p.ID) + "/image"
			}
			products = append(products, p)
		}

		if err := rows.Err(); err != nil {
			c.JSON(500, gin.H{"error": "rows error"})
			return
		}

		c.JSON(http.StatusOK, products)
	// Post /api/products - загрузка в БД товаров (только для админки, пропускает через authMiddleware перед обработкой )*/
	case "POST":
		var req models.ProductRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}
		if req.Name == "" || req.Price == "" {
			c.JSON(400, gin.H{"error": "name and price are required"})
			return
		}

		_, err := db.DB.Exec(
			"INSERT INTO Products (name, price, category) VALUES ($1, $2, $3)",
			req.Name, req.Price, req.Category,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": "database error"})
			return
		}

		c.JSON(201, gin.H{
			"name":     req.Name,
			"price":    req.Price,
			"category": req.Category,
		})
	}
}
