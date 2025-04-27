// Package userViews provides user interface functions for the PikaClean application
// focused on customer-related operations like order creation, viewing order history,
// and managing user profiles. It handles user interactions through the command line
// interface for customer-focused operations.
package userViews

import (
	utils "teamdev/cmd/cmdUtils"
	"teamdev/cmd/views/stringConst"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// login handles the user authentication process through the command line interface.
// It prompts the user for email and password credentials, then validates them
// against the database. This function is typically the entry point for user sessions
// and is called before allowing access to customer-specific functionality.
//
// Parameters:
//   - services: Service container providing access to business logic services,
//     particularly the UserService for authentication
//
// Returns:
//   - *models.User: Authenticated user entity containing profile information
//     if login was successful, or nil if authentication failed
//   - error: Any error that occurred during the login process,
//     such as invalid credentials or database connectivity issues
func login(services registry.Services) (*models.User, error) {
	var email = utils.EndlessReadWord(stringConst.EmailRequest)
	var password = utils.EndlessReadWord(stringConst.PasswordRequest)

	client, err := services.UserService.Login(email, password)
	if err != nil {
		return nil, err
	}

	return client, nil
}
