package workerViews

import (
	"fmt"
	"lab3/cmd/modelTables"
	"lab3/internal/models"
	"lab3/internal/registry"
)

func assignWorker(services registry.Services, order *models.Order) error {
	workers, err := services.WorkerService.GetWorkersByRole(models.MasterRole)
	if err != nil {
		return err
	}

	err = modelTables.Workers(services, workers)
	if err != nil {
		return err
	}

	var workerNumber int
	for {
		fmt.Print("Введите номер работника, чтобы назначить его на заказ или 0, чтобы выйти\n")

		_, err = fmt.Scanf("%d", &workerNumber)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if workerNumber == 0 {
			return nil
		}

		if workerNumber < 1 || workerNumber > len(workers) {
			fmt.Println("Неверный номер")
			continue
		}

		order.WorkerID = workers[workerNumber-1].ID
		_, err = services.OrderService.Update(order.ID, order.Status, order.Rate, order.WorkerID)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Работник назначен")
			return nil
		}
	}

}
