// Package repository_interfaces defines the contract interfaces for data persistence
// operations. It specifies the required functionality that any repository implementation
// must satisfy to interact with the application's domain models.
package repository_interfaces

import (
	"github.com/google/uuid"
	"teamdev/internal/models"
)

// IWorkerRepository defines the contract for worker data persistence operations.
// Any implementation of this interface must provide methods for creating, retrieving,
// updating, and deleting worker entities, as well as specialized worker queries.
type IWorkerRepository interface {
	// Create adds a new worker record to the data store.
	//
	// Parameters:
	//   - worker: Worker entity to be persisted
	//
	// Returns:
	//   - *models.Worker: Created worker with assigned ID
	//   - error: Error if creation fails
	Create(worker *models.Worker) (*models.Worker, error)

	// Update modifies an existing worker record in the data store.
	//
	// Parameters:
	//   - worker: Worker entity with updated values
	//
	// Returns:
	//   - *models.Worker: Updated worker data
	//   - error: Error if update fails
	Update(worker *models.Worker) (*models.Worker, error)

	// Delete removes a worker record from the data store by ID.
	//
	// Parameters:
	//   - id: UUID of the worker to delete
	//
	// Returns:
	//   - error: Error if deletion fails
	Delete(id uuid.UUID) error

	// GetWorkerByID retrieves a worker by unique identifier.
	//
	// Parameters:
	//   - id: UUID of the worker to retrieve
	//
	// Returns:
	//   - *models.Worker: Retrieved worker entity
	//   - error: Error if retrieval fails or worker not found
	GetWorkerByID(id uuid.UUID) (*models.Worker, error)

	// GetAllWorkers retrieves all workers from the data store.
	//
	// Returns:
	//   - []models.Worker: Slice of all worker entities
	//   - error: Error if retrieval fails
	GetAllWorkers() ([]models.Worker, error)

	// GetWorkerByEmail retrieves a worker by email address.
	//
	// Parameters:
	//   - email: Email address to search for
	//
	// Returns:
	//   - *models.Worker: Retrieved worker entity
	//   - error: Error if retrieval fails or worker not found
	GetWorkerByEmail(email string) (*models.Worker, error)

	// GetWorkersByRole retrieves workers filtered by role.
	//
	// Parameters:
	//   - role: Role ID to filter by
	//
	// Returns:
	//   - []models.Worker: Slice of workers with the specified role
	//   - error: Error if retrieval fails
	GetWorkersByRole(role int) ([]models.Worker, error)

	// GetAverageOrderRate calculates the average rating for completed orders
	// assigned to a specific worker.
	//
	// Parameters:
	//   - worker: Worker to calculate average rating for
	//
	// Returns:
	//   - float64: Average rating value (0.0-5.0)
	//   - error: Error if calculation fails
	GetAverageOrderRate(worker *models.Worker) (float64, error)
}
