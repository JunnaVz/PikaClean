// Package service_interfaces defines the contract interfaces for the business logic layer
// of the PikaClean application. These interfaces abstract the implementation details
// and provide a clear API for the application's services.
package service_interfaces

import "teamdev/internal/models"

// ICategoryService defines the contract for category management operations.
// This interface provides methods for retrieving, creating, updating, and deleting
// categories, as well as accessing tasks within categories.
type ICategoryService interface {
	// GetAll retrieves all available service categories.
	//
	// Returns:
	//   - []models.Category: Slice of all category entities
	//   - error: Error if retrieval fails
	GetAll() ([]models.Category, error)

	// GetTasksInCategory retrieves all cleaning tasks belonging to a specific category.
	//
	// Parameters:
	//   - id: Numeric ID of the category to retrieve tasks for
	//
	// Returns:
	//   - []models.Task: Slice of task entities in the specified category
	//   - error: Error if retrieval fails
	GetTasksInCategory(id int) ([]models.Task, error)

	// GetByID retrieves a category by its unique identifier.
	//
	// Parameters:
	//   - id: Numeric ID of the category to retrieve
	//
	// Returns:
	//   - *models.Category: Retrieved category entity
	//   - error: Error if retrieval fails or category not found
	GetByID(id int) (*models.Category, error)

	// Create adds a new service category with the specified name.
	//
	// Parameters:
	//   - name: Name for the new category
	//
	// Returns:
	//   - *models.Category: Created category with assigned ID
	//   - error: Error if creation fails or validation fails
	Create(name string) (*models.Category, error)

	// Update modifies an existing category's properties.
	//
	// Parameters:
	//   - category: Category entity with updated values
	//
	// Returns:
	//   - *models.Category: Updated category data
	//   - error: Error if update fails or validation fails
	Update(category *models.Category) (*models.Category, error)

	// Delete removes a category by its ID.
	// Note: This may fail if there are tasks still associated with the category.
	//
	// Parameters:
	//   - id: Numeric ID of the category to delete
	//
	// Returns:
	//   - error: Error if deletion fails
	Delete(id int) error
}
