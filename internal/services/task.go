// Package interfaces provides service layer implementations for the PikaClean application,
// handling the business logic between API controllers and data repositories.
// This package contains concrete implementations of the service interfaces defined
// in the service_interfaces package, with implementations for tasks, orders, users,
// workers, and categories.
package interfaces

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"teamdev/internal/models"
	"teamdev/internal/repository/repository_interfaces"
	"teamdev/internal/services/service_interfaces"
)

// TaskService implements the ITaskService interface and provides
// business logic for managing cleaning tasks in the application.
// It handles task creation, updates, deletion, and various retrieval methods.
type TaskService struct {
	TaskRepository repository_interfaces.ITaskRepository // Repository for persistent task operations
	logger         *log.Logger                           // Logger for recording service activity
}

// NewTaskService creates a new TaskService instance with the provided dependencies.
//
// Parameters:
//   - TaskRepository: Repository for task data access operations
//   - logger: Logger for recording service activity and errors
//
// Returns:
//   - service_interfaces.ITaskService: A fully initialized task service
func NewTaskService(TaskRepository repository_interfaces.ITaskRepository, logger *log.Logger) service_interfaces.ITaskService {
	return &TaskService{
		TaskRepository: TaskRepository,
		logger:         logger,
	}
}

// Create adds a new cleaning task to the system with the provided details.
//
// Parameters:
//   - name: Descriptive name of the cleaning task
//   - price: Cost per unit of the task
//   - category: Category ID the task belongs to
//
// Returns:
//   - *models.Task: Created task with assigned ID if successful
//   - error: Validation or persistence errors if they occur
func (t TaskService) Create(name string, price float64, category int) (*models.Task, error) {
	//if !validName(name) || !validPrice(price) || !validCategory(category) {
	//	t.logger.Error("SERVICE: Invalid input")
	//	return nil, fmt.Errorf("SERVICE: Invalid input")
	//}
	if !validName(name) {
		t.logger.Error("SERVICE: Invalid name")
		return nil, fmt.Errorf("SERVICE: Invalid name")
	}
	if !validPrice(price) {
		t.logger.Error("SERVICE: Invalid price")
		return nil, fmt.Errorf("SERVICE: Invalid price")
	}
	if !validCategory(category) {
		t.logger.Error("SERVICE: Invalid category")
		return nil, fmt.Errorf("SERVICE: Invalid category")
	}

	task := &models.Task{
		Name:           name,
		PricePerSingle: price,
		Category:       category,
	}

	task, err := t.TaskRepository.Create(task)
	if err != nil {
		t.logger.Error("SERVICE: CreateNewTask method failed", "error", err)
		return nil, err
	}

	t.logger.Info("SERVICE: Successfully created new task", "task", task)
	return task, nil
}

// Update modifies an existing task with new information.
//
// Parameters:
//   - taskID: UUID of the task to update
//   - category: New category ID for the task
//   - name: New name for the task
//   - price: New price per unit for the task
//
// Returns:
//   - *models.Task: Updated task after changes
//   - error: Validation or persistence errors if they occur
func (t TaskService) Update(taskID uuid.UUID, category int, name string, price float64) (*models.Task, error) {
	task, err := t.GetTaskByID(taskID)
	if err != nil {
		t.logger.Error("SERVICE: GetTaskByID method failed", "id", taskID, "error", err)
		return nil, err
	}

	//if !validCategory(category) || !validName(name) || !validPrice(price) {
	//	t.logger.Error("SERVICE: Invalid input")
	//	return nil, fmt.Errorf("SERVICE: Invalid input")
	//} else {

	if !validCategory(category) {
		t.logger.Error("SERVICE: Invalid category")
		return nil, fmt.Errorf("SERVICE: Invalid category")
	}
	if !validName(name) {
		t.logger.Error("SERVICE: Invalid name")
		return nil, fmt.Errorf("SERVICE: Invalid name")
	}
	if !validPrice(price) {
		t.logger.Error("SERVICE: Invalid price")
		return nil, fmt.Errorf("SERVICE: Invalid price")
	}

	task.Category = category
	task.Name = name
	task.PricePerSingle = price

	updatedTask, err := t.TaskRepository.Update(task)
	if err != nil {
		t.logger.Error("SERVICE: UpdateTask method failed", "error", err)
		return nil, err
	}

	t.logger.Info("SERVICE: Successfully updated task price", "task", task)
	return updatedTask, nil
}

// Delete removes a task from the system.
//
// Parameters:
//   - taskID: UUID of the task to delete
//
// Returns:
//   - error: Any errors that occur during the deletion process
func (t TaskService) Delete(taskID uuid.UUID) error {
	_, err := t.GetTaskByID(taskID)
	if err != nil {
		t.logger.Error("SERVICE: GetTaskByID method failed", "id", taskID, "error", err)
		return err
	}

	err = t.TaskRepository.Delete(taskID)
	if err != nil {
		t.logger.Error("SERVICE: DeleteTask method failed", "error", err)
		return err
	}

	t.logger.Info("SERVICE: Successfully deleted task", "task", taskID)
	return nil
}

// GetAllTasks retrieves all cleaning tasks in the system.
//
// Returns:
//   - []models.Task: Slice of all task entities
//   - error: Any retrieval errors
func (t TaskService) GetAllTasks() ([]models.Task, error) {
	tasks, err := t.TaskRepository.GetAllTasks()
	if err != nil {
		t.logger.Error("SERVICE: GetAllTasks method failed", "error", err)
		return nil, err
	}

	t.logger.Info("SERVICE: Successfully got all tasks", "tasks", tasks)
	return tasks, nil
}

// GetTaskByID retrieves a specific task by its unique identifier.
//
// Parameters:
//   - id: UUID of the task to retrieve
//
// Returns:
//   - *models.Task: Retrieved task entity
//   - error: Any retrieval errors
func (t TaskService) GetTaskByID(id uuid.UUID) (*models.Task, error) {
	task, err := t.TaskRepository.GetTaskByID(id)

	if err != nil {
		t.logger.Error("SERVICE: GetTaskByID method failed", "id", id, "error", err)
		return nil, err
	}

	t.logger.Info("SERVICE: Successfully got task with GetTaskByID", "id", id)
	return task, nil
}

// GetTasksInCategory retrieves all tasks belonging to a specific category.
//
// Parameters:
//   - category: Category ID to filter tasks by
//
// Returns:
//   - []models.Task: Slice of tasks in the specified category
//   - error: Validation or retrieval errors
func (t TaskService) GetTasksInCategory(category int) ([]models.Task, error) {
	print(category)
	if !validCategory(category) {
		t.logger.Error("SERVICE: Invalid category", "category", category)
		return nil, fmt.Errorf("SERVICE: Invalid category")
	}

	tasks, err := t.TaskRepository.GetTasksInCategory(category)
	if err != nil {
		t.logger.Error("SERVICE: GetTasksInCategory method failed", "error", err)
		return nil, err
	}

	t.logger.Info("SERVICE: Successfully got tasks in category", "category", category)
	return tasks, nil
}

// GetTaskByName retrieves a task by its name.
//
// Parameters:
//   - name: Name of the task to search for
//
// Returns:
//   - *models.Task: Retrieved task entity
//   - error: Any retrieval errors
func (t TaskService) GetTaskByName(name string) (*models.Task, error) {
	task, err := t.TaskRepository.GetTaskByName(name)

	if err != nil {
		t.logger.Error("SERVICE: GetTaskByName method failed", "name", name, "error", err)
		return nil, err
	}

	t.logger.Info("SERVICE: Successfully got task with GetTaskByName", "name", name)
	return task, nil
}
