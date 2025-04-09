package orderViews

import (
	"fmt"
	"lab3/internal/models"
	"lab3/internal/registry"
)

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
