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
	"lab3/internal/services/service_interfaces"
	mock_password_hash "lab3/tests/hasher_mocks"
	mock_repository_interfaces "lab3/tests/repository_mocks"
	"os"
	"testing"
)

type workerServiceFields struct {
	workerRepoMock *mock_repository_interfaces.MockIWorkerRepository
	logger         *log.Logger
	hash           *mock_password_hash.MockPasswordHash
}

func initWorkerServiceFields(ctrl *gomock.Controller) *workerServiceFields {
	workerRepoMock := mock_repository_interfaces.NewMockIWorkerRepository(ctrl)
	f, err := os.OpenFile("tests.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(f)

	return &workerServiceFields{
		workerRepoMock: workerRepoMock,
		hash:           mock_password_hash.NewMockPasswordHash(ctrl),
		logger:         logger,
	}
}

func initWorkerService(fields *workerServiceFields) service_interfaces.IWorkerService {
	return services.NewWorkerService(fields.workerRepoMock, fields.hash, fields.logger)
}

var testWorkerGetByID = []struct {
	testName  string
	inputData struct {
		id uuid.UUID
	}
	prepare     func(fields *workerServiceFields)
	checkOutput func(t *testing.T, worker *models.Worker, err error)
}{
	{
		testName: "Success",
		inputData: struct {
			id uuid.UUID
		}{
			id: uuid.New(),
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{
				Name:        "Test",
				Surname:     "Test",
				Email:       "test@gmail.com",
				Address:     "Test",
				PhoneNumber: "Test",
				Role:        1,
			}, nil)
		},
		checkOutput: func(t *testing.T, worker *models.Worker, err error) {
			assert.NoError(t, err)
			assert.Equal(t, "Test", worker.Name)
			assert.Equal(t, "Test", worker.Surname)
			assert.Equal(t, "test@gmail.com", worker.Email)
			assert.Equal(t, "Test", worker.Address)
			assert.Equal(t, "Test", worker.PhoneNumber)
			assert.Equal(t, 1, worker.Role)
		},
	},
	{
		testName: "worker not found",
		inputData: struct {
			id uuid.UUID
		}{
			id: uuid.New(),
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Nil(t, worker)
		},
	},
}

func TestWorkerService_GetWorkerByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initWorkerServiceFields(ctrl)
	service := initWorkerService(fields)

	for _, tt := range testWorkerGetByID {
		tt.prepare(fields)
		t.Run(tt.testName, func(t *testing.T) {
			worker, err := service.GetWorkerByID(tt.inputData.id)
			tt.checkOutput(t, worker, err)
		})
	}
}

var testWorkerGetAllWorkers = []struct {
	testName  string
	prepare   func(fields *workerServiceFields)
	checkFunc func(t *testing.T, workers []models.Worker, err error)
}{
	{
		testName: "Success",
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetAllWorkers().Return([]models.Worker{
				{
					Name:        "Test",
					Surname:     "Test",
					Email:       "test@gmail.com",
					Address:     "Test",
					PhoneNumber: "+79999999999",
					Role:        1,
				},
				{
					Name:        "Test 2",
					Surname:     "Test 2",
					Email:       "test2@gmail.com",
					Address:     "Test 2",
					PhoneNumber: "+79999999988",
					Role:        2,
				},
				{
					Name: "Test 3",
				},
				{
					Name: "Test 4",
				},
			}, nil)
		},
		checkFunc: func(t *testing.T, workers []models.Worker, err error) {
			assert.NoError(t, err)
			assert.Equal(t, 4, len(workers))
			assert.Equal(t, "Test", workers[0].Name)
			assert.Equal(t, "Test 2", workers[1].Name)
			assert.Equal(t, "test@gmail.com", workers[0].Email)
			assert.Equal(t, "test2@gmail.com", workers[1].Email)
			assert.Equal(t, "+79999999999", workers[0].PhoneNumber)
			assert.Equal(t, "+79999999988", workers[1].PhoneNumber)
		},
	},
	{
		testName: "empty list",
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetAllWorkers().Return([]models.Worker{}, nil)
		},
		checkFunc: func(t *testing.T, workers []models.Worker, err error) {
			assert.NoError(t, err)
			assert.Equal(t, 0, len(workers))
		},
	},
}

func TestWorkerService_GetAllWorkers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initWorkerServiceFields(ctrl)
	service := initWorkerService(fields)

	for _, tt := range testWorkerGetAllWorkers {
		tt.prepare(fields)
		t.Run(tt.testName, func(t *testing.T) {
			workers, err := service.GetAllWorkers()
			tt.checkFunc(t, workers, err)
		})
	}
}

var testWorkerDelete = []struct {
	testName  string
	inputData struct {
		id uuid.UUID
	}
	prepare   func(fields *workerServiceFields)
	checkFunc func(t *testing.T, err error)
}{
	{
		testName:  "Success",
		inputData: struct{ id uuid.UUID }{id: uuid.New()},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{}, nil)
			fields.workerRepoMock.EXPECT().Delete(gomock.Any()).Return(nil)
		},
		checkFunc: func(t *testing.T, err error) {
			assert.NoError(t, err)
		},
	},
	{
		testName:  "worker not found",
		inputData: struct{ id uuid.UUID }{id: uuid.New()},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkFunc: func(t *testing.T, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DoesNotExist, err)
		},
	},
}

func TestWorkerService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initWorkerServiceFields(ctrl)
	service := initWorkerService(fields)

	for _, tt := range testWorkerDelete {
		tt.prepare(fields)
		t.Run(tt.testName, func(t *testing.T) {
			err := service.Delete(tt.inputData.id)
			tt.checkFunc(t, err)
		})
	}
}

// ----------------------------------------
var testWorkerChangePassword = []struct {
	testName  string
	inputData struct {
		id          uuid.UUID
		name        string
		surname     string
		email       string
		address     string
		phoneNumber string
		role        int
		password    string
	}
	prepare   func(fields *workerServiceFields)
	checkFunc func(t *testing.T, worker *models.Worker, err error)
}{
	{
		testName: "Success",
		inputData: struct {
			id          uuid.UUID
			name        string
			surname     string
			email       string
			address     string
			phoneNumber string
			role        int
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "+79999999999",
			role:        1,
			password:    "password123", //change password from
		},
		prepare: func(fields *workerServiceFields) {
			var worker = &models.Worker{
				ID:          uuid.New(),
				Name:        "Test",
				Surname:     "Test",
				Email:       "test@gmail.com",
				Address:     "Test",
				PhoneNumber: "+79999999999",
				Role:        1,
				Password:    "password", //to
			}
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(worker, nil)
			fields.hash.EXPECT().GetHash(gomock.Any()).Return("hash", nil)
			fields.workerRepoMock.EXPECT().Update(worker).Return(worker, nil)
		},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.NoError(t, err)
		},
	},
	{
		testName: "worker not found",
		inputData: struct {
			id          uuid.UUID
			name        string
			surname     string
			email       string
			address     string
			phoneNumber string
			role        int
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "+79999999999",
			role:        1,
			password:    "password123",
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DoesNotExist, err)
		},
	},
	{
		testName: "invalid password",
		inputData: struct {
			id          uuid.UUID
			name        string
			surname     string
			email       string
			address     string
			phoneNumber string
			role        int
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "+79999999999",
			role:        1,
			password:    "123",
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{}, nil)
		},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
}

func TestWorkerService_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initWorkerServiceFields(ctrl)
	service := initWorkerService(fields)

	for _, tt := range testWorkerChangePassword {
		tt.prepare(fields)
		t.Run(tt.testName, func(t *testing.T) {
			worker, err := service.Update(tt.inputData.id, tt.inputData.name, tt.inputData.surname, tt.inputData.email, tt.inputData.address, tt.inputData.phoneNumber, tt.inputData.role, tt.inputData.password)
			tt.checkFunc(t, worker, err)
		})
	}
}

var testWorkerUpdateRole = []struct {
	testName  string
	inputData struct {
		id          uuid.UUID
		name        string
		surname     string
		email       string
		address     string
		phoneNumber string
		role        int
		password    string
	}
	prepare   func(fields *workerServiceFields)
	checkFunc func(t *testing.T, worker *models.Worker, err error)
}{
	{
		testName: "Success",
		inputData: struct {
			id          uuid.UUID
			name        string
			surname     string
			email       string
			address     string
			phoneNumber string
			role        int
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "+79999999999",
			role:        1,
			password:    "password123",
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{
				ID:          uuid.New(),
				Name:        "Test",
				Surname:     "Test",
				Email:       "test@gmail.com",
				Address:     "Test",
				PhoneNumber: "+79999999999",
				Role:        1,
				Password:    "password123",
			}, nil)

			fields.hash.EXPECT().GetHash(gomock.Any()).Return("hash", nil)
			fields.workerRepoMock.EXPECT().Update(gomock.Any()).Return(&models.Worker{
				ID:          uuid.New(),
				Name:        "Test",
				Surname:     "Test",
				Email:       "test@gmail.com",
				Address:     "Test",
				PhoneNumber: "+79999999999",
				Role:        2,
				Password:    "password123",
			}, nil)
		},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.NoError(t, err)
		},
	},
	{
		testName: "worker not found",
		inputData: struct {
			id          uuid.UUID
			name        string
			surname     string
			email       string
			address     string
			phoneNumber string
			role        int
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "+79999999999",
			role:        1,
			password:    "password123",
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Equal(t, repository_errors.DoesNotExist, err)
		},
	},
	{
		testName: "invalid role",
		inputData: struct {
			id          uuid.UUID
			name        string
			surname     string
			email       string
			address     string
			phoneNumber string
			role        int
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "+79999999999",
			role:        5, //invalid role
			password:    "password123",
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{}, nil)
		},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
}

func TestWorkerService_UpdateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initWorkerServiceFields(ctrl)
	service := initWorkerService(fields)

	for _, tt := range testWorkerUpdateRole {
		tt.prepare(fields)
		t.Run(tt.testName, func(t *testing.T) {
			worker, err := service.Update(tt.inputData.id, tt.inputData.name, tt.inputData.surname, tt.inputData.email, tt.inputData.address, tt.inputData.phoneNumber, tt.inputData.role, tt.inputData.password)
			tt.checkFunc(t, worker, err)
		})
	}
}

var testWorkerUpdatePersonalInformation = []struct {
	testName  string
	inputData struct {
		id          uuid.UUID
		name        string
		surname     string
		email       string
		address     string
		phoneNumber string
		role        int
		password    string
	}
	prepare   func(fields *workerServiceFields)
	checkFunc func(t *testing.T, worker *models.Worker, err error)
}{
	{
		testName: "Success",
		inputData: struct {
			id          uuid.UUID
			name        string
			surname     string
			email       string
			address     string
			phoneNumber string
			role        int
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "+79999999999",
			role:        1,
			password:    "password123",
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{
				ID:          uuid.New(),
				Name:        "Test 2",
				Surname:     "Test",
				Email:       "test@email.com",
				Address:     "Test",
				PhoneNumber: "+79999999999",
				Role:        1,
				Password:    "password123",
			}, nil)

			fields.hash.EXPECT().GetHash(gomock.Any()).Return("hash", nil)
			fields.workerRepoMock.EXPECT().Update(gomock.Any()).Return(&models.Worker{
				ID:          uuid.New(),
				Name:        "Test 2",
				Surname:     "Test",
				Email:       "test@email.com",
				Address:     "Test",
				PhoneNumber: "+79999999999",
				Role:        1,
				Password:    "password123",
			}, nil)
		},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.NoError(t, err)
			assert.Equal(t, "Test 2", worker.Name)
			assert.Equal(t, "Test", worker.Surname)
			assert.Equal(t, "test@email.com", worker.Email)
		},
	},
	{
		testName: "worker not found",
		inputData: struct {
			id          uuid.UUID
			name        string
			surname     string
			email       string
			address     string
			phoneNumber string
			role        int
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "+79999999999",
			role:        1,
			password:    "password123",
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Nil(t, worker)
			assert.Equal(t, repository_errors.DoesNotExist, err)
		},
	},
	{
		testName: "invalid name",
		inputData: struct {
			id          uuid.UUID
			name        string
			surname     string
			email       string
			address     string
			phoneNumber string
			role        int
			password    string
		}{
			id:          uuid.New(),
			name:        "", //invalid name
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "+79999999999",
			role:        1,
			password:    "password123",
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{}, nil)
		},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Nil(t, worker)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
	{
		testName: "invalid email",
		inputData: struct {
			id          uuid.UUID
			name        string
			surname     string
			email       string
			address     string
			phoneNumber string
			role        int
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "not-valid.com", //invalid email
			address:     "Test",
			phoneNumber: "+79999999999",
			role:        1,
			password:    "password123",
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{}, nil)
		},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Nil(t, worker)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
	{
		testName: "invalid address",
		inputData: struct {
			id          uuid.UUID
			name        string
			surname     string
			email       string
			address     string
			phoneNumber string
			role        int
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "", //invalid address
			phoneNumber: "+79999999999",
			role:        1,
			password:    "password123",
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{}, nil)
		},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Nil(t, worker)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)

		},
	},
	{
		testName: "invalid phone number",
		inputData: struct {
			id          uuid.UUID
			name        string
			surname     string
			email       string
			address     string
			phoneNumber string
			role        int
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "123", //invalid phone number
			role:        1,
			password:    "password123",
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByID(gomock.Any()).Return(&models.Worker{}, nil)
		},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Nil(t, worker)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
}

func TestWorkerService_UpdatePersonalInformation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initWorkerServiceFields(ctrl)
	service := initWorkerService(fields)

	for _, tt := range testWorkerUpdatePersonalInformation {
		tt.prepare(fields)
		t.Run(tt.testName, func(t *testing.T) {
			worker, err := service.Update(tt.inputData.id, tt.inputData.name, tt.inputData.surname, tt.inputData.email, tt.inputData.address, tt.inputData.phoneNumber, tt.inputData.role, tt.inputData.password)
			tt.checkFunc(t, worker, err)
		})
	}
}

var testWorkerCreate = []struct {
	testName  string
	inputData struct {
		worker   *models.Worker
		password string
	}
	prepare   func(fields *workerServiceFields)
	checkFunc func(t *testing.T, worker *models.Worker, err error)
}{
	{
		testName: "Success",
		inputData: struct {
			worker   *models.Worker
			password string
		}{
			worker: &models.Worker{
				Name:        "Test",
				Surname:     "Test",
				Email:       "test@email.com",
				Address:     "Test",
				PhoneNumber: "+79999999999",
				Role:        1,
			},
			password: "password123",
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByEmail(gomock.Any()).Return(nil, nil)
			fields.hash.EXPECT().GetHash(gomock.Any()).Return("hash", nil)
			fields.workerRepoMock.EXPECT().Create(gomock.Any()).Return(&models.Worker{
				ID:          uuid.New(),
				Name:        "Test",
				Surname:     "Test",
				Email:       "test@email.com",
				Address:     "Test",
				PhoneNumber: "+79999999999",
				Role:        1,
				Password:    "hash",
			}, nil)
		},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.NoError(t, err)
			assert.Equal(t, "Test", worker.Name)
			assert.Equal(t, "Test", worker.Surname)
			assert.Equal(t, "test@email.com", worker.Email)
		},
	},
	{
		testName: "worker already exists",
		inputData: struct {
			worker   *models.Worker
			password string
		}{
			worker: &models.Worker{
				Name:        "Test",
				Surname:     "Test",
				Email:       "test@email.com",
				Address:     "Test",
				PhoneNumber: "+79999999999",
				Role:        1,
			},
			password: "password123",
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByEmail(gomock.Any()).Return(&models.Worker{}, nil)
		},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Nil(t, worker)
			assert.Equal(t, fmt.Errorf("SERVICE: Worker with email already exists"), err)
		},
	},
	{
		testName: "invalid name",
		inputData: struct {
			worker   *models.Worker
			password string
		}{
			worker: &models.Worker{
				Name:        "",
				Surname:     "Test",
				Email:       "test@email.com",
				Address:     "Test",
				PhoneNumber: "+79999999999",
				Role:        1,
			},
			password: "password123",
		},
		prepare: func(fields *workerServiceFields) {},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Nil(t, worker)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
	{
		testName: "invalid email",
		inputData: struct {
			worker   *models.Worker
			password string
		}{
			worker: &models.Worker{
				Name:        "Test",
				Surname:     "Test",
				Email:       "not-valid.com",
				Address:     "Test",
				PhoneNumber: "+79999999999",
				Role:        1,
			},
			password: "password123",
		},
		prepare: func(fields *workerServiceFields) {},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Nil(t, worker)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
	{
		testName: "invalid address",
		inputData: struct {
			worker   *models.Worker
			password string
		}{
			worker: &models.Worker{
				Name:        "Test",
				Surname:     "Test",
				Email:       "test@email.com",
				Address:     "",
				PhoneNumber: "+79999999999",
				Role:        1,
			},
			password: "password123",
		},
		prepare: func(fields *workerServiceFields) {},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Nil(t, worker)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
	{
		testName: "invalid phone number",
		inputData: struct {
			worker   *models.Worker
			password string
		}{
			worker: &models.Worker{
				Name:        "Test",
				Surname:     "Test",
				Email:       "test@email.com",
				Address:     "Test",
				PhoneNumber: "123",
				Role:        1,
			},
			password: "password123",
		},

		prepare: func(fields *workerServiceFields) {},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Nil(t, worker)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},

	{
		testName: "invalid password",
		inputData: struct {
			worker   *models.Worker
			password string
		}{
			worker: &models.Worker{
				Name:        "Test",
				Surname:     "Test",
				Email:       "test@mail.com",
				Address:     "Test",
				PhoneNumber: "+79999999999",
				Role:        1,
			},
			password: "123",
		},
		prepare: func(fields *workerServiceFields) {},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Nil(t, worker)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
}

func TestWorkerServiceCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initWorkerServiceFields(ctrl)
	service := initWorkerService(fields)

	for _, tt := range testWorkerCreate {
		tt.prepare(fields)
		t.Run(tt.testName, func(t *testing.T) {
			worker, err := service.Create(tt.inputData.worker, tt.inputData.password)
			tt.checkFunc(t, worker, err)
		})
	}
}

var testWorkerLogin = []struct {
	testName  string
	inputData struct {
		email    string
		password string
	}
	prepare   func(fields *workerServiceFields)
	checkFunc func(t *testing.T, worker *models.Worker, err error)
}{
	{
		testName: "Success",
		inputData: struct {
			email    string
			password string
		}{
			email:    "test@email.com",
			password: "password123",
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByEmail(gomock.Any()).Return(&models.Worker{
				Name:     "Test",
				Surname:  "Test",
				Email:    "test@email.com",
				Address:  "Test",
				Role:     1,
				Password: "hash",
			}, nil)
			fields.hash.EXPECT().CompareHashAndPassword(gomock.Any(), gomock.Any()).Return(true)
		},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.NoError(t, err)
			assert.Equal(t, "Test", worker.Name)
			assert.Equal(t, "test@email.com", worker.Email)
			assert.Equal(t, "hash", worker.Password)
		},
	},
	{
		testName: "worker not found",
		inputData: struct {
			email    string
			password string
		}{
			email:    "not@found.com",
			password: "password123",
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByEmail(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Nil(t, worker)
			assert.Equal(t, repository_errors.DoesNotExist, err)
		},
	},
	{
		testName: "invalid password",
		inputData: struct {
			email    string
			password string
		}{
			email:    "test@email.com",
			password: "password123",
		},
		prepare: func(fields *workerServiceFields) {
			fields.workerRepoMock.EXPECT().GetWorkerByEmail(gomock.Any()).Return(&models.Worker{
				Name:    "Test",
				Surname: "Test",
				Email:   "test@email.com",
			}, nil)
			fields.hash.EXPECT().CompareHashAndPassword(gomock.Any(), gomock.Any()).Return(false)
		},
		checkFunc: func(t *testing.T, worker *models.Worker, err error) {
			assert.Error(t, err)
			assert.Nil(t, worker)
			assert.Equal(t, fmt.Errorf("SERVICE: Password is incorrect for worker with email"), err)
		},
	},
}

func TestWorkerServiceLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initWorkerServiceFields(ctrl)
	service := initWorkerService(fields)

	for _, tt := range testWorkerLogin {
		tt.prepare(fields)
		t.Run(tt.testName, func(t *testing.T) {
			worker, err := service.Login(tt.inputData.email, tt.inputData.password)
			tt.checkFunc(t, worker, err)
		})
	}
}
