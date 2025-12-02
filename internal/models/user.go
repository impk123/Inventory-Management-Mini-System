package models

import (
	"time"
)

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Email       string    `gorm:"uniqueIndex;not null" json:"email"`
	Password    string    `json:"-"`
	GoogleID    *string   `gorm:"uniqueIndex" json:"google_id,omitempty"`
	Name        string    `json:"name"`
	Role        string    `gorm:"default:user" json:"role"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	LastLoginAt time.Time `json:"last_login_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
