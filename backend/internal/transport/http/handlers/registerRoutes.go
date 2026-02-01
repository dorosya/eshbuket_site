package handlers

import (
	"database/sql"
	store "eshbuket/internal/repository/Store"
	"eshbuket/internal/repository/postgres"
	"eshbuket/internal/service/auth"
	"eshbuket/internal/service/order"
	"eshbuket/internal/service/product"
	"eshbuket/internal/transport/http/handlers"
	"eshbuket/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

// type handlers_init struct {
// 	service []interface{}
// 	handler []interface{}
// 	repo    []interface{}
// }

//TODO InitHandlers
//для этого каждый интерфейс должен реализовывать один общий интерфейс => переписывать все хендлеры...
//Слишком ммного заеба для всего ~5 хендлеров
//func InitHandlers(handler_list []interface{}) {}

type Handlers struct {
	ProductHandler *handlers.ProductHandler
	OrderHandler   *handlers.OrderHandler
	LoginHandler   *handlers.LoginHandler
	AuthService    *auth.AuthService
}

func CreateHandlers(db *sql.DB) *Handlers {
	OrderRepo := postgres.NewOrderRepository(db)
	ProductRepo := postgres.NewProductRepository(db)

	ordersService := order.NewOrderService(OrderRepo)
	productService := product.NewProductService(ProductRepo)
	authService := auth.NewAuthService(store.NewSessionStore())

	productHandler := handlers.NewProductHandler(*productService)
	orderHandler := handlers.NewOrderHandler(*ordersService)
	loginHandler := handlers.NewLoginHandler(*authService)

	return &Handlers{
		ProductHandler: productHandler,
		OrderHandler:   orderHandler,
		LoginHandler:   loginHandler,
		AuthService:    authService,
	}
}
func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	h := CreateHandlers(db)
	api := router.Group("/api")
	{
		api.POST("/login", h.LoginHandler)
		api.POST("/orders", h.OrderHandler)
		api.POST("/products", middleware.AuthMiddleware(h.AuthService), h.ProductHandler)
		api.GET("/products", h.ProductHandler)
		api.GET("/products/:id/image", GetProductImage)                                               // доделать
		api.POST("/products/:id/image", middleware.AuthMiddleware(h.AuthService), UploadProductImage) // доделать
	}
}
