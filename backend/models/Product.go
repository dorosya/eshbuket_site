package handlers

type Product struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Category string `json:"category,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}
