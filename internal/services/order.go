package interfaces

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"lab3/internal/models"
	"lab3/internal/repository/repository_errors"
	"lab3/internal/repository/repository_interfaces"
	"lab3/internal/services/service_interfaces"
	"time"
)

type OrderService struct {
	OrderRepository  repository_interfaces.IOrderRepository
	TaskRepository   repository_interfaces.ITaskRepository
	WorkerRepository repository_interfaces.IWorkerRepository
	UserRepository   repository_interfaces.IUserRepository
	logger           *log.Logger
}

func NewOrderService(orderRepository repository_interfaces.IOrderRepository, workerRepository repository_interfaces.IWorkerRepository, taskRepository repository_interfaces.ITaskRepository, userRepository repository_interfaces.IUserRepository, logger *log.Logger) service_interfaces.IOrderService {
	return &OrderService{
		OrderRepository:  orderRepository,
		TaskRepository:   taskRepository,
		WorkerRepository: workerRepository,
		UserRepository:   userRepository,
		logger:           logger,
	}
}

func orderIsCompleted(orderStatus int) bool {
	return orderStatus == models.CompletedOrderStatus || orderStatus == models.CancelledOrderStatus
}

func (o OrderService) checkTasksExistence(tasks []models.OrderedTask) (bool, error) {
	for _, task := range tasks {
		if task.Quantity <= 0 {
			o.logger.Error("SERVICE: Quantity is negative", "task", task)
			return false, fmt.Errorf("SERVICE: Quantity is negative")
		}

		_, err := o.TaskRepository.GetTaskByID(task.Task.ID)
		if errors.Is(err, repository_errors.DoesNotExist) {
			o.logger.Error("SERVICE: Task does not exist", "id", task.Task.ID)
			return false, fmt.Errorf("SERVICE: Task does not exist")
		} else if err != nil {
			o.logger.Error("SERVICE: GetTaskByID method failed", "id", task.Task.ID, "error", err)
			return false, err
		}
	}

	return true, nil
}

func (o OrderService) CreateOrder(userID uuid.UUID, address string, deadline time.Time, orderedTasks []models.OrderedTask) (*models.Order, error) {
	// checking if order is valid
	if !validAddress(address) || !validDeadline(deadline) || !validTasksNumber(orderedTasks) {
		o.logger.Error("SERVICE: Invalid input")
		return nil, fmt.Errorf("SERVICE: Invalid input")
	}

	if _, err := o.checkTasksExistence(orderedTasks); err != nil {
		o.logger.Error("SERVICE: CheckTasksExistence method failed", "orderedTasks", orderedTasks, "error", err)
		return nil, err
	}

	// checking if user exists
	_, err := o.UserRepository.GetUserByID(userID)
	if errors.Is(err, repository_errors.DoesNotExist) {
		o.logger.Error("SERVICE: User does not exist", "id", userID)
		return nil, fmt.Errorf("SERVICE: User does not exist")
	} else if err != nil {
		o.logger.Error("SERVICE: GetWorkerByID method failed", "id", userID, "error", err)
		return nil, err
	}

	// creating order
	var order = &models.Order{
		UserID:   userID,
		Status:   models.NewOrderStatus,
		Address:  address,
		Deadline: deadline,
	}

	order, err = o.OrderRepository.Create(order, orderedTasks)
	if err != nil {
		o.logger.Error("SERVICE: Create method failed", "order", order, "error", err)
		return nil, err
	}

	o.logger.Info("SERVICE: Successfully created order", "order", order)
	return order, nil
}

func (o OrderService) DeleteOrder(id uuid.UUID) error {
	order, err := o.OrderRepository.GetOrderByID(id)
	if err != nil {
		o.logger.Error("SERVICE: GetUserByID method failed", "id", id, "error", err)
		return err
	}

	tasksFromOrder, err := o.OrderRepository.GetTasksInOrder(order.ID)
	if err != nil {
		o.logger.Error("SERVICE: GetTasksInOrder method failed", "id", order.ID, "error", err)
		return err
	}

	for _, task := range tasksFromOrder {
		err = o.OrderRepository.RemoveTaskFromOrder(order.ID, task.ID)
		if err != nil {
			o.logger.Error("SERVICE: Delete method failed", "id", task.ID, "error", err)
			return err
		}
	}
	o.logger.Info("SERVICE: Successfully deleted tasks from order", "order_id", order.ID)

	err = o.OrderRepository.Delete(id)
	if err != nil {
		o.logger.Error("SERVICE: Delete method failed", "id", id, "error", err)
		return err
	}

	o.logger.Info("SERVICE: Successfully deleted order", "id", id)
	return nil
}

func (o OrderService) GetTasksInOrder(orderID uuid.UUID) ([]models.Task, error) {
	_, err := o.OrderRepository.GetOrderByID(orderID)
	if err != nil {
		o.logger.Error("SERVICE: GetOrderByID method failed", "id", orderID, "error", err)
		return nil, err
	}

	tasks, err := o.OrderRepository.GetTasksInOrder(orderID)
	if err != nil {
		o.logger.Error("SERVICE: GetTasksInOrder method failed", "order_id", orderID, "error", err)
		return nil, err
	}

	o.logger.Info("SERVICE: Successfully got tasks in order", "order_id", orderID)
	return tasks, nil
}

func (o OrderService) GetCurrentOrderByUserID(userID uuid.UUID) (*models.Order, error) {
	user, _ := o.UserRepository.GetUserByID(userID)
	if user == nil {
		o.logger.Error("SERVICE: GetUserByID method failed", "id", userID)
		return nil, fmt.Errorf("SERVICE: GetUserByID method failed")
	}

	order, err := o.OrderRepository.GetCurrentOrderByUserID(userID)
	if err != nil {
		o.logger.Error("SERVICE: GetCurrentOrderByUserID method failed", "id", userID, "error", err)
		return nil, err
	}

	o.logger.Info("SERVICE: Successfully got current order by user id", "user_id", userID)
	return order, nil
}

func (o OrderService) GetAllOrdersByUserID(userID uuid.UUID) ([]models.Order, error) {
	user, _ := o.UserRepository.GetUserByID(userID)
	if user == nil {
		o.logger.Error("SERVICE: GetUserByID method failed", "id", userID)
		return nil, fmt.Errorf("SERVICE: GetUserByID method failed")
	}

	orders, err := o.OrderRepository.GetAllOrdersByUserID(userID)
	if err != nil {
		o.logger.Error("SERVICE: GetAllOrdersByUserID method failed", "id", userID, "error", err)
		return nil, err
	}

	o.logger.Info("SERVICE: Successfully got all orders by user id", "user_id", userID)
	return orders, nil
}

func (o OrderService) Filter(params map[string]string) ([]models.Order, error) {
	orders, err := o.OrderRepository.Filter(params)
	if err != nil {
		o.logger.Error("SERVICE: Filter method failed", "params", params, "error", err)
		return nil, err
	}

	o.logger.Info("SERVICE: Successfully filtered orders", "params", params)
	return orders, nil
}

func (o OrderService) Update(orderID uuid.UUID, status int, rate int, workerID uuid.UUID) (*models.Order, error) {
	order, err := o.OrderRepository.GetOrderByID(orderID)
	if err != nil {
		o.logger.Error("SERVICE: GetOrderByID method failed", "id", orderID, "error", err)
		return nil, err
	}

	if workerID != uuid.Nil {
		_, err = o.WorkerRepository.GetWorkerByID(workerID)
		if err != nil {
			o.logger.Error("SERVICE: GetWorkerByID method failed", "id", workerID, "error", err)
			return nil, err
		}

		order.WorkerID = workerID
	} else {
		order.WorkerID = uuid.Nil
	}

	//for adding rate to a completed/cancelled order
	//if order.Rate == 0 && order.Status == models.CompletedOrderStatus || order.Status == models.CancelledOrderStatus {
	//	o.logger.Error("SERVICE: Order is already completed or cancelled", "order", order)
	//	return nil, fmt.Errorf("SERVICE: Order is already completed or cancelled")
	//}

	if !validStatus(status) {
		o.logger.Error("SERVICE: Invalid status", "status", status)
		return nil, fmt.Errorf("SERVICE: Invalid status")
	} else {
		order.Status = status
	}

	//for testing adding rate to an uncompleted order -> 0 = no status
	if !orderIsCompleted(status) && rate != 0 {
		o.logger.Error("SERVICE: Order is not completed", "order", order)
		return nil, fmt.Errorf("SERVICE: Order is not completed")
	}

	if !validRate(rate) {
		o.logger.Error("SERVICE: Rating is out of range", "rate", rate)
		return nil, fmt.Errorf("SERVICE: Rating is out of range")
	} else {
		order.Rate = rate
	}

	order, err = o.OrderRepository.Update(order)
	if err != nil {
		o.logger.Error("SERVICE: Update method failed", "order", order, "error", err)
		return nil, err
	}

	o.logger.Info("SERVICE: Successfully changed order status", "order_id", orderID, "status", status)
	return order, nil
}

func (o OrderService) AddTask(orderID uuid.UUID, taskID uuid.UUID) error {
	order, err := o.OrderRepository.GetOrderByID(orderID)
	if err != nil {
		o.logger.Error("SERVICE: GetOrderByID method failed", "id", orderID, "error", err)
		return err
	}

	attachedTasks, err := o.OrderRepository.GetTasksInOrder(order.ID)

	_, err = o.TaskRepository.GetTaskByID(taskID)
	if err != nil {
		o.logger.Error("SERVICE: GetTaskByID method failed", "id", taskID, "error", err)
		return err
	}

	if taskIsAttachedToOrder(taskID, attachedTasks) {
		o.logger.Error("SERVICE: Task is already attached to order", "order_id", orderID, "task_id", taskID)
		return fmt.Errorf("SERVICE: Task is already attached to order")
	}

	err = o.OrderRepository.AddTaskToOrder(order.ID, taskID)
	if err != nil {
		o.logger.Error("SERVICE: AddTaskToOrder method failed", "order_id", order.ID, "task_id", taskID, "error", err)
		return err
	}

	o.logger.Info("SERVICE: Successfully added tasks to order", "order_id", orderID)
	return nil
}

func (o OrderService) RemoveTask(orderID uuid.UUID, taskID uuid.UUID) error {
	order, err := o.OrderRepository.GetOrderByID(orderID)
	if err != nil {
		o.logger.Error("SERVICE: GetOrderByID method failed", "id", orderID, "error", err)
		return err
	}

	_, err = o.TaskRepository.GetTaskByID(taskID)
	if err != nil {
		o.logger.Error("SERVICE: GetTaskByID method failed", "id", taskID, "error", err)
		return err
	}

	attachedTasks, err := o.OrderRepository.GetTasksInOrder(order.ID)
	if err != nil {
		o.logger.Error("SERVICE: GetTasksInOrder method failed", "order_id", order.ID, "error", err)
		return err
	}

	if !taskIsAttachedToOrder(taskID, attachedTasks) {
		o.logger.Error("SERVICE: Task is not attached to order", "order_id", orderID, "task_id", taskID)
		return fmt.Errorf("SERVICE: Task is not attached to order")
	}

	// remove task from order
	err = o.OrderRepository.RemoveTaskFromOrder(order.ID, taskID)
	if err != nil {
		o.logger.Error("SERVICE: RemoveTaskFromOrder method failed", "order_id", order.ID, "task_id", taskID, "error", err)
		return err
	}

	o.logger.Info("SERVICE: Successfully removed task from order", "order_id", orderID, "task_id", taskID)
	return nil
}

func (o OrderService) GetOrderByID(id uuid.UUID) (*models.Order, error) {
	order, err := o.OrderRepository.GetOrderByID(id)
	if err != nil {
		o.logger.Error("SERVICE: GetOrderByID method failed", "id", id, "error", err)
		return nil, err
	}

	o.logger.Info("SERVICE: Successfully got order by id", "id", id)
	return order, nil
}

func (o OrderService) IncrementTaskQuantity(id uuid.UUID, taskID uuid.UUID) (int, error) {
	_, err := o.OrderRepository.GetOrderByID(id)
	if err != nil {
		o.logger.Error("SERVICE: GetOrderByID method failed", "id", id, "error", err)
		return 0, err
	}

	_, err = o.TaskRepository.GetTaskByID(taskID)
	if err != nil {
		o.logger.Error("SERVICE: GetTaskByID method failed", "id", taskID, "error", err)
		return 0, err
	}

	quantity, err := o.OrderRepository.GetTaskQuantity(id, taskID)
	if err != nil {
		o.logger.Error("SERVICE: GetTaskQuantity method failed", "order_id", id, "task_id", taskID, "error", err)
		return 0, err
	}

	quantity++
	err = o.OrderRepository.UpdateTaskQuantity(id, taskID, quantity)
	if err != nil {
		o.logger.Error("SERVICE: UpdateTaskQuantity method failed", "order_id", id, "task_id", taskID, "error", err)
		return 0, err
	}

	o.logger.Info("SERVICE: Successfully incremented task quantity", "order_id", id, "task_id", taskID, "quantity", quantity)
	return quantity, nil
}

func (o OrderService) DecrementTaskQuantity(id uuid.UUID, taskID uuid.UUID) (int, error) {
	_, err := o.OrderRepository.GetOrderByID(id)
	if err != nil {
		o.logger.Error("SERVICE: GetOrderByID method failed", "id", id, "error", err)
		return 0, err
	}

	_, err = o.TaskRepository.GetTaskByID(taskID)
	if err != nil {
		o.logger.Error("SERVICE: GetTaskByID method failed", "id", taskID, "error", err)
		return 0, err
	}

	quantity, err := o.OrderRepository.GetTaskQuantity(id, taskID)
	if err != nil {
		o.logger.Error("SERVICE: GetTaskQuantity method failed", "order_id", id, "task_id", taskID, "error", err)
		return 0, err
	}

	if quantity == 0 {
		o.logger.Error("SERVICE: Quantity is already 0", "order_id", id, "task_id", taskID)
		return 0, fmt.Errorf("SERVICE: Quantity is already 0")
	}

	quantity--
	err = o.OrderRepository.UpdateTaskQuantity(id, taskID, quantity)
	if err != nil {
		o.logger.Error("SERVICE: UpdateTaskQuantity method failed", "order_id", id, "task_id", taskID, "error", err)
		return 0, err
	}

	o.logger.Info("SERVICE: Successfully decremented task quantity", "order_id", id, "task_id", taskID, "quantity", quantity)
	return quantity, nil
}

func (o OrderService) SetTaskQuantity(id uuid.UUID, taskID uuid.UUID, quantity int) error {
	if quantity < 0 {
		o.logger.Error("SERVICE: Quantity is negative", "order_id", id, "task_id", taskID, "quantity", quantity)
		return fmt.Errorf("SERVICE: Quantity is negative")
	}

	_, err := o.OrderRepository.GetOrderByID(id)
	if err != nil {
		o.logger.Error("SERVICE: GetOrderByID method failed", "id", id, "error", err)
		return err
	}

	_, err = o.TaskRepository.GetTaskByID(taskID)
	if err != nil {
		o.logger.Error("SERVICE: GetTaskByID method failed", "id", taskID, "error", err)
		return err
	}

	err = o.OrderRepository.UpdateTaskQuantity(id, taskID, quantity)
	if err != nil {
		o.logger.Error("SERVICE: UpdateTaskQuantity method failed", "order_id", id, "task_id", taskID, "error", err)
		return err
	}

	o.logger.Info("SERVICE: Successfully set task quantity", "order_id", id, "task_id", taskID, "quantity", quantity)
	return nil
}

func (o OrderService) GetTaskQuantity(orderID uuid.UUID, taskID uuid.UUID) (int, error) {
	_, err := o.OrderRepository.GetOrderByID(orderID)
	if err != nil {
		o.logger.Error("SERVICE: GetOrderByID method failed", "id", orderID, "error", err)
		return 0, err
	}

	_, err = o.TaskRepository.GetTaskByID(taskID)
	if err != nil {
		o.logger.Error("SERVICE: GetTaskByID method failed", "id", taskID, "error", err)
		return 0, err
	}

	quantity, err := o.OrderRepository.GetTaskQuantity(orderID, taskID)
	if err != nil {
		o.logger.Error("SERVICE: GetTaskQuantity method failed", "order_id", orderID, "task_id", taskID, "error", err)
		return 0, err
	}

	o.logger.Info("SERVICE: Successfully got task quantity", "order_id", orderID, "task_id", taskID, "quantity", quantity)
	return quantity, nil
}

func (o OrderService) GetTotalPrice(orderID uuid.UUID) (float64, error) {
	orders, err := o.OrderRepository.GetTasksInOrder(orderID)
	if err != nil {
		o.logger.Error("SERVICE: GetTasksInOrder method failed", "order_id", orderID, "error", err)
		return 0, err
	}

	var sum float64 = 0
	for _, task := range orders {
		var quantity int
		quantity, err = o.OrderRepository.GetTaskQuantity(orderID, task.ID)
		if err != nil {
			o.logger.Error("SERVICE: GetTaskQuantity method failed", "order_id", orderID, "task_id", task.ID, "error", err)
			return 0, err
		}
		sum += task.PricePerSingle * float64(quantity)
	}

	o.logger.Info("SERVICE: Successfully got total price", "order_id", orderID, "total_price", sum)
	return sum, nil
}
