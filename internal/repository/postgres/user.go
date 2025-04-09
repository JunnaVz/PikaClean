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

// UserDB структура для представления пользователя в базе данных
type UserDB struct {
	ID          uuid.UUID `db:"id"`           // Идентификатор пользователя
	Name        string    `db:"name"`         // Имя пользователя
	Surname     string    `db:"surname"`      // Фамилия пользователя
	Address     string    `db:"address"`      // Адрес пользователя
	PhoneNumber string    `db:"phone_number"` // Номер телефона пользователя
	Email       string    `db:"email"`        // Email пользователя
	Password    string    `db:"password"`     // Пароль пользователя
}

// UserRepository структура репозитория для работы с пользователями
type UserRepository struct {
	db *sqlx.DB // Экземпляр базы данных
}

// NewUserRepository создает новый репозиторий для работы с пользователями
func NewUserRepository(db *sqlx.DB) repository_interfaces.IUserRepository {
	return &UserRepository{db: db}
}

// copyUserResultToModel преобразует данные из базы данных в модель пользователя
func copyUserResultToModel(userDB *UserDB) *models.User {
	return &models.User{
		ID:          userDB.ID,
		Name:        userDB.Name,
		Surname:     userDB.Surname,
		Address:     userDB.Address,
		PhoneNumber: userDB.PhoneNumber,
		Email:       userDB.Email,
		Password:    userDB.Password,
	}
}

// Create добавляет нового пользователя в базу данных
func (u UserRepository) Create(user *models.User) (*models.User, error) {
	query := `INSERT INTO users(name, surname, address, phone_number, email, password) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

	var userID uuid.UUID
	// Выполнение запроса на добавление нового пользователя в базу данных
	err := u.db.QueryRow(query, user.Name, user.Surname, user.Address, user.PhoneNumber, user.Email, user.Password).Scan(&userID)

	if err != nil {
		return nil, repository_errors.InsertError // Ошибка вставки данных в базу
	}

	// Возврат нового пользователя с его ID
	return &models.User{
		ID:          userID,
		Name:        user.Name,
		Surname:     user.Surname,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Password:    user.Password,
	}, nil
}

// Delete удаляет пользователя из базы данных
func (u UserRepository) Delete(id uuid.UUID) error {
	// Начинаем транзакцию для удаления данных
	tx, err := u.db.Begin()
	if err != nil {
		return repository_errors.TransactionBeginError // Ошибка при начале транзакции
	}

	// Удаляем все заказы, связанные с пользователем
	_, err = tx.Exec(`DELETE FROM orders WHERE user_id = $1;`, id)
	if err != nil {
		_ = tx.Rollback()                    // Откат транзакции в случае ошибки
		return repository_errors.DeleteError // Ошибка удаления данных
	}

	// Удаляем пользователя
	_, err = tx.Exec(`DELETE FROM users WHERE id = $1;`, id)
	if err != nil {
		_ = tx.Rollback()                    // Откат транзакции в случае ошибки
		return repository_errors.DeleteError // Ошибка удаления данных
	}

	// Подтверждаем транзакцию
	err = tx.Commit()
	if err != nil {
		return repository_errors.TransactionCommitError // Ошибка при коммите транзакции
	}

	return nil
}

// Update обновляет данные пользователя в базе данных
func (u UserRepository) Update(user *models.User) (*models.User, error) {
	query := `UPDATE users SET name = $1, surname = $2, email = $3, phone_number = $4, address = $5, password = $6 WHERE users.id = $7 RETURNING id, name, surname, address, phone_number, email, password;`

	var updatedUser models.User
	// Выполнение запроса на обновление данных пользователя
	err := u.db.QueryRow(query, user.Name, user.Surname, user.Email, user.PhoneNumber, user.Address, user.Password, user.ID).Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Surname, &updatedUser.Address, &updatedUser.PhoneNumber, &updatedUser.Email, &updatedUser.Password)
	if err != nil {
		return nil, repository_errors.UpdateError // Ошибка обновления данных
	}

	return &updatedUser, nil
}

// GetUserByID извлекает пользователя из базы данных по его ID
func (u UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	query := `SELECT * FROM users WHERE id = $1;`
	userDB := &UserDB{}
	// Выполнение запроса для получения данных пользователя по его ID
	err := u.db.Get(userDB, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository_errors.DoesNotExist // Если пользователь не найден
	} else if err != nil {
		return nil, repository_errors.SelectError // Ошибка при выполнении запроса
	}

	// Преобразование данных из базы данных в модель пользователя
	userModels := copyUserResultToModel(userDB)

	return userModels, nil
}

// GetUserByEmail извлекает пользователя из базы данных по его email
func (u UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email = $1;`
	userDB := &UserDB{}
	// Выполнение запроса для получения данных пользователя по email
	err := u.db.Get(userDB, query, email)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository_errors.DoesNotExist // Если пользователь не найден
	} else if err != nil {
		return nil, repository_errors.SelectError // Ошибка при выполнении запроса
	}

	// Преобразование данных из базы данных в модель пользователя
	userModels := copyUserResultToModel(userDB)

	return userModels, nil
}

// GetAllUsers извлекает всех пользователей из базы данных
func (u UserRepository) GetAllUsers() ([]models.User, error) {
	query := `SELECT name, surname, address, phone_number, email FROM users;`
	var userDB []UserDB

	// Выполнение запроса для получения всех пользователей
	err := u.db.Select(&userDB, query)

	if err != nil {
		return nil, repository_errors.SelectError // Ошибка при выполнении запроса
	}

	var userModels []models.User
	// Преобразование всех записей в модели пользователей
	for i := range userDB {
		user := copyUserResultToModel(&userDB[i])
		userModels = append(userModels, *user)
	}

	return userModels, nil
}
