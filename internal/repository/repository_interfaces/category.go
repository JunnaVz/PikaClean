// Package repository_interfaces defines the contract interfaces for data persistence
// operations. It specifies the required functionality that any repository implementation
// must satisfy to interact with the application's domain models.
package repository_interfaces

import "teamdev/internal/models"

// ICategoryRepository defines the contract for category data persistence operations.
// Any implementation of this interface must provide methods for creating, retrieving,
// updating, and deleting category entities.
type ICategoryRepository interface {
	// GetAll retrieves all categories from the data store.
	//
	// Returns:
	//   - []models.Category: Slice of all category entities
	//   - error: Error if retrieval fails
	GetAll() ([]models.Category, error)

	// GetByID retrieves a category by its unique identifier.
	//
	// Parameters:
	//   - id: Numeric ID of the category to retrieve
	//
	// Returns:
	//   - *models.Category: Retrieved category entity
	//   - error: Error if retrieval fails or category not found
	GetByID(id int) (*models.Category, error)

	// Create adds a new category record to the data store.
	//
	// Parameters:
	//   - category: Category entity to be persisted
	//
	// Returns:
	//   - *models.Category: Created category with assigned ID
	//   - error: Error if creation fails
	Create(category *models.Category) (*models.Category, error)

	// Update modifies an existing category record in the data store.
	//
	// Parameters:
	//   - category: Category entity with updated values
	//
	// Returns:
	//   - *models.Category: Updated category data
	//   - error: Error if update fails
	Update(category *models.Category) (*models.Category, error)

	// Delete removes a category record from the data store by ID.
	//
	// Parameters:
	//   - id: Numeric ID of the category to delete
	//
	// Returns:
	//   - error: Error if deletion fails
	Delete(id int) error
}
