// Package userViews provides user interface functions for the PikaClean application
// focused on customer-related operations like order creation, viewing order history,
// and managing user profiles. It handles user interactions through the command line
// interface for customer-focused operations.
package userViews

import (
	"fmt"
	utils "teamdev/cmd/cmdUtils"
	"teamdev/cmd/views/stringConst"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// registration handles the user registration process through the command line interface.
// It prompts the user for personal information including email, password, name, surname,
// phone number and address, then creates a new user account in the system. After successful
// registration, it displays a confirmation message.
//
// Parameters:
//   - services: Service container providing access to business logic services,
//     particularly the UserService for account creation
//
// Returns:
//   - *models.User: Newly created user entity containing profile information
//     if registration was successful, or nil if registration failed
//   - error: Any error that occurred during the registration process,
//     such as invalid data format or duplicate email address
func registration(services registry.Services) (*models.User, error) {
	var user *models.User
	var err error

	var email = utils.EndlessReadWord(stringConst.EmailRequest)
	var password = utils.EndlessReadWord(stringConst.PasswordRequest)
	var name = utils.EndlessReadWord(stringConst.NameRequest)
	var surname = utils.EndlessReadWord(stringConst.SurnameRequest)
	var phoneNumber = utils.EndlessReadWord(stringConst.PhoneRequest)
	var address = utils.EndlessReadRow(stringConst.AddressRequest)

	user, err = services.UserService.Register(&models.User{
		Email:       email,
		Name:        name,
		Surname:     surname,
		PhoneNumber: phoneNumber,
		Address:     address,
	}, password)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Пользователь %s %s успешно зарегистрирован\n\n\n", user.Name, user.Surname)

	return user, nil
}
