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

// WorkerDB представляет структуру работника, используемую для взаимодействия с базой данных.
type WorkerDB struct {
	ID          uuid.UUID `db:"id"`           // Идентификатор работника
	Name        string    `db:"name"`         // Имя работника
	Surname     string    `db:"surname"`      // Фамилия работника
	Address     string    `db:"address"`      // Адрес работника
	PhoneNumber string    `db:"phone_number"` // Номер телефона работника
	Email       string    `db:"email"`        // Электронная почта работника
	Role        int       `db:"role"`         // Роль работника (например, администратор, пользователь и т.д.)
	Password    string    `db:"password"`     // Пароль работника
}

// WorkerRepository представляет репозиторий для работы с сущностью работника в базе данных.
type WorkerRepository struct {
	db *sqlx.DB // Объект базы данных для выполнения запросов
}

// NewWorkerRepository создает новый экземпляр репозитория для работы с данными работников.
func NewWorkerRepository(db *sqlx.DB) repository_interfaces.IWorkerRepository {
	return &WorkerRepository{db: db}
}

// copyWorkerResultToModel копирует данные из структуры WorkerDB в структуру Worker для использования в бизнес-логике.
func copyWorkerResultToModel(workerDB *WorkerDB) *models.Worker {
	return &models.Worker{
		ID:          workerDB.ID,
		Name:        workerDB.Name,
		Surname:     workerDB.Surname,
		Address:     workerDB.Address,
		PhoneNumber: workerDB.PhoneNumber,
		Email:       workerDB.Email,
		Role:        workerDB.Role,
		Password:    workerDB.Password,
	}
}

// Create создает нового работника в базе данных и возвращает его с присвоенным ID.
func (w WorkerRepository) Create(worker *models.Worker) (*models.Worker, error) {
	// SQL запрос для добавления работника
	query := `INSERT INTO workers(name, surname, address, phone_number, email, role, password) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

	var workerID uuid.UUID
	// Выполнение запроса и получение нового ID
	err := w.db.QueryRow(query, worker.Name, worker.Surname, worker.Address, worker.PhoneNumber, worker.Email, worker.Role, worker.Password).Scan(&workerID)

	if err != nil {
		return nil, repository_errors.InsertError
	}

	// Возвращаем нового работника с присвоенным ID
	return &models.Worker{
		ID:          workerID,
		Name:        worker.Name,
		Surname:     worker.Surname,
		Address:     worker.Address,
		PhoneNumber: worker.PhoneNumber,
		Email:       worker.Email,
		Role:        worker.Role,
		Password:    worker.Password,
	}, nil
}

// Update обновляет данные работника в базе данных.
func (w WorkerRepository) Update(worker *models.Worker) (*models.Worker, error) {
	// SQL запрос для обновления данных работника
	query := `UPDATE workers SET name = $1, surname = $2, address = $3, phone_number = $4, email = $5, role = $6, password = $7 WHERE workers.id = $8 RETURNING id, name, surname, address, phone_number, email, role, password;`

	var updatedWorker models.Worker
	// Выполнение запроса и получение обновленных данных
	err := w.db.QueryRow(query, worker.Name, worker.Surname, worker.Address, worker.PhoneNumber, worker.Email, worker.Role, worker.Password, worker.ID).Scan(&updatedWorker.ID, &updatedWorker.Name, &updatedWorker.Surname, &updatedWorker.Address, &updatedWorker.PhoneNumber, &updatedWorker.Email, &updatedWorker.Role, &updatedWorker.Password)
	if err != nil {
		return nil, repository_errors.UpdateError
	}
	return &updatedWorker, nil
}

// Delete удаляет работника из базы данных по его ID.
func (w WorkerRepository) Delete(id uuid.UUID) error {
	// SQL запрос для удаления работника
	query := `DELETE FROM workers WHERE id = $1;`
	_, err := w.db.Exec(query, id)

	if err != nil {
		return repository_errors.DeleteError
	}

	return nil
}

// GetWorkerByID получает работника по его ID из базы данных.
func (w WorkerRepository) GetWorkerByID(id uuid.UUID) (*models.Worker, error) {
	// SQL запрос для получения данных работника по ID
	query := `SELECT * FROM workers WHERE id = $1;`
	workerDB := &WorkerDB{}
	err := w.db.Get(workerDB, query, id)

	// Проверка на отсутствие данных
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository_errors.DoesNotExist
	} else if err != nil {
		return nil, repository_errors.SelectError
	}

	// Преобразование данных из WorkerDB в модель Worker
	workerModels := copyWorkerResultToModel(workerDB)

	return workerModels, nil
}

// GetAllWorkers получает всех работников из базы данных.
func (w WorkerRepository) GetAllWorkers() ([]models.Worker, error) {
	// SQL запрос для получения всех работников
	query := `SELECT id, name, surname, address, phone_number, email, role FROM workers;`
	var workerDB []WorkerDB

	err := w.db.Select(&workerDB, query)

	if err != nil {
		return nil, repository_errors.SelectError
	}

	// Преобразование всех полученных данных в модели Worker
	var workerModels []models.Worker
	for i := range workerDB {
		worker := copyWorkerResultToModel(&workerDB[i])
		workerModels = append(workerModels, *worker)
	}

	return workerModels, nil
}

// GetWorkerByEmail получает работника по его email.
func (w WorkerRepository) GetWorkerByEmail(email string) (*models.Worker, error) {
	// SQL запрос для получения работника по email
	query := `SELECT * FROM workers WHERE email = $1;`
	workerDB := &WorkerDB{}
	err := w.db.Get(workerDB, query, email)

	// Проверка на отсутствие данных
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository_errors.DoesNotExist
	} else if err != nil {
		return nil, repository_errors.SelectError
	}

	// Преобразование данных из WorkerDB в модель Worker
	workerModels := copyWorkerResultToModel(workerDB)

	return workerModels, nil
}

// GetWorkersByRole получает всех работников по заданной роли.
func (w WorkerRepository) GetWorkersByRole(role int) ([]models.Worker, error) {
	// SQL запрос для получения работников по роли
	query := `SELECT * FROM workers WHERE role = $1;`
	var workerDB []WorkerDB

	err := w.db.Select(&workerDB, query, role)

	if err != nil {
		return nil, repository_errors.SelectError
	}

	// Преобразование всех полученных данных в модели Worker
	var workerModels []models.Worker
	for i := range workerDB {
		worker := copyWorkerResultToModel(&workerDB[i])
		workerModels = append(workerModels, *worker)
	}

	return workerModels, nil
}

// GetAverageOrderRate получает среднюю оценку работника по его заказам.
func (w WorkerRepository) GetAverageOrderRate(worker *models.Worker) (float64, error) {
	// SQL запрос для получения средней оценки заказов
	query := `SELECT AVG(rate) FROM orders WHERE worker_id = $1 AND status = 3 AND rate != 0;`
	var averageRate float64

	err := w.db.Get(&averageRate, query, worker.ID)

	if err != nil {
		return 0, repository_errors.SelectError
	}

	return averageRate, nil
}
