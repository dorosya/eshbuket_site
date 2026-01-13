package handlers

type ProductRequest struct {
	Name     string `json:"name" binding:"required"`
	Price    string `json:"price" binding:"required"`
	Category string `json:"category"`
}
