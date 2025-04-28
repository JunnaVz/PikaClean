// Package workerViews provides user interface functions for the PikaClean application
// focused on worker-related operations like assignment to orders, order management,
// and worker profile management. It handles worker interactions through the command line
// interface for operations performed by cleaning staff and administrators.
package workerViews

import (
	"fmt"
	"teamdev/cmd/modelTables"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// assignWorker handles the process of assigning a worker to a cleaning order.
// It retrieves available workers with the Master role, displays them in a table,
// and prompts the administrator to select a worker to assign to the given order.
// The function updates the order with the selected worker's ID.
//
// Parameters:
//   - services: Service container providing access to business logic services,
//     particularly WorkerService and OrderService
//   - order: Order entity to which a worker will be assigned
//
// Returns:
//   - error: Any error that occurred during the assignment process,
//     such as database errors or display errors
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
