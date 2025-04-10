package interfaces

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"lab3/internal/models"
	"lab3/internal/repository/repository_interfaces"
	"lab3/internal/services/service_interfaces"
)

type TaskService struct {
	TaskRepository repository_interfaces.ITaskRepository
	logger         *log.Logger
}

func NewTaskService(TaskRepository repository_interfaces.ITaskRepository, logger *log.Logger) service_interfaces.ITaskService {
	return &TaskService{
		TaskRepository: TaskRepository,
		logger:         logger,
	}
}

func (t TaskService) Create(name string, price float64, category int) (*models.Task, error) {
	if !validName(name) || !validPrice(price) || !validCategory(category) {
		t.logger.Error("SERVICE: Invalid input")
		return nil, fmt.Errorf("SERVICE: Invalid input")
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

func (t TaskService) Update(taskID uuid.UUID, category int, name string, price float64) (*models.Task, error) {
	task, err := t.GetTaskByID(taskID)
	if err != nil {
		t.logger.Error("SERVICE: GetTaskByID method failed", "id", taskID, "error", err)
		return nil, err
	}

	if !validCategory(category) || !validName(name) || !validPrice(price) {
		t.logger.Error("SERVICE: Invalid input")
		return nil, fmt.Errorf("SERVICE: Invalid input")
	} else {
		task.Category = category
		task.Name = name
		task.PricePerSingle = price
	}

	updatedTask, err := t.TaskRepository.Update(task)
	if err != nil {
		t.logger.Error("SERVICE: UpdateTask method failed", "error", err)
		return nil, err
	}

	t.logger.Info("SERVICE: Successfully updated task price", "task", task)
	return updatedTask, nil
}

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

func (t TaskService) GetAllTasks() ([]models.Task, error) {
	tasks, err := t.TaskRepository.GetAllTasks()
	if err != nil {
		t.logger.Error("SERVICE: GetAllTasks method failed", "error", err)
		return nil, err
	}

	t.logger.Info("SERVICE: Successfully got all tasks", "tasks", tasks)
	return tasks, nil
}

func (t TaskService) GetTaskByID(id uuid.UUID) (*models.Task, error) {
	task, err := t.TaskRepository.GetTaskByID(id)

	if err != nil {
		t.logger.Error("SERVICE: GetTaskByID method failed", "id", id, "error", err)
		return nil, err
	}

	t.logger.Info("SERVICE: Successfully got task with GetTaskByID", "id", id)
	return task, nil
}

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

func (t TaskService) GetTaskByName(name string) (*models.Task, error) {
	task, err := t.TaskRepository.GetTaskByName(name)

	if err != nil {
		t.logger.Error("SERVICE: GetTaskByName method failed", "name", name, "error", err)
		return nil, err
	}

	t.logger.Info("SERVICE: Successfully got task with GetTaskByName", "name", name)
	return task, nil
}
