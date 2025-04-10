package interfaces

import (
	"lab3/internal/models"
	"lab3/internal/repository/repository_interfaces"

	"github.com/charmbracelet/log"
)

type CategoryService struct {
	CategoryRepository repository_interfaces.ICategoryRepository
	TaskRepository     repository_interfaces.ITaskRepository
	logger             *log.Logger
}

func NewCategoryService(CategoryRepository repository_interfaces.ICategoryRepository, TaskRepository repository_interfaces.ITaskRepository, logger *log.Logger) *CategoryService {
	return &CategoryService{
		CategoryRepository: CategoryRepository,
		TaskRepository:     TaskRepository,
		logger:             logger,
	}
}

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

func (c *CategoryService) Update(category *models.Category) (*models.Category, error) {
	category, err := c.CategoryRepository.Update(category)
	if err != nil {
		c.logger.Error("Error updating category")
		return nil, err
	}

	return category, nil
}

func (c *CategoryService) Delete(id int) error {
	err := c.CategoryRepository.Delete(id)
	if err != nil {
		c.logger.Error("Error deleting category")
	}
	return err
}

func (c *CategoryService) GetAll() ([]models.Category, error) {
	categories, err := c.CategoryRepository.GetAll()
	if err != nil {
		c.logger.Error("Error getting all categories")
		return nil, err
	}

	return categories, nil
}

func (c *CategoryService) GetByID(id int) (*models.Category, error) {
	category, err := c.CategoryRepository.GetByID(id)
	if err != nil {
		c.logger.Error("Error getting category by id")
		return nil, err
	}

	return category, nil
}

func (c *CategoryService) GetTasksInCategory(id int) ([]models.Task, error) {
	tasks, err := c.TaskRepository.GetTasksInCategory(id)
	if err != nil {
		c.logger.Error("Error getting tasks in category")
		return nil, err
	}

	return tasks, nil
}
