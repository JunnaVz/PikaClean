package workerViews

import (
	"fmt"
	"lab3/cmd/modelTables"
	"lab3/internal/models"
	"lab3/internal/registry"
)

func getAllWorkers(services registry.Services, manager *models.Worker) error {
	workers, err := services.WorkerService.GetAllWorkers()

	if err != nil {
		return err
	}

	err = modelTables.Workers(services, workers)
	if err != nil {
		return err
	}

	var action int
	for {
		fmt.Print("Введите номер работника, чтобы изменить его профиль или 0, чтобы выйти\n")

		_, err = fmt.Scanf("%d", &action)
		if err != nil {
			fmt.Println(err)
		}

		if action == 0 {
			return nil
		}

		if action < 1 || action > len(workers) {
			fmt.Println("Неверный номер")
		}

		err = Update(services, workers[action-1].ID, manager)

	}
}
