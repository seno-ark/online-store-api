package entity

import (
	"time"
)

type Payment struct {
	ID              string     `json:"id"`
	OrderID         string     `json:"order_id"`
	PaymentMethod   string     `json:"payment_method"`
	PaymentProvider string     `json:"payment_provider"`
	BillAmount      int64      `json:"bill_amount"`
	PaidAmount      int64      `json:"paid_amount"`
	Status          string     `json:"status"`
	TransactionID   string     `json:"transaction_id"`
	PaidAt          *time.Time `json:"paid_at"`
	Log             *string    `json:"log"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type PaymentWebHookTransactionDetails struct {
	OrderID     string `json:"order_id" validate:"required,uuid"`
	GrossAmount int64  `json:"gross_amount" validate:"required"`
}

type PaymentWebHookUserDetails struct {
	FullName string `json:"full_name" validate:"required"`
	Email    int64  `json:"email"`
}

type InPaymentWebHook struct {
	TransactionID                    string                           `json:"transaction_id" validate:"required"`
	PaymentAmount                    int64                            `json:"payment_amount" validate:"required,min=1"`
	Status                           string                           `json:"status" validate:"required"`
	PaymentWebHookTransactionDetails PaymentWebHookTransactionDetails `json:"transaction_details"`
	PaymentWebHookUserDetails        PaymentWebHookUserDetails        `json:"user_details"`
}
