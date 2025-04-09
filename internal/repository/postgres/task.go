package postgres

import (
	"database/sql"
	"errors"
	"lab3/internal/models"
	"lab3/internal/repository/repository_errors"
	"lab3/internal/repository/repository_interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TaskDB struct {
	ID             uuid.UUID `db:"id"`
	Name           string    `db:"name"`
	PricePerSingle float64   `db:"price_per_single"`
	Category       int       `db:"category"`
}

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) repository_interfaces.ITaskRepository {
	return &TaskRepository{db: db}
}

func copyTaskResultToModel(taskDB *TaskDB) *models.Task {
	return &models.Task{
		ID:             taskDB.ID,
		Name:           taskDB.Name,
		PricePerSingle: taskDB.PricePerSingle,
		Category:       taskDB.Category,
	}
}

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

func (t TaskRepository) Update(task *models.Task) (*models.Task, error) {
	query := `UPDATE tasks SET name = $1, price_per_single = $2, category = $3 WHERE tasks.id = $4 RETURNING id, name, price_per_single, category;`

	var updatedTask models.Task
	err := t.db.QueryRow(query, task.Name, task.PricePerSingle, task.Category, task.ID).Scan(&updatedTask.ID, &updatedTask.Name, &updatedTask.PricePerSingle, &updatedTask.Category)
	if err != nil {
		return nil, repository_errors.UpdateError
	}
	return &updatedTask, nil
}

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
