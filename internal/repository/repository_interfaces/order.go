package repository_interfaces

import (
	"github.com/google/uuid"
	"lab3/internal/models"
)

type IOrderRepository interface {
	Create(order *models.Order, orderedTasks []models.OrderedTask) (*models.Order, error)
	Delete(id uuid.UUID) error
	Update(order *models.Order) (*models.Order, error)
	GetOrderByID(id uuid.UUID) (*models.Order, error)
	GetTasksInOrder(id uuid.UUID) ([]models.Task, error)
	GetCurrentOrderByUserID(id uuid.UUID) (*models.Order, error)
	GetAllOrdersByUserID(id uuid.UUID) ([]models.Order, error)
	AddTaskToOrder(orderID uuid.UUID, taskID uuid.UUID) error
	RemoveTaskFromOrder(orderID uuid.UUID, taskID uuid.UUID) error
	UpdateTaskQuantity(orderID uuid.UUID, taskID uuid.UUID, quantity int) error
	GetTaskQuantity(orderID uuid.UUID, taskID uuid.UUID) (int, error)
	Filter(params map[string]string) ([]models.Order, error)
}
