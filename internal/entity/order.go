package entity

import "time"

type Order struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	Status          string    `json:"status" db:"status"`
	OtherCost       int64     `json:"other_cost" db:"other_cost"`
	TotalCost       int64     `json:"total_cost" db:"total_cost"`
	ShipmentAddress string    `json:"shipment_address" db:"shipment_address"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type OrderItem struct {
	ID           string    `json:"id" db:"id"`
	OrderID      string    `json:"order_id" db:"order_id"`
	ProductID    string    `json:"product_id" db:"product_id"`
	Qty          int64     `json:"qty" db:"qty"`
	ProductPrice int64     `json:"product_price" db:"product_price"`
	Notes        string    `json:"notes" db:"notes"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type createOrderItem struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Qty       int64  `json:"qty" validate:"required,min=1"`
	Notes     string `json:"notes"`
}
type InCreateOrder struct {
	ShipmentAddress string            `json:"shipment_address"`
	Items           []createOrderItem `json:"items" validate:"required,max=5"`
}

type InGetListOrder struct {
	Limit  int
	Offset int
}

type InGetListOrderItem struct {
	Limit  int
	Offset int
}

type OutGetOrderItem struct {
	ID                 string    `json:"id" db:"id"`
	OrderID            string    `json:"order_id" db:"order_id"`
	ProductID          string    `json:"product_id" db:"product_id"`
	Qty                int64     `json:"qty" db:"qty"`
	ProductPrice       int64     `json:"product_price" db:"product_price"`
	Notes              string    `json:"notes" db:"notes"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	ProductName        string    `json:"product_name" db:"product_name"`
	ProductDescription string    `json:"product_description" db:"product_description"`
}
