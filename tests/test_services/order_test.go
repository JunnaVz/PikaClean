package test_services

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"lab3/internal/models"
	"lab3/internal/repository/repository_errors"
	services "lab3/internal/services"
	"lab3/internal/services/service_errors"
	"lab3/internal/services/service_interfaces"
	mock_repository_interfaces "lab3/tests/repository_mocks"
	"os"
	"testing"
	"time"
)

type orderServiceFields struct {
	orderRepoMock  *mock_repository_interfaces.MockIOrderRepository
	taskRepoMock   *mock_repository_interfaces.MockITaskRepository
	workerRepoMock *mock_repository_interfaces.MockIWorkerRepository
	userRepoMock   *mock_repository_interfaces.MockIUserRepository
	logger         *log.Logger
}

func initOrderServiceFields(ctrl *gomock.Controller) *orderServiceFields {
	orderRepoMock := mock_repository_interfaces.NewMockIOrderRepository(ctrl)
	taskRepoMock := mock_repository_interfaces.NewMockITaskRepository(ctrl)
	workerRepoMock := mock_repository_interfaces.NewMockIWorkerRepository(ctrl)
	userRepoMock := mock_repository_interfaces.NewMockIUserRepository(ctrl)

	f, err := os.OpenFile("tests.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(f)

	return &orderServiceFields{
		orderRepoMock:  orderRepoMock,
		taskRepoMock:   taskRepoMock,
		workerRepoMock: workerRepoMock,
		userRepoMock:   userRepoMock,
		logger:         logger,
	}
}

func initOrderService(fields *orderServiceFields) service_interfaces.IOrderService {
	return services.NewOrderService(fields.orderRepoMock, fields.workerRepoMock, fields.taskRepoMock, fields.userRepoMock, fields.logger)
}

var testOrderServiceCreate = []struct {
	testName  string
	inputData struct {
		userID   uuid.UUID
		address  string
		deadline time.Time
		tasks    []models.Task
	}
	prepare     func(fields *orderServiceFields)
	checkOutput func(t *testing.T, order *models.Order, err error)
}{
	{
		testName: "create order success test",
		inputData: struct {
			userID   uuid.UUID
			address  string
			deadline time.Time
			tasks    []models.Task
		}{
			uuid.New(),
			"address",
			time.Now().AddDate(0, 0, 1),
			[]models.Task{{ID: uuid.New()}, {ID: uuid.New()}},
		},
		prepare: func(fields *orderServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{}, nil).Times(2)
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(&models.User{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&models.Order{ID: uuid.New()}, nil)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, order)
		},
	},
	{
		testName: "empty tasks list",
		inputData: struct {
			userID   uuid.UUID
			address  string
			deadline time.Time
			tasks    []models.Task
		}{
			uuid.New(),
			"address",
			time.Now().AddDate(0, 0, 1),
			[]models.Task{},
		},
		prepare: func(fields *orderServiceFields) {},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Nil(t, order)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
	{
		testName: "invalid address",
		inputData: struct {
			userID   uuid.UUID
			address  string
			deadline time.Time
			tasks    []models.Task
		}{
			uuid.New(),
			"",
			time.Now().AddDate(0, 0, 1),
			[]models.Task{{ID: uuid.New()}, {ID: uuid.New()}},
		},
		prepare: func(fields *orderServiceFields) {},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Nil(t, order)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
	{
		testName: "invalid deadline",
		inputData: struct {
			userID   uuid.UUID
			address  string
			deadline time.Time
			tasks    []models.Task
		}{
			uuid.New(),
			"address",
			time.Now().AddDate(0, 0, -1),
			[]models.Task{{ID: uuid.New()}, {ID: uuid.New()}},
		},
		prepare: func(fields *orderServiceFields) {},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Nil(t, order)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
	{
		testName: "one of the tasks not found",
		inputData: struct {
			userID   uuid.UUID
			address  string
			deadline time.Time
			tasks    []models.Task
		}{
			uuid.New(),
			"address",
			time.Now().AddDate(0, 0, 1),
			[]models.Task{{ID: uuid.New()}, {ID: uuid.New()}},
		},
		prepare: func(fields *orderServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(nil, service_errors.InvalidReference)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Nil(t, order)
			assert.Equal(t, service_errors.InvalidReference, err)
		},
	},
	{
		testName: "user not found",
		inputData: struct {
			userID   uuid.UUID
			address  string
			deadline time.Time
			tasks    []models.Task
		}{
			uuid.New(),
			"address",
			time.Now().AddDate(0, 0, 1),
			[]models.Task{{ID: uuid.New()}, {ID: uuid.New()}},
		},
		prepare: func(fields *orderServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{}, nil).Times(2)
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(nil, service_errors.InvalidReference)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Nil(t, order)
			assert.Equal(t, service_errors.InvalidReference, err)
		},
	},
	{
		testName: "order creation error",
		inputData: struct {
			userID   uuid.UUID
			address  string
			deadline time.Time
			tasks    []models.Task
		}{
			uuid.New(),
			"address",
			time.Now().AddDate(0, 0, 1),
			[]models.Task{{ID: uuid.New()}, {ID: uuid.New()}},
		},
		prepare: func(fields *orderServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{}, nil).Times(2)
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(&models.User{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, repository_errors.InsertError)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Nil(t, order)
			assert.Equal(t, repository_errors.InsertError, err)
		},
	},
}

func TestOrderService_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initOrderServiceFields(ctrl)
	orderService := initOrderService(fields)

	for _, tt := range testOrderServiceCreate {
		t.Run(tt.testName, func(t *testing.T) {
			tt.prepare(fields)
			orderedTasks := make([]models.OrderedTask, len(tt.inputData.tasks))
			for i, task := range tt.inputData.tasks {
				orderedTasks[i] = models.OrderedTask{
					Task: &task,
					// Set the quantity or any other fields as needed
					Quantity: 1,
				}
			}
			order, err := orderService.CreateOrder(tt.inputData.userID, tt.inputData.address, tt.inputData.deadline, orderedTasks)
			tt.checkOutput(t, order, err)
		})
	}
}

var testOrderServiceDelete = []struct {
	testName  string
	inputData struct {
		orderID uuid.UUID
	}
	prepare     func(fields *orderServiceFields)
	checkOutput func(t *testing.T, err error)
}{
	{
		testName: "delete order success test",
		inputData: struct {
			orderID uuid.UUID
		}{
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.UUID{}}, nil)
			fields.orderRepoMock.EXPECT().GetTasksInOrder(gomock.Any()).Return([]models.Task{{ID: uuid.UUID{}}}, nil)
			fields.orderRepoMock.EXPECT().RemoveTaskFromOrder(gomock.Any(), gomock.Any()).Return(nil)
			fields.orderRepoMock.EXPECT().Delete(gomock.Any()).Return(nil)
		},
		checkOutput: func(t *testing.T, err error) {
			assert.NoError(t, err)
		},
	},
	{
		testName: "order not found",
		inputData: struct {
			orderID uuid.UUID
		}{
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DoesNotExist, err)
		},
	},
	{
		testName: "delete order error",
		inputData: struct {
			orderID uuid.UUID
		}{
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.UUID{}}, nil)
			fields.orderRepoMock.EXPECT().GetTasksInOrder(gomock.Any()).Return([]models.Task{{ID: uuid.UUID{}}}, nil)
			fields.orderRepoMock.EXPECT().RemoveTaskFromOrder(gomock.Any(), gomock.Any()).Return(nil)
			fields.orderRepoMock.EXPECT().Delete(gomock.Any()).Return(repository_errors.DeleteError)
		},
		checkOutput: func(t *testing.T, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DeleteError, err)
		},
	},
}

func TestOrderService_DeleteOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initOrderServiceFields(ctrl)
	orderService := initOrderService(fields)

	for _, tt := range testOrderServiceDelete {
		t.Run(tt.testName, func(t *testing.T) {
			tt.prepare(fields)
			err := orderService.DeleteOrder(tt.inputData.orderID)
			tt.checkOutput(t, err)
		})
	}
}

var testOrderServiceGetTasksInOrder = []struct {
	testName  string
	inputData struct {
		orderID uuid.UUID
	}
	prepare     func(fields *orderServiceFields)
	checkOutput func(t *testing.T, tasks []models.Task, err error)
}{
	{
		testName: "get tasks in order success test",
		inputData: struct {
			orderID uuid.UUID
		}{
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().GetTasksInOrder(gomock.Any()).Return([]models.Task{{ID: uuid.New()}, {ID: uuid.New()}}, nil)
		},
		checkOutput: func(t *testing.T, tasks []models.Task, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, tasks)
			assert.Len(t, tasks, 2)
		},
	},
	{
		testName: "order not found",
		inputData: struct {
			orderID uuid.UUID
		}{
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, tasks []models.Task, err error) {
			assert.Error(t, err)
			assert.Nil(t, tasks)
			assert.Equal(t, repository_errors.DoesNotExist, err)
		},
	},
	{
		testName: "get tasks in order error",
		inputData: struct {
			orderID uuid.UUID
		}{
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().GetTasksInOrder(gomock.Any()).Return(nil, repository_errors.SelectError)
		},
		checkOutput: func(t *testing.T, tasks []models.Task, err error) {
			assert.Error(t, err)
			assert.Nil(t, tasks)
			assert.Equal(t, repository_errors.SelectError, err)
		},
	},
}

func TestOrderService_GetTasksInOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initOrderServiceFields(ctrl)
	orderService := initOrderService(fields)

	for _, tt := range testOrderServiceGetTasksInOrder {
		t.Run(tt.testName, func(t *testing.T) {
			tt.prepare(fields)
			tasks, err := orderService.GetTasksInOrder(tt.inputData.orderID)
			tt.checkOutput(t, tasks, err)
		})
	}
}

var testOrderServiceGetCurrentOrderByUserID = []struct {
	testName  string
	inputData struct {
		userID uuid.UUID
	}
	prepare     func(fields *orderServiceFields)
	checkOutput func(t *testing.T, order *models.Order, err error)
}{
	{
		testName: "get current order by user id success test",
		inputData: struct {
			userID uuid.UUID
		}{
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(&models.User{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().GetCurrentOrderByUserID(gomock.Any()).Return(&models.Order{ID: uuid.New()}, nil)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, order)
		},
	},
	{
		testName: "user not found",
		inputData: struct {
			userID uuid.UUID
		}{
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Nil(t, order)
			assert.Equal(t, fmt.Errorf("SERVICE: GetUserByID method failed"), err)
		},
	},
	{
		testName: "get current order by user id error",
		inputData: struct {
			userID uuid.UUID
		}{
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(&models.User{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().GetCurrentOrderByUserID(gomock.Any()).Return(nil, repository_errors.SelectError)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Nil(t, order)
			assert.Equal(t, repository_errors.SelectError, err)
		},
	},
}

func TestOrderService_GetCurrentOrderByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initOrderServiceFields(ctrl)
	orderService := initOrderService(fields)

	for _, tt := range testOrderServiceGetCurrentOrderByUserID {
		t.Run(tt.testName, func(t *testing.T) {
			tt.prepare(fields)
			order, err := orderService.GetCurrentOrderByUserID(tt.inputData.userID)
			tt.checkOutput(t, order, err)
		})
	}
}

var testOrderServiceGetAllOrdersByUserID = []struct {
	testName  string
	inputData struct {
		userID uuid.UUID
	}
	prepare     func(fields *orderServiceFields)
	checkOutput func(t *testing.T, orders []models.Order, err error)
}{
	{
		testName: "get all orders by user id success test",
		inputData: struct {
			userID uuid.UUID
		}{
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(&models.User{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().GetAllOrdersByUserID(gomock.Any()).Return([]models.Order{{ID: uuid.New()}, {ID: uuid.New()}}, nil)
		},
		checkOutput: func(t *testing.T, orders []models.Order, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, orders)
			assert.Len(t, orders, 2)
		},
	},
	{
		testName: "user not found",
		inputData: struct {
			userID uuid.UUID
		}{
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, orders []models.Order, err error) {
			assert.Error(t, err)
			assert.Nil(t, orders)
			assert.Equal(t, fmt.Errorf("SERVICE: GetUserByID method failed"), err)
		},
	},
	{
		testName: "get all orders by user id error",
		inputData: struct {
			userID uuid.UUID
		}{
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(&models.User{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().GetAllOrdersByUserID(gomock.Any()).Return(nil, repository_errors.SelectError)
		},
		checkOutput: func(t *testing.T, orders []models.Order, err error) {
			assert.Error(t, err)
			assert.Nil(t, orders)
			assert.Equal(t, repository_errors.SelectError, err)
		},
	},
}

func TestOrderService_GetAllOrdersByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initOrderServiceFields(ctrl)
	orderService := initOrderService(fields)

	for _, tt := range testOrderServiceGetAllOrdersByUserID {
		t.Run(tt.testName, func(t *testing.T) {
			tt.prepare(fields)
			orders, err := orderService.GetAllOrdersByUserID(tt.inputData.userID)
			tt.checkOutput(t, orders, err)
		})
	}
}

var testOrderServiceChangeOrderStatus = []struct {
	testName  string
	inputData struct {
		orderID  uuid.UUID
		status   int
		rate     int
		workerID uuid.UUID
	}
	prepare     func(fields *orderServiceFields)
	checkOutput func(t *testing.T, order *models.Order, err error)
}{
	{
		testName: "Success",
		inputData: struct {
			orderID  uuid.UUID
			status   int
			rate     int
			workerID uuid.UUID
		}{
			uuid.New(),
			models.NewOrderStatus,
			0,
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{
				ID:       uuid.New(),
				Status:   models.NewOrderStatus,
				Rate:     0,
				WorkerID: uuid.New(),
			}, nil)
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().Update(gomock.Any()).Return(&models.Order{
				ID:       uuid.New(),
				Status:   models.InProgressOrderStatus,
				Rate:     0,
				WorkerID: uuid.New(),
			}, nil)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.NoError(t, err)
		},
	},
	{
		testName: "order not found",
		inputData: struct {
			orderID  uuid.UUID
			status   int
			rate     int
			workerID uuid.UUID
		}{
			uuid.New(),
			models.InProgressOrderStatus,
			5,
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DoesNotExist, err)
		},
	},
	{
		testName: "already completed order",
		inputData: struct {
			orderID  uuid.UUID
			status   int
			rate     int
			workerID uuid.UUID
		}{
			uuid.New(),
			models.CompletedOrderStatus,
			0,
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{
				ID:       uuid.New(),
				Status:   models.CompletedOrderStatus,
				Rate:     0,
				WorkerID: uuid.New(),
			}, nil)
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{ID: uuid.New()}, nil)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Equal(t, fmt.Errorf("SERVICE: Order is already completed or cancelled"), err)
		},
	},
	{
		testName: "invalid status",
		inputData: struct {
			orderID  uuid.UUID
			status   int
			rate     int
			workerID uuid.UUID
		}{
			uuid.New(),
			0,
			0,
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{
				ID:       uuid.New(),
				Status:   0,
				Rate:     0,
				WorkerID: uuid.New(),
			}, nil)
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{ID: uuid.New()}, nil)

		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid status"), err)
		},
	},
	{
		testName: "change order status error",
		inputData: struct {
			orderID  uuid.UUID
			status   int
			rate     int
			workerID uuid.UUID
		}{
			uuid.New(),
			models.InProgressOrderStatus,
			0,
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{
				ID:       uuid.New(),
				Status:   models.NewOrderStatus,
				Rate:     0,
				WorkerID: uuid.New(),
			}, nil)
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{ID: uuid.New()}, nil)

			fields.orderRepoMock.EXPECT().Update(gomock.Any()).Return(nil, repository_errors.UpdateError)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.UpdateError, err)
		},
	},
}

func TestOrderService_ChangeOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initOrderServiceFields(ctrl)
	orderService := initOrderService(fields)

	for _, tt := range testOrderServiceChangeOrderStatus {
		t.Run(tt.testName, func(t *testing.T) {
			tt.prepare(fields)
			order, err := orderService.Update(tt.inputData.orderID, tt.inputData.status, tt.inputData.rate, tt.inputData.workerID)
			tt.checkOutput(t, order, err)
		})
	}
}

var testOrderServiceRateOrder = []struct {
	testName  string
	inputData struct {
		orderID  uuid.UUID
		status   int
		rate     int
		workerID uuid.UUID
	}
	prepare     func(fields *orderServiceFields)
	checkOutput func(t *testing.T, order *models.Order, err error)
}{
	{
		testName: "Success",
		inputData: struct {
			orderID  uuid.UUID
			status   int
			rate     int
			workerID uuid.UUID
		}{
			uuid.New(),
			models.InProgressOrderStatus,
			0,
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{
				ID:       uuid.New(),
				Status:   models.InProgressOrderStatus,
				Rate:     0,
				WorkerID: uuid.New(),
			}, nil)
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().Update(gomock.Any()).Return(&models.Order{
				ID:       uuid.New(),
				Status:   models.CompletedOrderStatus,
				Rate:     5,
				WorkerID: uuid.New(),
			}, nil)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.NoError(t, err)
		},
	},
	{
		testName: "order not found",
		inputData: struct {
			orderID  uuid.UUID
			status   int
			rate     int
			workerID uuid.UUID
		}{
			uuid.New(),
			models.InProgressOrderStatus,
			5,
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DoesNotExist, err)
		},
	},
	{
		testName: "invalid rating",
		inputData: struct {
			orderID  uuid.UUID
			status   int
			rate     int
			workerID uuid.UUID
		}{
			uuid.New(),
			models.CompletedOrderStatus,
			7,
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{
				ID:       uuid.New(),
				Status:   models.CompletedOrderStatus,
				Rate:     7,
				WorkerID: uuid.New(),
			}, nil)
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{ID: uuid.New()}, nil)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Equal(t, fmt.Errorf("SERVICE: Rating is out of range"), err)
		},
	},
	{
		testName: "order is not completed",
		inputData: struct {
			orderID  uuid.UUID
			status   int
			rate     int
			workerID uuid.UUID
		}{
			uuid.New(),
			models.InProgressOrderStatus,
			5,
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{
				ID:       uuid.New(),
				Status:   models.InProgressOrderStatus,
				Rate:     5,
				WorkerID: uuid.New(),
			}, nil)
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{ID: uuid.New()}, nil)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Equal(t, fmt.Errorf("SERVICE: Order is not completed"), err)
		},
	},
	{
		testName: "rate order error",
		inputData: struct {
			orderID  uuid.UUID
			status   int
			rate     int
			workerID uuid.UUID
		}{
			uuid.New(),
			models.CompletedOrderStatus,
			5,
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{
				ID:       uuid.New(),
				Status:   models.CompletedOrderStatus,
				Rate:     5,
				WorkerID: uuid.New(),
			}, nil)
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().Update(gomock.Any()).Return(nil, repository_errors.UpdateError)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.UpdateError, err)
		},
	},
}

func TestOrderService_RateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initOrderServiceFields(ctrl)
	orderService := initOrderService(fields)

	for _, tt := range testOrderServiceRateOrder {
		t.Run(tt.testName, func(t *testing.T) {
			tt.prepare(fields)
			order, err := orderService.Update(tt.inputData.orderID, tt.inputData.status, tt.inputData.rate, tt.inputData.workerID)
			tt.checkOutput(t, order, err)
		})
	}
}

var testOrderServiceAttachWorkerToOrder = []struct {
	testName  string
	inputData struct {
		orderID  uuid.UUID
		status   int
		rate     int
		workerID uuid.UUID
	}
	prepare     func(fields *orderServiceFields)
	checkOutput func(t *testing.T, order *models.Order, err error)
}{
	{
		testName: "Success",
		inputData: struct {
			orderID  uuid.UUID
			status   int
			rate     int
			workerID uuid.UUID
		}{
			uuid.New(),
			models.InProgressOrderStatus,
			0,
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{
				ID:       uuid.New(),
				Status:   models.InProgressOrderStatus,
				Rate:     0,
				WorkerID: uuid.New(),
			}, nil)
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().Update(gomock.Any()).Return(&models.Order{}, nil)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.NoError(t, err)
		},
	},
	{
		testName: "order not found",
		inputData: struct {
			orderID  uuid.UUID
			status   int
			rate     int
			workerID uuid.UUID
		}{
			uuid.New(),
			models.InProgressOrderStatus,
			5,
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DoesNotExist, err)
		},
	},
	{
		testName: "worker not found",
		inputData: struct {
			orderID  uuid.UUID
			status   int
			rate     int
			workerID uuid.UUID
		}{
			uuid.New(),
			models.NewOrderStatus,
			0,
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{
				ID:       uuid.New(),
				Status:   models.NewOrderStatus,
				Rate:     0,
				WorkerID: uuid.New(),
			}, nil)
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DoesNotExist, err)
		},
	},
	{
		testName: "attach worker to order error",
		inputData: struct {
			orderID  uuid.UUID
			status   int
			rate     int
			workerID uuid.UUID
		}{
			uuid.New(),
			models.NewOrderStatus,
			0,
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{
				ID:       uuid.New(),
				Status:   models.NewOrderStatus,
				Rate:     0,
				WorkerID: uuid.New(),
			}, nil)
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().Update(gomock.Any()).Return(nil, repository_errors.UpdateError)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.UpdateError, err)
		},
	},
}

func TestOrderService_AttachWorkerToOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initOrderServiceFields(ctrl)
	orderService := initOrderService(fields)

	for _, tt := range testOrderServiceAttachWorkerToOrder {
		t.Run(tt.testName, func(t *testing.T) {
			tt.prepare(fields)
			order, err := orderService.Update(tt.inputData.orderID, tt.inputData.status, tt.inputData.rate, tt.inputData.workerID)
			tt.checkOutput(t, order, err)
		})
	}
}

var testOrderServiceDetachWorkerFromOrder = []struct {
	testName  string
	inputData struct {
		orderID  uuid.UUID
		status   int
		rate     int
		workerID uuid.UUID
	}
	prepare     func(fields *orderServiceFields)
	checkOutput func(t *testing.T, order *models.Order, err error)
}{
	{
		testName: "Success",
		inputData: struct {
			orderID  uuid.UUID
			status   int
			rate     int
			workerID uuid.UUID
		}{
			uuid.New(),
			models.NewOrderStatus,
			0,
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{
				ID:       uuid.New(),
				Status:   models.NewOrderStatus,
				Rate:     0,
				WorkerID: uuid.New(),
			}, nil)
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().Update(gomock.Any()).Return(&models.Order{}, nil)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.NoError(t, err)
		},
	},
	{
		testName: "order not found",
		inputData: struct {
			orderID  uuid.UUID
			status   int
			rate     int
			workerID uuid.UUID
		}{
			uuid.New(),
			models.InProgressOrderStatus,
			5,
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DoesNotExist, err)
		},
	},
	{
		testName: "detach worker from order error",
		inputData: struct {
			orderID  uuid.UUID
			status   int
			rate     int
			workerID uuid.UUID
		}{
			uuid.New(),
			models.NewOrderStatus,
			0,
			uuid.New(),
		},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{
				ID:       uuid.New(),
				Status:   models.NewOrderStatus,
				Rate:     0,
				WorkerID: uuid.New(),
			}, nil)
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().Update(gomock.Any()).Return(nil, repository_errors.UpdateError)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.UpdateError, err)
		},
	},
}

func TestOrderService_DetachWorkerFromOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initOrderServiceFields(ctrl)
	orderService := initOrderService(fields)

	for _, tt := range testOrderServiceDetachWorkerFromOrder {
		t.Run(tt.testName, func(t *testing.T) {
			tt.prepare(fields)
			order, err := orderService.Update(tt.inputData.orderID, tt.inputData.status, tt.inputData.rate, tt.inputData.workerID)
			tt.checkOutput(t, order, err)
		})
	}
}

var testOrderServiceGetOrderByID = []struct {
	testName    string
	inputData   struct{ orderID uuid.UUID }
	prepare     func(fields *orderServiceFields)
	checkOutput func(t *testing.T, order *models.Order, err error)
}{
	{
		testName:  "get order by id success test",
		inputData: struct{ orderID uuid.UUID }{uuid.New()},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.New()}, nil)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, order)
		},
	},
	{
		testName:  "order not found",
		inputData: struct{ orderID uuid.UUID }{uuid.New()},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Nil(t, order)
			assert.Equal(t, repository_errors.DoesNotExist, err)
		},
	},
	{
		testName:  "get order by id error",
		inputData: struct{ orderID uuid.UUID }{uuid.New()},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(nil, repository_errors.SelectError)
		},
		checkOutput: func(t *testing.T, order *models.Order, err error) {
			assert.Error(t, err)
			assert.Nil(t, order)
			assert.Equal(t, repository_errors.SelectError, err)
		},
	},
}

func TestOrderService_GetOrderByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initOrderServiceFields(ctrl)
	orderService := initOrderService(fields)

	for _, tt := range testOrderServiceGetOrderByID {
		t.Run(tt.testName, func(t *testing.T) {
			tt.prepare(fields)
			order, err := orderService.GetOrderByID(tt.inputData.orderID)
			tt.checkOutput(t, order, err)
		})
	}
}

var testOrderServiceIncrementTaskQuantity = []struct {
	testName  string
	inputData struct {
		orderID uuid.UUID
		taskID  uuid.UUID
	}
	prepare     func(fields *orderServiceFields)
	checkOutput func(t *testing.T, quantity int, err error)
}{
	{
		testName: "increment task quantity success test",
		inputData: struct {
			orderID uuid.UUID
			taskID  uuid.UUID
		}{uuid.New(), uuid.New()},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.New(), Status: models.NewOrderStatus}, nil)
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().GetTaskQuantity(gomock.Any(), gomock.Any()).Return(1, nil)
			fields.orderRepoMock.EXPECT().UpdateTaskQuantity(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		},
		checkOutput: func(t *testing.T, quantity int, err error) {
			assert.NoError(t, err)
			assert.Equal(t, 2, quantity)
		},
	},
	{
		testName: "order not found",
		inputData: struct {
			orderID uuid.UUID
			taskID  uuid.UUID
		}{uuid.New(), uuid.New()},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, quantity int, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DoesNotExist, err)
			assert.Equal(t, 0, quantity)
		},
	},
	{
		testName: "task not found",
		inputData: struct {
			orderID uuid.UUID
			taskID  uuid.UUID
		}{uuid.New(), uuid.New()},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.New(), Status: models.NewOrderStatus}, nil)
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, quantity int, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DoesNotExist, err)
			assert.Equal(t, 0, quantity)
		},
	},
	{
		testName: "increment task quantity error",
		inputData: struct {
			orderID uuid.UUID
			taskID  uuid.UUID
		}{uuid.New(), uuid.New()},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.New(), Status: models.NewOrderStatus}, nil)
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().GetTaskQuantity(gomock.Any(), gomock.Any()).Return(0, repository_errors.SelectError)
		},
		checkOutput: func(t *testing.T, quantity int, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.SelectError, err)
			assert.Equal(t, 0, quantity)
		},
	},
	{
		testName: "task quantity is not incremented",
		inputData: struct {
			orderID uuid.UUID
			taskID  uuid.UUID
		}{uuid.New(), uuid.New()},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.New(), Status: models.NewOrderStatus}, nil)
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().GetTaskQuantity(gomock.Any(), gomock.Any()).Return(1, nil)
			fields.orderRepoMock.EXPECT().UpdateTaskQuantity(gomock.Any(), gomock.Any(), gomock.Any()).Return(repository_errors.UpdateError)
		},
		checkOutput: func(t *testing.T, quantity int, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.UpdateError, err)
			assert.Equal(t, 0, quantity)
		},
	},
}

func TestOrderService_IncrementTaskQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initOrderServiceFields(ctrl)
	orderService := initOrderService(fields)

	for _, tt := range testOrderServiceIncrementTaskQuantity {
		t.Run(tt.testName, func(t *testing.T) {
			tt.prepare(fields)
			quantity, err := orderService.IncrementTaskQuantity(tt.inputData.orderID, tt.inputData.taskID)
			tt.checkOutput(t, quantity, err)
		})
	}
}

var testOrderServiceDecrementTaskQuantity = []struct {
	testName  string
	inputData struct {
		orderID uuid.UUID
		taskID  uuid.UUID
	}
	prepare     func(fields *orderServiceFields)
	checkOutput func(t *testing.T, quantity int, err error)
}{
	{
		testName: "decrement task quantity success test",
		inputData: struct {
			orderID uuid.UUID
			taskID  uuid.UUID
		}{uuid.New(), uuid.New()},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.New(), Status: models.NewOrderStatus}, nil)
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().GetTaskQuantity(gomock.Any(), gomock.Any()).Return(2, nil)
			fields.orderRepoMock.EXPECT().UpdateTaskQuantity(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		},
		checkOutput: func(t *testing.T, quantity int, err error) {
			assert.NoError(t, err)
			assert.Equal(t, 1, quantity)
		},
	},
	{
		testName: "order not found",
		inputData: struct {
			orderID uuid.UUID
			taskID  uuid.UUID
		}{uuid.New(), uuid.New()},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, quantity int, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DoesNotExist, err)
			assert.Equal(t, 0, quantity)
		},
	},
	{
		testName: "task not found",
		inputData: struct {
			orderID uuid.UUID
			taskID  uuid.UUID
		}{uuid.New(), uuid.New()},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.New(), Status: models.NewOrderStatus}, nil)
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, quantity int, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DoesNotExist, err)
			assert.Equal(t, 0, quantity)
		},
	},
	{
		testName: "negative result",
		inputData: struct {
			orderID uuid.UUID
			taskID  uuid.UUID
		}{uuid.New(), uuid.New()},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.New(), Status: models.NewOrderStatus}, nil)
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().GetTaskQuantity(gomock.Any(), gomock.Any()).Return(0, nil)
		},
		checkOutput: func(t *testing.T, quantity int, err error) {
			assert.Error(t, err)
			assert.Equal(t, fmt.Errorf("SERVICE: Quantity is already 0"), err)
			assert.Equal(t, 0, quantity)
		},
	},
	{
		testName: "decrement task quantity error",
		inputData: struct {
			orderID uuid.UUID
			taskID  uuid.UUID
		}{uuid.New(), uuid.New()},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.New(), Status: models.NewOrderStatus}, nil)
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().GetTaskQuantity(gomock.Any(), gomock.Any()).Return(1, nil)
			fields.orderRepoMock.EXPECT().UpdateTaskQuantity(gomock.Any(), gomock.Any(), gomock.Any()).Return(repository_errors.UpdateError)
		},
		checkOutput: func(t *testing.T, quantity int, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.UpdateError, err)
			assert.Equal(t, 0, quantity)
		},
	},
}

func TestOrderService_DecrementTaskQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initOrderServiceFields(ctrl)
	orderService := initOrderService(fields)

	for _, tt := range testOrderServiceDecrementTaskQuantity {
		t.Run(tt.testName, func(t *testing.T) {
			tt.prepare(fields)
			quantity, err := orderService.DecrementTaskQuantity(tt.inputData.orderID, tt.inputData.taskID)
			tt.checkOutput(t, quantity, err)
		})
	}
}

var testOrderSetTaskQuantity = []struct {
	testName  string
	inputData struct {
		orderID  uuid.UUID
		taskID   uuid.UUID
		quantity int
	}
	prepare     func(fields *orderServiceFields)
	checkOutput func(t *testing.T, err error)
}{
	{
		testName: "set task quantity success test",
		inputData: struct {
			orderID  uuid.UUID
			taskID   uuid.UUID
			quantity int
		}{uuid.New(), uuid.New(), 5},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.New(), Status: models.NewOrderStatus}, nil)
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().UpdateTaskQuantity(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		},
		checkOutput: func(t *testing.T, err error) {
			assert.NoError(t, err)
		},
	},
	{
		testName: "order not found",
		inputData: struct {
			orderID  uuid.UUID
			taskID   uuid.UUID
			quantity int
		}{uuid.New(), uuid.New(), 5},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DoesNotExist, err)
		},
	},
	{
		testName: "task not found",
		inputData: struct {
			orderID  uuid.UUID
			taskID   uuid.UUID
			quantity int
		}{uuid.New(), uuid.New(), 5},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.New(), Status: models.NewOrderStatus}, nil)
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DoesNotExist, err)
		},
	},
	{
		testName: "negative quantity",
		inputData: struct {
			orderID  uuid.UUID
			taskID   uuid.UUID
			quantity int
		}{uuid.New(), uuid.New(), -1},
		prepare: func(fields *orderServiceFields) {},
		checkOutput: func(t *testing.T, err error) {
			assert.Error(t, err)
			assert.Equal(t, fmt.Errorf("SERVICE: Quantity is negative"), err)
		},
	},
	{
		testName: "set task quantity error",
		inputData: struct {
			orderID  uuid.UUID
			taskID   uuid.UUID
			quantity int
		}{uuid.New(), uuid.New(), 5},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.New(), Status: models.NewOrderStatus}, nil)
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().UpdateTaskQuantity(gomock.Any(), gomock.Any(), gomock.Any()).Return(repository_errors.UpdateError)
		},
		checkOutput: func(t *testing.T, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.UpdateError, err)
		},
	},
}

func TestOrderService_SetTaskQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initOrderServiceFields(ctrl)
	orderService := initOrderService(fields)

	for _, tt := range testOrderSetTaskQuantity {
		t.Run(tt.testName, func(t *testing.T) {
			tt.prepare(fields)
			err := orderService.SetTaskQuantity(tt.inputData.orderID, tt.inputData.taskID, tt.inputData.quantity)
			tt.checkOutput(t, err)
		})
	}
}

var testOrderServiceGetTaskQuantity = []struct {
	testName  string
	inputData struct {
		orderID uuid.UUID
		taskID  uuid.UUID
	}
	prepare     func(fields *orderServiceFields)
	checkOutput func(t *testing.T, quantity int, err error)
}{
	{
		testName: "get task quantity success test",
		inputData: struct {
			orderID uuid.UUID
			taskID  uuid.UUID
		}{uuid.New(), uuid.New()},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.New()}, nil)
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().GetTaskQuantity(gomock.Any(), gomock.Any()).Return(5, nil)
		},
		checkOutput: func(t *testing.T, quantity int, err error) {
			assert.NoError(t, err)
			assert.Equal(t, 5, quantity)
		},
	},
	{
		testName: "order not found",
		inputData: struct {
			orderID uuid.UUID
			taskID  uuid.UUID
		}{uuid.New(), uuid.New()},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, quantity int, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DoesNotExist, err)
			assert.Equal(t, 0, quantity)
		},
	},
	{
		testName: "task not found",
		inputData: struct {
			orderID uuid.UUID
			taskID  uuid.UUID
		}{uuid.New(), uuid.New()},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.New()}, nil)
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, quantity int, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DoesNotExist, err)
			assert.Equal(t, 0, quantity)
		},
	},
	{
		testName: "get task quantity error",
		inputData: struct {
			orderID uuid.UUID
			taskID  uuid.UUID
		}{uuid.New(), uuid.New()},
		prepare: func(fields *orderServiceFields) {
			fields.orderRepoMock.EXPECT().GetOrderByID(gomock.Any()).Return(&models.Order{ID: uuid.New()}, nil)
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{ID: uuid.New()}, nil)
			fields.orderRepoMock.EXPECT().GetTaskQuantity(gomock.Any(), gomock.Any()).Return(0, repository_errors.SelectError)
		},
		checkOutput: func(t *testing.T, quantity int, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.SelectError, err)
			assert.Equal(t, 0, quantity)
		},
	},
}

func TestOrderService_GetTaskQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initOrderServiceFields(ctrl)
	orderService := initOrderService(fields)

	for _, tt := range testOrderServiceGetTaskQuantity {
		t.Run(tt.testName, func(t *testing.T) {
			tt.prepare(fields)
			quantity, err := orderService.GetTaskQuantity(tt.inputData.orderID, tt.inputData.taskID)
			tt.checkOutput(t, quantity, err)
		})
	}
}
