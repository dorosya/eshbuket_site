package main

import (
	DB "eshbuket/database"
	"eshbuket/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	DB.Connect()
	router := gin.Default()
	handlers.RegisterRoutes(router)
	router.Run()
}
