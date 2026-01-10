package main

import (
	"eshbuket/handlers"
	"eshbuket/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/login", handlers.LoginHandler)
	router.POST("/orders", handlers.OrdersHandler)
	router.POST("/products", middleware.AuthMiddleware, handlers.ProductsHandler)
	router.GET("/products", handlers.ProductsHandler)
	router.Run()
}
