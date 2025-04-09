package test_repositories

import (
	"lab3/internal/models"
	"lab3/internal/repository/mongodb"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	"context"
)

func mongoCreateUser(fields *mongodb.MongoConnection) *models.User {
	user, _ := mongodb.CreateUserRepository(fields).Create(&models.User{
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

func mongoCreateWorker(fields *mongodb.MongoConnection) *models.Worker {
	worker, _ := mongodb.CreateWorkerRepository(fields).Create(&models.Worker{
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

func mongoCreateTasks(fields *mongodb.MongoConnection) []models.OrderedTask {
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
		createdTask, _ := mongodb.CreateTaskRepository(fields).Create(task.Task)
		createdTasks = append(createdTasks, models.OrderedTask{
			Task:     createdTask,
			Quantity: tasksModels[i].Quantity,
		})
	}
	return createdTasks
}
func TestMongoOrderRepositoryCreate(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}

	for _, test := range testOrderRepositoryCreateSuccess {
		orderRepository := mongodb.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := mongoCreateUser(&fields)
			worker := mongoCreateWorker(&fields)
			tasks := mongoCreateTasks(&fields)

			test.InputData.Order.UserID = user.ID
			test.InputData.Order.WorkerID = worker.ID

			createdOrder, err := orderRepository.Create(test.InputData.Order, tasks)
			test.CheckOutput(t, test.InputData.Order, createdOrder, err)
		})
	}
}

func TestMongoOrderRepositoryGetByID(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}

	for _, test := range testOrderRepositoryGetByIDSuccess {
		orderRepository := mongodb.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := mongoCreateUser(&fields)
			worker := mongoCreateWorker(&fields)
			createdTasks := mongoCreateTasks(&fields)

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

func TestMongoOrderRepositoryDelete(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}

	for _, test := range testOrderRepositoryDeleteSuccess {
		orderRepository := mongodb.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := mongoCreateUser(&fields)
			worker := mongoCreateWorker(&fields)
			tasks := mongoCreateTasks(&fields)

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
		orderRepository := mongodb.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			err := orderRepository.Delete(uuid.New())
			test.CheckOutput(t, err)
		})
	}
}

func TestMongoOrderRepositoryUpdate(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}

	for _, test := range testOrderRepositoryUpdateSuccess {
		orderRepository := mongodb.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := mongoCreateUser(&fields)
			worker := mongoCreateWorker(&fields)
			tasks := mongoCreateTasks(&fields)

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

func TestMongoOrderRepositoryGetTasksInOrder(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}

	for _, test := range testOrderRepositoryGetTasksInOrderSuccess {
		orderRepository := mongodb.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := mongoCreateUser(&fields)
			worker := mongoCreateWorker(&fields)
			createdTasks := mongoCreateTasks(&fields)

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

func TestMongoOrderRepositoryGetCurrentOrderByUserID(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}

	for _, test := range testOrderRepositoryGetCurrentOrderByUserIDSuccess {
		orderRepository := mongodb.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := mongoCreateUser(&fields)
			worker := mongoCreateWorker(&fields)
			createdTasks := mongoCreateTasks(&fields)

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

func TestMongoOrderRepositoryGetAllOrdersByUserID(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}

	for _, test := range testOrderRepositoryGetAllOrdersByUserIDSuccess {
		orderRepository := mongodb.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := mongoCreateUser(&fields)
			worker := mongoCreateWorker(&fields)
			createdTasks := mongoCreateTasks(&fields)

			createdOrder, _ := orderRepository.Create(&models.Order{
				ID:       uuid.New(),
				WorkerID: worker.ID,
				UserID:   user.ID,
				Status:   1,
				Address:  "Address",
				Deadline: time.Now().AddDate(0, 0, 1),
				Rate:     0,
			}, createdTasks)

			createdTasks2 := mongoCreateTasks(&fields)

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

func TestMongoOrderRepositoryAddTaskToOrder(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}

	for _, test := range testOrderRepositoryAddTaskToOrderSuccess {
		orderRepository := mongodb.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := mongoCreateUser(&fields)
			worker := mongoCreateWorker(&fields)
			createdTasks := mongoCreateTasks(&fields)

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

			createdTask, _ := mongodb.CreateTaskRepository(&fields).Create(task)

			err := orderRepository.AddTaskToOrder(createdOrder.ID, createdTask.ID)
			test.CheckOutput(t, createdOrder, err)

			tasks, err := orderRepository.GetTasksInOrder(createdOrder.ID)
			require.NoError(t, err)
			require.Equal(t, 3, len(tasks))
		})
	}
}

func TestMongoOrderRepositoryRemoveTaskFromOrder(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}

	for _, test := range testOrderRepositoryRemoveTaskFromOrderSuccess {
		orderRepository := mongodb.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := mongoCreateUser(&fields)
			worker := mongoCreateWorker(&fields)
			createdTasks := mongoCreateTasks(&fields)

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

func TestMongoOrderRepositoryUpdateTaskQuantity(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}

	for _, test := range testOrderRepositoryUpdateTaskQuantitySuccess {
		orderRepository := mongodb.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := mongoCreateUser(&fields)
			worker := mongoCreateWorker(&fields)
			createdTasks := mongoCreateTasks(&fields)

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

func TestMongoOrderRepositoryGetTaskQuantity(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}

	for _, test := range testOrderRepositoryGetTaskQuantitySuccess {
		orderRepository := mongodb.CreateOrderRepository(&fields)
		t.Run(test.TestName, func(t *testing.T) {
			user := mongoCreateUser(&fields)
			worker := mongoCreateWorker(&fields)
			createdTasks := mongoCreateTasks(&fields)

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
