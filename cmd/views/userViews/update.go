// Package userViews provides user interface functions for the PikaClean application
// focused on customer-related operations like order creation, viewing order history,
// and managing user profiles. It handles user interactions through the command line
// interface for customer-focused operations.
package userViews

import (
	"fmt"
	"strings"
	"teamdev/cmd/cmdUtils"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// requestForChange prompts the user to update a specific profile field.
// It displays the current value of the field and allows the user to either
// provide a new value or keep the existing one by entering nothing.
//
// Parameters:
//   - fieldName: Name of the profile field being updated (for display to the user)
//   - fieldValue: Current value of the field
//   - word: Boolean indicating whether input should be treated as a single word (true)
//     or can contain multiple words/spaces (false)
//
// Returns:
//   - string: New value if the user entered one, or original value if left empty
func requestForChange(fieldName string, fieldValue string, word bool) string {
	fmt.Printf("Изменить %s (оставьте пустым, чтобы не менять): ", fieldName)

	input, err := cmdUtils.StringReader(word)
	strings.TrimSpace(input)

	if err != nil || len(input) == 0 {
		return fieldValue
	}

	return input
}

// Update handles the profile update process for a user through the command line interface.
// It retrieves the user's current profile data and prompts for updates to each field,
// then submits the updated information to the user service.
//
// Parameters:
//   - services: Service container providing access to business logic services,
//     particularly the UserService for profile updates
//   - user: Current authenticated user whose profile will be updated
//
// Returns:
//   - error: Any error that occurred during the update process,
//     such as validation failures or database errors
func Update(services registry.Services, user *models.User) error {
	userFromDB, err := services.UserService.GetUserByID(user.ID)

	var email = requestForChange("email", userFromDB.Email, true)
	var password = requestForChange("пароль", userFromDB.Password, true)
	var name = requestForChange("имя", userFromDB.Name, true)
	var surname = requestForChange("фамилию", userFromDB.Surname, true)
	var phoneNumber = requestForChange("номер телефона", userFromDB.PhoneNumber, true)
	var address = requestForChange("адрес", userFromDB.Address, false)

	_, err = services.UserService.Update(user.ID, name, surname, email, address, phoneNumber, password)

	if err != nil {
		return err
	}

	return nil
}
