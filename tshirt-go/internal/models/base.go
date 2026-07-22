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

// Category represents the top-level grouping (e.g., "T-shirt", "Shoes")
type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"unique;not null;type:varchar(255)" json:"name"`
	Products  []Product      `gorm:"foreignKey:CategoryID" json:"products,omitempty"` // Has-Many relation
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// CreateCategoryDTO maps incoming payloads for adding taxonomy
type CreateCategoryDTO struct {
	Name string `json:"name"`
}
