package service_interfaces

import (
	"github.com/google/uuid"
	"lab3/internal/models"
)

type IUserService interface {
	GetUserByID(id uuid.UUID) (*models.User, error)
	Register(user *models.User, password string) (*models.User, error)
	Login(email, password string) (*models.User, error)
	Update(id uuid.UUID, name string, surname string, email string, address string, phoneNumber string, password string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}
