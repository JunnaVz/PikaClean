package workerViews

import (
	"fmt"
	"lab3/cmd/menu"
	"lab3/internal/models"
	"lab3/internal/registry"
)

func WorkerLoginMenu(services registry.Services) error {
	var m menu.Menu
	m.CreateMenu(
		[]menu.Item{
			{
				Name: "войти",
				Handler: func() error {
					worker, err := login(services)
					if err == nil {
						if worker.Role == models.ManagerRole {
							err = managerMainMenu(services, worker)
						} else if worker.Role == models.MasterRole {
							err = workerMainMenu(services, worker)
						} else {
							err = fmt.Errorf("")
						}
					}
					return err
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

func managerMainMenu(services registry.Services, worker *models.Worker) error {
	// Создание меню и добавление элементов
	var m menu.Menu
	m.CreateMenu(
		[]menu.Item{
			{
				Name: "Просмотреть профиль",
				Handler: func() error {
					return Get(services, worker)
				},
			},
			{
				Name: "Изменить профиль",
				Handler: func() error {
					return Update(services, worker.ID, worker)
				},
			},
			{
				Name: "Список работников",
				Handler: func() error {
					return getAllWorkers(services, worker)
				},
			},
			{
				Name: "Добавить работника",
				Handler: func() error {
					return create(services)
				},
			},
			{
				Name: "Посмотреть неназначенные заказы",
				Handler: func() error {
					return unassignedOrders(services)
				},
			},
			{
				Name: "Посмотреть заказы в работе",
				Handler: func() error {
					return inProgressOrders(services)
				},
			},
			{
				Name: "Посмотреть законченные заказы",
				Handler: func() error {
					return completedOrders(services)
				},
			},
			{
				Name: "База услуг",
				Handler: func() error {
					return managerTasks(services)
				},
			},
		})

	// Показать меню
	err := m.Menu()
	if err != nil {
		return err
	}

	return nil
}

func workerMainMenu(services registry.Services, worker *models.Worker) error {
	var m menu.Menu
	m.CreateMenu(
		[]menu.Item{
			{
				Name: "Просмотреть профиль",
				Handler: func() error {
					return Get(services, worker)
				},
			},
			{
				Name: "Изменить профиль",
				Handler: func() error {
					return Update(services, worker.ID, worker)
				},
			},
			{
				Name: "Посмотреть законченные заказы",
				Handler: func() error {
					return completedOrdersByWorker(services, worker)
				},
			},
			{
				Name: "Посмотреть заказы в работе",
				Handler: func() error {
					return inProgressOrdersByWorker(services, worker)
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
