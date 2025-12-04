package handlers

import "github.com/impk123/Inventory-Management-Mini-System/internal/models"

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

// ProductListResponse represents the response for getting all products
type ProductListResponse struct {
	Products []models.Product `json:"products"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
	Pages    int              `json:"pages"`
}

// StockHistoryResponse represents the response for getting stock history
type StockHistoryResponse struct {
	Stock     []models.Stock `json:"stock"`
	Movements []models.Stock `json:"movements"`
}

// CurrentStockResponse represents the response for getting current stock levels
type CurrentStockResponse struct {
	Products      []models.Product `json:"products"`
	Total         int64            `json:"total"`
	Page          int              `json:"page"`
	Limit         int              `json:"limit"`
	LowStock      []models.Product `json:"low_stock"`
	LowStockCount int              `json:"low_stock_count"`
}
