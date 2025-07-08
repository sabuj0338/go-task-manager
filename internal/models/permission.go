package models

import "time"

// Permission represents an action a user can perform.
type Permission struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"` // e.g., "Read Users"
	Code      string    `json:"code"` // e.g., "users:read"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
