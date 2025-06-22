package models

import (
	"time"
)

type User struct {
	ID          int       `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	Name        *string   `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserWriteModel struct {
	PhoneNumber string  `json:"phone_number"`
	Name        *string `json:"name,omitempty"`
}