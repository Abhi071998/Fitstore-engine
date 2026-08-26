// Package orders reads the customer-facing order tables (owned/migrated by
// fitstore-core via Prisma) so the admin backend can review pending orders.
// These models are read-mapped onto existing tables and must never be
// passed to AutoMigrate here — schema ownership stays with fitstore-core.
package orders

import (
	"time"

	"tshirt-store/internal/models"
)

// Order mirrors the Prisma `orders` model.
// NOTE: AdminComment maps to an `admin_comment` (nullable text) column that
// does not exist in the schema pasted from fitstore-core yet — it needs to
// be added there via a Prisma migration before reject-with-comment works.
type Order struct {
	ID               uint64     `gorm:"column:id;primaryKey" json:"id"`
	CustUserID       uint64     `gorm:"column:cust_user_id" json:"cust_user_id"`
	Status           string     `gorm:"column:status" json:"status"`
	ShippingName     *string    `gorm:"column:shipping_name" json:"shipping_name"`
	ShippingEmail    *string    `gorm:"column:shipping_email" json:"shipping_email"`
	ShippingAddress  *string    `gorm:"column:shipping_address" json:"shipping_address"`
	ShippingCity     *string    `gorm:"column:shipping_city" json:"shipping_city"`
	ShippingState    *string    `gorm:"column:shipping_state" json:"shipping_state"`
	ShippingPincode  *string    `gorm:"column:shipping_pincode" json:"shipping_pincode"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DecidedAt        *time.Time `gorm:"column:decided_at" json:"decided_at"`
	AdminComment     *string    `gorm:"column:admin_comment" json:"admin_comment"`

	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

func (Order) TableName() string { return "orders" }

// OrderItem mirrors the Prisma `order_items` model. unit_price is a
// snapshot of the product's selling price at submit time.
type OrderItem struct {
	ID            uint64    `gorm:"column:id;primaryKey" json:"id"`
	OrderID       uint64    `gorm:"column:order_id" json:"order_id"`
	ProductSizeID uint64    `gorm:"column:product_size_id" json:"product_size_id"`
	Quantity      uint64    `gorm:"column:quantity" json:"quantity"`
	UnitPrice     float64   `gorm:"column:unit_price" json:"unit_price"`
	Status        string    `gorm:"column:status" json:"status"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`

	// ProductSize (and its Product) live in the shared catalog tables this
	// service already owns — joined in for display, not part of fitstore-core's schema.
	ProductSize *models.ProductSize `gorm:"foreignKey:ProductSizeID;references:ID" json:"product_size,omitempty"`
}

func (OrderItem) TableName() string { return "order_items" }
