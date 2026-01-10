package handlers

import "github.com/gin-gonic/gin"

// /api/products - endpoint со списком всех товаров. Парсится из bd и загружает туда же
func ProductsHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		// Get /api/products - select всех товаров из БД по категории
	} else if c.Request.Method == "POST" {
		// Post /api/products - загрузка в БД товаров (только для админки, пропускает через authMiddleware перед обработкой )*/
	}
}
