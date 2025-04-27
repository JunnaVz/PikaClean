// Package orderViews provides view functions for managing cleaning service orders
// in the PikaClean application. It contains functions that encapsulate business
// logic for common order-related operations and UI workflows.
package orderViews

import (
	"fmt"
	"teamdev/cmd/modelTables"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// GetUnassignedOrder displays details of an order that hasn't been assigned to a worker yet
// and provides options for managing it. It shows all tasks in the order and presents
// a menu for canceling the order or assigning a worker to it.
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - order: The unassigned order to be viewed and potentially modified
//
// Returns:
//   - error: Any error that occurred during task retrieval or order modification,
//     or nil if the operation was successful
func GetUnassignedOrder(services registry.Services, order *models.Order) error {
	tasks, err := services.OrderService.GetTasksInOrder(order.ID)
	if err != nil {
		return err
	}

	fmt.Printf("\nУслуги в заказе:\n")
	for i, task := range tasks {
		taskAmount, _ := services.OrderService.GetTaskQuantity(order.ID, task.ID)
		fmt.Printf("%d.\t%s\t%d\n", i+1, task.Name, taskAmount)
	}

	fmt.Printf("\n-----------\n1 -- отменить заказ\n2 -- назначит работника\n0 -- выход\n\n")

	for {
		var action int
		_, err = fmt.Scanf("%d", &action)
		if err != nil {
			fmt.Println(err)
		}

		if action == 0 {
			return nil
		}

		if action == 1 {
			return CancelOrder(services, order)
		} else if action == 2 {
			return assignWorker(services, order)
		}
	}
}

// assignWorker handles the worker assignment process for an unassigned order.
// It displays available workers with the Master role, allows selecting one by number,
// and updates the order with the selected worker ID.
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - order: The order to which a worker should be assigned
//
// Returns:
//   - error: Any error that occurred during worker retrieval or order update,
//     or nil if the operation was successful
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
