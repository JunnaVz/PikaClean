// Package repository_interfaces defines the contract interfaces for data persistence
// operations. It specifies the required functionality that any repository implementation
// must satisfy to interact with the application's domain models.
package repository_interfaces

import (
	"github.com/google/uuid"
	"teamdev/internal/models"
)

// IUserRepository defines the contract for user data persistence operations.
// Any implementation of this interface must provide methods for creating, retrieving,
// updating, and deleting user entities.
type IUserRepository interface {
	// Create adds a new user record to the data store.
	//
	// Parameters:
	//   - user: User entity to be persisted
	//
	// Returns:
	//   - *models.User: Created user with assigned ID
	//   - error: Error if creation fails
	Create(user *models.User) (*models.User, error)

	// Update modifies an existing user record in the data store.
	//
	// Parameters:
	//   - user: User entity with updated values
	//
	// Returns:
	//   - *models.User: Updated user data
	//   - error: Error if update fails
	Update(user *models.User) (*models.User, error)

	// Delete removes a user record from the data store by ID.
	// This typically includes removing all associated data like orders.
	//
	// Parameters:
	//   - id: UUID of the user to delete
	//
	// Returns:
	//   - error: Error if deletion fails
	Delete(id uuid.UUID) error

	// GetUserByID retrieves a user by unique identifier.
	//
	// Parameters:
	//   - id: UUID of the user to retrieve
	//
	// Returns:
	//   - *models.User: Retrieved user entity
	//   - error: Error if retrieval fails or user not found
	GetUserByID(id uuid.UUID) (*models.User, error)

	// GetAllUsers retrieves all users from the data store.
	// For security reasons, this method typically does not return password data.
	//
	// Returns:
	//   - []models.User: Slice of all user entities
	//   - error: Error if retrieval fails
	GetAllUsers() ([]models.User, error)

	// GetUserByEmail retrieves a user by email address.
	// This is commonly used for authentication purposes.
	//
	// Parameters:
	//   - email: Email address to search for
	//
	// Returns:
	//   - *models.User: Retrieved user entity
	//   - error: Error if retrieval fails or user not found
	GetUserByEmail(email string) (*models.User, error)
}
