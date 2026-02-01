package handlers

import (
	"eshbuket/internal/service/product"
	"eshbuket/internal/transport/http/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service product.ProductService
}

func NewProductHandler(PS product.ProductService) *ProductHandler {
	return &ProductHandler{PS}
}

func (Handler *ProductHandler) ProductsHandler(c *gin.Context) {
	switch c.Request.Method {
	// Get /api/products - select всех товаров из БД по категории
	case http.MethodGet:
		products, err := Handler.service.ProductGetService(c.Request.Context(), c.Query("category"))

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, products)
	// Post /api/products - загрузка в БД товаров (только для админки, пропускает через authMiddleware перед обработкой )*/
	case http.MethodPost:
		var req dto.ProductRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err := Handler.service.ProductPostService(c.Request.Context(), req)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, gin.H{
			"name":     req.Name,
			"price":    req.Price,
			"category": req.Category,
		})
	}
}
