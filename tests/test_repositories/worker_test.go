package test_repositories

import (
	"context"
	"fmt"
	"lab3/internal/models"
	"lab3/internal/repository/postgres"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

var testWorkerRepositoryCreateSuccess = []struct {
	TestName  string
	InputData struct {
		worker *models.Worker
	}
	CheckOutput func(t *testing.T, inputData struct{ worker *models.Worker }, createdWorker *models.Worker, err error)
}{
	{
		TestName: "create success test",
		InputData: struct {
			worker *models.Worker
		}{
			&models.Worker{
				ID:          uuid.New(),
				Name:        "First Name",
				Surname:     "Last Name",
				Address:     "Address",
				PhoneNumber: "+79999999999",
				Email:       "test@email.com",
				Password:    "hashed_password",
				Role:        1,
			},
		},
		CheckOutput: func(t *testing.T, inputData struct{ worker *models.Worker }, createdWorker *models.Worker, err error) {
			require.NoError(t, err)
			require.Equal(t, inputData.worker.Name, createdWorker.Name)
			require.Equal(t, inputData.worker.Surname, createdWorker.Surname)
			require.Equal(t, inputData.worker.Address, createdWorker.Address)
			require.Equal(t, inputData.worker.PhoneNumber, createdWorker.PhoneNumber)
			require.Equal(t, inputData.worker.Email, createdWorker.Email)
			require.Equal(t, inputData.worker.Password, createdWorker.Password)
			require.Equal(t, inputData.worker.Role, createdWorker.Role)
		},
	},
}

func TestWorkerRepositoryCreate(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}

	for _, test := range testWorkerRepositoryCreateSuccess {
		workerRepository := postgres.CreateWorkerRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			createdWorker, err := workerRepository.Create(test.InputData.worker)
			test.CheckOutput(t, test.InputData, createdWorker, err)
		})
	}
}

var testWorkerRepositoryGetByIDSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdWorker *models.Worker, receivedWorker *models.Worker, err error)
}{
	{
		TestName: "get by id success test",
		CheckOutput: func(t *testing.T, createdWorker *models.Worker, receivedWorker *models.Worker, err error) {
			require.NoError(t, err)
			require.Equal(t, createdWorker.ID, receivedWorker.ID)
		},
	},
}

func TestWorkerRepositoryGetByID(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}
	workerRepository := postgres.CreateWorkerRepository(&fields)

	for _, test := range testWorkerRepositoryGetByIDSuccess {
		t.Run(test.TestName, func(t *testing.T) {
			createdWorker, err := workerRepository.Create(&models.Worker{
				Name:        "First Name",
				Surname:     "Last Name",
				Address:     "Address",
				PhoneNumber: "+79999999999",
				Email:       "email@test.ru",
				Password:    "hashed_password",
				Role:        1,
			})
			require.NoError(t, err)

			receivedWorker, err := workerRepository.GetWorkerByID(createdWorker.ID)
			test.CheckOutput(t, createdWorker, receivedWorker, err)
		})
	}
}

var testWorkerRepositoryGetByEmailSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdWorker *models.Worker, receivedWorker *models.Worker, err error)
}{
	{
		TestName: "get by email success test",
		CheckOutput: func(t *testing.T, createdWorker *models.Worker, receivedWorker *models.Worker, err error) {
			require.NoError(t, err)
			require.Equal(t, createdWorker.Email, receivedWorker.Email)
		},
	},
}

func TestWorkerRepositoryGetByEmail(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}
	workerRepository := postgres.CreateWorkerRepository(&fields)

	for _, test := range testWorkerRepositoryGetByEmailSuccess {
		t.Run(test.TestName, func(t *testing.T) {
			createdWorker, err := workerRepository.Create(&models.Worker{
				Name:        "First Name",
				Surname:     "Last Name",
				Address:     "Address",
				PhoneNumber: "+79999999999",
				Email:       "email@email.com",
				Password:    "hashed_password",
				Role:        1,
			})
			require.NoError(t, err)

			receivedWorker, err := workerRepository.GetWorkerByEmail(createdWorker.Email)
			test.CheckOutput(t, createdWorker, receivedWorker, err)
		})
	}
}

var testWorkerRepositoryUpdateSuccess = []struct {
	TestName  string
	InputData struct {
		worker *models.Worker
	}
	CheckOutput func(t *testing.T, inputData struct{ worker *models.Worker }, updatedWorker *models.Worker, err error)
}{
	{
		TestName: "update success test",
		InputData: struct {
			worker *models.Worker
		}{
			&models.Worker{
				ID:          uuid.New(),
				Name:        "First Name",
				Surname:     "Last Name",
				Address:     "Address",
				PhoneNumber: "+79999999999",
				Email:       "email@email.com",
				Password:    "hashed_password",
				Role:        1,
			},
		},
		CheckOutput: func(t *testing.T, inputData struct{ worker *models.Worker }, updatedWorker *models.Worker, err error) {
			require.NoError(t, err)
			require.NotEqual(t, inputData.worker.Name, updatedWorker.Name)
			require.NotEqual(t, inputData.worker.Surname, updatedWorker.Surname)
			require.NotEqual(t, inputData.worker.Address, updatedWorker.Address)
			require.NotEqual(t, inputData.worker.PhoneNumber, updatedWorker.PhoneNumber)
			require.NotEqual(t, inputData.worker.Email, updatedWorker.Email)
			require.NotEqual(t, inputData.worker.Password, updatedWorker.Password)
			require.NotEqual(t, inputData.worker.Role, updatedWorker.Role)
		},
	},
}

func TestWorkerRepositoryUpdate(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}
	workerRepository := postgres.CreateWorkerRepository(&fields)

	for _, test := range testWorkerRepositoryUpdateSuccess {
		t.Run(test.TestName, func(t *testing.T) {
			createdWorker, err := workerRepository.Create(test.InputData.worker)
			require.NoError(t, err)

			createdWorker.Name = "Updated Name"
			createdWorker.Surname = "Updated Surname"
			createdWorker.Address = "Updated Address"
			createdWorker.PhoneNumber = "+79999999998"
			createdWorker.Email = "email1@email.com"
			createdWorker.Password = "updated_hashed_password"
			createdWorker.Role = 2

			updatedWorker, err := workerRepository.Update(createdWorker)
			test.CheckOutput(t, test.InputData, updatedWorker, err)
		})
	}
}

var testWorkerRepositoryDeleteSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdWorker *models.Worker, err error)
}{
	{
		TestName: "delete success test",
		CheckOutput: func(t *testing.T, createdWorker *models.Worker, err error) {
			require.NoError(t, err)
		},
	},
}

var testWorkerRepositoryDeleteFail = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdWorker *models.Worker, err error)
}{
	{
		TestName: "delete fail test",
		CheckOutput: func(t *testing.T, createdWorker *models.Worker, err error) {
			require.Nil(t, err)
		},
	},
}

func TestWorkerRepositoryDelete(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}
	workerRepository := postgres.CreateWorkerRepository(&fields)

	for _, test := range testWorkerRepositoryDeleteSuccess {
		t.Run(test.TestName, func(t *testing.T) {
			createdWorker, err := workerRepository.Create(&models.Worker{
				Name:        "First Name",
				Surname:     "Last Name",
				Address:     "Address",
				PhoneNumber: "+79999999999",
				Email:       "test@email.com",
				Password:    "hashed_password",
				Role:        1,
			})
			require.NoError(t, err)

			err = workerRepository.Delete(createdWorker.ID)
			test.CheckOutput(t, createdWorker, err)

			_, err = workerRepository.GetWorkerByID(createdWorker.ID)
			require.Error(t, err)
		})
	}

	for _, test := range testWorkerRepositoryDeleteFail {
		t.Run(test.TestName, func(t *testing.T) {
			createdWorker, err := workerRepository.Create(&models.Worker{
				Name:        "First Name",
				Surname:     "Last Name",
				Address:     "Address",
				PhoneNumber: "+79999999999",
				Email:       "test@email.com",
				Password:    "hashed_password",
				Role:        1,
			})
			require.NoError(t, err)

			err = workerRepository.Delete(uuid.New())
			test.CheckOutput(t, createdWorker, err)
		})
	}
}

var testWorkerRepositoryGetAllSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdWorkers []models.Worker, receivedWorkers []models.Worker, err error)
}{
	{
		TestName: "get all success test",
		CheckOutput: func(t *testing.T, createdWorkers []models.Worker, receivedWorkers []models.Worker, err error) {
			require.NoError(t, err)
			require.Equal(t, len(createdWorkers), len(receivedWorkers))
		},
	},
}

func TestWorkerRepositoryGetAll(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}
	workerRepository := postgres.CreateWorkerRepository(&fields)

	for i, test := range testWorkerRepositoryGetAllSuccess {
		t.Run(test.TestName, func(t *testing.T) {
			createdWorkers := []models.Worker{
				{
					Name:        fmt.Sprintf("First Name %d", i+1),
					Surname:     fmt.Sprintf("Last Name %d", i+1),
					Address:     fmt.Sprintf("Address   %d", i+1),
					PhoneNumber: fmt.Sprintf("+7999999999%d", i),
					Email:       fmt.Sprintf("test%d@email.com", i),
					Password:    "hashed_password",
					Role:        1,
				},
			}
			for _, worker := range createdWorkers {
				_, err := workerRepository.Create(&worker)
				require.NoError(t, err)
			}

			receivedWorkers, err := workerRepository.GetAllWorkers()
			test.CheckOutput(t, createdWorkers, receivedWorkers, err)
		})
	}
}
