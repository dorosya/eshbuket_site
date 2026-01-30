package handlers

import (
	"eshbuket/internal/service/order"
	"eshbuket/internal/transport/http/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service order.OrderService
}

func NewOrderHandler(OS order.OrderService) *OrderHandler {
	return &OrderHandler{OS}
}

// OrdersHandler обрабатывает POST /api/orders
func (handler *OrderHandler) OrdersHandler(c *gin.Context) {
	var req dto.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if len(req.Products) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no products in order"})
		return
	}

	orderID, totalPrice, err := handler.service.OrderServiceFunc(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(201, gin.H{"order_id": orderID, "total_price": totalPrice})
}
