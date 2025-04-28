// Package workerViews provides user interface functions for the PikaClean application
// focused on worker-related operations like authentication, profile management,
// and administrative worker functions. It handles worker interactions through the command line
// interface for operations performed by cleaning staff and administrators.
package workerViews

import (
	utils "teamdev/cmd/cmdUtils"
	"teamdev/cmd/views/stringConst"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// login handles the worker authentication process through the command line interface.
// It prompts for email and password credentials, then validates them against the database.
// If successful, the authenticated worker entity is returned for use in subsequent operations.
//
// Parameters:
//   - services: Service container providing access to business logic services,
//     particularly the WorkerService for authentication
//
// Returns:
//   - *models.Worker: The authenticated worker entity if login is successful
//   - error: Any error that occurred during the login process,
//     such as invalid credentials or database errors
func login(services registry.Services) (*models.Worker, error) {
	var email = utils.EndlessReadWord(stringConst.EmailRequest)
	var password = utils.EndlessReadWord(stringConst.PasswordRequest)

	worker, err := services.WorkerService.Login(email, password)
	if err != nil {
		return nil, err
	}

	return worker, nil
}
