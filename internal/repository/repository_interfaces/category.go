package repository_interfaces

import "lab3/internal/models"

type ICategoryRepository interface {
	GetAll() ([]models.Category, error)
	GetByID(id int) (*models.Category, error)
	Create(category *models.Category) (*models.Category, error)
	Update(category *models.Category) (*models.Category, error)
	Delete(id int) error
}
