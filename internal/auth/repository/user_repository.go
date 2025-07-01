package repository

import (
	"database/sql"
	"errors"

	"github.com/sabuj0338/go-task-manager/internal/models"
	"github.com/sabuj0338/go-task-manager/pkg/database"
)

func GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, name, email, password, role, verified, created_at, updated_at FROM users WHERE email = ?`
	row := database.DB.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Verified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id int) (*models.User, error) {
	query := `SELECT id, name, email, password, role, verified, created_at, updated_at FROM users WHERE id = ?`
	row := database.DB.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Verified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *models.User) error {
	query := `INSERT INTO users (name, email, password, role, verified, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())`
	_, err := database.DB.Exec(query, user.Name, user.Email, user.Password, user.Role, user.Verified)
	return err
}

func UpdateUserMFASecret(userID uint, secret string) error {
	query := `UPDATE users SET mfa_secret = ?, mfa_enabled = true, updated_at = NOW() WHERE id = ?`
	_, err := database.DB.Exec(query, secret, userID)
	return err
}

func GetUserMFA(userID uint) (string, bool, error) {
	var secret string
	var enabled bool
	query := `SELECT mfa_secret, mfa_enabled FROM users WHERE id = ?`
	err := database.DB.QueryRow(query, userID).Scan(&secret, &enabled)
	return secret, enabled, err
}
