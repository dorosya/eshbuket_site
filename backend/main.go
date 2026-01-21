package main

import (
	db "eshbuket/database"
	"eshbuket/handlers"
	"eshbuket/models"
	"log"
	"os"

	"github.com/gin-contrib/cors"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db.Connect(dbHost, dbPort, dbUser, dbPassword, dbName)
	err := db.InitSchema(db.DB)
	if err != nil {
		log.Println("Не получилось создать таблицы для ДБ")
	}

	env := os.Getenv("APP_ENV")
	router := models.NewRouter(models.Config{env})
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))
	handlers.RegisterRoutes(router)

	// Подгрузка фронтенда через
	// router.Static("/static", "./frontend")
	// router.GET("/", func(c *gin.Context) {
	// 	c.File("./frontend/index.html")
	// })
	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}
