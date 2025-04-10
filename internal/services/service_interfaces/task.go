package service_interfaces

import (
	"github.com/google/uuid"
	"lab3/internal/models"
)

type ITaskService interface {
	Create(name string, price float64, category int) (*models.Task, error)
	Update(taskID uuid.UUID, category int, name string, price float64) (*models.Task, error)
	Delete(taskID uuid.UUID) error
	GetAllTasks() ([]models.Task, error)
	GetTaskByID(id uuid.UUID) (*models.Task, error)
	GetTasksInCategory(category int) ([]models.Task, error)
	GetTaskByName(name string) (*models.Task, error)
}
