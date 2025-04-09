package taskViews

import (
	"fmt"
	"lab3/cmd/modelTables"
	"lab3/internal/models"
	"lab3/internal/registry"
)

func AllTasks(services registry.Services) error {
	tasks, err := services.TaskService.GetAllTasks()
	if err != nil {
		return err
	}
	return modelTables.Tasks(tasks)
}

func TasksByCategory(service registry.Services, category int) ([]models.Task, error) {
	tasks, err := service.TaskService.GetTasksInCategory(category)
	if err != nil {
		return nil, err
	}

	return tasks, err
}

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
