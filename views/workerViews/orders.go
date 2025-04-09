package workerViews

import (
	"fmt"
	"lab3/cmd/modelTables"
	"lab3/cmd/views/orderViews"
	"lab3/internal/models"
	"lab3/internal/registry"
)

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

func validateOrderNumber(orderNumber int, orders []models.Order) bool {
	orderNumber--
	return orderNumber >= 0 && orderNumber < len(orders)
}

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
