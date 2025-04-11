package workerViews

import (
	utils "teamdev/cmd/cmdUtils"
	"teamdev/cmd/views/stringConst"
	"teamdev/internal/models"
	"teamdev/internal/registry"
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
