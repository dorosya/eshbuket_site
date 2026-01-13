package handlers

import (
	"eshbuket/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/login", LoginHandler)
		api.POST("/orders", OrdersHandler)
		api.POST("/products", middleware.AuthMiddleware, ProductsHandler)
		api.GET("/products", ProductsHandler)
	}
}
