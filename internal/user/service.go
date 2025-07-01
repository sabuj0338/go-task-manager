package user

import (
	"errors"

	"github.com/sabuj0338/go-task-manager/internal/user/repository"
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

func UpdateUserById(id int, dto UpdateUserDTO) error {
	return repository.UpdateById(id, dto.Name, dto.Email)
}

func DeleteUserById(id int) error {
	return repository.RemoveById(id)
}
