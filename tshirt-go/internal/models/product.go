package models

import (
	"time"

	"gorm.io/gorm"
)

// ProductSize handles inventory breakdown per apparel size (e.g., S: 10, M: 25)
type ProductSize struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProductID uint           `gorm:"not null;index" json:"product_id"`
	Size      string         `gorm:"type:varchar(10);not null" json:"size"` // e.g., "S", "M", "L", "XL", "XXL"
	Stock     int            `gorm:"not null;default:0" json:"stock"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Product represents the complete apparel product record
type Product struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"not null;type:varchar(255)" json:"name"`
	Brand       string `gorm:"type:varchar(100);default:'FITstore'" json:"brand"`
	Description string `gorm:"type:text" json:"description"`
	ProductCode string `gorm:"unique;not null;type:varchar(100)" json:"product_code"` // e.g., "STYLE-43054612"
	SKU         string `gorm:"unique;not null;type:varchar(100)" json:"sku"`          // Stock Keeping Unit

	// Pricing attributes
	MRP          float64 `gorm:"not null;type:decimal(10,2)" json:"mrp"`           // Original list price
	SellingPrice float64 `gorm:"not null;type:decimal(10,2)" json:"selling_price"` // Discounted offer price
	Discount     int     `gorm:"default:0" json:"discount_percentage"`             // e.g., 20% off

	// Media & Visuals (Stored as JSON array strings or postgres text array)
	Images string `gorm:"type:text" json:"images"` // Comma-separated or JSON array of image URLs

	// Metadata & Specifications (Key-Value structured metadata)
	// Example: {"Fabric": "100% Cotton", "Fit": "Slim Fit", "Pattern": "Striped", "Sleeve": "Long Sleeve"}
	Specifications string `gorm:"type:jsonb" json:"specifications,omitempty"`

	// Category Relationship
	CategoryID uint      `gorm:"not null;index" json:"category_id"`
	Category   *Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"category,omitempty"`

	// Inventory Breakdown across Sizes
	Sizes []ProductSize `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;" json:"sizes,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
