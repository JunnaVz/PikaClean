// Package userViews provides user interface functions for the PikaClean application
// focused on customer-related operations like order creation, viewing order history,
// and managing user profiles. It handles user interactions through the command line
// interface for customer-focused operations.
package userViews

import (
	"fmt"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// Get retrieves and displays detailed information about the current user.
// It fetches the user's profile from the database and presents the information
// in a formatted output including email, name, surname, phone number, and address.
// This function is typically used in the account management section of the application.
//
// Parameters:
//   - service: Service container providing access to business logic services
//   - user: Current authenticated user whose information will be displayed
//
// Returns:
//   - error: Any error that occurred during the retrieval of user information,
//     or nil if the operation was successful
func Get(service registry.Services, user *models.User) error {
	userFromDB, err := service.UserService.GetUserByID(user.ID)
	if err != nil {
		return err
	}

	fmt.Print("\nUser info:\n")
	fmt.Printf("Email: %s\nИмя: %s\nФамилия: %s\nТелефон: %s\nАдрес: %s\n", userFromDB.Email, userFromDB.Name, userFromDB.Surname, userFromDB.PhoneNumber, userFromDB.Address)
	fmt.Print("----------------\n")
	return nil
}
