package task

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sabuj0338/go-task-manager/pkg/response"
)

var validate = validator.New()

func CreateTaskHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var dto CreateTaskDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}
	if err := validate.Struct(dto); err != nil {
		return response.ValidationErrorResponse(c, err)
	}
	if err := Create(userID, dto); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create task", err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Task created", nil)
}

// GetTasksHandler retrieves a paginated list of tasks.
func GetTasksHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	// Parse pagination query parameters, providing sensible defaults.
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	tasks, meta, err := GetAll(userID, page, limit)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get tasks", err.Error())
	}
	return response.SuccessWithMeta(c, fiber.StatusOK, "Tasks retrieved", tasks, meta)
}

func GetTaskHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid task id")
	}

	task, err := GetByID(uint(id), userID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Task not found", err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Task retrieved", task)
}

func UpdateTaskHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid task id")
	}

	var dto UpdateTaskDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	if err := validate.Struct(dto); err != nil {
		return response.ValidationErrorResponse(c, err)
	}

	if err := Update(uint(id), userID, dto); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update task", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Task updated", nil)
}

func DeleteTaskHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid task id")
	}

	if err := Delete(uint(id), userID); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete task", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Task deleted", nil)
}
