// Package service_interfaces defines the contract interfaces for the business logic layer
// of the PikaClean application. These interfaces abstract the implementation details
// and provide a clear API for the application's services.
package service_interfaces

import (
	"github.com/google/uuid"
	"teamdev/internal/models"
)

// IUserService defines the contract for user management operations.
// It provides methods for user authentication, registration, and profile management.
type IUserService interface {
	// GetUserByID retrieves a user by their unique identifier.
	//
	// Parameters:
	//   - id: UUID of the user to retrieve
	//
	// Returns:
	//   - *models.User: Retrieved user entity
	//   - error: Error if retrieval fails or user not found
	GetUserByID(id uuid.UUID) (*models.User, error)

	// Register creates a new user account in the system with the specified credentials.
	//
	// Parameters:
	//   - user: User information including name, contact details, etc.
	//   - password: Plain text password that will be hashed before storage
	//
	// Returns:
	//   - *models.User: Created user with assigned ID
	//   - error: Error if registration fails or validation fails
	Register(user *models.User, password string) (*models.User, error)

	// Login authenticates a user with email and password credentials.
	//
	// Parameters:
	//   - email: User's email address
	//   - password: Plain text password to validate
	//
	// Returns:
	//   - *models.User: Authenticated user data
	//   - error: Error if authentication fails or credentials are invalid
	Login(email, password string) (*models.User, error)

	// Update modifies an existing user's profile information.
	//
	// Parameters:
	//   - id: UUID of the user to update
	//   - name: New first name
	//   - surname: New last name
	//   - email: New email address
	//   - address: New physical address
	//   - phoneNumber: New contact phone number
	//   - password: New password (will be hashed before storage)
	//
	// Returns:
	//   - *models.User: Updated user data
	//   - error: Error if update fails or validation fails
	Update(id uuid.UUID, name string, surname string, email string, address string, phoneNumber string, password string) (*models.User, error)

	// GetUserByEmail retrieves a user by their email address.
	//
	// Parameters:
	//   - email: Email address to search for
	//
	// Returns:
	//   - *models.User: Retrieved user entity
	//   - error: Error if retrieval fails or user not found
	GetUserByEmail(email string) (*models.User, error)
}
