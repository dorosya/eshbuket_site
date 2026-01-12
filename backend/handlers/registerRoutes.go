package handlers

import (
	"eshbuket/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.POST("/login", LoginHandler)
	router.POST("/orders", OrdersHandler)
	router.POST("/products", middleware.AuthMiddleware, ProductsHandler)
	router.GET("/products", ProductsHandler)
}
