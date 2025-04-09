package postgres

import (
	"PikaClean/internal/models"
	"PikaClean/internal/repository/repository_errors"
	"PikaClean/internal/repository/repository_interfaces"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// TaskDB - структура для представления задачи в базе данных
type TaskDB struct {
	ID             uuid.UUID `db:"id"`               // Уникальный идентификатор задачи
	Name           string    `db:"name"`             // Название задачи
	PricePerSingle float64   `db:"price_per_single"` // Цена за единицу задачи
	Category       int       `db:"category"`         // Категория задачи
}

// TaskRepository - структура для взаимодействия с репозиторием задач
type TaskRepository struct {
	db *sqlx.DB // Соединение с базой данных
}

// NewTaskRepository - конструктор для создания нового репозитория задач
func NewTaskRepository(db *sqlx.DB) repository_interfaces.ITaskRepository {
	return &TaskRepository{db: db}
}

// copyTaskResultToModel - функция для преобразования данных задачи из базы в модель
func copyTaskResultToModel(taskDB *TaskDB) *models.Task {
	return &models.Task{
		ID:             taskDB.ID,
		Name:           taskDB.Name,
		PricePerSingle: taskDB.PricePerSingle,
		Category:       taskDB.Category,
	}
}

// Create - метод для создания новой задачи в базе данных
func (t TaskRepository) Create(task *models.Task) (*models.Task, error) {
	// SQL-запрос на добавление задачи
	query := `INSERT INTO tasks(name, price_per_single, category) VALUES ($1, $2, $3) RETURNING id;`

	var taskID uuid.UUID
	// Выполнение запроса и получение ID новой задачи
	err := t.db.QueryRow(query, task.Name, task.PricePerSingle, task.Category).Scan(&taskID)

	if err != nil {
		return nil, repository_errors.InsertError // Возвращаем ошибку, если добавление не удалось
	}

	// Возвращаем модель задачи с установленным ID
	return &models.Task{
		ID:             taskID,
		Name:           task.Name,
		PricePerSingle: task.PricePerSingle,
		Category:       task.Category,
	}, nil
}

// Delete - метод для удаления задачи из базы данных
func (t TaskRepository) Delete(id uuid.UUID) error {
	// SQL-запрос на удаление задачи
	query := `DELETE FROM tasks WHERE id = $1;`
	result, err := t.db.Exec(query, id)

	if err != nil {
		return repository_errors.DeleteError // Возвращаем ошибку, если удаление не удалось
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err // Возвращаем ошибку, если не удалось получить количество затронутых строк
	}

	if rowsAffected == 0 {
		return errors.New("no task found to delete") // Возвращаем ошибку, если задача не найдена для удаления
	}

	return nil
}

// Update - метод для обновления данных задачи в базе
func (t TaskRepository) Update(task *models.Task) (*models.Task, error) {
	// SQL-запрос на обновление данных задачи
	query := `UPDATE tasks SET name = $1, price_per_single = $2, category = $3 WHERE tasks.id = $4 RETURNING id, name, price_per_single, category;`

	var updatedTask models.Task
	// Выполнение запроса и получение обновленных данных задачи
	err := t.db.QueryRow(query, task.Name, task.PricePerSingle, task.Category, task.ID).Scan(&updatedTask.ID, &updatedTask.Name, &updatedTask.PricePerSingle, &updatedTask.Category)
	if err != nil {
		return nil, repository_errors.UpdateError // Возвращаем ошибку, если обновление не удалось
	}

	return &updatedTask, nil
}

// GetTaskByID - метод для получения задачи по ID
func (t TaskRepository) GetTaskByID(id uuid.UUID) (*models.Task, error) {
	// SQL-запрос на получение задачи по ID
	query := `SELECT * FROM tasks WHERE id = $1;`
	taskDB := &TaskDB{}
	err := t.db.Get(taskDB, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository_errors.DoesNotExist // Возвращаем ошибку, если задача не найдена
	} else if err != nil {
		return nil, repository_errors.SelectError // Возвращаем ошибку, если произошла ошибка при запросе
	}

	// Преобразуем данные из базы в модель задачи
	taskModels := copyTaskResultToModel(taskDB)

	return taskModels, nil
}

// GetTaskByName - метод для получения задачи по имени
func (t TaskRepository) GetTaskByName(name string) (*models.Task, error) {
	// SQL-запрос на получение задачи по имени
	query := `SELECT * FROM tasks WHERE name = $1 LIMIT 1;`
	taskDB := &TaskDB{}
	err := t.db.Get(taskDB, query, name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository_errors.DoesNotExist // Возвращаем ошибку, если задача не найдена
	} else if err != nil {
		return nil, repository_errors.SelectError // Возвращаем ошибку, если произошла ошибка при запросе
	}

	// Преобразуем данные из базы в модель задачи
	return copyTaskResultToModel(taskDB), nil
}

// GetAllTasks - метод для получения всех задач
func (t TaskRepository) GetAllTasks() ([]models.Task, error) {
	// SQL-запрос на получение всех задач
	query := `SELECT id, name, price_per_single, category FROM tasks;`
	var taskDB []TaskDB

	// Выполнение запроса и получение всех данных задач
	err := t.db.Select(&taskDB, query)

	if err != nil {
		return nil, repository_errors.SelectError // Возвращаем ошибку, если не удалось выполнить запрос
	}

	// Преобразуем данные из базы в модели задач
	var taskModels []models.Task
	for i := range taskDB {
		task := copyTaskResultToModel(&taskDB[i])
		taskModels = append(taskModels, *task)
	}

	return taskModels, nil
}

// GetTasksInCategory - метод для получения задач по категории
func (t TaskRepository) GetTasksInCategory(category int) ([]models.Task, error) {
	// SQL-запрос на получение задач по категории
	query := `SELECT * FROM tasks WHERE category = $1;`
	var taskDB []TaskDB

	// Выполнение запроса и получение данных задач по категории
	err := t.db.Select(&taskDB, query, category)

	if err != nil {
		return nil, repository_errors.SelectError // Возвращаем ошибку, если не удалось выполнить запрос
	}

	// Преобразуем данные из базы в модели задач
	var taskModels []models.Task
	for i := range taskDB {
		task := copyTaskResultToModel(&taskDB[i])
		taskModels = append(taskModels, *task)
	}

	return taskModels, nil
}
