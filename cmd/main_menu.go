// Package cmd contains command-line interface components for the PikaClean application.
// It provides the entry points, menu systems and command handlers that form the
// user interface layer of the application.
package cmd

import (
	"fmt"
	"teamdev/cmd/menu"
	"teamdev/cmd/views/taskViews"
	"teamdev/cmd/views/userViews"
	"teamdev/cmd/views/workerViews"
	"teamdev/internal/registry"
)

// RunMenu displays and handles the main application menu which serves as the entry point
// for all user interactions. It presents three primary paths: customer login, worker login,
// or viewing services as a guest. The menu uses the menu package to create an interactive
// console-based navigation system.
//
// This function creates the top-level menu with three options:
//   - "клиент" (customer): Redirects to the customer login menu
//   - "работник" (worker): Redirects to the worker login menu
//   - "гость, посмотреть цены" (guest, view prices): Shows available services and pricing
//
// Parameters:
//   - a: Pointer to registry.Services containing all application business logic services
//
// Returns:
//   - error: Any error that occurred during menu navigation or handling
func RunMenu(a *registry.Services) error {
	fmt.Print("Кто вы?\n")
	var m menu.Menu
	m.CreateMenu(
		[]menu.Item{
			{
				Name: "клиент",
				Handler: func() error {
					return userViews.UserLoginMenu(*a)
				},
			},
			{
				Name: "работник",
				Handler: func() error {
					return workerViews.WorkerLoginMenu(*a)
				},
			},
			{
				Name: "гость, посмотреть цены",
				Handler: func() error {
					return taskViews.AllTasks(*a)
				},
			},
		},
	)

	// Показать меню
	err := m.Menu()
	if err != nil {
		return err
	}

	return nil
}
