// Package repository_interfaces defines the contract interfaces for data persistence
// operations. It specifies the required functionality that any repository implementation
// must satisfy to interact with the application's domain models.
package repository_interfaces

import (
	"github.com/google/uuid"
	"teamdev/internal/models"
)

// ITaskRepository defines the contract for task data persistence operations.
// Any implementation of this interface must provide methods for creating, retrieving,
// updating, and deleting task entities, as well as specialized task queries.
type ITaskRepository interface {
	// Create adds a new task record to the data store.
	//
	// Parameters:
	//   - task: Task entity to be persisted
	//
	// Returns:
	//   - *models.Task: Created task with assigned ID
	//   - error: Error if creation fails
	Create(task *models.Task) (*models.Task, error)

	// Delete removes a task record from the data store by ID.
	//
	// Parameters:
	//   - id: UUID of the task to delete
	//
	// Returns:
	//   - error: Error if deletion fails
	Delete(id uuid.UUID) error

	// Update modifies an existing task record in the data store.
	//
	// Parameters:
	//   - task: Task entity with updated values
	//
	// Returns:
	//   - *models.Task: Updated task data
	//   - error: Error if update fails
	Update(task *models.Task) (*models.Task, error)

	// GetTaskByID retrieves a task by unique identifier.
	//
	// Parameters:
	//   - id: UUID of the task to retrieve
	//
	// Returns:
	//   - *models.Task: Retrieved task entity
	//   - error: Error if retrieval fails or task not found
	GetTaskByID(id uuid.UUID) (*models.Task, error)

	// GetAllTasks retrieves all tasks from the data store.
	//
	// Returns:
	//   - []models.Task: Slice of all task entities
	//   - error: Error if retrieval fails
	GetAllTasks() ([]models.Task, error)

	// GetTasksInCategory retrieves all tasks belonging to a specific category.
	//
	// Parameters:
	//   - category: Category ID to filter by
	//
	// Returns:
	//   - []models.Task: Slice of tasks in the specified category
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
