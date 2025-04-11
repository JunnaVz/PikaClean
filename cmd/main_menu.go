package cmd

import (
	"fmt"
	"teamdev/cmd/menu"
	"teamdev/cmd/views/taskViews"
	"teamdev/cmd/views/userViews"
	"teamdev/cmd/views/workerViews"
	"teamdev/internal/registry"
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
