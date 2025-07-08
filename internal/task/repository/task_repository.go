package repository

import (
	"github.com/sabuj0338/go-task-manager/internal/models"
	"github.com/sabuj0338/go-task-manager/pkg/database"
)

func CreateTask(task *models.Task) error {
	return database.DB.Create(task).Error
}

func GetTasks(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := database.DB.Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}

func GetTaskByID(id uint, userID uint) (*models.Task, error) {
	var task models.Task
	err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error
	return &task, err
}

func UpdateTask(id uint, userID uint, title string, description string, completed bool) error {
	return database.DB.Model(&models.Task{}).Where("id = ? AND user_id = ?", id, userID).
		Updates(models.Task{Title: title, Description: description, Completed: completed}).Error
}

func DeleteTask(id uint, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Task{}).Error
}
