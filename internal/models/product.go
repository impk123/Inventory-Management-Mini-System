package models

import (
	"time"
)

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SKU         string    `gorm:"uniqueIndex;not null" json:"sku"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	UnitPrice   float64   `gorm:"not null" json:"unit_price"`
	CostPrice   float64   `json:"cost_price"`
	Quantity    int       `gorm:"not null;default:0" json:"quantity"`
	MinQuantity int       `gorm:"default:10" json:"min_quantity"`
	MaxQuantity int       `json:"max_quantity"`
	Location    string    `json:"location"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedBy   uint      `json:"created_by"`
	UpdatedBy   uint      `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:CreatedBy" json:"-"`
}
