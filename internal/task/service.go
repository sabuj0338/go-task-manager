package task

import (
	"github.com/sabuj0338/go-task-manager/internal/models"
	"github.com/sabuj0338/go-task-manager/internal/task/repository"
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

func GetAll(userID uint) ([]models.Task, error) {
	return repository.GetTasks(userID)
}

func GetByID(id uint, userID uint) (*models.Task, error) {
	return repository.GetTaskByID(id, userID)
}

func Update(id uint, userID uint, dto UpdateTaskDTO) error {
	return repository.UpdateTask(id, userID, dto.Title, dto.Description, *dto.Completed)
}

func Delete(id uint, userID uint) error {
	return repository.DeleteTask(id, userID)
}
