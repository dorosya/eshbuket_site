package main

import (
	db "eshbuket/database"
	"eshbuket/handlers"
	"log"
	"os"

	"github.com/gin-contrib/cors"
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
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	router.Use(cors.New(config))
	// gin.SetMode(gin.TestMode) //Тестмод для проходки тестов на Post /api/products (для отключения авторизации)
	handlers.RegisterRoutes(router)

	// Подгрузка фронтенда через
	// router.Static("/static", "./frontend")
	// router.GET("/", func(c *gin.Context) {
	// 	c.File("./frontend/index.html")
	// })
	router.Run()
}
