package repository

import (
	"github.com/sabuj0338/go-task-manager/internal/models"
	"github.com/sabuj0338/go-task-manager/pkg/database"
)

func CreateTask(task *models.Task) error {
	query := `INSERT INTO tasks (title, description, completed, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())`
	_, err := database.DB.Exec(query, task.Title, task.Description, task.Completed, task.UserID)
	return err
}

func GetTasks(userID uint) ([]models.Task, error) {
	rows, err := database.DB.Query(`SELECT id, title, description, completed, user_id, created_at, updated_at FROM tasks WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.UserID, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func GetTaskByID(id uint, userID uint) (*models.Task, error) {
	row := database.DB.QueryRow(`SELECT id, title, description, completed, user_id, created_at, updated_at FROM tasks WHERE id = ? AND user_id = ?`, id, userID)
	var task models.Task
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.UserID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func UpdateTask(id uint, userID uint, title string, description string, completed bool) error {
	query := `UPDATE tasks SET title = ?, description = ?, completed = ?, updated_at = NOW() WHERE id = ? AND user_id = ?`
	_, err := database.DB.Exec(query, title, description, completed, id, userID)
	return err
}

func DeleteTask(id uint, userID uint) error {
	_, err := database.DB.Exec(`DELETE FROM tasks WHERE id = ? AND user_id = ?`, id, userID)
	return err
}
