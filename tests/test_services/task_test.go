package test_services

import (
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
)

type taskServiceFields struct {
	taskRepoMock *mock_repository_interfaces.MockITaskRepository
	logger       *log.Logger
}

func initTaskServiceFields(ctrl *gomock.Controller) *taskServiceFields {
	taskRepoMock := mock_repository_interfaces.NewMockITaskRepository(ctrl)
	f, err := os.OpenFile("tests.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(f)

	return &taskServiceFields{
		taskRepoMock: taskRepoMock,
		logger:       logger,
	}
}

func initTaskService(fields *taskServiceFields) service_interfaces.ITaskService {
	return services.NewTaskService(fields.taskRepoMock, fields.logger)
}

var testTaskCreateSuccess = []struct {
	testName  string
	inputData struct {
		name     string
		price    float64
		category int
	}
	prepare     func(fields *taskServiceFields)
	checkOutput func(t *testing.T, task *models.Task, err error)
}{
	{
		testName: "basic create",
		inputData: struct {
			name     string
			price    float64
			category int
		}{name: "Test Task", price: 100.0, category: 1},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().Create(gomock.Any()).Return(&models.Task{
				ID:             uuid.New(),
				Name:           "Test Task",
				PricePerSingle: 100.0,
				Category:       1,
			}, nil)
		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.NoError(t, err)
			assert.Equal(t, "Test Task", task.Name)
			assert.Equal(t, 100.0, task.PricePerSingle)
			assert.Equal(t, 1, task.Category)
		},
	},
	{
		testName: "cyrillic name",
		inputData: struct {
			name     string
			price    float64
			category int
		}{
			name:     "Тестовое задание",
			price:    100.0,
			category: 1,
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().Create(gomock.Any()).Return(&models.Task{
				ID:             uuid.New(),
				Name:           "Тестовое задание",
				PricePerSingle: 100.0,
				Category:       1,
			}, nil)
		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.NoError(t, err)
			assert.Equal(t, "Тестовое задание", task.Name)
			assert.Equal(t, 100.0, task.PricePerSingle)
			assert.Equal(t, 1, task.Category)
		},
	},
}

var testTaskCreateFail = []struct {
	testName  string
	inputData struct {
		name     string
		price    float64
		category int
	}
	prepare     func(fields *taskServiceFields)
	checkOutput func(t *testing.T, task *models.Task, err error)
}{
	{
		testName: "empty name",
		inputData: struct {
			name     string
			price    float64
			category int
		}{name: "", price: 100.0, category: 1},
		prepare: func(fields *taskServiceFields) {},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.Error(t, err)
			assert.Nil(t, task)
		},
	},
	{
		testName: "negative price",
		inputData: struct {
			name     string
			price    float64
			category int
		}{name: "Test Task", price: -1.0, category: 1},
		prepare: func(fields *taskServiceFields) {},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.Error(t, err)
			assert.Nil(t, task)
		},
	},
	{
		testName: "invalid category",
		inputData: struct {
			name     string
			price    float64
			category int
		}{
			name:     "Test Task",
			price:    100.0,
			category: 0,
		},
		prepare: func(fields *taskServiceFields) {},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.Error(t, err)
			assert.Nil(t, task)
		},
	},
}

func TestTaskServiceCreate(t *testing.T) {
	for _, tt := range testTaskCreateSuccess {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := initTaskServiceFields(ctrl)
			taskService := initTaskService(fields)

			tt.prepare(fields)

			task, err := taskService.Create(tt.inputData.name, tt.inputData.price, tt.inputData.category)
			tt.checkOutput(t, task, err)
		})
	}

	for _, tt := range testTaskCreateFail {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := initTaskServiceFields(ctrl)
			taskService := initTaskService(fields)

			tt.prepare(fields)

			task, err := taskService.Create(tt.inputData.name, tt.inputData.price, tt.inputData.category)
			tt.checkOutput(t, task, err)
		})
	}
}

var testTaskChangeCategorySuccess = []struct {
	testName  string
	inputData struct {
		taskID   uuid.UUID
		category int
		name     string
		price    float64
	}
	prepare     func(fields *taskServiceFields)
	checkOutput func(t *testing.T, task *models.Task, err error)
}{
	{
		testName: "Success",
		inputData: struct {
			taskID   uuid.UUID
			category int
			name     string
			price    float64
		}{
			taskID:   uuid.New(),
			category: 2,
			name:     "Test Task",
			price:    100.0,
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{
				ID:             uuid.New(),
				Name:           "Test Task",
				PricePerSingle: 100.0,
				Category:       1, //category from 1
			}, nil)
			fields.taskRepoMock.EXPECT().Update(gomock.Any()).Return(&models.Task{
				ID:             uuid.New(),
				Name:           "Test Task",
				PricePerSingle: 100.0,
				Category:       2, //to category 2
			}, nil)
		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.NoError(t, err)
		},
	},
	{
		testName: "same category",
		inputData: struct {
			taskID   uuid.UUID
			category int
			name     string
			price    float64
		}{
			taskID:   uuid.New(),
			category: 1,
			name:     "Test Task",
			price:    100.0,
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{
				ID:             uuid.New(),
				Name:           "Test Task",
				PricePerSingle: 100.0,
				Category:       1, //no diff
			}, nil)
			fields.taskRepoMock.EXPECT().Update(gomock.Any()).Return(&models.Task{
				ID:             uuid.New(),
				Name:           "Test Task",
				PricePerSingle: 100.0,
				Category:       1, //no diff
			}, nil)
		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.NoError(t, err)
		},
	},
}

var testTaskChangeCategoryFail = []struct {
	testName  string
	inputData struct {
		taskID   uuid.UUID
		category int
		name     string
		price    float64
	}
	prepare     func(fields *taskServiceFields)
	checkOutput func(t *testing.T, task *models.Task, err error)
}{
	{
		testName: "invalid category",
		inputData: struct {
			taskID   uuid.UUID
			category int
			name     string
			price    float64
		}{
			taskID:   uuid.New(),
			category: 6,
			name:     "Test Task",
			price:    100.0,
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(nil, service_errors.InvalidCategory)
		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.Error(t, err)
		},
	},
	{
		testName: "task not found",
		inputData: struct {
			taskID   uuid.UUID
			category int
			name     string
			price    float64
		}{
			taskID:   uuid.New(),
			category: 1,
			name:     "Test Task",
			price:    100.0,
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.Error(t, err)
		},
	},
}

func TestTaskServiceChangeCategory(t *testing.T) {
	for _, tt := range testTaskChangeCategorySuccess {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := initTaskServiceFields(ctrl)
			taskService := initTaskService(fields)

			tt.prepare(fields)

			task, err := taskService.Update(tt.inputData.taskID, tt.inputData.category, tt.inputData.name, tt.inputData.price)
			tt.checkOutput(t, task, err)
		})
	}

	for _, tt := range testTaskChangeCategoryFail {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := initTaskServiceFields(ctrl)
			taskService := initTaskService(fields)

			tt.prepare(fields)

			task, err := taskService.Update(tt.inputData.taskID, tt.inputData.category, tt.inputData.name, tt.inputData.price)
			tt.checkOutput(t, task, err)
		})
	}
}

var testTaskChangeNameSuccess = []struct {
	testName  string
	inputData struct {
		taskID   uuid.UUID
		category int
		name     string
		price    float64
	}
	prepare     func(fields *taskServiceFields)
	checkOutput func(t *testing.T, task *models.Task, err error)
}{
	{
		testName: "basic change",
		inputData: struct {
			taskID   uuid.UUID
			category int
			name     string
			price    float64
		}{
			taskID:   uuid.New(),
			category: 1,
			name:     "New Task Name",
			price:    100.0,
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{
				ID:             uuid.New(),
				Name:           "Old Task Name", //change name from
				PricePerSingle: 100.0,
				Category:       1,
			}, nil)
			fields.taskRepoMock.EXPECT().Update(gomock.Any()).Return(&models.Task{
				ID:             uuid.New(),
				Name:           "New Task Name", //to
				PricePerSingle: 100.0,
				Category:       1,
			}, nil)
		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.NoError(t, err)
		},
	},
	{
		testName: "same name",
		inputData: struct {
			taskID   uuid.UUID
			category int
			name     string
			price    float64
		}{
			taskID:   uuid.New(),
			category: 1,
			name:     "Old Task Name",
			price:    100.0,
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{
				ID:             uuid.New(),
				Name:           "Old Task Name", //no diff
				PricePerSingle: 100.0,
				Category:       1,
			}, nil)
			fields.taskRepoMock.EXPECT().Update(gomock.Any()).Return(&models.Task{
				ID:             uuid.New(),
				Name:           "Old Task Name", //no diff
				PricePerSingle: 100.0,
				Category:       1,
			}, nil)
		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.NoError(t, err)
		},
	},
}

var testTaskChangeNameFail = []struct {
	testName  string
	inputData struct {
		taskID   uuid.UUID
		category int
		name     string
		price    float64
	}
	prepare     func(fields *taskServiceFields)
	checkOutput func(t *testing.T, task *models.Task, err error)
}{
	{
		testName: "empty name",
		inputData: struct {
			taskID   uuid.UUID
			category int
			name     string
			price    float64
		}{
			taskID:   uuid.New(),
			category: 1,
			name:     "",
			price:    100.0,
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(nil, service_errors.InvalidName)

		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.Error(t, err)
		},
	},
	{
		testName: "task not found",
		inputData: struct {
			taskID   uuid.UUID
			category int
			name     string
			price    float64
		}{
			taskID:   uuid.New(),
			category: 1,
			name:     "New Task Name",
			price:    100.0,
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.Error(t, err)
		},
	},
}

func TestTaskServiceChangeName(t *testing.T) {
	for _, tt := range testTaskChangeNameSuccess {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := initTaskServiceFields(ctrl)
			taskService := initTaskService(fields)

			tt.prepare(fields)

			task, err := taskService.Update(tt.inputData.taskID, tt.inputData.category, tt.inputData.name, tt.inputData.price)
			tt.checkOutput(t, task, err)
		})
	}

	for _, tt := range testTaskChangeNameFail {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := initTaskServiceFields(ctrl)
			taskService := initTaskService(fields)

			tt.prepare(fields)

			task, err := taskService.Update(tt.inputData.taskID, tt.inputData.category, tt.inputData.name, tt.inputData.price)
			tt.checkOutput(t, task, err)
		})
	}
}

var testTaskChangePriceSuccess = []struct {
	testName  string
	inputData struct {
		taskID   uuid.UUID
		category int
		name     string
		price    float64
	}
	prepare     func(fields *taskServiceFields)
	checkOutput func(t *testing.T, task *models.Task, err error)
}{
	{
		testName: "basic change",
		inputData: struct {
			taskID   uuid.UUID
			category int
			name     string
			price    float64
		}{
			taskID:   uuid.New(),
			category: 1,
			name:     "Test Task",
			price:    200.0,
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{
				ID:             uuid.New(),
				Name:           "Test Task",
				PricePerSingle: 100.0, //change price from
				Category:       1,
			}, nil)
			fields.taskRepoMock.EXPECT().Update(gomock.Any()).Return(&models.Task{
				ID:             uuid.New(),
				Name:           "Test Task",
				PricePerSingle: 200.0, //to
				Category:       1,
			}, nil)
		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.NoError(t, err)
		},
	},
	{
		testName: "same price",
		inputData: struct {
			taskID   uuid.UUID
			category int
			name     string
			price    float64
		}{
			taskID:   uuid.New(),
			category: 1,
			name:     "Test Task",
			price:    100.0,
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{
				ID:             uuid.New(),
				Name:           "Test Task",
				PricePerSingle: 100.0, //no diff
				Category:       1,
			}, nil)
			fields.taskRepoMock.EXPECT().Update(gomock.Any()).Return(&models.Task{
				ID:             uuid.New(),
				Name:           "Test Task",
				PricePerSingle: 100.0, //no diff
				Category:       1,
			}, nil)
		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.NoError(t, err)
		},
	},
}

var testTaskChangePriceFail = []struct {
	testName  string
	inputData struct {
		taskID   uuid.UUID
		category int
		name     string
		price    float64
	}
	prepare     func(fields *taskServiceFields)
	checkOutput func(t *testing.T, task *models.Task, err error)
}{
	{
		testName: "negative price",
		inputData: struct {
			taskID   uuid.UUID
			category int
			name     string
			price    float64
		}{
			taskID:   uuid.New(),
			category: 1,
			name:     "Test Task",
			price:    -1.0,
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(nil, service_errors.InvalidPrice)
		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.Error(t, err)
		},
	},
	{
		testName: "task not found",
		inputData: struct {
			taskID   uuid.UUID
			category int
			name     string
			price    float64
		}{
			taskID:   uuid.New(),
			category: 1,
			name:     "Test Task",
			price:    200.0,
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.Error(t, err)
		},
	},
}

func TestTaskServiceChangePrice(t *testing.T) {
	for _, tt := range testTaskChangePriceSuccess {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := initTaskServiceFields(ctrl)
			taskService := initTaskService(fields)

			tt.prepare(fields)

			task, err := taskService.Update(tt.inputData.taskID, tt.inputData.category, tt.inputData.name, tt.inputData.price)
			tt.checkOutput(t, task, err)
		})
	}

	for _, tt := range testTaskChangePriceFail {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := initTaskServiceFields(ctrl)
			taskService := initTaskService(fields)

			tt.prepare(fields)

			task, err := taskService.Update(tt.inputData.taskID, tt.inputData.category, tt.inputData.name, tt.inputData.price)
			tt.checkOutput(t, task, err)
		})
	}
}

var testTaskDeleteSuccess = []struct {
	testName  string
	inputData struct {
		taskID uuid.UUID
	}
	prepare     func(fields *taskServiceFields)
	checkOutput func(t *testing.T, err error)
}{
	{
		testName: "basic delete",
		inputData: struct {
			taskID uuid.UUID
		}{
			taskID: uuid.New(),
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{
				ID:             uuid.New(),
				Name:           "Test Task",
				PricePerSingle: 100.0,
				Category:       1,
			}, nil)
			fields.taskRepoMock.EXPECT().Delete(gomock.Any()).Return(nil)
		},
		checkOutput: func(t *testing.T, err error) {
			assert.NoError(t, err)
		},
	},
}

var testTaskDeleteFail = []struct {
	testName  string
	inputData struct {
		taskID uuid.UUID
	}
	prepare     func(fields *taskServiceFields)
	checkOutput func(t *testing.T, err error)
}{
	{
		testName: "task not found",
		inputData: struct {
			taskID uuid.UUID
		}{
			taskID: uuid.New(),
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, err error) {
			assert.Error(t, err)
		},
	},
}

func TestTaskServiceDelete(t *testing.T) {
	for _, tt := range testTaskDeleteSuccess {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := initTaskServiceFields(ctrl)
			taskService := initTaskService(fields)

			tt.prepare(fields)

			err := taskService.Delete(tt.inputData.taskID)
			tt.checkOutput(t, err)
		})
	}

	for _, tt := range testTaskDeleteFail {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := initTaskServiceFields(ctrl)
			taskService := initTaskService(fields)

			tt.prepare(fields)

			err := taskService.Delete(tt.inputData.taskID)
			tt.checkOutput(t, err)
		})
	}
}

var testTaskGetAllSuccess = []struct {
	testName    string
	prepare     func(fields *taskServiceFields)
	checkOutput func(t *testing.T, tasks []models.Task, err error)
}{
	{
		testName: "basic get all",
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetAllTasks().Return([]models.Task{
				{
					ID:             uuid.New(),
					Name:           "Test Task 1",
					PricePerSingle: 100.0,
					Category:       1,
				},
				{
					ID:             uuid.New(),
					Name:           "Test Task 2",
					PricePerSingle: 200.0,
					Category:       2,
				},
			}, nil)
		},
		checkOutput: func(t *testing.T, tasks []models.Task, err error) {
			assert.NoError(t, err)
			assert.Len(t, tasks, 2)
			assert.Equal(t, "Test Task 1", tasks[0].Name)
			assert.Equal(t, 100.0, tasks[0].PricePerSingle)
			assert.Equal(t, 1, tasks[0].Category)
			assert.Equal(t, "Test Task 2", tasks[1].Name)
			assert.Equal(t, 200.0, tasks[1].PricePerSingle)
			assert.Equal(t, 2, tasks[1].Category)
		},
	},

	{
		testName: "empty list",
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetAllTasks().Return([]models.Task{}, nil)
		},
		checkOutput: func(t *testing.T, tasks []models.Task, err error) {
			assert.NoError(t, err)
			assert.Len(t, tasks, 0)
		},
	},
}

func TestTaskServiceGetAll(t *testing.T) {
	for _, tt := range testTaskGetAllSuccess {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := initTaskServiceFields(ctrl)
			taskService := initTaskService(fields)

			tt.prepare(fields)

			tasks, err := taskService.GetAllTasks()
			tt.checkOutput(t, tasks, err)
		})
	}
}

var testTaskGetByIDSuccess = []struct {
	testName  string
	inputData struct {
		taskID uuid.UUID
	}
	prepare     func(fields *taskServiceFields)
	checkOutput func(t *testing.T, task *models.Task, err error)
}{
	{
		testName: "basic get by id",
		inputData: struct {
			taskID uuid.UUID
		}{
			taskID: uuid.New(),
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(&models.Task{
				ID:             uuid.New(),
				Name:           "Test Task",
				PricePerSingle: 100.0,
				Category:       1,
			}, nil)
		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.NoError(t, err)
			assert.Equal(t, "Test Task", task.Name)
			assert.Equal(t, 100.0, task.PricePerSingle)
			assert.Equal(t, 1, task.Category)
		},
	},
}

var testTaskGetByIDFail = []struct {
	testName  string
	inputData struct {
		taskID uuid.UUID
	}
	prepare     func(fields *taskServiceFields)
	checkOutput func(t *testing.T, task *models.Task, err error)
}{
	{
		testName: "task not found",
		inputData: struct {
			taskID uuid.UUID
		}{
			taskID: uuid.New(),
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.Error(t, err)
			assert.Nil(t, task)
		},
	},
}

func TestTaskServiceGetByID(t *testing.T) {
	for _, tt := range testTaskGetByIDSuccess {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := initTaskServiceFields(ctrl)
			taskService := initTaskService(fields)

			tt.prepare(fields)

			task, err := taskService.GetTaskByID(tt.inputData.taskID)
			tt.checkOutput(t, task, err)
		})
	}

	for _, tt := range testTaskGetByIDFail {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := initTaskServiceFields(ctrl)
			taskService := initTaskService(fields)

			tt.prepare(fields)

			task, err := taskService.GetTaskByID(tt.inputData.taskID)
			tt.checkOutput(t, task, err)
		})
	}
}

var testTaskGetByNameSuccess = []struct {
	testName  string
	inputData struct {
		name string
	}
	prepare     func(fields *taskServiceFields)
	checkOutput func(t *testing.T, task *models.Task, err error)
}{
	{
		testName: "basic get by name",
		inputData: struct {
			name string
		}{
			name: "Test Task",
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByName(gomock.Any()).Return(&models.Task{
				ID:             uuid.New(),
				Name:           "Test Task",
				PricePerSingle: 100.0,
				Category:       1,
			}, nil)
		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.NoError(t, err)
			assert.Equal(t, "Test Task", task.Name)
			assert.Equal(t, 100.0, task.PricePerSingle)
			assert.Equal(t, 1, task.Category)
		},
	},
}

var testTaskGetByNameFail = []struct {
	testName  string
	inputData struct {
		name string
	}
	prepare     func(fields *taskServiceFields)
	checkOutput func(t *testing.T, task *models.Task, err error)
}{
	{
		testName: "task not found",
		inputData: struct {
			name string
		}{
			name: "Test Task",
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTaskByName(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, task *models.Task, err error) {
			assert.Error(t, err)
			assert.Nil(t, task)
		},
	},
}

func TestTaskServiceGetByName(t *testing.T) {
	for _, tt := range testTaskGetByNameSuccess {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := initTaskServiceFields(ctrl)
			taskService := initTaskService(fields)

			tt.prepare(fields)

			task, err := taskService.GetTaskByName(tt.inputData.name)
			tt.checkOutput(t, task, err)
		})
	}

	for _, tt := range testTaskGetByNameFail {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := initTaskServiceFields(ctrl)
			taskService := initTaskService(fields)

			tt.prepare(fields)

			task, err := taskService.GetTaskByName(tt.inputData.name)
			tt.checkOutput(t, task, err)
		})
	}
}

var testTaskGetInCategorySuccess = []struct {
	testName  string
	inputData struct {
		category int
	}
	prepare     func(fields *taskServiceFields)
	checkOutput func(t *testing.T, tasks []models.Task, err error)
}{
	{
		testName: "basic get in category",
		inputData: struct {
			category int
		}{
			category: 1,
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTasksInCategory(gomock.Any()).Return([]models.Task{
				{
					ID:             uuid.New(),
					Name:           "Test Task 1",
					PricePerSingle: 100.0,
					Category:       1,
				},
				{
					ID:             uuid.New(),
					Name:           "Test Task 2",
					PricePerSingle: 200.0,
					Category:       1,
				},
			}, nil)
		},
		checkOutput: func(t *testing.T, tasks []models.Task, err error) {
			assert.NoError(t, err)
			assert.Len(t, tasks, 2)
			assert.Equal(t, "Test Task 1", tasks[0].Name)
			assert.Equal(t, 100.0, tasks[0].PricePerSingle)
			assert.Equal(t, 1, tasks[0].Category)
			assert.Equal(t, "Test Task 2", tasks[1].Name)
			assert.Equal(t, 200.0, tasks[1].PricePerSingle)
			assert.Equal(t, 1, tasks[1].Category)
		},
	},
	{
		testName: "empty list",
		inputData: struct {
			category int
		}{
			category: 1,
		},
		prepare: func(fields *taskServiceFields) {
			fields.taskRepoMock.EXPECT().GetTasksInCategory(gomock.Any()).Return([]models.Task{}, nil)
		},
		checkOutput: func(t *testing.T, tasks []models.Task, err error) {
			assert.NoError(t, err)
			assert.Len(t, tasks, 0)
		},
	},
}

var testTaskGetInCategoryFail = []struct {
	testName  string
	inputData struct {
		category int
	}
	prepare     func(fields *taskServiceFields)
	checkOutput func(t *testing.T, tasks []models.Task, err error)
}{
	{
		testName: "invalid category",
		inputData: struct {
			category int
		}{
			category: 0,
		},
		prepare: func(fields *taskServiceFields) {},
		checkOutput: func(t *testing.T, tasks []models.Task, err error) {
			assert.Error(t, err)
		},
	},
}

func TestTaskServiceGetInCategory(t *testing.T) {
	for _, tt := range testTaskGetInCategorySuccess {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := initTaskServiceFields(ctrl)
			taskService := initTaskService(fields)

			tt.prepare(fields)

			tasks, err := taskService.GetTasksInCategory(tt.inputData.category)
			tt.checkOutput(t, tasks, err)
		})
	}

	for _, tt := range testTaskGetInCategoryFail {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := initTaskServiceFields(ctrl)
			taskService := initTaskService(fields)

			tt.prepare(fields)

			tasks, err := taskService.GetTasksInCategory(tt.inputData.category)
			tt.checkOutput(t, tasks, err)
		})
	}
}
