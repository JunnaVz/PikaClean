// Package repository_interfaces defines the contract interfaces for data persistence
// operations. It specifies the required functionality that any repository implementation
// must satisfy to interact with the application's domain models.
package repository_interfaces

import (
	"github.com/google/uuid"
	"teamdev/internal/models"
)

// IOrderRepository defines the contract for order data persistence operations.
// Any implementation of this interface must provide methods for creating, retrieving,
// updating, and deleting order entities, as well as managing order-task relationships.
type IOrderRepository interface {
	// Create adds a new order record to the data store along with its associated tasks.
	//
	// Parameters:
	//   - order: Order entity to be persisted
	//   - orderedTasks: Slice of tasks associated with the order and their quantities
	//
	// Returns:
	//   - *models.Order: Created order with assigned ID
	//   - error: Error if creation fails
	Create(order *models.Order, orderedTasks []models.OrderedTask) (*models.Order, error)

	// Delete removes an order record and all associated task relationships from the data store.
	//
	// Parameters:
	//   - id: UUID of the order to delete
	//
	// Returns:
	//   - error: Error if deletion fails
	Delete(id uuid.UUID) error

	// Update modifies an existing order record in the data store.
	//
	// Parameters:
	//   - order: Order entity with updated values
	//
	// Returns:
	//   - *models.Order: Updated order data
	//   - error: Error if update fails
	Update(order *models.Order) (*models.Order, error)

	// GetOrderByID retrieves an order by unique identifier.
	//
	// Parameters:
	//   - id: UUID of the order to retrieve
	//
	// Returns:
	//   - *models.Order: Retrieved order entity
	//   - error: Error if retrieval fails or order not found
	GetOrderByID(id uuid.UUID) (*models.Order, error)

	// GetTasksInOrder retrieves all tasks associated with a specific order.
	//
	// Parameters:
	//   - id: UUID of the order to retrieve tasks for
	//
	// Returns:
	//   - []models.Task: Slice of task entities associated with the order
	//   - error: Error if retrieval fails
	GetTasksInOrder(id uuid.UUID) ([]models.Task, error)

	// GetCurrentOrderByUserID retrieves the most recent order for a specific user.
	//
	// Parameters:
	//   - id: UUID of the user to retrieve the current order for
	//
	// Returns:
	//   - *models.Order: Retrieved order entity
	//   - error: Error if retrieval fails or no orders found
	GetCurrentOrderByUserID(id uuid.UUID) (*models.Order, error)

	// GetAllOrdersByUserID retrieves all orders for a specific user.
	//
	// Parameters:
	//   - id: UUID of the user to retrieve orders for
	//
	// Returns:
	//   - []models.Order: Slice of order entities for the specified user
	//   - error: Error if retrieval fails
	GetAllOrdersByUserID(id uuid.UUID) ([]models.Order, error)

	// AddTaskToOrder associates a task with an order.
	//
	// Parameters:
	//   - orderID: UUID of the order
	//   - taskID: UUID of the task to add
	//
	// Returns:
	//   - error: Error if association fails
	AddTaskToOrder(orderID uuid.UUID, taskID uuid.UUID) error

	// RemoveTaskFromOrder removes a task association from an order.
	//
	// Parameters:
	//   - orderID: UUID of the order
	//   - taskID: UUID of the task to remove
	//
	// Returns:
	//   - error: Error if removal fails
	RemoveTaskFromOrder(orderID uuid.UUID, taskID uuid.UUID) error

	// UpdateTaskQuantity updates the quantity of a specific task in an order.
	//
	// Parameters:
	//   - orderID: UUID of the order
	//   - taskID: UUID of the task
	//   - quantity: New quantity value
	//
	// Returns:
	//   - error: Error if update fails
	UpdateTaskQuantity(orderID uuid.UUID, taskID uuid.UUID, quantity int) error

	// GetTaskQuantity retrieves the quantity of a specific task in an order.
	//
	// Parameters:
	//   - orderID: UUID of the order
	//   - taskID: UUID of the task
	//
	// Returns:
	//   - int: Quantity of the task in the order
	//   - error: Error if retrieval fails
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
}
