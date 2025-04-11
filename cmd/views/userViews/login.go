package userViews

import (
	utils "teamdev/cmd/cmdUtils"
	"teamdev/cmd/views/stringConst"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

func login(services registry.Services) (*models.User, error) {
	var email = utils.EndlessReadWord(stringConst.EmailRequest)
	var password = utils.EndlessReadWord(stringConst.PasswordRequest)

	client, err := services.UserService.Login(email, password)
	if err != nil {
		return nil, err
	}

	return client, nil
}
