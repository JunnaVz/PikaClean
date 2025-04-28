// Package postgres provides repository implementations for data persistence
// using PostgreSQL database. It includes repositories for managing workers,
// users, tasks, orders, and categories.
package postgres

import (
	"database/sql"
	"errors"
	"teamdev/internal/models"
	"teamdev/internal/repository/repository_errors"
	"teamdev/internal/repository/repository_interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// WorkerDB represents a worker entity as stored in the PostgreSQL database.
// It maps directly to the columns in the workers table.
type WorkerDB struct {
	ID          uuid.UUID `db:"id"`           // Unique identifier for the worker
	Name        string    `db:"name"`         // First name of the worker
	Surname     string    `db:"surname"`      // Last name of the worker
	Address     string    `db:"address"`      // Physical address of the worker
	PhoneNumber string    `db:"phone_number"` // Contact phone number
	Email       string    `db:"email"`        // Email address, used as username for login
	Role        int       `db:"role"`         // Role identifier (determines permissions)
	Password    string    `db:"password"`     // Hashed password for authentication
}

// WorkerRepository implements the IWorkerRepository interface for PostgreSQL.
// It provides methods for creating, updating, and retrieving worker records.
type WorkerRepository struct {
	db *sqlx.DB // Database connection
}

// NewWorkerRepository creates a new WorkerRepository instance with the provided
// database connection.
//
// Parameters:
//   - db: An initialized sqlx.DB connection to PostgreSQL
//
// Returns:
//   - repository_interfaces.IWorkerRepository: Repository implementation
func NewWorkerRepository(db *sqlx.DB) repository_interfaces.IWorkerRepository {
	return &WorkerRepository{db: db}
}

// copyWorkerResultToModel converts a WorkerDB database entity to a models.Worker domain entity.
//
// Parameters:
//   - workerDB: Database entity to convert
//
// Returns:
//   - *models.Worker: Corresponding domain entity
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

// Create inserts a new worker record into the database.
//
// Parameters:
//   - worker: Worker entity to be created
//
// Returns:
//   - *models.Worker: Created worker with assigned ID
//   - error: repository_errors.InsertError if the operation fails
func (w WorkerRepository) Create(worker *models.Worker) (*models.Worker, error) {
	query := `INSERT INTO workers(name, surname, address, phone_number, email, role, password) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

	var workerID uuid.UUID
	err := w.db.QueryRow(query, worker.Name, worker.Surname, worker.Address, worker.PhoneNumber, worker.Email, worker.Role, worker.Password).Scan(&workerID)

	if err != nil {
		return nil, repository_errors.InsertError
	}

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

// Update modifies an existing worker record in the database.
//
// Parameters:
//   - worker: Worker entity with updated values
//
// Returns:
//   - *models.Worker: Updated worker after the operation
//   - error: repository_errors.UpdateError if the operation fails
func (w WorkerRepository) Update(worker *models.Worker) (*models.Worker, error) {
	query := `UPDATE workers SET name = $1, surname = $2, address = $3, phone_number = $4, email = $5, role = $6, password = $7 WHERE workers.id = $8 RETURNING id, name, surname, address, phone_number, email, role, password;`

	var updatedWorker models.Worker
	err := w.db.QueryRow(query, worker.Name, worker.Surname, worker.Address, worker.PhoneNumber, worker.Email, worker.Role, worker.Password, worker.ID).Scan(&updatedWorker.ID, &updatedWorker.Name, &updatedWorker.Surname, &updatedWorker.Address, &updatedWorker.PhoneNumber, &updatedWorker.Email, &updatedWorker.Role, &updatedWorker.Password)
	if err != nil {
		return nil, repository_errors.UpdateError
	}
	return &updatedWorker, nil
}

// Delete removes a worker record from the database by ID.
//
// Parameters:
//   - id: UUID of the worker to delete
//
// Returns:
//   - error: repository_errors.DeleteError if the operation fails
func (w WorkerRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM workers WHERE id = $1;`
	_, err := w.db.Exec(query, id)

	if err != nil {
		return repository_errors.DeleteError
	}

	return nil
}

// GetWorkerByID retrieves a worker by their unique identifier.
//
// Parameters:
//   - id: UUID of the worker to retrieve
//
// Returns:
//   - *models.Worker: Retrieved worker entity
//   - error: repository_errors.DoesNotExist if no worker found,
//     repository_errors.SelectError for other failures
func (w WorkerRepository) GetWorkerByID(id uuid.UUID) (*models.Worker, error) {
	query := `SELECT * FROM workers WHERE id = $1;`
	workerDB := &WorkerDB{}
	err := w.db.Get(workerDB, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository_errors.DoesNotExist
	} else if err != nil {
		return nil, repository_errors.SelectError
	}

	workerModels := copyWorkerResultToModel(workerDB)

	return workerModels, nil
}

// GetAllWorkers retrieves all workers from the database.
//
// Returns:
//   - []models.Worker: Slice of all worker entities
//   - error: repository_errors.SelectError if the operation fails
func (w WorkerRepository) GetAllWorkers() ([]models.Worker, error) {
	query := `SELECT id, name, surname, address, phone_number, email, role FROM workers;`
	var workerDB []WorkerDB

	err := w.db.Select(&workerDB, query)

	if err != nil {
		return nil, repository_errors.SelectError
	}

	var workerModels []models.Worker
	for i := range workerDB {
		worker := copyWorkerResultToModel(&workerDB[i])
		workerModels = append(workerModels, *worker)
	}

	return workerModels, nil
}

// GetWorkerByEmail retrieves a worker by their email address.
//
// Parameters:
//   - email: Email address to search for
//
// Returns:
//   - *models.Worker: Retrieved worker entity
//   - error: repository_errors.DoesNotExist if no worker found,
//     repository_errors.SelectError for other failures
func (w WorkerRepository) GetWorkerByEmail(email string) (*models.Worker, error) {
	query := `SELECT * FROM workers WHERE email = $1;`
	workerDB := &WorkerDB{}
	err := w.db.Get(workerDB, query, email)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository_errors.DoesNotExist
	} else if err != nil {
		return nil, repository_errors.SelectError
	}

	workerModels := copyWorkerResultToModel(workerDB)

	return workerModels, nil
}

// GetWorkersByRole retrieves all workers with a specific role.
//
// Parameters:
//   - role: Role identifier to filter by
//
// Returns:
//   - []models.Worker: Slice of worker entities with the specified role
//   - error: repository_errors.SelectError if the operation fails
func (w WorkerRepository) GetWorkersByRole(role int) ([]models.Worker, error) {
	query := `SELECT * FROM workers WHERE role = $1;`
	var workerDB []WorkerDB

	err := w.db.Select(&workerDB, query, role)

	if err != nil {
		return nil, repository_errors.SelectError
	}

	var workerModels []models.Worker
	for i := range workerDB {
		worker := copyWorkerResultToModel(&workerDB[i])
		workerModels = append(workerModels, *worker)
	}

	return workerModels, nil
}

// GetAverageOrderRate calculates the average rating for completed orders
// assigned to a specific worker.
//
// Parameters:
//   - worker: Worker entity to calculate average rating for
//
// Returns:
//   - float64: Average rating value (0.0-5.0)
//   - error: repository_errors.SelectError if the operation fails
func (w WorkerRepository) GetAverageOrderRate(worker *models.Worker) (float64, error) {
	query := `SELECT AVG(rate) FROM orders WHERE worker_id = $1 AND status = 3 AND rate != 0;`
	var averageRate float64

	err := w.db.Get(&averageRate, query, worker.ID)

	if err != nil {
		return 0, repository_errors.SelectError
	}

	return averageRate, nil
}
