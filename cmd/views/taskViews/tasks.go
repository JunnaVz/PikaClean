// Package taskViews provides view functions for managing cleaning tasks
// in the PikaClean application. It contains functions that encapsulate user interaction
// for creating, updating, and displaying cleaning service tasks.
package taskViews

import (
	"fmt"
	"teamdev/cmd/modelTables"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// AllTasks displays all available cleaning tasks in a tabular format.
// It retrieves all tasks from the service layer and passes them to the
// modelTables formatter for display.
//
// Parameters:
//   - services: Service container providing access to business logic services
//
// Returns:
//   - error: Any error that occurred during task retrieval or display,
//     or nil if the operation was successful
func AllTasks(services registry.Services) error {
	tasks, err := services.TaskService.GetAllTasks()
	if err != nil {
		return err
	}
	return modelTables.Tasks(tasks)
}

// TasksByCategory retrieves tasks belonging to a specific category.
// It returns the list of tasks that match the specified category ID.
//
// Parameters:
//   - service: Service container providing access to business logic services
//   - category: ID of the category to filter tasks by
//
// Returns:
//   - []models.Task: List of tasks in the specified category
//   - error: Any error that occurred during task retrieval,
//     or nil if the operation was successful
func TasksByCategory(service registry.Services, category int) ([]models.Task, error) {
	tasks, err := service.TaskService.GetTasksInCategory(category)
	if err != nil {
		return nil, err
	}

	return tasks, err
}

// ChooseTaskCategory displays a list of available task categories and
// prompts the user to select one. It validates the selection and returns
// the selected category ID.
//
// Returns:
//   - int: The ID of the selected task category
func ChooseTaskCategory() int {
	fmt.Println("Выберите категорию задачи:")

	for i, category := range models.TaskCategories {
		fmt.Printf("%d. %s\n", i+1, category)
	}

	var category int
	for {
		fmt.Scanf("%d", &category)
		//print(len(models.TaskCategories), category)
		if category < 1 || category > len(models.TaskCategories) {
			fmt.Println("Неверный номер категории")
		} else {
			return category
		}
	}
}

// Tasks provides a menu-driven interface for viewing tasks, allowing the user
// to view all tasks or filter them by category. It returns the selected list
// of tasks after successful display.
//
// Parameters:
//   - services: Service container providing access to business logic services
//
// Returns:
//   - []models.Task: List of tasks based on the user's selection
//   - error: Any error that occurred during task retrieval or display,
//     or nil if the operation was successful
func Tasks(services registry.Services) ([]models.Task, error) {
	const menu = "1 -- просмотреть все услуги \n2 -- смотреть по категории\nВыберите действие: "
	var action int
	var tasks []models.Task

	for {
		fmt.Printf(menu)

		_, err := fmt.Scanf("%d", &action)
		if err != nil {
			return nil, err
		}

		switch action {
		case 1:
			tasks, err = services.TaskService.GetAllTasks()
			err = AllTasks(services)
		case 2:
			category := ChooseTaskCategory()
			tasks, err = TasksByCategory(services, category)
			err = modelTables.Tasks(tasks)
		default:
			fmt.Println("Такого пункта в меню нету")
		}

		if err == nil {
			return tasks, nil
		}
	}
}
