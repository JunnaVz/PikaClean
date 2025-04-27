// Package orderViews provides view functions for managing cleaning service orders
// in the PikaClean application. It contains functions that encapsulate business
// logic for common order-related operations and UI workflows.
package orderViews

import (
	"fmt"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// GetTasksInOrder displays all tasks included in a specific order along with their quantities.
// It retrieves tasks from the order service and displays them in a formatted list.
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - order: The order whose tasks should be displayed
//
// Returns:
//   - error: Any error that occurred during task retrieval,
//     or nil if the operation was successful
func GetTasksInOrder(services registry.Services, order *models.Order) error {
	tasks, err := services.OrderService.GetTasksInOrder(order.ID)
	if err != nil {
		return err
	}

	fmt.Printf("\nУслуги в заказе:\n")
	for i, task := range tasks {
		taskAmount, _ := services.OrderService.GetTaskQuantity(order.ID, task.ID)
		fmt.Printf("%d.\t%s\t%d\n", i+1, task.Name, taskAmount)
	}

	return nil
}
