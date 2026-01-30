package handlers

import (
	"eshbuket/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/login", LoginHandler)
		api.POST("/orders", OrdersHandler)
		api.POST("/products", middleware.AuthMiddleware, ProductsHandler)
		api.GET("/products", ProductsHandler)
		api.GET("/products/:id/image", GetProductImage)
		api.POST("/products/:id/image", middleware.AuthMiddleware, UploadProductImage)
	}
}
