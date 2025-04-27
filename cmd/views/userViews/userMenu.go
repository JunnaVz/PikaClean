// Package userViews provides user interface functions for the PikaClean application
// focused on customer-related operations like order creation, viewing order history,
// and managing user profiles. It handles user interactions through the command line
// interface for customer-focused operations.
package userViews

import (
	"teamdev/cmd/menu"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// UserLoginMenu displays the initial menu for customer authentication,
// allowing users to either login with existing credentials or register
// as a new user. Upon successful authentication, it directs to the main
// user menu.
//
// Parameters:
//   - services: Service container providing access to business logic services
//
// Returns:
//   - error: Any error that occurred during menu operation,
//     or nil if menu exited normally
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

// userMainMenu displays the primary navigation menu for authenticated customers.
// It offers options for profile management, order creation, and order history viewing.
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - user: Currently authenticated user whose context is used for operations
//
// Returns:
//   - error: Any error that occurred during menu operation,
//     or nil if menu exited normally
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
