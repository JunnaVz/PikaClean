package interfaces

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"lab3/internal/models"
	"lab3/internal/repository/repository_interfaces"
	"lab3/internal/services/service_interfaces"
	"lab3/password_hash"
)

type WorkerService struct {
	WorkerRepository repository_interfaces.IWorkerRepository
	hash             password_hash.PasswordHash
	logger           *log.Logger
}

func NewWorkerService(WorkerRepository repository_interfaces.IWorkerRepository, hash password_hash.PasswordHash, logger *log.Logger) service_interfaces.IWorkerService {
	return &WorkerService{
		WorkerRepository: WorkerRepository,
		hash:             hash,
		logger:           logger,
	}
}

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

func (w WorkerService) GetWorkerByID(id uuid.UUID) (*models.Worker, error) {
	worker, err := w.WorkerRepository.GetWorkerByID(id)

	if err != nil {
		w.logger.Error("SERVICE: GetWorkerByID method failed", "id", id, "error", err)
		return nil, err
	}

	w.logger.Info("SERVICE: Successfully got user with GetWorkerByID", "id", id)
	return worker, nil
}

func (w WorkerService) GetAllWorkers() ([]models.Worker, error) {
	workers, err := w.WorkerRepository.GetAllWorkers()

	if err != nil {
		w.logger.Error("SERVICE: GetAllWorkers method failed", "error", err)
		return nil, err
	}

	w.logger.Info("SERVICE: Successfully got all workers")
	return workers, nil
}

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

func (w WorkerService) GetWorkersByRole(role int) ([]models.Worker, error) {
	workers, err := w.WorkerRepository.GetWorkersByRole(role)

	if err != nil {
		w.logger.Error("SERVICE: GetWorkersByRole method failed", "error", err)
		return nil, err
	}

	w.logger.Info("SERVICE: Successfully got workers by role")
	return workers, nil

}

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
