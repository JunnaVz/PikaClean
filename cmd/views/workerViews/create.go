// Package workerViews provides user interface functions for the PikaClean application
// focused on worker-related operations like worker registration, profile management,
// and administrative worker functions. It handles worker interactions through the command line
// interface for operations performed by managers and administrators.
package workerViews

import (
	"fmt"
	utils "teamdev/cmd/cmdUtils"
	"teamdev/cmd/views/stringConst"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// create handles the worker registration process through the command line interface.
// It prompts for worker details including personal information and role,
// then submits the information to create a new worker account in the system.
// The function supports creation of both manager and master (cleaner) roles.
//
// Parameters:
//   - services: Service container providing access to business logic services,
//     particularly the WorkerService for worker creation
//
// Returns:
//   - error: Any error that occurred during the worker creation process,
//     such as validation failures or database errors
func create(services registry.Services) error {
	var worker *models.Worker
	var err error

	var email = utils.EndlessReadWord(stringConst.EmailRequest)
	var password = utils.EndlessReadWord(stringConst.PasswordRequest)
	var name = utils.EndlessReadWord(stringConst.NameRequest)
	var surname = utils.EndlessReadWord(stringConst.SurnameRequest)
	var phoneNumber = utils.EndlessReadWord(stringConst.PhoneRequest)
	var address = utils.EndlessReadRow(stringConst.AddressRequest)
	var roleStr = utils.EndlessReadWord(stringConst.RoleRequest)
	var role int

	switch roleStr {
	case "1":
		role = models.ManagerRole
	default:
		role = models.MasterRole
	}

	worker, err = services.WorkerService.Create(&models.Worker{
		Email:       email,
		Name:        name,
		Surname:     surname,
		PhoneNumber: phoneNumber,
		Address:     address,
		Role:        role,
	}, password)

	if err != nil {
		return err
	}

	fmt.Printf("%s %s %s успешно зарегистрирован\n\n\n", worker.DisplayRole(), worker.Name, worker.Surname)

	return nil
}
