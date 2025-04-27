// Package postgres provides repository implementations for data persistence
// using PostgreSQL database. It includes repositories for managing workers,
// users, tasks, orders, and categories.
package postgres

import (
	"database/sql"
	"errors"
	"teamdev/internal/models"
	"teamdev/internal/repository/repository_errors"
	"teamdev/internal/repository/repository_interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// TaskDB represents a task entity as stored in the PostgreSQL database.
// It maps directly to the columns in the tasks table.
type TaskDB struct {
	ID             uuid.UUID `db:"id"`               // Unique identifier for the task
	Name           string    `db:"name"`             // Descriptive name of the cleaning task
	PricePerSingle float64   `db:"price_per_single"` // Cost per unit of the task
	Category       int       `db:"category"`         // Category ID the task belongs to
}

// TaskRepository implements the ITaskRepository interface for PostgreSQL.
// It provides methods for creating, updating, and retrieving task records.
type TaskRepository struct {
	db *sqlx.DB // Database connection
}

// NewTaskRepository creates a new TaskRepository instance with the provided
// database connection.
//
// Parameters:
//   - db: An initialized sqlx.DB connection to PostgreSQL
//
// Returns:
//   - repository_interfaces.ITaskRepository: Repository implementation
func NewTaskRepository(db *sqlx.DB) repository_interfaces.ITaskRepository {
	return &TaskRepository{db: db}
}

// copyTaskResultToModel converts a TaskDB database entity to a models.Task domain entity.
//
// Parameters:
//   - taskDB: Database entity to convert
//
// Returns:
//   - *models.Task: Corresponding domain entity
func copyTaskResultToModel(taskDB *TaskDB) *models.Task {
	return &models.Task{
		ID:             taskDB.ID,
		Name:           taskDB.Name,
		PricePerSingle: taskDB.PricePerSingle,
		Category:       taskDB.Category,
	}
}

// Create inserts a new task record into the database.
//
// Parameters:
//   - task: Task entity to be created
//
// Returns:
//   - *models.Task: Created task with assigned ID
//   - error: repository_errors.InsertError if the operation fails
func (t TaskRepository) Create(task *models.Task) (*models.Task, error) {
	query := `INSERT INTO tasks(name, price_per_single, category) VALUES ($1, $2, $3) RETURNING id;`

	var taskID uuid.UUID
	err := t.db.QueryRow(query, task.Name, task.PricePerSingle, task.Category).Scan(&taskID)

	if err != nil {
		return nil, repository_errors.InsertError
	}

	return &models.Task{
		ID:             taskID,
		Name:           task.Name,
		PricePerSingle: task.PricePerSingle,
		Category:       task.Category,
	}, nil
}

// Delete removes a task record from the database by ID.
//
// Parameters:
//   - id: UUID of the task to delete
//
// Returns:
//   - error: repository_errors.DeleteError if the operation fails,
//     or a custom error if the task doesn't exist
func (t TaskRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1;`
	result, err := t.db.Exec(query, id)

	if err != nil {
		return repository_errors.DeleteError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no task found to delete")
	}

	return nil
}

// Update modifies an existing task record in the database.
//
// Parameters:
//   - task: Task entity with updated values
//
// Returns:
//   - *models.Task: Updated task after the operation
//   - error: repository_errors.UpdateError if the operation fails
func (t TaskRepository) Update(task *models.Task) (*models.Task, error) {
	query := `UPDATE tasks SET name = $1, price_per_single = $2, category = $3 WHERE tasks.id = $4 RETURNING id, name, price_per_single, category;`

	var updatedTask models.Task
	err := t.db.QueryRow(query, task.Name, task.PricePerSingle, task.Category, task.ID).Scan(&updatedTask.ID, &updatedTask.Name, &updatedTask.PricePerSingle, &updatedTask.Category)
	if err != nil {
		return nil, repository_errors.UpdateError
	}
	return &updatedTask, nil
}

// GetTaskByID retrieves a task by its unique identifier.
//
// Parameters:
//   - id: UUID of the task to retrieve
//
// Returns:
//   - *models.Task: Retrieved task entity
//   - error: repository_errors.DoesNotExist if no task found,
//     repository_errors.SelectError for other failures
func (t TaskRepository) GetTaskByID(id uuid.UUID) (*models.Task, error) {
	query := `SELECT * FROM tasks WHERE id = $1;`
	taskDB := &TaskDB{}
	err := t.db.Get(taskDB, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository_errors.DoesNotExist
	} else if err != nil {
		return nil, repository_errors.SelectError
	}

	taskModels := copyTaskResultToModel(taskDB)

	return taskModels, nil
}

// GetTaskByName retrieves a task by its name.
//
// Parameters:
//   - name: Name of the task to search for
//
// Returns:
//   - *models.Task: Retrieved task entity
//   - error: repository_errors.DoesNotExist if no task found,
//     repository_errors.SelectError for other failures
func (t TaskRepository) GetTaskByName(name string) (*models.Task, error) {
	query := `SELECT * FROM tasks WHERE name = $1 LIMIT 1;`
	taskDB := &TaskDB{}
	err := t.db.Get(taskDB, query, name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository_errors.DoesNotExist
	} else if err != nil {
		return nil, repository_errors.SelectError
	}

	return copyTaskResultToModel(taskDB), nil
}

// GetAllTasks retrieves all tasks from the database.
//
// Returns:
//   - []models.Task: Slice of all task entities
//   - error: repository_errors.SelectError if the operation fails
func (t TaskRepository) GetAllTasks() ([]models.Task, error) {
	query := `SELECT id, name, price_per_single, category FROM tasks;`
	var taskDB []TaskDB

	err := t.db.Select(&taskDB, query)

	if err != nil {
		return nil, repository_errors.SelectError
	}

	var taskModels []models.Task
	for i := range taskDB {
		user := copyTaskResultToModel(&taskDB[i])
		taskModels = append(taskModels, *user)
	}

	return taskModels, nil
}

// GetTasksInCategory retrieves all tasks belonging to a specific category.
//
// Parameters:
//   - category: Category ID to filter by
//
// Returns:
//   - []models.Task: Slice of task entities in the specified category
//   - error: repository_errors.SelectError if the operation fails
func (t TaskRepository) GetTasksInCategory(category int) ([]models.Task, error) {
	query := `SELECT * FROM tasks WHERE category = $1;`
	var taskDB []TaskDB

	err := t.db.Select(&taskDB, query, category)

	if err != nil {
		return nil, repository_errors.SelectError
	}

	var taskModels []models.Task
	for i := range taskDB {
		task := copyTaskResultToModel(&taskDB[i])
		taskModels = append(taskModels, *task)
	}

	return taskModels, nil
}
