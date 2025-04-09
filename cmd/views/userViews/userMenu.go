package userViews

import (
	"lab3/cmd/menu"
	"lab3/internal/models"
	"lab3/internal/registry"
)

func UserLoginMenu(services registry.Services) error {
	var m menu.Menu
	m.CreateMenu(
		[]menu.Item{
			{
				Name: "войти",
				Handler: func() error {
					user, err := login(services)
					if err == nil {
						return userMainMenu(services, user)
					}
					return err
				},
			},
			{
				Name: "зарегистрироваться",
				Handler: func() error {
					user, err := registration(services)
					if err == nil {
						return userMainMenu(services, user)
					}
					return err
				},
			},
		},
	)

	err := m.Menu()
	if err != nil {
		return err
	}

	return nil
}

func userMainMenu(services registry.Services, user *models.User) error {
	var m menu.Menu
	m.CreateMenu(
		[]menu.Item{
			{
				Name: "просмотреть профиль",
				Handler: func() error {
					return Get(services, user)
				},
			},
			{
				Name: "изменить профиль",
				Handler: func() error {
					return Update(services, user)
				},
			},
			{
				Name: "создать заказ",
				Handler: func() error {
					return createOrder(services, user)
				},
			},
			{
				Name: "посмотреть законченные заказы",
				Handler: func() error {
					return getCompletedOrders(services, user)
				},
			},
			{
				Name: "посмотреть заказы в работе",
				Handler: func() error {
					return getOrdersInWork(services, user)

				},
			},
		},
	)

	err := m.Menu()
	if err != nil {
		return err
	}

	return nil
}
