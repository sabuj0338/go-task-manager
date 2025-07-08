package models

import "time"

// Role is a collection of permissions.
type Role struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
