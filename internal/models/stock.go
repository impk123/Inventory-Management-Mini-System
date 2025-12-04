package models

import (
	"time"
)

type Stock struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProductID   uint      `gorm:"not null" json:"product_id"`
	Type        string    `gorm:"not null" json:"type"` // IN, OUT, ADJUST
	Quantity    int       `gorm:"not null" json:"quantity"`
	OldQuantity int       `json:"old_quantity"`
	NewQuantity int       `json:"new_quantity"`
	Reference   string    `json:"reference"` // PO number, Sales order, etc.
	Notes       string    `json:"notes"`
	CreatedBy   uint      `json:"created_by"`
	UpdatedBy   uint      `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Product     Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	User        User      `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
}

type StockUpdateRequest struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Type      string `json:"type" binding:"required,oneof=IN OUT ADJUST"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
	Notes     string `json:"notes"`
}
