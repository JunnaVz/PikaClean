package test_repositories

import (
	"context"
	"lab3/internal/models"
	"lab3/internal/repository/postgres"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

var testTaskRepositoryCreateSuccess = []struct {
	TestName  string
	InputData struct {
		ID       uuid.UUID
		Name     string
		Price    float64
		Category int
	}
	CheckOutput func(t *testing.T, inputData struct {
		ID       uuid.UUID
		Name     string
		Price    float64
		Category int
	}, createdTask *models.Task, err error)
}{
	{
		TestName: "create success test",
		InputData: struct {
			ID       uuid.UUID
			Name     string
			Price    float64
			Category int
		}{
			uuid.New(),
			"TaskName",
			100.0,
			1,
		},
		CheckOutput: func(t *testing.T, inputData struct {
			ID       uuid.UUID
			Name     string
			Price    float64
			Category int
		}, createdTask *models.Task, err error) {
			require.NoError(t, err)
			require.Equal(t, inputData.Name, createdTask.Name)
			require.Equal(t, inputData.Price, createdTask.PricePerSingle)
			require.Equal(t, inputData.Category, createdTask.Category)
		},
	},
}

func TestTaskRepositoryCreate(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}

	for _, test := range testTaskRepositoryCreateSuccess {
		taskRepository := postgres.CreateTaskRepository(&fields)

		createdTask, err := taskRepository.Create(
			&models.Task{
				Name:           test.InputData.Name,
				PricePerSingle: test.InputData.Price,
				Category:       test.InputData.Category,
			},
		)

		test.CheckOutput(t, test.InputData, createdTask, err)
	}
}

var testTaskRepositoryGetByIDSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdTask *models.Task, receivedTask *models.Task, err error)
}{
	{
		TestName: "get by id success test",
		CheckOutput: func(t *testing.T, createdTask *models.Task, receivedTask *models.Task, err error) {
			require.NoError(t, err)
			require.Equal(t, createdTask.ID, receivedTask.ID)
		},
	},
}

func TestTaskRepositoryGetByID(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}
	taskRepository := postgres.CreateTaskRepository(&fields)

	for _, test := range testTaskRepositoryGetByIDSuccess {
		createdTask, err := taskRepository.Create(
			&models.Task{
				Name:           "TaskName",
				PricePerSingle: 100.0,
				Category:       1,
			},
		)

		receivedTask, err := taskRepository.GetTaskByID(createdTask.ID)
		test.CheckOutput(t, createdTask, receivedTask, err)
	}
}

var testTaskRepositoryDeleteSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "delete success test",
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

var testTaskRepositoryDeleteFailure = []struct {
	TestName    string
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "delete non-existent task test",
		CheckOutput: func(t *testing.T, err error) {
			require.Error(t, err)
		},
	},
}

func TestTaskRepositoryDelete(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}
	taskRepository := postgres.CreateTaskRepository(&fields)

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

var testTaskRepositoryUpdateSuccess = []struct {
	TestName  string
	InputData struct {
		Task *models.Task
	}
	CheckOutput func(t *testing.T, createdTask *models.Task, updatedTask *models.Task, err error)
}{
	{
		TestName: "update success test",
		InputData: struct {
			Task *models.Task
		}{
			&models.Task{
				Name:           "TaskName",
				PricePerSingle: 100.0,
				Category:       1,
			},
		},
		CheckOutput: func(t *testing.T, createdTask *models.Task, updatedTask *models.Task, err error) {
			require.NoError(t, err)
			require.Equal(t, createdTask.ID, updatedTask.ID)
			require.NotEqual(t, createdTask.Name, updatedTask.Name)
			require.NotEqual(t, createdTask.PricePerSingle, updatedTask.PricePerSingle)
			require.NotEqual(t, createdTask.Category, updatedTask.Category)
		},
	},
}

func TestTaskRepositoryUpdate(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}
	taskRepository := postgres.CreateTaskRepository(&fields)

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

var testTaskRepositoryGetTaskByName = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdTask *models.Task, receivedTask *models.Task, err error)
}{
	{
		TestName: "get by name success test",
		CheckOutput: func(t *testing.T, createdTask *models.Task, receivedTask *models.Task, err error) {
			require.NoError(t, err)
			require.Equal(t, createdTask.ID, receivedTask.ID)
		},
	},
}

func TestTaskRepositoryGetTaskByName(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}
	taskRepository := postgres.CreateTaskRepository(&fields)

	for _, test := range testTaskRepositoryGetTaskByName {
		createdTask, err := taskRepository.Create(
			&models.Task{
				Name:           "TaskName",
				PricePerSingle: 100.0,
				Category:       1,
			},
		)

		receivedTask, err := taskRepository.GetTaskByName(createdTask.Name)
		test.CheckOutput(t, createdTask, receivedTask, err)
	}
}

var testTaskRepositoryGetAllTasks = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdTasks []models.Task, receivedTasks []models.Task, err error)
}{
	{
		TestName: "get all tasks success test",
		CheckOutput: func(t *testing.T, createdTasks []models.Task, receivedTasks []models.Task, err error) {
			require.NoError(t, err)
			require.Equal(t, len(createdTasks), len(receivedTasks))
		},
	},
}

func TestTaskRepositoryGetAllTasks(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}
	taskRepository := postgres.CreateTaskRepository(&fields)

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

var testTaskRepositoryGetAllTasksByCategory = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdTasks []models.Task, receivedTasks []models.Task, err error)
}{
	{
		TestName: "get all tasks by category success test",
		CheckOutput: func(t *testing.T, createdTasks []models.Task, receivedTasks []models.Task, err error) {
			require.NoError(t, err)
			require.Equal(t, len(createdTasks), len(receivedTasks))
			for _, task := range receivedTasks {
				require.Equal(t, createdTasks[0].Category, task.Category)
			}
		},
	},
}

func TestTaskRepositoryGetAllTasksByCategory(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}
	taskRepository := postgres.CreateTaskRepository(&fields)

	for _, test := range testTaskRepositoryGetAllTasksByCategory {
		createdTasks := []models.Task{
			{
				Name:           "TaskName1",
				PricePerSingle: 100.0,
				Category:       1,
			},
			{
				Name:           "TaskName2",
				PricePerSingle: 200.0,
				Category:       1,
			},
			{
				Name:           "TaskName3",
				PricePerSingle: 300.0,
				Category:       1,
			},
		}

		_ = []models.Task{
			{
				Name:           "TaskName4",
				PricePerSingle: 400.0,
				Category:       2,
			},
			{
				Name:           "TaskName5",
				PricePerSingle: 500.0,
				Category:       3,
			},
		}

		for i := range createdTasks {
			_, err := taskRepository.Create(&createdTasks[i])
			require.NoError(t, err)
		}

		receivedTasks, err := taskRepository.GetTasksInCategory(1)
		test.CheckOutput(t, createdTasks, receivedTasks, err)
	}
}
