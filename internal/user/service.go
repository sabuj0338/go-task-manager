package user

import (
	"errors"

	auth_repo "github.com/sabuj0338/go-task-manager/internal/auth/repository"
	"github.com/sabuj0338/go-task-manager/internal/models"
	"github.com/sabuj0338/go-task-manager/internal/user/repository"
	"github.com/sabuj0338/go-task-manager/internal/utils"
)

func GetAllUsers() (interface{}, error) {
	return repository.FindAll()
}

func GetUserByID(id int) (interface{}, error) {
	user, err := repository.FindById(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func CreateNewUser(dto CreateUserDTO) error {
	existingUser, _ := repository.GetUserByEmail(dto.Email)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	// Find the default "user" role. Admins might want to assign other roles,
	// so the DTO should ideally contain role information in the future.
	// For now, we'll assign the 'user' role by default.
	userRole, err := auth_repo.GetRoleByName("user")
	if err != nil {
		return errors.New("default user role not found")
	}

	hashed, _ := utils.HashPassword(dto.Password)

	user := &models.User{
		Name:          dto.Name,
		Email:         dto.Email,
		Password:      hashed,
		EmailVerified: false, // Users created by an admin are not verified by default
	}

	if err := repository.CreateUser(user); err != nil {
		return err
	}
	// Assign the role to the newly created user.
	return auth_repo.AssignRoleToUser(user.ID, userRole.ID)
}

func UpdateUserById(id int, dto UpdateUserDTO) error {
	return repository.UpdateById(id, dto.Name, dto.Email)
}

func DeleteUserById(id int) error {
	return repository.RemoveById(id)
}
