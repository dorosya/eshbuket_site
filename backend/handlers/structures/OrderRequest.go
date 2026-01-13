package handlers

type OrderProductRequest struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required"`
}

type OrderRequest struct {
	ContactData string                `json:"contact_data" binding:"required"`
	Comment     string                `json:"comment"`
	Products    []OrderProductRequest `json:"products" binding:"required,dive"`
}
