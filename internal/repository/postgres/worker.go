package postgres

import (
	"database/sql"
	"errors"
	"lab3/internal/models"
	"lab3/internal/repository/repository_errors"
	"lab3/internal/repository/repository_interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type WorkerDB struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Surname     string    `db:"surname"`
	Address     string    `db:"address"`
	PhoneNumber string    `db:"phone_number"`
	Email       string    `db:"email"`
	Role        int       `db:"role"`
	Password    string    `db:"password"`
}

type WorkerRepository struct {
	db *sqlx.DB
}

func NewWorkerRepository(db *sqlx.DB) repository_interfaces.IWorkerRepository {
	return &WorkerRepository{db: db}
}

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

func (w WorkerRepository) Update(worker *models.Worker) (*models.Worker, error) {
	query := `UPDATE workers SET name = $1, surname = $2, address = $3, phone_number = $4, email = $5, role = $6, password = $7 WHERE workers.id = $8 RETURNING id, name, surname, address, phone_number, email, role, password;`

	var updatedWorker models.Worker
	err := w.db.QueryRow(query, worker.Name, worker.Surname, worker.Address, worker.PhoneNumber, worker.Email, worker.Role, worker.Password, worker.ID).Scan(&updatedWorker.ID, &updatedWorker.Name, &updatedWorker.Surname, &updatedWorker.Address, &updatedWorker.PhoneNumber, &updatedWorker.Email, &updatedWorker.Role, &updatedWorker.Password)
	if err != nil {
		return nil, repository_errors.UpdateError
	}
	return &updatedWorker, nil
}

func (w WorkerRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM workers WHERE id = $1;`
	_, err := w.db.Exec(query, id)

	if err != nil {
		return repository_errors.DeleteError
	}

	return nil
}

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

func (w WorkerRepository) GetAverageOrderRate(worker *models.Worker) (float64, error) {
	query := `SELECT AVG(rate) FROM orders WHERE worker_id = $1 AND status = 3 AND rate != 0;`
	var averageRate float64

	err := w.db.Get(&averageRate, query, worker.ID)

	if err != nil {
		return 0, repository_errors.SelectError
	}

	return averageRate, nil

}
