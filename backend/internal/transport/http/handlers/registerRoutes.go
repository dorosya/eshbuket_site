package handlers

import (
	"database/sql"
	store "eshbuket/internal/repository/Store"
	"eshbuket/internal/repository/postgres"
	"eshbuket/internal/service/auth"
	"eshbuket/internal/service/order"
	"eshbuket/internal/service/product"
	"eshbuket/internal/service/productimage"
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
	ProductHandler *ProductHandler
	OrderHandler   *OrderHandler
	LoginHandler   *LoginHandler
	ImageHandler   *ImageHandler
	AuthService    *auth.AuthService
}

func CreateHandlers(db *sql.DB) *Handlers {
	OrderRepo := postgres.NewOrderRepository(db)
	ProductRepo := postgres.NewProductRepository(db)
	ImageRepo := postgres.NewProductImageRepository(db)

	ordersService := order.NewOrderService(OrderRepo)
	productService := product.NewProductService(ProductRepo)
	imageService := productimage.NewService(ImageRepo, "")
	authService := auth.NewAuthService(store.NewSessionStore())

	productHandler := NewProductHandler(*productService)
	orderHandler := NewOrderHandler(*ordersService)
	loginHandler := NewLoginHandler(*authService)
	imageHandler := NewImageHandler(imageService)

	return &Handlers{
		ProductHandler: productHandler,
		OrderHandler:   orderHandler,
		LoginHandler:   loginHandler,
		ImageHandler:   imageHandler,
		AuthService:    authService,
	}
}
func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	h := CreateHandlers(db)
	api := router.Group("/api")
	{
		api.POST("/login", h.LoginHandler.LoginHandler)
		api.POST("/orders", h.OrderHandler.OrdersHandler)
		api.POST("/products", middleware.AuthMiddleware(h.AuthService), h.ProductHandler.ProductsHandler)
		api.GET("/products", h.ProductHandler.ProductsHandler)
		api.GET("/products/:id/image", h.ImageHandler.GetProductImage)
		api.POST("/products/:id/image", middleware.AuthMiddleware(h.AuthService), h.ImageHandler.UploadProductImage)
	}
}
