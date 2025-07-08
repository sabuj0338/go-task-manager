package models

import "time"

type User struct {
	// ID        uint      `json:"id"`
	// Name      string    `json:"name"`
	// Email     string    `json:"email"`
	// Phone     string    `json:"phone"`
	// Password  string    `json:"-"`
	// Role      string    `json:"role"`
	// Verified  bool      `json:"verified"`
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone,omitempty"`
	Password      string    `json:"-"`
	Roles         []Role    `json:"roles,omitempty" gorm:"many2many:user_roles;"`
	EmailVerified bool      `json:"email_verified"`
	MFASecret     *string   `json:"-"`
	MFAEnabled    bool      `json:"mfa_enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
