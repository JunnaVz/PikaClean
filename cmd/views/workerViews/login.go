package workerViews

import (
	utils "lab3/cmd/cmdUtils"
	"lab3/cmd/views/stringConst"
	"lab3/internal/models"
	"lab3/internal/registry"
)

func login(services registry.Services) (*models.Worker, error) {
	var email = utils.EndlessReadWord(stringConst.EmailRequest)
	var password = utils.EndlessReadWord(stringConst.PasswordRequest)

	worker, err := services.WorkerService.Login(email, password)
	if err != nil {
		return nil, err
	}

	return worker, nil
}
