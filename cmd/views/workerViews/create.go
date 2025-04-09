package workerViews

import (
	"fmt"
	utils "lab3/cmd/cmdUtils"
	"lab3/cmd/views/stringConst"
	"lab3/internal/models"
	"lab3/internal/registry"
)

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
