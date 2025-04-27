// Package postgres provides repository implementations for data persistence
// using PostgreSQL database. It includes repositories for managing workers,
// users, tasks, orders, and categories.
package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"teamdev/internal/models"
	"teamdev/internal/repository/repository_errors"
	"teamdev/internal/repository/repository_interfaces"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// OrderDB represents an order entity as stored in the PostgreSQL database.
// It maps directly to the columns in the orders table.
type OrderDB struct {
	ID           uuid.UUID `db:"id"`            // Unique identifier for the order
	WorkerID     uuid.UUID `db:"worker_id"`     // ID of the worker assigned to the order
	UserID       uuid.UUID `db:"user_id"`       // ID of the user who created the order
	Status       int       `db:"status"`        // Current status of the order (numeric code)
	Address      string    `db:"address"`       // Location where the cleaning service should be performed
	CreationDate time.Time `db:"creation_date"` // When the order was created
	Deadline     time.Time `db:"deadline"`      // When the order should be completed
	Rate         int       `db:"rate"`          // Customer satisfaction rating (0-5)
}

// OrderRepository implements the IOrderRepository interface for PostgreSQL.
// It provides methods for creating, updating, and retrieving order records.
type OrderRepository struct {
	db *sqlx.DB // Database connection
}

// NewOrderRepository creates a new OrderRepository instance with the provided
// database connection.
//
// Parameters:
//   - db: An initialized sqlx.DB connection to PostgreSQL
//
// Returns:
//   - repository_interfaces.IOrderRepository: Repository implementation
func NewOrderRepository(db *sqlx.DB) repository_interfaces.IOrderRepository {
	return &OrderRepository{db: db}
}

// copyOrderResultToModel converts an OrderDB database entity to a models.Order domain entity.
//
// Parameters:
//   - orderDB: Database entity to convert
//
// Returns:
//   - *models.Order: Corresponding domain entity
func copyOrderResultToModel(orderDB *OrderDB) *models.Order {
	return &models.Order{
		ID:           orderDB.ID,
		WorkerID:     orderDB.WorkerID,
		UserID:       orderDB.UserID,
		Status:       orderDB.Status,
		Address:      orderDB.Address,
		CreationDate: orderDB.CreationDate,
		Deadline:     orderDB.Deadline,
		Rate:         orderDB.Rate,
	}
}

// Create inserts a new order record into the database along with its associated tasks.
// The operation is performed within a transaction to ensure data consistency.
//
// Parameters:
//   - order: Order entity to be created
//   - orderedTasks: Slice of tasks associated with the order and their quantities
//
// Returns:
//   - *models.Order: Created order with assigned ID
//   - error: repository_errors.TransactionBeginError, repository_errors.TransactionRollbackError,
//     repository_errors.InsertError, or repository_errors.TransactionCommitError if the operation fails
func (o OrderRepository) Create(order *models.Order, orderedTasks []models.OrderedTask) (*models.Order, error) {
	transaction, err := o.db.Begin()
	if err != nil {
		return nil, repository_errors.TransactionBeginError
	}

	query := `INSERT INTO orders(user_id, status, address, deadline) VALUES ($1, $2, $3, $4) RETURNING id;`

	err = transaction.QueryRow(query, order.UserID, order.Status, order.Address, order.Deadline).Scan(&order.ID)

	if err != nil {
		err = transaction.Rollback()
		if err != nil {
			return nil, repository_errors.TransactionRollbackError
		}
		return nil, repository_errors.InsertError
	}

	for _, task := range orderedTasks {
		query = `INSERT INTO order_contains_tasks(order_id, task_id, quantity) VALUES ($1, $2, $3);`
		_, err = transaction.Exec(query, order.ID, task.Task.ID, task.Quantity)
		if err != nil {
			err = transaction.Rollback()
			if err != nil {
				return nil, repository_errors.TransactionRollbackError
			}
			return nil, repository_errors.InsertError
		}
	}

	err = transaction.Commit()
	if err != nil {
		return nil, repository_errors.TransactionCommitError
	}

	return order, nil
}

// Delete removes an order record and all associated task relationships from the database.
// The operation is performed within a transaction to ensure data consistency.
//
// Parameters:
//   - id: UUID of the order to delete
//
// Returns:
//   - error: repository_errors.TransactionBeginError, repository_errors.TransactionRollbackError,
//     repository_errors.DeleteError, repository_errors.TransactionCommitError,
//     or a custom error if no order was found to delete
func (o OrderRepository) Delete(id uuid.UUID) error {
	// Start a new transaction
	tx, err := o.db.Begin()
	if err != nil {
		return repository_errors.TransactionBeginError
	}

	// Delete the records in the order_contains_tasks table that reference the order
	_, err = tx.Exec(`DELETE FROM order_contains_tasks WHERE order_id = $1;`, id)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return repository_errors.TransactionRollbackError
		}
		return repository_errors.DeleteError
	}

	// Delete the order
	result, err := tx.Exec(`DELETE FROM orders WHERE id = $1;`, id)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return repository_errors.TransactionRollbackError
		}
		return repository_errors.DeleteError
	}

	// Check if the order was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return repository_errors.TransactionRollbackError
		}
		return repository_errors.DeleteError
	}

	if rowsAffected == 0 {
		return errors.New("no order found to delete")
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return repository_errors.TransactionCommitError
	}

	return nil
}

// Update modifies an existing order record in the database.
// It handles NULL worker IDs by using interface{} to pass NULL to the database when appropriate.
//
// Parameters:
//   - order: Order entity with updated values
//
// Returns:
//   - *models.Order: Updated order after the operation
//   - error: repository_errors.UpdateError if the operation fails
func (o OrderRepository) Update(order *models.Order) (*models.Order, error) {
	query := `UPDATE orders SET worker_id = $1, user_id = $2, status = $3, address = $4, creation_date = $5, deadline = $6, rate = $7 WHERE id = $8 RETURNING id, worker_id, user_id, status, address, creation_date, deadline, rate;`

	var workerID interface{}
	if order.WorkerID != uuid.Nil {
		workerID = order.WorkerID
	}

	var updatedOrder models.Order
	err := o.db.QueryRow(query, workerID, order.UserID, order.Status, order.Address, order.CreationDate, order.Deadline, order.Rate, order.ID).Scan(&updatedOrder.ID, &updatedOrder.WorkerID, &updatedOrder.UserID, &updatedOrder.Status, &updatedOrder.Address, &updatedOrder.CreationDate, &updatedOrder.Deadline, &updatedOrder.Rate)
	if err != nil {
		return nil, repository_errors.UpdateError
	}
	return &updatedOrder, nil
}

// GetOrderByID retrieves an order by its unique identifier.
//
// Parameters:
//   - id: UUID of the order to retrieve
//
// Returns:
//   - *models.Order: Retrieved order entity
//   - error: repository_errors.DoesNotExist if no order found,
//     repository_errors.SelectError for other failures
func (o OrderRepository) GetOrderByID(id uuid.UUID) (*models.Order, error) {
	query := `SELECT * FROM orders WHERE id = $1;`
	orderDB := &OrderDB{}
	err := o.db.Get(orderDB, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository_errors.DoesNotExist
	} else if err != nil {
		return nil, repository_errors.SelectError
	}

	orderModels := copyOrderResultToModel(orderDB)

	return orderModels, nil
}

// GetTasksInOrder retrieves all tasks associated with a specific order.
//
// Parameters:
//   - id: UUID of the order to retrieve tasks for
//
// Returns:
//   - []models.Task: Slice of task entities associated with the order
//   - error: repository_errors.SelectError if the operation fails
func (o OrderRepository) GetTasksInOrder(id uuid.UUID) ([]models.Task, error) {
	query := `SELECT * FROM tasks WHERE id IN (SELECT task_id FROM order_contains_tasks WHERE order_id = $1);`
	var tasksDB []TaskDB
	err := o.db.Select(&tasksDB, query, id)
	if err != nil {
		return nil, repository_errors.SelectError
	}

	var taskModels []models.Task

	for i := range tasksDB {
		order := copyTaskResultToModel(&tasksDB[i])
		taskModels = append(taskModels, *order)
	}

	return taskModels, nil
}

// GetCurrentOrderByUserID retrieves the most recent order for a specific user.
//
// Parameters:
//   - id: UUID of the user to retrieve the current order for
//
// Returns:
//   - *models.Order: Retrieved order entity
//   - error: repository_errors.DoesNotExist if no order found,
//     repository_errors.SelectError for other failures
func (o OrderRepository) GetCurrentOrderByUserID(id uuid.UUID) (*models.Order, error) {
	query := `SELECT * FROM orders WHERE user_id = $1 ORDER BY creation_date DESC LIMIT 1;`
	orderDB := &OrderDB{}
	err := o.db.Get(orderDB, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository_errors.DoesNotExist
	} else if err != nil {
		return nil, repository_errors.SelectError
	}

	orderModels := copyOrderResultToModel(orderDB)

	return orderModels, nil
}

// GetAllOrdersByUserID retrieves all orders for a specific user.
//
// Parameters:
//   - id: UUID of the user to retrieve orders for
//
// Returns:
//   - []models.Order: Slice of order entities for the specified user
//   - error: repository_errors.SelectError if the operation fails
func (o OrderRepository) GetAllOrdersByUserID(id uuid.UUID) ([]models.Order, error) {
	query := `SELECT * FROM orders WHERE user_id = $1;`
	var orderDB []OrderDB

	err := o.db.Select(&orderDB, query, id)

	if err != nil {
		return nil, repository_errors.SelectError
	}

	var orderModels []models.Order
	for i := range orderDB {
		order := copyOrderResultToModel(&orderDB[i])
		orderModels = append(orderModels, *order)
	}

	return orderModels, nil
}

// Filter retrieves orders matching the specified criteria.
// Supports flexible filtering by multiple fields and multiple values per field.
//
// Parameters:
//   - params: Map of field names to filter values
//     (values can be comma-separated for OR conditions)
//
// Returns:
//   - []models.Order: Slice of order entities matching the filter criteria
//   - error: repository_errors.SelectError if the operation fails
func (o OrderRepository) Filter(params map[string]string) ([]models.Order, error) {
	var query strings.Builder
	query.WriteString("SELECT * FROM orders")

	if len(params) > 0 {
		query.WriteString(" WHERE ")
		for field, value := range params {
			// Разделяем значения по запятой
			values := strings.Split(value, ",")
			if len(values) > 1 {
				// Если есть несколько значений, создаем условие SQL с OR
				query.WriteString("(")
				for _, v := range values {
					if v == "null" {
						query.WriteString(fmt.Sprintf("%s IS NULL OR ", field))
					} else if field == "status" {
						query.WriteString(fmt.Sprintf("%s = %s OR ", field, v))
					} else {
						query.WriteString(fmt.Sprintf("%s::text LIKE '%s' OR ", field, v))
					}
				}
				// Удаляем последний " OR "
				str := query.String()[:query.Len()-4]
				query.Reset()
				query.WriteString(str)
				query.WriteString(") AND ")
			} else {
				if value == "null" {
					query.WriteString(fmt.Sprintf("%s IS NULL AND ", field))
				} else if value == "not null" {
					query.WriteString(fmt.Sprintf("%s IS NOT NULL AND ", field))
				} else if field == "status" {
					query.WriteString(fmt.Sprintf("%s = %s AND ", field, value))
				} else {
					query.WriteString(fmt.Sprintf("%s::text LIKE '%s' AND ", field, value))
				}
			}
		}
		// Удаляем последний " AND "
		str := query.String()[:query.Len()-5]
		query.Reset()
		query.WriteString(str)
	}

	var orderDB []OrderDB
	err := o.db.Select(&orderDB, query.String())

	if err != nil {
		return nil, repository_errors.SelectError
	}

	var orderModels []models.Order
	for i := range orderDB {
		order := copyOrderResultToModel(&orderDB[i])
		orderModels = append(orderModels, *order)
	}

	return orderModels, nil
}

// AddTaskToOrder associates a task with an order.
//
// Parameters:
//   - orderID: UUID of the order
//   - taskID: UUID of the task to add
//
// Returns:
//   - error: repository_errors.InsertError if the operation fails
func (o OrderRepository) AddTaskToOrder(orderID uuid.UUID, taskID uuid.UUID) error {
	query := `INSERT INTO order_contains_tasks(order_id, task_id) VALUES ($1, $2);`
	_, err := o.db.Exec(query, orderID, taskID)

	if err != nil {
		return repository_errors.InsertError
	}

	return nil
}

// RemoveTaskFromOrder removes a task association from an order.
//
// Parameters:
//   - orderID: UUID of the order
//   - taskID: UUID of the task to remove
//
// Returns:
//   - error: repository_errors.DeleteError if the operation fails
func (o OrderRepository) RemoveTaskFromOrder(orderID uuid.UUID, taskID uuid.UUID) error {
	query := `DELETE FROM order_contains_tasks WHERE order_id = $1 AND task_id = $2;`
	_, err := o.db.Exec(query, orderID, taskID)

	if err != nil {
		return repository_errors.DeleteError
	}

	return nil
}

// UpdateTaskQuantity updates the quantity of a specific task in an order.
//
// Parameters:
//   - orderID: UUID of the order
//   - taskID: UUID of the task
//   - quantity: New quantity value
//
// Returns:
//   - error: repository_errors.UpdateError if the operation fails
func (o OrderRepository) UpdateTaskQuantity(orderID uuid.UUID, taskID uuid.UUID, quantity int) error {
	query := `UPDATE order_contains_tasks SET quantity = $1 WHERE order_id = $2 AND task_id = $3;`
	_, err := o.db.Exec(query, quantity, orderID, taskID)

	if err != nil {
		return repository_errors.UpdateError
	}

	return nil
}

// GetTaskQuantity retrieves the quantity of a specific task in an order.
//
// Parameters:
//   - orderID: UUID of the order
//   - taskID: UUID of the task
//
// Returns:
//   - int: Quantity of the task in the order
//   - error: repository_errors.DoesNotExist if the task is not in the order,
//     repository_errors.SelectError for other failures
func (o OrderRepository) GetTaskQuantity(orderID uuid.UUID, taskID uuid.UUID) (int, error) {
	query := `SELECT quantity FROM order_contains_tasks WHERE order_id = $1 AND task_id = $2 LIMIT 1;`
	var quantity int

	err := o.db.Get(&quantity, query, orderID, taskID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, repository_errors.DoesNotExist
	} else if err != nil {
		return 0, repository_errors.SelectError
	}

	return quantity, nil
}
