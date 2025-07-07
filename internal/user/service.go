package user

import (
	"errors"

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

	hashed, _ := utils.HashPassword(dto.Password)

	user := &models.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: hashed,
		Role:     "user",
		Verified: false,
	}

	return repository.CreateUser(user)
}

func UpdateUserById(id int, dto UpdateUserDTO) error {
	return repository.UpdateById(id, dto.Name, dto.Email)
}

func DeleteUserById(id int) error {
	return repository.RemoveById(id)
}
