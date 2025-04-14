package test_repositories

import (
	"context"
	"teamdev/internal/models"
	"teamdev/internal/repository/postgres"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func createUser(fields *postgres.PostgresConnection) *models.User {
	user, _ := postgres.CreateUserRepository(fields).Create(&models.User{
		ID:          uuid.New(),
		Name:        "Name",
		Surname:     "Surname",
		Address:     "Address",
		PhoneNumber: "+79999999999",
		Email:       "user@email.com",
		Password:    "hashed_password",
	})
	return user
}

func createWorker(fields *postgres.PostgresConnection) *models.Worker {
	worker, _ := postgres.CreateWorkerRepository(fields).Create(&models.Worker{
		ID:          uuid.New(),
		Name:        "Worker Name",
		Surname:     "Worker Surname",
		Address:     "Worker Address",
		PhoneNumber: "+79999999998",
		Email:       "worker@email.com",
		Password:    "hashed_password",
		Role:        1,
	})
	return worker
}

func createTasks(fields *postgres.PostgresConnection) []models.OrderedTask {
	tasksModels := []models.OrderedTask{
		{
			Task: &models.Task{
				ID:             uuid.New(),
				Name:           "Task Name",
				PricePerSingle: 100,
				Category:       1,
			},
			Quantity: 2,
		},
		{
			Task: &models.Task{
				ID:             uuid.New(),
				Name:           "Task Name 2",
				PricePerSingle: 200,
				Category:       2,
			},
			Quantity: 1,
		},
	}
	createdTasks := make([]models.OrderedTask, 0)
	for i, task := range tasksModels {
		createdTask, _ := postgres.CreateTaskRepository(fields).Create(task.Task)
		createdTasks = append(createdTasks, models.OrderedTask{
			Task:     createdTask,
			Quantity: tasksModels[i].Quantity,
		})
	}
	return createdTasks
}

var testOrderRepositoryCreateSuccess = []struct {
	TestName  string
	InputData struct {
		Order        *models.Order
		OrderedTasks []models.OrderedTask
	}
	CheckOutput func(t *testing.T, inputData *models.Order, createdOrder *models.Order, err error)
}{
	{
		TestName: "create success test",
		InputData: struct {
			Order        *models.Order
			OrderedTasks []models.OrderedTask
		}{
			&models.Order{
				ID:       uuid.New(),
				WorkerID: uuid.New(),
				UserID:   uuid.New(),
				Status:   1,
				Address:  "Address",
				Deadline: time.Now().AddDate(0, 0, 1),
				Rate:     5,
			},
			[]models.OrderedTask{
				{
					Task: &models.Task{
						ID:             uuid.New(),
						Name:           "Task Name",
						PricePerSingle: 100,
						Category:       1,
					},
					Quantity: 2,
				},
				{
					Task: &models.Task{
						ID:             uuid.New(),
						Name:           "Task Name 2",
						PricePerSingle: 200,
						Category:       2,
					},
					Quantity: 1,
				},
			},
		},
		CheckOutput: func(t *testing.T, inputData *models.Order, createdOrder *models.Order, err error) {
			if err != nil {
				println(err.Error())
			}
			require.NoError(t, err)
			require.Equal(t, inputData.WorkerID, createdOrder.WorkerID)
			require.Equal(t, inputData.UserID, createdOrder.UserID)
			require.Equal(t, inputData.Status, createdOrder.Status)
			require.Equal(t, inputData.Address, createdOrder.Address)
			require.Equal(t, inputData.Deadline, createdOrder.Deadline)
			require.Equal(t, inputData.Rate, createdOrder.Rate)
		},
	},
}

func TestOrderRepositoryCreate(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}

	for _, test := range testOrderRepositoryCreateSuccess {
		orderRepository := postgres.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := createUser(&fields)
			worker := createWorker(&fields)
			tasks := createTasks(&fields)

			test.InputData.Order.UserID = user.ID
			test.InputData.Order.WorkerID = worker.ID

			createdOrder, err := orderRepository.Create(test.InputData.Order, tasks)
			test.CheckOutput(t, test.InputData.Order, createdOrder, err)
		})
	}
}

var testOrderRepositoryGetByIDSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdOrder *models.Order, receivedOrder *models.Order, err error)
}{
	{
		TestName: "get by id success test",
		CheckOutput: func(t *testing.T, createdOrder *models.Order, receivedOrder *models.Order, err error) {
			require.NoError(t, err)
			require.Equal(t, createdOrder.ID, receivedOrder.ID)
		},
	},
}

func TestOrderRepositoryGetByID(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}

	for _, test := range testOrderRepositoryGetByIDSuccess {
		orderRepository := postgres.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := createUser(&fields)
			worker := createWorker(&fields)
			createdTasks := createTasks(&fields)

			createdOrder, err := orderRepository.Create(&models.Order{
				ID:       uuid.New(),
				WorkerID: worker.ID,
				UserID:   user.ID,
				Status:   1,
				Address:  "Address",
				Deadline: time.Now().AddDate(0, 0, 1),
				Rate:     0,
			}, createdTasks)

			receivedOrder, err := orderRepository.GetOrderByID(createdOrder.ID)
			test.CheckOutput(t, createdOrder, receivedOrder, err)
		})
	}
}

var testOrderRepositoryDeleteSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdOrder *models.Order, err error)
}{
	{
		TestName: "delete success test",
		CheckOutput: func(t *testing.T, createdOrder *models.Order, err error) {
			require.NoError(t, err)
		},
	},
}

var testOrderRepositoryDeleteFailure = []struct {
	TestName    string
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "delete non-existent order test",
		CheckOutput: func(t *testing.T, err error) {
			require.Error(t, err)
		},
	},
}

func TestOrderRepositoryDelete(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}

	for _, test := range testOrderRepositoryDeleteSuccess {
		orderRepository := postgres.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := createUser(&fields)
			worker := createWorker(&fields)
			tasks := createTasks(&fields)

			createdOrder, err := orderRepository.Create(&models.Order{
				ID:       uuid.New(),
				WorkerID: worker.ID,
				UserID:   user.ID,
				Status:   1,
				Address:  "Address",
				Deadline: time.Now().AddDate(0, 0, 1),
				Rate:     0,
			}, tasks)

			err = orderRepository.Delete(createdOrder.ID)
			test.CheckOutput(t, createdOrder, err)

			_, err = orderRepository.GetOrderByID(createdOrder.ID)
			require.Error(t, err)
		})
	}

	for _, test := range testOrderRepositoryDeleteFailure {
		orderRepository := postgres.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			err := orderRepository.Delete(uuid.New())
			test.CheckOutput(t, err)
		})
	}
}

var testOrderRepositoryUpdateSuccess = []struct {
	TestName string

	CheckOutput func(t *testing.T, createdOrder *models.Order, updatedOrder *models.Order, err error)
}{
	{
		TestName: "update success test",
		CheckOutput: func(t *testing.T, createdOrder *models.Order, updatedOrder *models.Order, err error) {
			require.NoError(t, err)
			require.Equal(t, createdOrder.ID, updatedOrder.ID)
			require.NotEqual(t, createdOrder.Status, updatedOrder.Status)
			require.NotEqual(t, createdOrder.Address, updatedOrder.Address)
		},
	},
}

func TestOrderRepositoryUpdate(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}

	for _, test := range testOrderRepositoryUpdateSuccess {
		orderRepository := postgres.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := createUser(&fields)
			worker := createWorker(&fields)
			tasks := createTasks(&fields)

			createdOrder, _ := orderRepository.Create(&models.Order{
				ID:       uuid.New(),
				WorkerID: worker.ID,
				UserID:   user.ID,
				Status:   1,
				Address:  "Address",
				Deadline: time.Now().AddDate(0, 0, 1),
				Rate:     0,
			}, tasks)

			updatedOrder, err := orderRepository.Update(&models.Order{
				ID:       createdOrder.ID,
				WorkerID: createdOrder.WorkerID,
				UserID:   createdOrder.UserID,
				Status:   2,
				Address:  "New Address",
				Deadline: createdOrder.Deadline,
				Rate:     createdOrder.Rate,
			})
			test.CheckOutput(t, createdOrder, updatedOrder, err)
		})
	}
}

var testOrderRepositoryGetTasksInOrderSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdTasks []models.OrderedTask, receivedTasks []models.Task, err error)
}{
	{
		TestName: "get tasks in order success test",
		CheckOutput: func(t *testing.T, createdTasks []models.OrderedTask, receivedTasks []models.Task, err error) {
			require.NoError(t, err)
			require.Equal(t, len(createdTasks), len(receivedTasks))
		},
	},
}

func TestOrderRepositoryGetTasksInOrder(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}

	for _, test := range testOrderRepositoryGetTasksInOrderSuccess {
		orderRepository := postgres.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := createUser(&fields)
			worker := createWorker(&fields)
			createdTasks := createTasks(&fields)

			createdOrder, _ := orderRepository.Create(&models.Order{
				ID:       uuid.New(),
				WorkerID: worker.ID,
				UserID:   user.ID,
				Status:   1,
				Address:  "Address",
				Deadline: time.Now().AddDate(0, 0, 1),
				Rate:     0,
			}, createdTasks)

			receivedTasks, err := orderRepository.GetTasksInOrder(createdOrder.ID)
			test.CheckOutput(t, createdTasks, receivedTasks, err)
		})
	}
}

var testOrderRepositoryGetCurrentOrderByUserIDSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdOrder *models.Order, receivedOrder *models.Order, err error)
}{
	{
		TestName: "get current order by user id success test",
		CheckOutput: func(t *testing.T, createdOrder *models.Order, receivedOrder *models.Order, err error) {
			require.NoError(t, err)
			require.Equal(t, createdOrder.ID, receivedOrder.ID)
		},
	},
}

func TestOrderRepositoryGetCurrentOrderByUserID(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}

	for _, test := range testOrderRepositoryGetCurrentOrderByUserIDSuccess {
		orderRepository := postgres.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := createUser(&fields)
			worker := createWorker(&fields)
			createdTasks := createTasks(&fields)

			createdOrder, _ := orderRepository.Create(&models.Order{
				ID:       uuid.New(),
				WorkerID: worker.ID,
				UserID:   user.ID,
				Status:   1,
				Address:  "Address",
				Deadline: time.Now().AddDate(0, 0, 1),
				Rate:     0,
			}, createdTasks)

			receivedOrder, err := orderRepository.GetCurrentOrderByUserID(user.ID)
			test.CheckOutput(t, createdOrder, receivedOrder, err)
		})
	}
}

var testOrderRepositoryGetAllOrdersByUserIDSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdOrders []models.Order, receivedOrders []models.Order, err error)
}{
	{
		TestName: "get all orders by user id success test",
		CheckOutput: func(t *testing.T, createdOrders []models.Order, receivedOrders []models.Order, err error) {
			require.NoError(t, err)
			require.Equal(t, len(createdOrders), len(receivedOrders))
		},
	},
}

func TestOrderRepositoryGetAllOrdersByUserID(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}

	for _, test := range testOrderRepositoryGetAllOrdersByUserIDSuccess {
		orderRepository := postgres.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := createUser(&fields)
			worker := createWorker(&fields)
			createdTasks := createTasks(&fields)

			createdOrder, _ := orderRepository.Create(&models.Order{
				ID:       uuid.New(),
				WorkerID: worker.ID,
				UserID:   user.ID,
				Status:   1,
				Address:  "Address",
				Deadline: time.Now().AddDate(0, 0, 1),
				Rate:     0,
			}, createdTasks)

			createdTasks2 := createTasks(&fields)

			createdOrder2, _ := orderRepository.Create(&models.Order{
				ID:       uuid.New(),
				WorkerID: worker.ID,
				UserID:   user.ID,
				Status:   1,
				Address:  "Address",
				Deadline: time.Now().AddDate(0, 0, 1),
				Rate:     0,
			}, createdTasks2)

			receivedOrders, err := orderRepository.GetAllOrdersByUserID(user.ID)
			test.CheckOutput(t, []models.Order{*createdOrder, *createdOrder2}, receivedOrders, err)
		})
	}
}

var testOrderRepositoryAddTaskToOrderSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdOrder *models.Order, err error)
}{
	{
		TestName: "add task to order success test",
		CheckOutput: func(t *testing.T, createdOrder *models.Order, err error) {
			require.NoError(t, err)
		},
	},
}

func TestOrderRepositoryAddTaskToOrder(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}

	for _, test := range testOrderRepositoryAddTaskToOrderSuccess {
		orderRepository := postgres.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := createUser(&fields)
			worker := createWorker(&fields)
			createdTasks := createTasks(&fields)

			createdOrder, _ := orderRepository.Create(&models.Order{
				ID:       uuid.New(),
				WorkerID: worker.ID,
				UserID:   user.ID,
				Status:   1,
				Address:  "Address",
				Deadline: time.Now().AddDate(0, 0, 1),
				Rate:     0,
			}, createdTasks)

			task := &models.Task{
				ID:             uuid.New(),
				Name:           "Task Name",
				PricePerSingle: 100,
				Category:       1,
			}

			createdTask, _ := postgres.CreateTaskRepository(&fields).Create(task)

			err := orderRepository.AddTaskToOrder(createdOrder.ID, createdTask.ID)
			test.CheckOutput(t, createdOrder, err)

			tasks, err := orderRepository.GetTasksInOrder(createdOrder.ID)
			require.NoError(t, err)
			require.Equal(t, 3, len(tasks))
		})
	}
}

var testOrderRepositoryRemoveTaskFromOrderSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdOrder *models.Order, err error)
}{
	{
		TestName: "remove task from order success test",
		CheckOutput: func(t *testing.T, createdOrder *models.Order, err error) {
			require.NoError(t, err)
		},
	},
}

func TestOrderRepositoryRemoveTaskFromOrder(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}

	for _, test := range testOrderRepositoryRemoveTaskFromOrderSuccess {
		orderRepository := postgres.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := createUser(&fields)
			worker := createWorker(&fields)
			createdTasks := createTasks(&fields)

			createdOrder, _ := orderRepository.Create(&models.Order{
				ID:       uuid.New(),
				WorkerID: worker.ID,
				UserID:   user.ID,
				Status:   1,
				Address:  "Address",
				Deadline: time.Now().AddDate(0, 0, 1),
				Rate:     0,
			}, createdTasks)

			err := orderRepository.RemoveTaskFromOrder(createdOrder.ID, createdTasks[0].Task.ID)
			test.CheckOutput(t, createdOrder, err)

			tasks, err := orderRepository.GetTasksInOrder(createdOrder.ID)
			require.NoError(t, err)
			require.Equal(t, 1, len(tasks))
		})
	}
}

var testOrderRepositoryUpdateTaskQuantitySuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "update task quantity success test",
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

func TestOrderRepositoryUpdateTaskQuantity(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}

	for _, test := range testOrderRepositoryUpdateTaskQuantitySuccess {
		orderRepository := postgres.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := createUser(&fields)
			worker := createWorker(&fields)
			createdTasks := createTasks(&fields)

			createdOrder, _ := orderRepository.Create(&models.Order{
				ID:       uuid.New(),
				WorkerID: worker.ID,
				UserID:   user.ID,
				Status:   1,
				Address:  "Address",
				Deadline: time.Now().AddDate(0, 0, 1),
				Rate:     0,
			}, createdTasks)

			err := orderRepository.UpdateTaskQuantity(createdOrder.ID, createdTasks[0].Task.ID, 5)
			test.CheckOutput(t, err)

			quantityTask, err := orderRepository.GetTaskQuantity(createdOrder.ID, createdTasks[0].Task.ID)
			require.NoError(t, err)
			require.Equal(t, 5, quantityTask)
		})
	}
}

var testOrderRepositoryGetTaskQuantitySuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, quantity int, err error)
}{
	{
		TestName: "get task quantity success test",
		CheckOutput: func(t *testing.T, quantity int, err error) {
			require.NoError(t, err)
			require.Equal(t, 5, quantity)
		},
	},
}

func TestOrderRepositoryGetTaskQuantity(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}

	for _, test := range testOrderRepositoryGetTaskQuantitySuccess {
		orderRepository := postgres.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := createUser(&fields)
			worker := createWorker(&fields)
			createdTasks := createTasks(&fields)

			createdOrder, _ := orderRepository.Create(&models.Order{
				ID:       uuid.New(),
				WorkerID: worker.ID,
				UserID:   user.ID,
				Status:   1,
				Address:  "Address",
				Deadline: time.Now().AddDate(0, 0, 1),
				Rate:     0,
			}, createdTasks)

			err := orderRepository.UpdateTaskQuantity(createdOrder.ID, createdTasks[0].Task.ID, 5)
			require.NoError(t, err)

			quantityTask, err := orderRepository.GetTaskQuantity(createdOrder.ID, createdTasks[0].Task.ID)
			test.CheckOutput(t, quantityTask, err)
		})
	}
}
