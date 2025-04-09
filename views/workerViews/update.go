package workerViews

import (
	"fmt"
	"github.com/google/uuid"
	"lab3/cmd/cmdUtils"
	"lab3/internal/models"
	"lab3/internal/registry"
	"strings"
)

func requestForChange(fieldName string, fieldValue string, word bool) string {
	fmt.Printf("Изменить %s (оставьте пустым, чтобы не менять): ", fieldName)

	input, _ := cmdUtils.StringReader(word)
	strings.TrimSpace(input)

	if len(input) == 0 {
		return fieldValue
	}

	return input
}

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
