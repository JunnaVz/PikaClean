package repository_interfaces

import (
	"github.com/google/uuid"
	"lab3/internal/models"
)

type IWorkerRepository interface {
	Create(worker *models.Worker) (*models.Worker, error)
	Update(worker *models.Worker) (*models.Worker, error)
	Delete(id uuid.UUID) error
	GetWorkerByID(id uuid.UUID) (*models.Worker, error)
	GetAllWorkers() ([]models.Worker, error)
	GetWorkerByEmail(email string) (*models.Worker, error)

	GetWorkersByRole(role int) ([]models.Worker, error)
	GetAverageOrderRate(worker *models.Worker) (float64, error)
}
