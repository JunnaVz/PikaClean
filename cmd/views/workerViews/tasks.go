// Package workerViews provides user interface functions for the PikaClean application
// focused on worker-related operations including task management for managers.
// This file contains functionality for viewing, editing, and managing cleaning tasks.
package workerViews

import (
	"fmt"
	"teamdev/cmd/menu"
	"teamdev/cmd/modelTables"
	"teamdev/cmd/views/taskViews"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// pickTaskForEditing displays a list of tasks and allows a manager to select
// and edit individual tasks. It provides an interactive interface for task modification.
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - tasks: Slice of tasks to display and make available for editing
//
// Returns:
//   - error: Any error that occurred during operation
func pickTaskForEditing(services registry.Services, tasks []models.Task) error {
	var err error
	var taskID int

	for {
		err = modelTables.Tasks(tasks)
		if err != nil {
			return err
		}

		fmt.Printf("Выберите услуги для заказа. Введите 0, чтобы вернуться обратно.\n")

		fmt.Scanf("%d", &taskID)
		if taskID == 0 {
			return nil
		}

		if taskID > 0 && taskID <= len(tasks) {
			updatedTask, updErr := taskViews.Update(services, tasks[taskID-1])
			if updErr != nil {
				fmt.Println(updErr.Error())
			} else {
				tasks[taskID-1] = *updatedTask
			}
		} else {
			fmt.Println("Неверный номер услуги")
		}
	}
}

// managerTasks displays a menu with options for task management operations
// that are available to managers. It allows viewing all tasks, viewing tasks by category,
// and creating new tasks.
//
// Parameters:
//   - services: Service container providing access to business logic services
//
// Returns:
//   - error: Any error that occurred during operation
func managerTasks(services registry.Services) error {
	var m menu.Menu
	m.CreateMenu(
		[]menu.Item{
			{
				Name: "Просмотреть все услуги",
				Handler: func() error {
					tasks, err := services.TaskService.GetAllTasks()
					if err != nil {
						fmt.Println(err.Error())
					}
					return pickTaskForEditing(services, tasks)
				},
			},
			{
				Name: "Просмотреть по категории",
				Handler: func() error {
					category := taskViews.ChooseTaskCategory()
					tasks, err := taskViews.TasksByCategory(services, category)
					if err != nil {
						fmt.Println(err.Error())
					}
					return pickTaskForEditing(services, tasks)
				},
			},
			{
				Name: "Создать новую услугу",
				Handler: func() error {
					return taskViews.Create(services)
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
