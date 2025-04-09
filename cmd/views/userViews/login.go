package userViews

import (
	utils "lab3/cmd/cmdUtils"
	"lab3/cmd/views/stringConst"
	"lab3/internal/models"
	"lab3/internal/registry"
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
