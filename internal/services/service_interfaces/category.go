package service_interfaces

import "lab3/internal/models"

type ICategoryService interface {
	GetAll() ([]models.Category, error)
	GetTasksInCategory(id int) ([]models.Task, error)
	GetByID(id int) (*models.Category, error)
	Create(name string) (*models.Category, error)
	Update(category *models.Category) (*models.Category, error)
	Delete(id int) error
}
