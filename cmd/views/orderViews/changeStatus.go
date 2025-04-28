// Package orderViews provides view functions for managing cleaning service orders
// in the PikaClean application. It contains functions that encapsulate business
// logic for common order-related operations and UI workflows.
package orderViews

import (
	"fmt"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// OrderMenuChangeStatus displays a menu for reviewing order tasks and changing the order status.
// It lists all tasks in the order with their quantities, and provides options for status management.
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - order: The order to be viewed and potentially modified
//
// Returns:
//   - error: Any error that occurred during task retrieval or status update,
//     or nil if the operation was successful
func OrderMenuChangeStatus(services registry.Services, order *models.Order) error {
	tasks, err := services.OrderService.GetTasksInOrder(order.ID)
	if err != nil {
		return err
	}

	fmt.Printf("\nУслуги в заказе:\n")
	for i, task := range tasks {
		taskAmount, _ := services.OrderService.GetTaskQuantity(order.ID, task.ID)
		fmt.Printf("%d.\t%s\t%d\n", i+1, task.Name, taskAmount)
	}

	fmt.Printf("\n-----------\n1 -- изменить статус заказа\n0 -- выход\n\n")

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
			return changeStatus(services, order)
		}
	}
}

// changeStatus handles the actual status change operation after a user selects this option.
// It displays available status options, validates user input, and updates the order's status.
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - order: The order whose status should be changed
//
// Returns:
//   - error: Any error that occurred during input processing or status update,
//     or nil if the operation was successful
func changeStatus(services registry.Services, order *models.Order) error {
	fmt.Print("Введите новый статус заказа:\n2 -- в работе\n3 -- выполнен\n0 -- выход\n\n")

	var newStatus int
	_, err := fmt.Scanf("%d", &newStatus)
	if err != nil {
		return err
	}

	if newStatus == 0 {
		return nil
	}

	if newStatus < 2 || newStatus > 3 {
		fmt.Println("Неверный статус заказа")
		return nil
	}

	_, err = services.OrderService.Update(order.ID, newStatus, order.Rate, order.WorkerID)
	if err != nil {
		return err
	}

	return nil
}
