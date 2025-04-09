package userViews

import (
	"fmt"
	"lab3/internal/models"
	"lab3/internal/registry"
)

func Get(service registry.Services, user *models.User) error {
	userFromDB, err := service.UserService.GetUserByID(user.ID)
	if err != nil {
		return err
	}

	fmt.Print("\nUser info:\n")
	fmt.Printf("Email: %s\nИмя: %s\nФамилия: %s\nТелефон: %s\nАдрес: %s\n", userFromDB.Email, userFromDB.Name, userFromDB.Surname, userFromDB.PhoneNumber, userFromDB.Address)
	fmt.Print("----------------\n")
	return nil
}
