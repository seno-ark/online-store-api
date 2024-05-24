package entity

import (
	"time"
)

type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"-" db:"email"`
	Password  string    `json:"-" db:"password"`
	FullName  string    `json:"full_name" db:"full_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

type InUserRegister struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
}

type InUserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
