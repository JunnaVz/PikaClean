package repository_interfaces

import (
	"github.com/google/uuid"
	"lab3/internal/models"
)

type ITaskRepository interface {
	Create(task *models.Task) (*models.Task, error)
	Delete(id uuid.UUID) error
	Update(task *models.Task) (*models.Task, error)
	GetTaskByID(id uuid.UUID) (*models.Task, error)
	GetAllTasks() ([]models.Task, error)
	GetTasksInCategory(category int) ([]models.Task, error)
	GetTaskByName(name string) (*models.Task, error)
}
