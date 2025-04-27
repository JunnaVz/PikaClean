// Package workerViews provides user interface functions for the PikaClean application
// focused on worker-related operations like order assignment, order management,
// and administrative worker functions. It handles worker interactions through the command line
// interface for operations performed by cleaning staff and administrators.
package workerViews

import (
	"fmt"
	"teamdev/cmd/modelTables"
	"teamdev/cmd/views/orderViews"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// getOrderNumber reads an order number from standard input.
// It repeatedly prompts until a valid integer is entered.
//
// Returns:
//   - int: The valid order number entered by the user
func getOrderNumber() int {
	var orderNumber int
	for {
		_, err := fmt.Scanf("%d", &orderNumber)
		if err == nil {
			return orderNumber
		}
		fmt.Println(err)
	}
}

// validateOrderNumber checks if the entered order number is within
// the valid range for the available orders.
//
// Parameters:
//   - orderNumber: The order number entered by the user (1-based index)
//   - orders: Slice of available orders to validate against
//
// Returns:
//   - bool: True if the order number is valid, false otherwise
func validateOrderNumber(orderNumber int, orders []models.Order) bool {
	orderNumber--
	return orderNumber >= 0 && orderNumber < len(orders)
}

// unassignedOrders displays all unassigned orders and allows
// a manager to assign workers to them.
//
// Parameters:
//   - services: Service container providing access to business logic services
//
// Returns:
//   - error: Any error that occurred during operation
func unassignedOrders(services registry.Services) error {
	params := map[string]string{
		"worker_id": "null",
	}

	orders, err := services.OrderService.Filter(params)

	if err != nil {
		return err
	}

	err = modelTables.Orders(orders)
	if err != nil {
		return err
	}

	fmt.Printf("\n-----------\n" +
		"Введите номер заказа, чтобы назначить работника\n" +
		"Введите 0, чтобы выйти\n\n")

	orderNumber := getOrderNumber()

	if orderNumber == 0 {
		return nil
	}

	if !validateOrderNumber(orderNumber, orders) {
		fmt.Println("Неверный номер заказа")
		return nil
	}

	return orderViews.GetUnassignedOrder(services, &orders[orderNumber-1])
}

// completedOrders displays all completed orders and allows
// viewing their details.
//
// Parameters:
//   - services: Service container providing access to business logic services
//
// Returns:
//   - error: Any error that occurred during operation
func completedOrders(services registry.Services) error {
	params := map[string]string{
		"status": "3",
	}

	orders, err := services.OrderService.Filter(params)

	if err != nil {
		return err
	}

	err = modelTables.Orders(orders)
	if err != nil {
		return err
	}

	fmt.Printf("\n-----------\n" +
		"Введите номер заказа, чтобы просмотреть его содержимое\n" +
		"Введите 0, чтобы выйти\n\n")

	orderNumber := getOrderNumber()

	if orderNumber == 0 {
		return nil
	}

	if !validateOrderNumber(orderNumber, orders) {
		fmt.Println("Неверный номер заказа")
		return nil
	}

	err = orderViews.GetTasksInOrder(services, &orders[orderNumber-1])
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Нажмите Enter, чтобы продолжить")
	fmt.Scanln()

	return nil
}

// inProgressOrders displays all in-progress orders (status 1 or 2)
// and allows viewing their details and cancellation.
//
// Parameters:
//   - services: Service container providing access to business logic services
//
// Returns:
//   - error: Any error that occurred during operation
func inProgressOrders(services registry.Services) error {
	params := map[string]string{
		"status": "1,2",
	}

	orders, err := services.OrderService.Filter(params)

	if err != nil {
		return err
	}

	err = modelTables.Orders(orders)
	if err != nil {
		return err
	}

	fmt.Printf("\n-----------\n" +
		"Введите номер заказа, чтобы просмотреть его содержимое\n" +
		"Введите 0, чтобы выйти\n\n")

	orderNumber := getOrderNumber()

	if orderNumber == 0 {
		return nil
	}

	if !validateOrderNumber(orderNumber, orders) {
		fmt.Println("Неверный номер заказа")
		return nil
	}

	err = orderViews.GetTasksInOrder(services, &orders[orderNumber-1])
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\n-----------\n" +
		"Введите 1, чтобы отменить заказ\n\n" +
		"Введите 0, чтобы выйти\n\n")

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
			err = orderViews.CancelOrder(services, &orders[orderNumber-1])
			if err != nil {
				return err
			}

			return nil
		}
	}
}

// completedOrdersByWorker displays all completed orders assigned
// to a specific worker and allows viewing their details.
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - worker: The worker whose completed orders should be displayed
//
// Returns:
//   - error: Any error that occurred during operation
func completedOrdersByWorker(services registry.Services, worker *models.Worker) error {
	params := map[string]string{
		"status":    "3",
		"worker_id": worker.ID.String(),
	}

	orders, err := services.OrderService.Filter(params)

	if err != nil {
		return err
	}

	err = modelTables.Orders(orders)
	if err != nil {
		return err
	}

	fmt.Printf("\n-----------\n" +
		"Введите номер заказа, чтобы просмотреть его содержимое\n" +
		"Введите 0, чтобы выйти\n\n")

	orderNumber := getOrderNumber()

	if orderNumber == 0 {
		return nil
	}

	if !validateOrderNumber(orderNumber, orders) {
		fmt.Println("Неверный номер заказа")
		return nil
	}

	err = orderViews.GetTasksInOrder(services, &orders[orderNumber-1])
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Нажмите Enter, чтобы продолжить")
	fmt.Scanln()

	return nil
}

// inProgressOrdersByWorker displays all in-progress orders assigned
// to a specific worker and allows changing their status.
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - worker: The worker whose in-progress orders should be displayed
//
// Returns:
//   - error: Any error that occurred during operation
func inProgressOrdersByWorker(services registry.Services, worker *models.Worker) error {
	params := map[string]string{
		"status":    "1,2",
		"worker_id": worker.ID.String(),
	}

	orders, err := services.OrderService.Filter(params)

	if err != nil {
		return err
	}

	err = modelTables.Orders(orders)
	if err != nil {
		return err
	}

	fmt.Printf("\n-----------\n" +
		"Введите номер заказа, чтобы просмотреть его содержимое\n" +
		"Введите 0, чтобы выйти\n\n")

	orderNumber := getOrderNumber()

	if orderNumber == 0 {
		return nil
	}

	if !validateOrderNumber(orderNumber, orders) {
		fmt.Println("Неверный номер заказа")
		return nil
	}

	return orderViews.OrderMenuChangeStatus(services, &orders[orderNumber-1])
}
