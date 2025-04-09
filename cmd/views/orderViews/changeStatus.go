package orderViews

import (
	"fmt"
	"lab3/internal/models"
	"lab3/internal/registry"
)

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
