// Package userViews provides user interface functions for the PikaClean application
// focused on customer-related operations like order creation, viewing order history,
// and managing user profiles. It handles user interactions through the command line
// interface for customer-focused operations.
package userViews

import (
	"fmt"
	utils "teamdev/cmd/cmdUtils"
	"teamdev/cmd/views/taskViews"
	"teamdev/internal/models"
	"teamdev/internal/registry"
	"time"
)

// orderSumPrice calculates the total price of an order by summing the cost of all
// ordered tasks, taking into account their quantities and individual prices.
//
// Parameters:
//   - orderedTasks: Slice of ordered tasks with their quantities
//
// Returns:
//   - float64: Total price of the order in the local currency
func orderSumPrice(orderedTasks []models.OrderedTask) float64 {
	var sum float64
	for _, task := range orderedTasks {
		sum += task.Task.PricePerSingle * float64(task.Quantity)
	}
	return sum
}

// addTaskToCart adds a task to the order cart, increasing quantity if the task
// already exists or appending it as a new item if not. This function ensures
// that duplicate tasks are consolidated with increased quantities rather than
// appearing multiple times in the cart.
//
// Parameters:
//   - orderedTask: Task to be added with its quantity
//   - tasks: Current list of tasks in the cart
//
// Returns:
//   - []models.OrderedTask: Updated list of tasks in the cart
func addTaskToCart(orderedTask models.OrderedTask, tasks []models.OrderedTask) []models.OrderedTask {
	for i, task := range tasks {
		if task.Task.Name == orderedTask.Task.Name {
			tasks[i].Quantity += orderedTask.Quantity
			return tasks
		}
	}

	return append(tasks, orderedTask)
}

// createOrder guides users through the process of creating a new cleaning order.
// It collects the delivery address, deadline, and allows selection of multiple
// cleaning tasks with quantities. The function validates input at each step
// and displays an order summary upon successful creation.
//
// Parameters:
//   - service: Service container providing access to business logic services
//   - user: Current authenticated user creating the order
//
// Returns:
//   - error: Any error that occurred during order creation,
//     or nil if the operation was successful
func createOrder(service registry.Services, user *models.User) error {
	var yesno string
	var address string

	fmt.Printf("Адрес заказа совпадает с Вашим?: (y/n) ")
	fmt.Scanf("%s", &yesno)
	if yesno == "n" {
		address = utils.EndlessReadWord("Введите адрес заказа: ")
	} else {
		address = user.Address
	}

	var err error
	const dateLayout = "2006-01-02"
	var deadline time.Time
	for {
		deadline, err = time.Parse(dateLayout, utils.EndlessReadWord("Введите крайний срок выполнения: (yyyy-mm-dd) "))
		if err != nil {
			fmt.Println("Неверный формат даты")
		} else {
			break
		}
	}

	var tasks []models.Task
	var orderedTasks []models.OrderedTask

	tasks, err = taskViews.Tasks(service)
	if err != nil {
		return err
	}

	fmt.Printf("Выберите услуги для заказа. Чтобы остановить выбор, введите пустую строку. Введите <, чтобы выбрать другую категорию.\n")
	for {
		var taskNum int
		var amount int

		fmt.Println("Введите номер услуги и количество через пробел (1 1): ")
		_, err = fmt.Scanf("%d %d", &taskNum, &amount)
		if err != nil {
			if err.Error() == "unexpected newline" {
				if len(orderedTasks) == 0 {
					fmt.Println("Не выбрано ни одной услуги")
					continue
				}
				break
			} else if err.Error() == "expected integer" {
				fmt.Scanln()
				tasks, err = taskViews.Tasks(service)
				if err != nil {
					return err
				}
			}
			continue
		}

		if taskNum < 1 || taskNum > len(tasks) {
			fmt.Println("Неверный номер услуги")
			continue
		}

		orderedTasks = addTaskToCart(models.OrderedTask{Task: &tasks[taskNum-1], Quantity: amount}, orderedTasks)
	}

	_, err = service.OrderService.CreateOrder(user.ID, address, deadline, orderedTasks)

	if err == nil {
		fmt.Println("Заказ успешно создан\nДобавлены следующие услуги:")
		for i, task := range orderedTasks {
			fmt.Printf("%d. %s %d\n", i+1, task.Task.Name, task.Quantity)
		}
		fmt.Printf("Адрес: %s\nКрайний срок: %s\n", address, deadline.Format(dateLayout))
		fmt.Printf("Стоимость заказа: %.2f рублей\n", orderSumPrice(orderedTasks))
		fmt.Printf("Ожидайте звонка оператора\n-------------------\n")
	}

	return err
}
