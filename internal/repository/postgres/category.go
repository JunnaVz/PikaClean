// Package postgres provides repository implementations for data persistence
// using PostgreSQL database. It includes repositories for managing workers,
// users, tasks, orders, and categories.
package postgres

import (
	"teamdev/internal/models"

	"github.com/jmoiron/sqlx"
)

// Category represents a category entity as stored in the PostgreSQL database.
// It maps directly to the columns in the categories table.
type Category struct {
	ID   int    // Unique identifier for the category
	Name string // Descriptive name of the category
}

// CategoryRepository implements the ICategoryRepository interface for PostgreSQL.
// It provides methods for creating, retrieving, updating and deleting category records.
type CategoryRepository struct {
	db *sqlx.DB // Database connection
}

// NewCategoryRepository creates a new CategoryRepository instance with the provided
// database connection.
//
// Parameters:
//   - db: An initialized sqlx.DB connection to PostgreSQL
//
// Returns:
//   - *CategoryRepository: Repository implementation
func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// GetAll retrieves all categories from the database.
//
// Returns:
//   - []models.Category: Slice of all category entities
//   - error: Database error if the operation fails
func (c CategoryRepository) GetAll() ([]models.Category, error) {
	var categories []Category
	err := c.db.Select(&categories, "SELECT * FROM categories")
	if err != nil {
		return nil, err
	}

	var categoryModels []models.Category
	for i := range categories {
		categoryModel := models.Category{
			ID:   categories[i].ID,
			Name: categories[i].Name,
		}

		categoryModels = append(categoryModels, categoryModel)
	}
	return categoryModels, nil
}

// GetByID retrieves a category by its unique identifier.
//
// Parameters:
//   - id: ID of the category to retrieve
//
// Returns:
//   - *models.Category: Retrieved category entity
//   - error: Database error if the operation fails or no category found
func (c CategoryRepository) GetByID(id int) (*models.Category, error) {
	var category Category
	err := c.db.Get(&category, "SELECT * FROM categories WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &models.Category{
		ID:   category.ID,
		Name: category.Name,
	}, nil
}

// Create inserts a new category record into the database.
//
// Parameters:
//   - category: Category entity to be created
//
// Returns:
//   - *models.Category: Created category with assigned ID
//   - error: Database error if the operation fails
func (c CategoryRepository) Create(category *models.Category) (*models.Category, error) {
	query := `INSERT INTO categories(name) VALUES ($1) RETURNING id;`

	var categoryID int
	err := c.db.QueryRow(query, category.Name).Scan(&categoryID)

	if err != nil {
		return nil, err
	}

	return &models.Category{
		ID:   categoryID,
		Name: category.Name,
	}, nil
}

// Update modifies an existing category record in the database.
//
// Parameters:
//   - category: Category entity with updated values
//
// Returns:
//   - *models.Category: Updated category after the operation
//   - error: Database error if the operation fails
func (c CategoryRepository) Update(category *models.Category) (*models.Category, error) {
	query := `UPDATE categories SET name = $2 WHERE id = $1 RETURNING id;`

	var categoryID int
	err := c.db.QueryRow(query, category.ID, category.Name).Scan(&categoryID)

	if err != nil {
		return nil, err
	}

	return &models.Category{
		ID:   categoryID,
		Name: category.Name,
	}, nil
}

// Delete removes a category record from the database by ID.
//
// Parameters:
//   - id: ID of the category to delete
//
// Returns:
//   - error: Database error if the operation fails
func (c CategoryRepository) Delete(id int) error {
	_, err := c.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
