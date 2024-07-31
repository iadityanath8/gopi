package models

import "gorm.io/gorm"

type RequestBody struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type RequestBodyDel struct {
	ProductID uint `json:"product_id"`
}

type Product struct {
	gorm.Model
	Name  string `json:"name"`
	Price uint   `json:"price"`
	Stock uint   `json:"stock"`
}

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex" json:"username"`
	Password string `json:"password"`
}

// READ THIS CAREFULLY IN HERE OK
type CartItem struct {
	gorm.Model
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
	CartID    uint `json:"cart_id"`
}

type Cart struct {
	gorm.Model
	UserID uint       `json:"user_id"`
	Items  []CartItem `json:"items" gorm:"foreignkey:CartID"`
}

type CartItemResponse struct {
	ProductId uint    `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  uint    `json:"quantity"`
}
