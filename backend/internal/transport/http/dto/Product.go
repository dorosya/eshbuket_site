package models

type ProductResponse struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Category string `json:"category,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

type ProductRequest struct {
	Name     string `json:"name" binding:"required"`
	Price    string `json:"price" binding:"required"`
	Category string `json:"category"`
}
