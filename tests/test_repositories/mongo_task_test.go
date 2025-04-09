package test_repositories

import (
	"context"
	"fmt"
	"lab3/internal/models"
	"lab3/internal/repository/mongodb"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func TestTaskMongoRepositoryCreate(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}

	for _, test := range testTaskRepositoryCreateSuccess {
		taskRepository := mongodb.CreateTaskRepository(&fields)

		_, err := taskRepository.Create(
			&models.Task{
				Name:           test.InputData.Name,
				PricePerSingle: test.InputData.Price,
				Category:       test.InputData.Category,
			},
		)

		createdTask, _ := taskRepository.GetTaskByName(test.InputData.Name)

		test.CheckOutput(t, test.InputData, createdTask, err)
	}
}

func TestTaskMongoRepositoryGetByID(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}
	taskRepository := mongodb.CreateTaskRepository(&fields)

	for _, test := range testTaskRepositoryGetByIDSuccess {
		createdTask, err := taskRepository.Create(
			&models.Task{
				Name:           "TaskName",
				PricePerSingle: 100.0,
				Category:       1,
			},
		)

		if err != nil {
			t.Error(err)
		}

		tasks, err := taskRepository.GetAllTasks()
		fmt.Println(tasks)
		receivedTask, err := taskRepository.GetTaskByID(createdTask.ID)
		test.CheckOutput(t, createdTask, receivedTask, err)
	}
}

func TestTaskMongoRepositoryDelete(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}
	taskRepository := mongodb.CreateTaskRepository(&fields)

	for _, test := range testTaskRepositoryDeleteSuccess {
		createdTask, err := taskRepository.Create(
			&models.Task{
				Name:           "TaskName",
				PricePerSingle: 100.0,
				Category:       1,
			},
		)

		err = taskRepository.Delete(createdTask.ID)
		test.CheckOutput(t, err)

		_, err = taskRepository.GetTaskByID(createdTask.ID)
		require.Error(t, err)
	}

	for _, test := range testTaskRepositoryDeleteFailure {
		t.Run(test.TestName, func(t *testing.T) {
			err := taskRepository.Delete(uuid.New())
			test.CheckOutput(t, err)
		})
	}
}

func TestTaskMongoRepositoryUpdate(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}
	taskRepository := mongodb.CreateTaskRepository(&fields)

	for _, test := range testTaskRepositoryUpdateSuccess {
		createdTask, err := taskRepository.Create(test.InputData.Task)

		updatedTask, err := taskRepository.Update(
			&models.Task{
				ID:             createdTask.ID,
				Name:           "UpdatedTaskName",
				PricePerSingle: 200.0,
				Category:       2,
			},
		)

		test.CheckOutput(t, createdTask, updatedTask, err)
	}
}

func TestTaskMongoRepositoryGetTaskByName(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}
	taskRepository := mongodb.CreateTaskRepository(&fields)

	for _, test := range testTaskRepositoryGetTaskByName {
		createdTask, err := taskRepository.Create(
			&models.Task{
				ID:             uuid.New(),
				Name:           "TaskName1",
				PricePerSingle: 100.0,
				Category:       1,
			},
		)

		receivedTask, _ := taskRepository.GetTaskByName(createdTask.Name)
		test.CheckOutput(t, createdTask, receivedTask, err)
	}
}

func TestTaskMongoRepositoryGetAllTasks(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}
	taskRepository := mongodb.CreateTaskRepository(&fields)

	for _, test := range testTaskRepositoryGetAllTasks {
		createdTasks := []models.Task{
			{
				Name:           "TaskName1",
				PricePerSingle: 100.0,
				Category:       1,
			},
			{
				Name:           "TaskName2",
				PricePerSingle: 200.0,
				Category:       2,
			},
			{
				Name:           "TaskName3",
				PricePerSingle: 300.0,
				Category:       3,
			},
		}

		for i := range createdTasks {
			_, err := taskRepository.Create(&createdTasks[i])
			require.NoError(t, err)
		}

		receivedTasks, err := taskRepository.GetAllTasks()
		test.CheckOutput(t, createdTasks, receivedTasks, err)
	}
}
