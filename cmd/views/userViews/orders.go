// Package userViews provides user interface functions for the PikaClean application
// focused on customer-related operations like order creation, viewing order history,
// and managing user profiles. It handles user interactions through the command line
// interface for customer-focused operations.
package userViews

import (
	"fmt"
	"teamdev/cmd/modelTables"
	"teamdev/cmd/views/orderViews"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// getOrderNumber reads an order number from the command line input.
// It continues to prompt until valid integer input is received.
//
// Returns:
//   - int: The order number entered by the user
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

// validateOrderNumber checks if the provided order number is valid
// within the range of available orders.
//
// Parameters:
//   - orderNumber: The order number to validate (1-based indexing from user perspective)
//   - orders: The slice of orders to validate against
//
// Returns:
//   - bool: true if the order number is valid, false otherwise
func validateOrderNumber(orderNumber int, orders []models.Order) bool {
	orderNumber--
	return orderNumber >= 0 && orderNumber < len(orders)
}

// getCompletedOrders displays a list of completed orders for the current user
// and allows them to rate these orders. Only orders with status 3 (completed)
// are displayed. The user can select an order by number to change its rating.
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - user: Current authenticated user whose completed orders will be displayed
//
// Returns:
//   - error: Any error that occurred during the operation,
//     or nil if the operation was successful
func getCompletedOrders(services registry.Services, user *models.User) error {
	params := map[string]string{
		"status":  "3",
		"user_id": user.ID.String(),
	}

	orders, err := services.OrderService.Filter(params)

	if err != nil {
		return err
	}

	err = modelTables.Orders(orders)
	if err != nil {
		return err
	}

	for {
		fmt.Printf("\n-----------\n" +
			"Введите номер заказа, чтобы изменить оценку его оценку\n" +
			"Введите 0, чтобы выйти\n\n")
		var orderNumber = getOrderNumber()

		if orderNumber == 0 {
			return nil
		}

		if !validateOrderNumber(orderNumber, orders) {
			fmt.Println("Неверный номер заказа")
			continue
		}

		err = rateOrder(services, &orders[orderNumber-1])
		if err != nil {
			return err
		}
	}
}

// getOrdersInWork displays a list of orders currently in progress for the current user
// and allows them to view order details or cancel orders. Only orders with status 1 or 2
// (in progress) are displayed. The user can select an order by number to view its tasks
// or cancel it.
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - user: Current authenticated user whose in-progress orders will be displayed
//
// Returns:
//   - error: Any error that occurred during the operation,
//     or nil if the operation was successful
func getOrdersInWork(services registry.Services, user *models.User) error {
	params := map[string]string{
		"status":  "1,2",
		"user_id": user.ID.String(),
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

	for {
		var orderNumber int
		_, err = fmt.Scanf("%d", &orderNumber)
		if err != nil {
			return err
		}

		if orderNumber == 0 {
			return nil
		}

		if !validateOrderNumber(orderNumber, orders) {
			fmt.Println("Неверный номер заказа")
			continue
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
}

// rateOrder allows a user to provide a satisfaction rating for a completed order.
// The rating is a numeric value that's stored with the order and can be used
// to evaluate worker performance.
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - order: The order to be rated
//
// Returns:
//   - error: Any error that occurred during the operation,
//     or nil if the operation was successful
func rateOrder(services registry.Services, order *models.Order) error {
	fmt.Printf("Введите оценку заказа: ")
	var rate int
	_, err := fmt.Scanf("%d", &rate)
	if err != nil {
		return err
	}

	order.Rate = rate
	_, err = services.OrderService.Update(order.ID, order.Status, rate, order.WorkerID)
	if err != nil {
		return err
	}

	return nil
}
