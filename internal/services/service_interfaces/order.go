// Package service_interfaces defines the contract interfaces for the business logic layer
// of the PikaClean application. These interfaces abstract the implementation details
// and provide a clear API for the application's services.
package service_interfaces

import (
	"github.com/google/uuid"
	"teamdev/internal/models"
	"time"
)

// IOrderService defines the contract for order management operations.
// This interface provides methods for creating, updating, and managing cleaning orders,
// including task assignment, pricing calculations, and order status management.
type IOrderService interface {
	// CreateOrder creates a new cleaning service order for a customer.
	//
	// Parameters:
	//   - userID: UUID of the customer placing the order
	//   - address: Location where the cleaning service should be performed
	//   - deadline: When the order should be completed by
	//   - orderedTasks: Slice of tasks with their quantities to be included in the order
	//
	// Returns:
	//   - *models.Order: Created order with assigned ID and initial status
	//   - error: Error if creation fails or validation fails
	CreateOrder(userID uuid.UUID, address string, deadline time.Time, orderedTasks []models.OrderedTask) (*models.Order, error)

	// DeleteOrder removes an order and its associated task relationships.
	//
	// Parameters:
	//   - id: UUID of the order to delete
	//
	// Returns:
	//   - error: Error if deletion fails
	DeleteOrder(id uuid.UUID) error

	// GetTasksInOrder retrieves all cleaning tasks associated with a specific order.
	//
	// Parameters:
	//   - orderID: UUID of the order to retrieve tasks for
	//
	// Returns:
	//   - []models.Task: Slice of task entities in the specified order
	//   - error: Error if retrieval fails
	GetTasksInOrder(orderID uuid.UUID) ([]models.Task, error)

	// GetOrderByID retrieves an order by its unique identifier.
	//
	// Parameters:
	//   - id: UUID of the order to retrieve
	//
	// Returns:
	//   - *models.Order: Retrieved order entity
	//   - error: Error if retrieval fails or order not found
	GetOrderByID(id uuid.UUID) (*models.Order, error)

	// GetCurrentOrderByUserID retrieves the most recent order for a specific user.
	//
	// Parameters:
	//   - userID: UUID of the user to retrieve the current order for
	//
	// Returns:
	//   - *models.Order: Retrieved order entity
	//   - error: Error if retrieval fails or no orders exist
	GetCurrentOrderByUserID(userID uuid.UUID) (*models.Order, error)

	// GetAllOrdersByUserID retrieves all orders for a specific user.
	//
	// Parameters:
	//   - userID: UUID of the user to retrieve orders for
	//
	// Returns:
	//   - []models.Order: Slice of order entities for the specified user
	//   - error: Error if retrieval fails
	GetAllOrdersByUserID(userID uuid.UUID) ([]models.Order, error)

	// Update modifies an existing order's status, rating, or worker assignment.
	//
	// Parameters:
	//   - orderID: UUID of the order to update
	//   - status: New status code for the order
	//   - rate: Customer satisfaction rating (0-5)
	//   - workerID: UUID of the worker to assign to the order
	//
	// Returns:
	//   - *models.Order: Updated order data
	//   - error: Error if update fails or validation fails
	Update(orderID uuid.UUID, status int, rate int, workerID uuid.UUID) (*models.Order, error)

	// AddTask associates a new task with an existing order.
	//
	// Parameters:
	//   - orderID: UUID of the order
	//   - tasksID: UUID of the task to add
	//
	// Returns:
	//   - error: Error if addition fails or task is already in the order
	AddTask(orderID uuid.UUID, tasksID uuid.UUID) error

	// RemoveTask removes a task association from an order.
	//
	// Parameters:
	//   - orderID: UUID of the order
	//   - taskID: UUID of the task to remove
	//
	// Returns:
	//   - error: Error if removal fails or task is not in the order
	RemoveTask(orderID uuid.UUID, taskID uuid.UUID) error

	// IncrementTaskQuantity increases the quantity of a specific task in an order by one.
	//
	// Parameters:
	//   - id: UUID of the order
	//   - taskID: UUID of the task
	//
	// Returns:
	//   - int: New quantity after increment
	//   - error: Error if update fails or task is not in the order
	IncrementTaskQuantity(id uuid.UUID, taskID uuid.UUID) (int, error)

	// DecrementTaskQuantity decreases the quantity of a specific task in an order by one.
	//
	// Parameters:
	//   - id: UUID of the order
	//   - taskID: UUID of the task
	//
	// Returns:
	//   - int: New quantity after decrement
	//   - error: Error if update fails, task is not in the order, or quantity would become negative
	DecrementTaskQuantity(id uuid.UUID, taskID uuid.UUID) (int, error)

	// SetTaskQuantity sets the exact quantity of a specific task in an order.
	//
	// Parameters:
	//   - id: UUID of the order
	//   - taskID: UUID of the task
	//   - quantity: New quantity value
	//
	// Returns:
	//   - error: Error if update fails, task is not in the order, or quantity is invalid
	SetTaskQuantity(id uuid.UUID, taskID uuid.UUID, quantity int) error

	// GetTaskQuantity retrieves the quantity of a specific task in an order.
	//
	// Parameters:
	//   - orderID: UUID of the order
	//   - taskID: UUID of the task
	//
	// Returns:
	//   - int: Quantity of the task in the order
	//   - error: Error if retrieval fails or task is not in the order
	GetTaskQuantity(orderID uuid.UUID, taskID uuid.UUID) (int, error)

	// Filter retrieves orders matching the specified criteria.
	//
	// Parameters:
	//   - params: Map of field names to filter values
	//
	// Returns:
	//   - []models.Order: Slice of order entities matching the filter criteria
	//   - error: Error if filtering fails
	Filter(params map[string]string) ([]models.Order, error)

	// GetTotalPrice calculates the total price for an order based on tasks and quantities.
	//
	// Parameters:
	//   - orderID: UUID of the order to calculate price for
	//
	// Returns:
	//   - float64: Total price of the order
	//   - error: Error if calculation fails
	GetTotalPrice(orderID uuid.UUID) (float64, error)
}
