package service_interfaces

import (
	"github.com/google/uuid"
	"lab3/internal/models"
)

type IWorkerService interface {
	Login(email, password string) (*models.Worker, error)
	Create(worker *models.Worker, password string) (*models.Worker, error)
	Delete(id uuid.UUID) error
	GetWorkerByID(id uuid.UUID) (*models.Worker, error)
	GetAllWorkers() ([]models.Worker, error)
	Update(id uuid.UUID, name string, surname string, email string, address string, phoneNumber string, role int, password string) (*models.Worker, error)

	GetWorkersByRole(role int) ([]models.Worker, error)
	GetAverageOrderRate(worker *models.Worker) (float64, error)
}
