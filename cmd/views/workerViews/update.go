// Package workerViews provides user interface functions for the PikaClean application
// focused on worker-related operations including profile viewing, updating, and management.
// This file contains functionality for updating worker profile information.
package workerViews

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"teamdev/cmd/cmdUtils"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// requestForChange prompts the user to input a new value for a specific worker profile field.
// If the user provides no input (empty string), the original value is maintained.
//
// Parameters:
//   - fieldName: The name of the field being updated (displayed in the prompt)
//   - fieldValue: The current value of the field
//   - word: Boolean flag indicating whether to read a single word (true) or multiple words (false)
//
// Returns:
//   - string: The new value if provided, or the original value if input was empty
func requestForChange(fieldName string, fieldValue string, word bool) string {
	fmt.Printf("Изменить %s (оставьте пустым, чтобы не менять): ", fieldName)

	input, _ := cmdUtils.StringReader(word)
	strings.TrimSpace(input)

	if len(input) == 0 {
		return fieldValue
	}

	return input
}

// Update modifies a worker's profile information based on user input.
// It allows changing various fields including email, password, name, surname,
// phone number, address, and role (if the editor has manager privileges).
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - workerID: UUID of the worker whose profile is being updated
//   - editor: The worker who is performing the update operation
//
// Returns:
//   - error: Any error that occurred during the update operation
func Update(services registry.Services, workerID uuid.UUID, editor *models.Worker) error {
	worker, err := services.WorkerService.GetWorkerByID(workerID)

	if err != nil {
		return err
	}

	var email = requestForChange("email", worker.Email, true)
	var password = requestForChange("пароль", worker.Password, true)
	var name = requestForChange("имя", worker.Name, true)
	var surname = requestForChange("фамилию", worker.Surname, true)
	var phoneNumber = requestForChange("номер телефона", worker.PhoneNumber, true)
	var address = requestForChange("адрес", worker.Address, false)
	var role int

	if editor.Role == models.ManagerRole {
		roleStr := requestForChange("роль (1 - менеджер, 2 - мастер)", worker.DisplayRole(), true)
		switch roleStr {
		case "1":
			role = models.ManagerRole
		case "2":
			role = models.MasterRole
		default:
			role = worker.Role
		}
	} else {
		role = worker.Role
	}

	_, err = services.WorkerService.Update(worker.ID, name, surname, email, address, phoneNumber, role, password)

	if err != nil {
		return err
	}

	return nil
}
