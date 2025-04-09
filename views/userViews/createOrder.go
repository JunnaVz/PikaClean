package userViews

import (
	"fmt"
	utils "lab3/cmd/cmdUtils"
	"lab3/cmd/views/taskViews"
	"lab3/internal/models"
	"lab3/internal/registry"
	"time"
)

func orderSumPrice(orderedTasks []models.OrderedTask) float64 {
	var sum float64
	for _, task := range orderedTasks {
		sum += task.Task.PricePerSingle * float64(task.Quantity)
	}
	return sum
}

func addTaskToCart(orderedTask models.OrderedTask, tasks []models.OrderedTask) []models.OrderedTask {
	for i, task := range tasks {
		if task.Task.Name == orderedTask.Task.Name {
			tasks[i].Quantity += orderedTask.Quantity
			return tasks
		}
	}

	return append(tasks, orderedTask)
}

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
		_, err = fmt.Scanf("%d %d", &taskNum, &amount) // Использование адресов переменных
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
