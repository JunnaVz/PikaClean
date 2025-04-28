// Package workerViews provides user interface functions for the PikaClean application
// focused on worker-related operations including profile viewing, updating, and management.
// This file contains functionality for worker menu navigation and interaction.
package workerViews

import (
	"fmt"
	"teamdev/cmd/menu"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// WorkerLoginMenu displays the login menu for worker authentication.
// It provides an entry point for workers to authenticate and access the system.
//
// Parameters:
//   - services: Service container providing access to business logic services
//
// Returns:
//   - error: Any error that occurred during the login process
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

// managerMainMenu displays the main navigation menu for managers.
// It provides access to manager-specific functionality including worker management,
// order tracking, and task management.
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - worker: The authenticated manager worker object
//
// Returns:
//   - error: Any error that occurred during menu navigation
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

// workerMainMenu displays the main navigation menu for regular workers (masters).
// It provides access to worker-specific functionality including profile management
// and order tracking.
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - worker: The authenticated worker object
//
// Returns:
//   - error: Any error that occurred during menu navigation
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
