package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the administrative account schema in the system
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"unique;not null;type:varchar(255)" json:"email"`
	Name      string         `gorm:"not null;type:varchar(255)" json:"name"`
	Password  string         `gorm:"not null;type:varchar(255)" json:"-"` // "-" hides it from JSON responses
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Product represents the t-shirt inventory item schema
type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;type:varchar(255)" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"not null;type:decimal(10,2)" json:"price"`
	SKU         string         `gorm:"unique;not null;type:varchar(100)" json:"sku"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// =========================================================================
// 📥 REQUEST DTOs (Data Transfer Objects)
// =========================================================================

// RegisterDTO maps incoming sign-up payloads
type RegisterDTO struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// LoginDTO maps incoming sign-in credentials
type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
