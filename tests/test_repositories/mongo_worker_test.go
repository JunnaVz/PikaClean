package test_repositories

import (
	"fmt"
	"lab3/internal/models"
	"lab3/internal/repository/mongodb"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	"context"
)

func TestMongoWorkerRepositoryCreate(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}

	for _, test := range testWorkerRepositoryCreateSuccess {
		workerRepository := mongodb.CreateWorkerRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			createdWorker, err := workerRepository.Create(test.InputData.worker)
			test.CheckOutput(t, test.InputData, createdWorker, err)
		})
	}
}

func TestMongoWorkerRepositoryGetByID(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}
	workerRepository := mongodb.CreateWorkerRepository(&fields)

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

func TestMongoWorkerRepositoryGetByEmail(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}
	workerRepository := mongodb.CreateWorkerRepository(&fields)

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

func TestMongoWorkerRepositoryUpdate(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}
	workerRepository := mongodb.CreateWorkerRepository(&fields)

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

func TestMongoWorkerRepositoryDelete(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}
	workerRepository := mongodb.CreateWorkerRepository(&fields)

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

func TestMongoWorkerRepositoryGetAll(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}
	workerRepository := mongodb.CreateWorkerRepository(&fields)

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
