package userViews

import (
	"fmt"
	"lab3/cmd/cmdUtils"
	"lab3/internal/models"
	"lab3/internal/registry"
	"strings"
)

func requestForChange(fieldName string, fieldValue string, word bool) string {
	fmt.Printf("Изменить %s (оставьте пустым, чтобы не менять): ", fieldName)

	input, err := cmdUtils.StringReader(word)
	strings.TrimSpace(input)

	if err != nil || len(input) == 0 {
		return fieldValue
	}

	return input
}

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
