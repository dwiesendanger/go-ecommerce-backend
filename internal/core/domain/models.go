package domain

import (
	"gorm.io/gorm"
)

// USER
type User struct {
	gorm.Model
	Email    string  `gorm:"uniqueIndex;not null" json:"email"`
	Password string  `json:"-"`
	Orders   []Order `json:"orders,omitempty"` // One-to-Many
	Cart     Cart    `json:"cart,omitempty"`   // One-to-One
}

// PRODUCT
type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	SKU         string  `gorm:"uniqueIndex" json:"sku"` // Stock Keeping Unit
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// CART
type Cart struct {
	gorm.Model
	UserID uint       `gorm:"uniqueIndex"`
	Items  []CartItem `json:"items"`
}

type CartItem struct {
	gorm.Model
	CartID    uint    `gorm:"index"`
	ProductID uint    `gorm:"index"`
	Product   Product `json:"product"`
	Quantity  int     `json:"quantity"`
}

// ORDER
type Order struct {
	gorm.Model
	UserID      uint        `gorm:"index"`
	OrderNumber string      `gorm:"uniqueIndex" json:"order_number"`
	Status      string      `gorm:"index;default:'pending'" json:"status"` // pending, paid, shipped, cancelled
	TotalAmount float64     `json:"total_amount"`
	Items       []OrderItem `json:"items"`
}

// ORDER ITEM
type OrderItem struct {
	gorm.Model
	OrderID   uint    `gorm:"index"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}
