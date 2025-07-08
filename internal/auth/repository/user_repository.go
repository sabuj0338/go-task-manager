package repository

import (
	"github.com/sabuj0338/go-task-manager/internal/models"
	"github.com/sabuj0338/go-task-manager/pkg/database"
)

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	// Use GORM to find the user. GORM handles 'record not found' errors.
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	// Fetch associated roles for the user. It's okay if the user has no roles.
	roles, err := GetUserRoles(user.ID)
	if err == nil {
		user.Roles = roles
	}

	return &user, nil
}

func GetUserByID(id int) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	// Fetch associated roles for the user.
	roles, err := GetUserRoles(user.ID)
	if err == nil {
		user.Roles = roles
	}

	return &user, nil
}

func CheckEmailExists(userID uint, email string) bool {
	var count int64
	database.DB.Model(&models.User{}).Where("email = ? AND id != ?", email, userID).Count(&count)
	return count > 0
}

func CreateUser(user *models.User) error {
	// GORM handles the insert and populates the user's ID.
	return database.DB.Create(user).Error
}

func UpdateUserMFASecret(userID uint, secret string) error {
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"mfa_secret":  secret,
		"mfa_enabled": true,
	}).Error
}

func GetUserMFA(userID uint) (string, bool, error) {
	var user models.User
	err := database.DB.Select("mfa_secret", "mfa_enabled").First(&user, userID).Error
	if err != nil {
		return "", false, err
	}
	var secretVal string
	if user.MFASecret != nil {
		secretVal = *user.MFASecret
	}
	return secretVal, user.MFAEnabled, nil
}

// DisableUserMFA clears the MFA secret and disables the MFA flag for a user.
func DisableUserMFA(userID uint) error {
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"mfa_secret":  nil, // Sets the column to NULL
		"mfa_enabled": false,
	}).Error
}

func UpdateUserPassword(userID uint, newPasswordHash string) error {
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Update("password", newPasswordHash).Error
}
