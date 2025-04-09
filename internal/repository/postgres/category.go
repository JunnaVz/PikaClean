package postgres

import (
	"PikaClean/internal/models"
	"github.com/jmoiron/sqlx"
)

// Category представляет структуру категории, используемую в базе данных.
type Category struct {
	ID   int    // Идентификатор категории
	Name string // Название категории
}

// CategoryRepository предоставляет методы для работы с категориями в базе данных.
type CategoryRepository struct {
	db *sqlx.DB // Подключение к базе данных
}

// NewCategoryRepository создает новый экземпляр репозитория для категорий.
func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// GetAll получает все категории из базы данных.
func (c CategoryRepository) GetAll() ([]models.Category, error) {
	// Выполняем запрос для получения всех категорий
	var categories []Category
	err := c.db.Select(&categories, "SELECT * FROM categories")
	if err != nil {
		return nil, err // Возвращаем ошибку, если запрос не удался
	}

	// Преобразуем данные из структуры Category в модель модели Category
	var categoryModels []models.Category
	for i := range categories {
		categoryModel := models.Category{
			ID:   categories[i].ID,
			Name: categories[i].Name,
		}
		categoryModels = append(categoryModels, categoryModel)
	}
	return categoryModels, nil // Возвращаем все категории
}

// GetByID получает категорию по ее идентификатору.
func (c CategoryRepository) GetByID(id int) (*models.Category, error) {
	// Выполняем запрос для получения категории по ID
	var category Category
	err := c.db.Get(&category, "SELECT * FROM categories WHERE id = $1", id)
	if err != nil {
		return nil, err // Возвращаем ошибку, если категория не найдена или произошла ошибка
	}

	// Возвращаем категорию в виде модели Category
	return &models.Category{
		ID:   category.ID,
		Name: category.Name,
	}, nil
}

// Create создает новую категорию в базе данных.
func (c CategoryRepository) Create(category *models.Category) (*models.Category, error) {
	// Вставляем новую категорию в базу данных
	query := `INSERT INTO categories(name) VALUES ($1) RETURNING id;`

	var categoryID int
	err := c.db.QueryRow(query, category.Name).Scan(&categoryID)
	if err != nil {
		return nil, err // Возвращаем ошибку, если не удалось вставить категорию
	}

	// Возвращаем созданную категорию с полученным ID
	return &models.Category{
		ID:   categoryID,
		Name: category.Name,
	}, nil
}

// Update обновляет информацию о категории в базе данных.
func (c CategoryRepository) Update(category *models.Category) (*models.Category, error) {
	// Обновляем категорию в базе данных по ID
	query := `UPDATE categories SET name = $2 WHERE id = $1 RETURNING id;`

	var categoryID int
	err := c.db.QueryRow(query, category.ID, category.Name).Scan(&categoryID)
	if err != nil {
		return nil, err // Возвращаем ошибку, если не удалось обновить категорию
	}

	// Возвращаем обновленную категорию
	return &models.Category{
		ID:   categoryID,
		Name: category.Name,
	}, nil
}

// Delete удаляет категорию по ее идентификатору.
func (c CategoryRepository) Delete(id int) error {
	// Выполняем запрос на удаление категории по ID
	_, err := c.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err // Возвращаем ошибку, если удаление не удалось
	}
	return nil // Возвращаем nil, если категория была успешно удалена
}
