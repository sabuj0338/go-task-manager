package repository

import (
	"github.com/sabuj0338/go-task-manager/internal/models"
	"github.com/sabuj0338/go-task-manager/pkg/database"
)

func FindAll() ([]models.User, error) {
	var users []models.User
	// Use GORM to find all users and preload their roles
	err := database.DB.Preload("Roles").Find(&users).Error
	return users, err
}

func FindById(id int) (*models.User, error) {
	var user models.User
	// Use GORM to find a user by ID and preload their roles
	if err := database.DB.Preload("Roles").First(&user, id).Error; err != nil {
		return nil, err // GORM handles ErrRecordNotFound
	}
	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).Preload("Roles").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *models.User) error {
	// GORM handles the insert and populates the user's ID.
	return database.DB.Create(user).Error
}

func UpdateById(id int, name string, email string) error {
	// Use GORM's Updates method
	return database.DB.Model(&models.User{}).Where("id = ?", id).Updates(models.User{Name: name, Email: email}).Error
}

func RemoveById(id int) error {
	// Use GORM's Delete method
	return database.DB.Where("id = ?", id).Delete(&models.User{}).Error
}
