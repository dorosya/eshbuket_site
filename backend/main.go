package main

import (
	db "eshbuket/database"
	"eshbuket/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env файл не найден")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db.Connect(dbHost, dbPort, dbUser, dbPassword, dbName)
	db.InitSchema(db.DB)

	router := gin.Default()
	handlers.RegisterRoutes(router)
	router.Run()
}
