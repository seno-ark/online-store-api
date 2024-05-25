package entity

import (
	"time"
)

type CartItem struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	ProductID string    `json:"product_id" db:"product_id"`
	Notes     string    `json:"notes" db:"notes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type InCreateCartItem struct {
	ProductID string `json:"product_id" validate:"required,uuid" example:"02a1a6a3-1c9c-4f46-ae18-162e2b0d7a9a"`
	Notes     string `json:"notes" example:"Nggak pakai sambal"`
}

type InGetListCartItem struct {
	Limit  int
	Offset int
}

/*
carts.id, carts.user_id, carts.product_id, carts.notes,
products.category_id, products.name AS product_name, products.description AS product_description,
products.price, products.stock,
categories.name AS category_name
*/
type OutGetCartItem struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	ProductID string    `json:"product_id" db:"product_id"`
	Notes     string    `json:"notes" db:"notes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	CategoryID         string `json:"category_id" db:"category_id"`
	ProductName        string `json:"product_name" db:"product_name"`
	ProductDescription string `json:"product_description" db:"product_description"`
	ProductPrice       int64  `json:"product_price" db:"product_price"`
	ProductStock       int64  `json:"product_stock" db:"product_stock"`

	CategoryName string `json:"category_name" db:"category_name"`
}
