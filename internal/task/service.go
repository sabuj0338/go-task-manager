package task

import (
	"math"

	"github.com/sabuj0338/go-task-manager/internal/models"
	"github.com/sabuj0338/go-task-manager/internal/task/repository"
	"github.com/sabuj0338/go-task-manager/pkg/response"
)

func Create(userID uint, dto CreateTaskDTO) error {
	task := &models.Task{
		Title:       dto.Title,
		Description: dto.Description,
		Completed:   false,
		UserID:      userID,
	}
	return repository.CreateTask(task)
}

// GetAll retrieves a paginated list of tasks for a user.
func GetAll(userID uint, page int, limit int, title string, sort string) ([]models.Task, *response.PaginationMeta, error) {
	// Set default values for pagination if they are not provided
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10 // A sensible default limit
	}

	// Calculate offset for the database query
	offset := (page - 1) * limit

	var sortBy = "created_at DESC"

	switch sort {
	case "title":
		sortBy = "title ASC"
	case "-title":
		sortBy = "title DESC"
	case "description":
		sortBy = "description ASC"
	case "-description":
		sortBy = "description DESC"
	case "completed":
		sortBy = "completed ASC"
	case "-completed":
		sortBy = "completed DESC"
	case "created_at":
		sortBy = "created_at ASC"
	case "-created_at":
		sortBy = "created_at DESC"
	}

	tasks, total, err := repository.GetTasks(userID, limit, offset, title, sortBy)
	if err != nil {
		return nil, nil, err
	}

	meta := &response.PaginationMeta{
		TotalItems:   total,
		TotalPages:   int(math.Ceil(float64(total) / float64(limit))),
		CurrentPage:  page,
		ItemsPerPage: limit,
	}

	return tasks, meta, nil
}

func GetByID(id uint, userID uint) (*models.Task, error) {
	return repository.GetTaskByID(id, userID)
}

func Update(id uint, userID uint, dto UpdateTaskDTO) error {
	// The current repository `UpdateTask` uses gorm.Updates(struct{...}), which
	// has a known issue where it ignores zero-value fields (like `false` for a boolean).
	// A more robust approach is to fetch the existing record and then apply changes.
	task, err := repository.GetTaskByID(id, userID)
	if err != nil {
		return err // Handles "not found" and other database errors.
	}

	// Only update the 'completed' status if it was explicitly provided in the request.
	if dto.Completed != nil {
		task.Completed = *dto.Completed
	}
	return repository.UpdateTask(id, userID, dto.Title, dto.Description, task.Completed)
}

func Delete(id uint, userID uint) error {
	return repository.DeleteTask(id, userID)
}
