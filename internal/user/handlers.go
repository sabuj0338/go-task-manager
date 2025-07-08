package user

import (
	"strconv"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sabuj0338/go-task-manager/pkg/response"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("strong_password", strongPasswordValidation)
}

func strongPasswordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}

func GetUsers(c *fiber.Ctx) error {
	// Parse pagination query parameters, providing sensible defaults.
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	users, meta, err := GetAllUsers(page, limit)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get users", err.Error())
	}
	return response.SuccessWithMeta(c, fiber.StatusOK, "Users retrieved", users, meta)
}

func GetUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user id")
	}
	user, err := GetUserByID(id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "User not found", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "User retrieved", user)
}

func CreateUser(c *fiber.Ctx) error {
	var dto CreateUserDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}
	if err := validate.Struct(dto); err != nil {
		return response.ValidationErrorResponse(c, err)
	}

	if err := CreateNewUser(dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "New user created successful", nil)
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user id")
	}

	var dto UpdateUserDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}
	if err := validate.Struct(dto); err != nil {
		return response.ValidationErrorResponse(c, err)
	}

	if err := UpdateUserById(id, dto); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update user", err.Error())
	}
	return response.Success(c, fiber.StatusOK, "User updated successfully", nil)
}

func DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user id")
	}
	if err := DeleteUserById(id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete user", err.Error())
	}
	return response.Success(c, fiber.StatusOK, "User deleted successfully", nil)
}
