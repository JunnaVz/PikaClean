package workerViews

import (
	"fmt"
	"lab3/cmd/menu"
	"lab3/cmd/modelTables"
	"lab3/cmd/views/taskViews"
	"lab3/internal/models"
	"lab3/internal/registry"
)

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
