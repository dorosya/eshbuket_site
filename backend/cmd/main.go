package main

import (
	"eshbuket/internal/repository/postgres"
	"eshbuket/internal/service/order"
	"eshbuket/internal/transport/http/dto"
	"eshbuket/internal/transport/http/handlers"
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
	db := postgres.Connect(dbHost, dbPort, dbUser, dbPassword, dbName)

	OrderRepo := postgres.NewOrderRepository(db)
	OrderService := order.NewOrderService(OrderRepo)
	OrderHandler := handlers.NewOrderHandler(OrderService)

	err := postgres.InitSchema()
	if err != nil {
		log.Println("Не получилось создать таблицы для ДБ")
	}

	env := os.Getenv("APP_ENV")
	router := dto.NewRouter(dto.Config{Env: env})
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
