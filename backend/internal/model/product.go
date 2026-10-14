package model

type Product struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Category          string `json:"category"`
	Description       string `json:"description"`
	Price             int64  `json:"price"`
	ImageURL          string `json:"imageUrl"`
	StockQuantity     int    `json:"stockQuantity"`
	LowStockThreshold int    `json:"lowStockThreshold"`
	IsActive          bool   `json:"isActive"`
	Available         bool   `json:"available"`
}
