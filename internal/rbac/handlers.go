package rbac

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sabuj0338/go-task-manager/pkg/response"
)

var validate = validator.New()

func CreateRoleHandler(c *fiber.Ctx) error {
	var dto CreateRoleDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}
	if err := validate.Struct(dto); err != nil {
		return response.ValidationErrorResponse(c, err)
	}

	if err := CreateRole(dto); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create role", err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Role created successfully", nil)
}

func CreatePermissionHandler(c *fiber.Ctx) error {
	var dto CreatePermissionDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}
	if err := validate.Struct(dto); err != nil {
		return response.ValidationErrorResponse(c, err)
	}

	if err := CreatePermission(dto); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create permission", err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Permission created successfully", nil)
}

func AssignPermissionToRoleHandler(c *fiber.Ctx) error {
	var dto AssignPermissionToRoleDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}
	if err := validate.Struct(dto); err != nil {
		return response.ValidationErrorResponse(c, err)
	}

	if err := AssignPermissionToRole(dto); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to assign permission to role", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Permission assigned to role successfully", nil)
}

func AssignRoleToUserHandler(c *fiber.Ctx) error {
	var dto AssignRoleToUserDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}
	if err := validate.Struct(dto); err != nil {
		return response.ValidationErrorResponse(c, err)
	}

	if err := AssignRoleToUser(dto); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to assign role to user", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Role assigned to user successfully", nil)
}
