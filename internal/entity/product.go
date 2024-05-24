package entity

import (
	"time"
)

type Product struct {
	ID          string    `json:"id" db:"id"`
	CategoryID  string    `json:"category_id" db:"category_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       int64     `json:"price" db:"price"`
	Stock       int64     `json:"stock" db:"stock"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type InGetListProduct struct {
	CategoryID string
	Limit      int
	Offset     int
}
