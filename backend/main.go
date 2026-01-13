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
		log.Println("Локальный .env файл не найден")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db.Connect(dbHost, dbPort, dbUser, dbPassword, dbName)
	err = db.InitSchema(db.DB)
	if err != nil {
		log.Println("Не получилось создать таблицы для ДБ")
	}

	router := gin.Default()
	gin.SetMode(gin.TestMode)
	handlers.RegisterRoutes(router)
	router.Run()
}
