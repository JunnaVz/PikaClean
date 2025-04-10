package service_interfaces

import (
	"github.com/google/uuid"
	"lab3/internal/models"
	"time"
)

type IOrderService interface {
	CreateOrder(userID uuid.UUID, address string, deadline time.Time, orderedTasks []models.OrderedTask) (*models.Order, error)
	DeleteOrder(id uuid.UUID) error
	GetTasksInOrder(orderID uuid.UUID) ([]models.Task, error)
	GetOrderByID(id uuid.UUID) (*models.Order, error)
	GetCurrentOrderByUserID(userID uuid.UUID) (*models.Order, error)
	GetAllOrdersByUserID(userID uuid.UUID) ([]models.Order, error)

	Update(orderID uuid.UUID, status int, rate int, workerID uuid.UUID) (*models.Order, error)

	AddTask(orderID uuid.UUID, tasksID uuid.UUID) error
	RemoveTask(orderID uuid.UUID, taskID uuid.UUID) error

	IncrementTaskQuantity(id uuid.UUID, taskID uuid.UUID) (int, error)
	DecrementTaskQuantity(id uuid.UUID, taskID uuid.UUID) (int, error)
	SetTaskQuantity(id uuid.UUID, taskID uuid.UUID, quantity int) error
	GetTaskQuantity(orderID uuid.UUID, taskID uuid.UUID) (int, error)

	Filter(params map[string]string) ([]models.Order, error)
	GetTotalPrice(orderID uuid.UUID) (float64, error)
}
