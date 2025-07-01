package repository

import (
	"database/sql"

	"github.com/sabuj0338/go-task-manager/internal/models"
	"github.com/sabuj0338/go-task-manager/pkg/database"
)

func FindAll() ([]models.User, error) {
	rows, err := database.DB.Query("SELECT id, name, email, role, verified, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Verified, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func FindById(id int) (*models.User, error) {
	row := database.DB.QueryRow("SELECT id, name, email, role, verified, created_at, updated_at FROM users WHERE id = ?", id)
	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Verified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func UpdateById(id int, name string, email string) error {
	query := "UPDATE users SET name = ?, email = ?, updated_at = NOW() WHERE id = ?"
	_, err := database.DB.Exec(query, name, email, id)
	return err
}

func RemoveById(id int) error {
	_, err := database.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}
