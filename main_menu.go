package cmd

import (
	"fmt"
	"lab3/cmd/menu"
	"lab3/cmd/views/taskViews"
	"lab3/cmd/views/userViews"
	"lab3/cmd/views/workerViews"
	"lab3/internal/registry"
)

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
