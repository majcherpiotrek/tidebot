package models

import "time"

type NotificationSubscription struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}