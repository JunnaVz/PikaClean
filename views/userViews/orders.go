package userViews

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
