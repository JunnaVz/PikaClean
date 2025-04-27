// Package interfaces implements the business logic services for the PikaClean application.
// This package contains concrete implementations of the service interfaces defined in service_interfaces.
package interfaces

import (
	"teamdev/internal/models"
	"teamdev/internal/repository/repository_interfaces"

	"github.com/charmbracelet/log"
)

// CategoryService implements the ICategoryService interface and provides
// business logic for category management operations.
type CategoryService struct {
	CategoryRepository repository_interfaces.ICategoryRepository // Repository for category data access
	TaskRepository     repository_interfaces.ITaskRepository     // Repository for task data access
	logger             *log.Logger                               // Logger for recording service events
}

// NewCategoryService creates a new CategoryService instance with the provided repositories
// and logger.
//
// Parameters:
//   - CategoryRepository: Repository for category data operations
//   - TaskRepository: Repository for task data operations
//   - logger: Logger for service operations
//
// Returns:
//   - *CategoryService: Initialized service implementation
func NewCategoryService(CategoryRepository repository_interfaces.ICategoryRepository, TaskRepository repository_interfaces.ITaskRepository, logger *log.Logger) *CategoryService {
	return &CategoryService{
		CategoryRepository: CategoryRepository,
		TaskRepository:     TaskRepository,
		logger:             logger,
	}
}

// Create adds a new cleaning service category with the specified name.
//
// Parameters:
//   - name: Name for the new category
//
// Returns:
//   - *models.Category: Created category with assigned ID
//   - error: Error if creation fails
func (c *CategoryService) Create(name string) (*models.Category, error) {
	category := &models.Category{
		Name: name,
	}

	category, err := c.CategoryRepository.Create(category)
	if err != nil {
		c.logger.Error("Error creating category")
		return nil, err
	}

	return category, nil
}

// Update modifies an existing category with new information.
//
// Parameters:
//   - category: Category with updated data (name)
//
// Returns:
//   - *models.Category: Updated category data
//   - error: Error if update fails
func (c *CategoryService) Update(category *models.Category) (*models.Category, error) {
	category, err := c.CategoryRepository.Update(category)
	if err != nil {
		c.logger.Error("Error updating category")
		return nil, err
	}

	return category, nil
}

// Delete removes a category from the system by ID.
//
// Parameters:
//   - id: ID of the category to delete
//
// Returns:
//   - error: Error if deletion fails
func (c *CategoryService) Delete(id int) error {
	err := c.CategoryRepository.Delete(id)
	if err != nil {
		c.logger.Error("Error deleting category")
	}
	return err
}

// GetAll retrieves all categories from the system.
//
// Returns:
//   - []models.Category: Slice of all category entities
//   - error: Error if retrieval fails
func (c *CategoryService) GetAll() ([]models.Category, error) {
	categories, err := c.CategoryRepository.GetAll()
	if err != nil {
		c.logger.Error("Error getting all categories")
		return nil, err
	}

	return categories, nil
}

// GetByID retrieves a category by its unique identifier.
//
// Parameters:
//   - id: ID of the category to retrieve
//
// Returns:
//   - *models.Category: Retrieved category entity
//   - error: Error if retrieval fails or category not found
func (c *CategoryService) GetByID(id int) (*models.Category, error) {
	category, err := c.CategoryRepository.GetByID(id)
	if err != nil {
		c.logger.Error("Error getting category by id")
		return nil, err
	}

	return category, nil
}

// GetTasksInCategory retrieves all tasks belonging to a specific category.
//
// Parameters:
//   - id: ID of the category to get tasks for
//
// Returns:
//   - []models.Task: Slice of task entities in the specified category
//   - error: Error if retrieval fails
func (c *CategoryService) GetTasksInCategory(id int) ([]models.Task, error) {
	tasks, err := c.TaskRepository.GetTasksInCategory(id)
	if err != nil {
		c.logger.Error("Error getting tasks in category")
		return nil, err
	}

	return tasks, nil
}
