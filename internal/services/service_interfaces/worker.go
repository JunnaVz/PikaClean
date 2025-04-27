// Package service_interfaces defines the contract interfaces for the business logic layer
// of the PikaClean application. These interfaces abstract the implementation details
// and provide a clear API for the application's services.
package service_interfaces

import (
	"github.com/google/uuid"
	"teamdev/internal/models"
)

// IWorkerService defines the contract for worker management operations.
// It provides methods for worker authentication, profile management, and performance metrics.
type IWorkerService interface {
	// Login authenticates a worker with email and password credentials.
	//
	// Parameters:
	//   - email: Worker's email address
	//   - password: Plain text password to validate
	//
	// Returns:
	//   - *models.Worker: Authenticated worker data
	//   - error: Error if authentication fails or credentials are invalid
	Login(email, password string) (*models.Worker, error)

	// Create registers a new worker account in the system with the specified credentials.
	//
	// Parameters:
	//   - worker: Worker information including name, contact details, role, etc.
	//   - password: Plain text password that will be hashed before storage
	//
	// Returns:
	//   - *models.Worker: Created worker with assigned ID
	//   - error: Error if registration fails or validation fails
	Create(worker *models.Worker, password string) (*models.Worker, error)

	// Delete removes a worker account from the system.
	//
	// Parameters:
	//   - id: UUID of the worker to delete
	//
	// Returns:
	//   - error: Error if deletion fails
	Delete(id uuid.UUID) error

	// GetWorkerByID retrieves a worker by their unique identifier.
	//
	// Parameters:
	//   - id: UUID of the worker to retrieve
	//
	// Returns:
	//   - *models.Worker: Retrieved worker entity
	//   - error: Error if retrieval fails or worker not found
	GetWorkerByID(id uuid.UUID) (*models.Worker, error)

	// GetAllWorkers retrieves all workers registered in the system.
	//
	// Returns:
	//   - []models.Worker: Slice of all worker entities
	//   - error: Error if retrieval fails
	GetAllWorkers() ([]models.Worker, error)

	// Update modifies an existing worker's profile information.
	//
	// Parameters:
	//   - id: UUID of the worker to update
	//   - name: New first name
	//   - surname: New last name
	//   - email: New email address
	//   - address: New physical address
	//   - phoneNumber: New contact phone number
	//   - role: New worker role (determines permissions)
	//   - password: New password (will be hashed before storage)
	//
	// Returns:
	//   - *models.Worker: Updated worker data
	//   - error: Error if update fails or validation fails
	Update(id uuid.UUID, name string, surname string, email string, address string, phoneNumber string, role int, password string) (*models.Worker, error)

	// GetWorkersByRole retrieves all workers with a specific role.
	//
	// Parameters:
	//   - role: Role identifier to filter by
	//
	// Returns:
	//   - []models.Worker: Slice of worker entities with the specified role
	//   - error: Error if retrieval fails
	GetWorkersByRole(role int) ([]models.Worker, error)

	// GetAverageOrderRate calculates the average customer satisfaction rating
	// for completed orders assigned to a specific worker.
	//
	// Parameters:
	//   - worker: Worker to calculate the average rating for
	//
	// Returns:
	//   - float64: Average rating value (0.0-5.0)
	//   - error: Error if calculation fails
	GetAverageOrderRate(worker *models.Worker) (float64, error)
}
