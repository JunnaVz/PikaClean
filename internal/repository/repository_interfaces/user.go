package repository_interfaces

import (
	"PikaClean/internal/models"
	"github.com/google/uuid"
)

type IUserRepository interface {
	Create(user *models.User) (*models.User, error)
	Update(user *models.User) (*models.User, error)
	Delete(id uuid.UUID) error
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}
