package postgres

import (
	"lab3/internal/models"

	"github.com/jmoiron/sqlx"
)

type Category struct {
	ID   int
	Name string
}

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

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

func (c CategoryRepository) Delete(id int) error {
	_, err := c.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
