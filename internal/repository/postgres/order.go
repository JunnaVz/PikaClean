package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"lab3/internal/models"
	"lab3/internal/repository/repository_errors"
	"lab3/internal/repository/repository_interfaces"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderDB struct {
	ID           uuid.UUID `db:"id"`
	WorkerID     uuid.UUID `db:"worker_id"`
	UserID       uuid.UUID `db:"user_id"`
	Status       int       `db:"status"`
	Address      string    `db:"address"`
	CreationDate time.Time `db:"creation_date"`
	Deadline     time.Time `db:"deadline"`
	Rate         int       `db:"rate"`
}

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) repository_interfaces.IOrderRepository {
	return &OrderRepository{db: db}
}

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

func (o OrderRepository) AddTaskToOrder(orderID uuid.UUID, taskID uuid.UUID) error {
	query := `INSERT INTO order_contains_tasks(order_id, task_id) VALUES ($1, $2);`
	_, err := o.db.Exec(query, orderID, taskID)

	if err != nil {
		return repository_errors.InsertError
	}

	return nil
}

func (o OrderRepository) RemoveTaskFromOrder(orderID uuid.UUID, taskID uuid.UUID) error {
	query := `DELETE FROM order_contains_tasks WHERE order_id = $1 AND task_id = $2;`
	_, err := o.db.Exec(query, orderID, taskID)

	if err != nil {
		return repository_errors.DeleteError
	}

	return nil
}

func (o OrderRepository) UpdateTaskQuantity(orderID uuid.UUID, taskID uuid.UUID, quantity int) error {
	query := `UPDATE order_contains_tasks SET quantity = $1 WHERE order_id = $2 AND task_id = $3;`
	_, err := o.db.Exec(query, quantity, orderID, taskID)

	if err != nil {
		return repository_errors.UpdateError
	}

	return nil
}

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
