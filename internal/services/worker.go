// Package interfaces provides service implementations for the business logic layer
// of the PikaClean application. This package contains the core functionality that
// connects the repository layer with the application's domain models.
package interfaces

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"teamdev/internal/models"
	"teamdev/internal/repository/repository_interfaces"
	"teamdev/internal/services/service_interfaces"
	"teamdev/password_hash"
)

// WorkerService implements the IWorkerService interface to handle worker-related
// business logic operations. It orchestrates worker data operations and enforces
// business rules related to worker management.
type WorkerService struct {
	WorkerRepository repository_interfaces.IWorkerRepository // Repository for worker persistence
	hash             password_hash.PasswordHash              // Password hashing utility
	logger           *log.Logger                             // Logger for tracking operations
}

// NewWorkerService creates and initializes a new WorkerService with the provided dependencies.
//
// Parameters:
//   - WorkerRepository: Repository for accessing worker data
//   - hash: Utility for password hashing and verification
//   - logger: Logger for operation tracking and error reporting
//
// Returns:
//   - service_interfaces.IWorkerService: Initialized worker service implementation
func NewWorkerService(WorkerRepository repository_interfaces.IWorkerRepository, hash password_hash.PasswordHash, logger *log.Logger) service_interfaces.IWorkerService {
	return &WorkerService{
		WorkerRepository: WorkerRepository,
		hash:             hash,
		logger:           logger,
	}
}

// checkIfWorkerWithEmailExists verifies if a worker with the specified email exists in the system.
//
// Parameters:
//   - email: Email address to check for existence
//
// Returns:
//   - *models.Worker: Worker model if found, nil if not found
//   - error: Error that occurred during repository operation, nil if successful
func (w WorkerService) checkIfWorkerWithEmailExists(email string) (*models.Worker, error) {
	w.logger.Info("SERVICE: Checking if worker with email exists", "email", email)
	tempWorker, err := w.WorkerRepository.GetWorkerByEmail(email)

	if err != nil && err.Error() == "GET operation has failed. Such row does not exist" {
		w.logger.Info("SERVICE: Worker with email does not exist", "email", email)
		return nil, nil
	} else if err != nil {
		w.logger.Error("SERVICE: GetWorkerByEmail method failed", "email", email, "error", err)
		return nil, err
	} else {
		w.logger.Info("SERVICE: Worker with email exists", "email", email)
		return tempWorker, nil
	}
}

// Login authenticates a worker using email and password credentials.
//
// Parameters:
//   - email: Worker's email address for identification
//   - password: Worker's password for verification
//
// Returns:
//   - *models.Worker: Authenticated worker if credentials are valid
//   - error: Authentication error or repository error, nil if successful
func (w WorkerService) Login(email, password string) (*models.Worker, error) {
	w.logger.Infof("SERVICE: Checking if worker with email %s exists", email)
	tempWorker, err := w.checkIfWorkerWithEmailExists(email)
	if err != nil {
		w.logger.Error("SERVICE: Error occurred during checking if worker with email exists")
		return nil, err
	} else if tempWorker == nil {
		w.logger.Info("SERVICE: Worker with email does not exist")
		return nil, fmt.Errorf("SERVICE: Worker with email does not exist")
	}

	w.logger.Infof("SERVICE: Checking if password is correct for worker with email %s", email)
	isPasswordCorrect := w.hash.CompareHashAndPassword(tempWorker.Password, password)
	if !isPasswordCorrect {
		w.logger.Info("SERVICE: Password is incorrect for worker with email")
		return nil, fmt.Errorf("SERVICE: Password is incorrect for worker with email")
	}

	w.logger.Info("SERVICE: Successfully logged in worker with email", "email", email)
	return tempWorker, nil
}

// Create registers a new worker in the system with validation of input data.
//
// Parameters:
//   - worker: Worker model with personal information to be registered
//   - password: Plain text password to be hashed and stored
//
// Returns:
//   - *models.Worker: Created worker with assigned ID if successful
//   - error: Validation error or repository error, nil if successful
func (w WorkerService) Create(worker *models.Worker, password string) (*models.Worker, error) {
	w.logger.Info("SERVICE: Validating data")
	if !validName(worker.Name) || !validName(worker.Surname) || !validEmail(worker.Email) || !validAddress(worker.Address) || !validPhoneNumber(worker.PhoneNumber) || !validRole(worker.Role) || !validPassword(password) {
		w.logger.Error("SERVICE: Invalid input")
		return nil, fmt.Errorf("SERVICE: Invalid input")
	}

	w.logger.Infof("SERVICE: Checking if worker with email %s exists", worker.Email)
	tempWorker, err := w.checkIfWorkerWithEmailExists(worker.Email)
	if err != nil {
		w.logger.Error("SERVICE: Error occurred during checking if worker with email exists")
		return nil, err
	} else if tempWorker != nil {
		w.logger.Info("SERVICE: Worker with email exists", "email", worker.Email)
		return nil, fmt.Errorf("SERVICE: Worker with email already exists")
	}

	w.logger.Infof("SERVICE: Creating new worker: %s %s", worker.Name, worker.Surname)
	hashedPassword, err := w.hash.GetHash(password)
	if err != nil {
		w.logger.Error("SERVICE: Error occurred during password hashing")
		return nil, err
	} else {
		worker.Password = hashedPassword
	}

	createdWorker, err := w.WorkerRepository.Create(worker)
	if err != nil {
		w.logger.Error("SERVICE: Create method failed", "error", err)
		return nil, err
	}

	w.logger.Info("SERVICE: Successfully created new user with ", "id", createdWorker.ID)

	return createdWorker, nil
}

// Delete removes a worker record from the system by ID.
//
// Parameters:
//   - id: UUID of the worker to be deleted
//
// Returns:
//   - error: Repository error if deletion fails, nil if successful
func (w WorkerService) Delete(id uuid.UUID) error {
	_, err := w.WorkerRepository.GetWorkerByID(id)
	if err != nil {
		w.logger.Error("SERVICE: GetWorkerByID method failed", "id", id, "error", err)
		return err
	}

	err = w.WorkerRepository.Delete(id)
	if err != nil {
		w.logger.Error("SERVICE: Delete method failed", "error", err)
	}

	w.logger.Info("SERVICE: Successfully deleted worker", "id", id)
	return nil
}

// GetWorkerByID retrieves a worker by their unique identifier.
//
// Parameters:
//   - id: UUID of the worker to retrieve
//
// Returns:
//   - *models.Worker: Retrieved worker if found
//   - error: Repository error if retrieval fails, nil if successful
func (w WorkerService) GetWorkerByID(id uuid.UUID) (*models.Worker, error) {
	worker, err := w.WorkerRepository.GetWorkerByID(id)

	if err != nil {
		w.logger.Error("SERVICE: GetWorkerByID method failed", "id", id, "error", err)
		return nil, err
	}

	w.logger.Info("SERVICE: Successfully got user with GetWorkerByID", "id", id)
	return worker, nil
}

// GetAllWorkers retrieves all workers registered in the system.
//
// Returns:
//   - []models.Worker: Slice of all worker records
//   - error: Repository error if retrieval fails, nil if successful
func (w WorkerService) GetAllWorkers() ([]models.Worker, error) {
	workers, err := w.WorkerRepository.GetAllWorkers()

	if err != nil {
		w.logger.Error("SERVICE: GetAllWorkers method failed", "error", err)
		return nil, err
	}

	w.logger.Info("SERVICE: Successfully got all workers")
	return workers, nil
}

// Update modifies a worker's information after validating the new data.
//
// Parameters:
//   - id: UUID of the worker to update
//   - name: New first name
//   - surname: New last name
//   - email: New email address
//   - address: New physical address
//   - phoneNumber: New contact phone number
//   - role: New role identifier
//   - password: New password (will be hashed if different from current)
//
// Returns:
//   - *models.Worker: Updated worker record
//   - error: Validation error or repository error, nil if successful
func (w WorkerService) Update(id uuid.UUID, name string, surname string, email string, address string, phoneNumber string, role int, password string) (*models.Worker, error) {
	worker, err := w.WorkerRepository.GetWorkerByID(id)
	if err != nil {
		w.logger.Error("SERVICE: GetUserByID method failed", "id", id, "error", err)
		return nil, err
	}

	if !validName(name) || !validName(surname) || !validEmail(email) || !validAddress(address) || !validPhoneNumber(phoneNumber) || !validRole(role) || !validPassword(password) {
		w.logger.Error("SERVICE: Invalid input")
		return nil, fmt.Errorf("SERVICE: Invalid input")
	} else {
		worker.Name = name
		worker.Surname = surname
		worker.Email = email
		worker.Address = address
		worker.PhoneNumber = phoneNumber
		worker.Role = role

		if password != worker.Password {
			hashedPassword, hashErr := w.hash.GetHash(password)
			if hashErr != nil {
				w.logger.Error("SERVICE: Error occurred during password hashing")
				return nil, hashErr
			} else {
				worker.Password = hashedPassword
			}
		}
	}

	worker, err = w.WorkerRepository.Update(worker)
	if err != nil {
		w.logger.Error("SERVICE: Update method failed", "error", err)
		return nil, err
	}

	w.logger.Info("SERVICE: Successfully updated worker personal information", "worker", worker)
	return worker, nil
}

// GetWorkersByRole retrieves all workers with a specific role.
//
// Parameters:
//   - role: Role identifier to filter by
//
// Returns:
//   - []models.Worker: Slice of workers with the specified role
//   - error: Repository error if retrieval fails, nil if successful
func (w WorkerService) GetWorkersByRole(role int) ([]models.Worker, error) {
	workers, err := w.WorkerRepository.GetWorkersByRole(role)

	if err != nil {
		w.logger.Error("SERVICE: GetWorkersByRole method failed", "error", err)
		return nil, err
	}

	w.logger.Info("SERVICE: Successfully got workers by role")
	return workers, nil
}

// GetAverageOrderRate calculates the average rating for a worker based on completed orders.
//
// Parameters:
//   - worker: Worker model to calculate average rating for
//
// Returns:
//   - float64: Average rating value (0.0-5.0)
//   - error: Repository error if calculation fails, nil if successful
func (w WorkerService) GetAverageOrderRate(worker *models.Worker) (float64, error) {
	_, err := w.WorkerRepository.GetWorkerByID(worker.ID)
	if err != nil {
		w.logger.Error("SERVICE: GetWorkerByID method failed", "id", worker.ID, "error", err)
		return 0, err
	}

	workerRate, err := w.WorkerRepository.GetAverageOrderRate(worker)

	if err != nil {
		w.logger.Error("SERVICE: GetAverageOrderRate method failed", "error", err)
		return 0, err
	}

	w.logger.Info("SERVICE: Successfully got average order rate for worker", "worker", worker)
	return workerRate, nil
}
