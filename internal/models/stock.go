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
	Reason      string    `json:"reason"`
	CreatedBy   uint      `json:"created_by"`
	UpdatedBy   uint      `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Product     Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	User        User      `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
}
