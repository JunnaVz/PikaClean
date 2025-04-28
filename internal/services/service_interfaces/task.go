// Package service_interfaces defines the contract interfaces for the business logic layer
// of the PikaClean application. These interfaces abstract the implementation details
// and provide a clear API for the application's services.
package service_interfaces

import (
	"github.com/google/uuid"
	"teamdev/internal/models"
)

// ITaskService defines the contract for cleaning task management operations.
// It provides methods for creating, retrieving, updating, and deleting tasks,
// as well as organizing them by categories.
type ITaskService interface {
	// Create adds a new cleaning task to the system.
	//
	// Parameters:
	//   - name: Descriptive name of the cleaning task
	//   - price: Cost per unit of the task
	//   - category: Category ID the task belongs to
	//
	// Returns:
	//   - *models.Task: Created task with assigned ID
	//   - error: Error if creation fails or validation fails
	Create(name string, price float64, category int) (*models.Task, error)

	// Update modifies an existing task's properties.
	//
	// Parameters:
	//   - taskID: UUID of the task to update
	//   - category: New category ID for the task
	//   - name: New name for the task
	//   - price: New price for the task
	//
	// Returns:
	//   - *models.Task: Updated task data
	//   - error: Error if update fails or validation fails
	Update(taskID uuid.UUID, category int, name string, price float64) (*models.Task, error)

	// Delete removes a task from the system.
	//
	// Parameters:
	//   - taskID: UUID of the task to delete
	//
	// Returns:
	//   - error: Error if deletion fails
	Delete(taskID uuid.UUID) error

	// GetAllTasks retrieves all available cleaning tasks.
	//
	// Returns:
	//   - []models.Task: Slice of all task entities
	//   - error: Error if retrieval fails
	GetAllTasks() ([]models.Task, error)

	// GetTaskByID retrieves a task by its unique identifier.
	//
	// Parameters:
	//   - id: UUID of the task to retrieve
	//
	// Returns:
	//   - *models.Task: Retrieved task entity
	//   - error: Error if retrieval fails or task not found
	GetTaskByID(id uuid.UUID) (*models.Task, error)

	// GetTasksInCategory retrieves all tasks belonging to a specific category.
	//
	// Parameters:
	//   - category: Category ID to filter by
	//
	// Returns:
	//   - []models.Task: Slice of task entities in the specified category
	//   - error: Error if retrieval fails
	GetTasksInCategory(category int) ([]models.Task, error)

	// GetTaskByName retrieves a task by its name.
	//
	// Parameters:
	//   - name: Name of the task to search for
	//
	// Returns:
	//   - *models.Task: Retrieved task entity
	//   - error: Error if retrieval fails or task not found
	GetTaskByName(name string) (*models.Task, error)
}
