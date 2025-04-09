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

type userServiceFields struct {
	userRepoMock *mock_repository_interfaces.MockIUserRepository
	logger       *log.Logger
	hash         *mock_password_hash.MockPasswordHash
}

func initUserServiceFields(ctrl *gomock.Controller) *userServiceFields {
	userRepoMock := mock_repository_interfaces.NewMockIUserRepository(ctrl)
	f, err := os.OpenFile("tests.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(f)

	return &userServiceFields{
		userRepoMock: userRepoMock,
		hash:         mock_password_hash.NewMockPasswordHash(ctrl),
		logger:       logger,
	}
}

func initUserService(fields *userServiceFields) service_interfaces.IUserService {
	return services.NewUserService(fields.userRepoMock, fields.hash, fields.logger)
}

var testUserGetByIDSuccess = []struct {
	testName  string
	inputData struct {
		id uuid.UUID
	}
	prepare     func(fields *userServiceFields)
	checkOutput func(t *testing.T, user *models.User, err error)
}{
	{
		testName: "basic get by id",
		inputData: struct {
			id uuid.UUID
		}{
			id: uuid.New(),
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(&models.User{}, nil)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, user)
		},
	},
}

var testUserGetByIDFail = []struct {
	testName  string
	inputData struct {
		id uuid.UUID
	}
	prepare     func(fields *userServiceFields)
	checkOutput func(t *testing.T, user *models.User, err error)
}{
	{
		testName: "user not found",
		inputData: struct {
			id uuid.UUID
		}{
			id: uuid.New(),
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
			assert.Nil(t, user)
		},
	},
}

func TestUserServiceGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initUserServiceFields(ctrl)
	service := initUserService(fields)

	for _, test := range testUserGetByIDSuccess {
		t.Run(test.testName, func(t *testing.T) {
			test.prepare(fields)
			user, err := service.GetUserByID(test.inputData.id)
			test.checkOutput(t, user, err)
		})
	}

	for _, test := range testUserGetByIDFail {
		t.Run(test.testName, func(t *testing.T) {
			test.prepare(fields)
			user, err := service.GetUserByID(test.inputData.id)
			test.checkOutput(t, user, err)
		})
	}
}

var testUserLoginSuccess = []struct {
	testName  string
	inputData struct {
		email    string
		password string
	}
	prepare     func(fields *userServiceFields)
	checkOutput func(t *testing.T, user *models.User, err error)
}{
	{
		testName: "basic login",
		inputData: struct {
			email    string
			password string
		}{
			email:    "test@gmail.com",
			password: "password123",
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByEmail(gomock.Any()).Return(&models.User{Password: "password123"}, nil)
			fields.hash.EXPECT().CompareHashAndPassword(gomock.Any(), gomock.Any()).Return(true)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, user)
		},
	},
}

var testUserLoginFail = []struct {
	testName  string
	inputData struct {
		email    string
		password string
	}
	prepare     func(fields *userServiceFields)
	checkOutput func(t *testing.T, user *models.User, err error)
}{
	{
		testName: "user not found",
		inputData: struct {
			email    string
			password string
		}{
			email:    "test@gmail.com",
			password: "password123",
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByEmail(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
			assert.Nil(t, user)
		},
	},
	{
		testName: "password mismatch",
		inputData: struct {
			email    string
			password string
		}{
			email:    "test@gmail.com",
			password: "passwordaaa",
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByEmail(gomock.Any()).Return(&models.User{Password: "password123"}, nil)
			fields.hash.EXPECT().CompareHashAndPassword(gomock.Any(), gomock.Any()).Return(false)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
			assert.Nil(t, user)
		},
	},
}

func TestUserServiceLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initUserServiceFields(ctrl)
	service := initUserService(fields)

	for _, test := range testUserLoginSuccess {
		t.Run(test.testName, func(t *testing.T) {
			test.prepare(fields)
			user, err := service.Login(test.inputData.email, test.inputData.password)
			test.checkOutput(t, user, err)
		})
	}

	for _, test := range testUserLoginFail {
		t.Run(test.testName, func(t *testing.T) {
			test.prepare(fields)
			user, err := service.Login(test.inputData.email, test.inputData.password)
			test.checkOutput(t, user, err)
		})
	}
}

var testUserChangePasswordSuccess = []struct {
	testName  string
	inputData struct {
		id          uuid.UUID
		name        string
		surname     string
		email       string
		address     string
		phoneNumber string
		password    string
	}
	prepare     func(fields *userServiceFields)
	checkOutput func(t *testing.T, user *models.User, err error)
}{
	{
		testName: "basic change password",
		inputData: struct {
			id          uuid.UUID
			name        string
			surname     string
			email       string
			address     string
			phoneNumber string
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "+79999999999",
			password:    "password123",
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(&models.User{
				Password: "password123",
			}, nil)
			fields.hash.EXPECT().GetHash(gomock.Any()).Return("password123", nil)
			fields.userRepoMock.EXPECT().Update(gomock.Any()).Return(&models.User{}, nil)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.NoError(t, err)
		},
	},
}

var testUserChangePasswordFail = []struct {
	testName  string
	inputData struct {
		id          uuid.UUID
		name        string
		surname     string
		email       string
		address     string
		phoneNumber string
		password    string
	}
	prepare     func(fields *userServiceFields)
	checkOutput func(t *testing.T, user *models.User, err error)
}{
	{
		testName: "user not found",
		inputData: struct {
			id          uuid.UUID
			name        string
			surname     string
			email       string
			address     string
			phoneNumber string
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "+79999999999",
			password:    "password123",
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
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
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "+79999999999",
			password:    "pass",
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(&models.User{}, nil)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
}

func TestUserServiceChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initUserServiceFields(ctrl)
	service := initUserService(fields)

	for _, test := range testUserChangePasswordSuccess {
		t.Run(test.testName, func(t *testing.T) {
			test.prepare(fields)
			user, err := service.Update(test.inputData.id, test.inputData.name, test.inputData.surname, test.inputData.email, test.inputData.address, test.inputData.phoneNumber, test.inputData.password)
			test.checkOutput(t, user, err)
		})
	}

	for _, test := range testUserChangePasswordFail {
		t.Run(test.testName, func(t *testing.T) {
			test.prepare(fields)
			user, err := service.Update(test.inputData.id, test.inputData.name, test.inputData.surname, test.inputData.email, test.inputData.address, test.inputData.phoneNumber, test.inputData.password)
			test.checkOutput(t, user, err)
		})
	}
}

var testUserRegisterSuccess = []struct {
	testName  string
	inputData struct {
		user     *models.User
		password string
	}
	prepare     func(fields *userServiceFields)
	checkOutput func(t *testing.T, user *models.User, err error)
}{
	{
		testName: "basic register",
		inputData: struct {
			user     *models.User
			password string
		}{
			user: &models.User{
				Email:       "test@gmail.com",
				Name:        "Test",
				Surname:     "Test",
				Address:     "Test",
				PhoneNumber: "+79999999999",
			},
			password: "password123",
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByEmail(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
			fields.hash.EXPECT().GetHash(gomock.Any()).Return("password123", nil)
			fields.userRepoMock.EXPECT().Create(gomock.Any()).Return(&models.User{}, nil)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, user)
		},
	},
}

var testUserRegisterFail = []struct {
	testName  string
	inputData struct {
		user     *models.User
		password string
	}
	prepare     func(fields *userServiceFields)
	checkOutput func(t *testing.T, user *models.User, err error)
}{
	{
		testName: "user already exists",
		inputData: struct {
			user     *models.User
			password string
		}{
			user: &models.User{
				Email:       "test@gmail.com",
				Name:        "Test",
				Surname:     "Test",
				Address:     "Test",
				PhoneNumber: "+79999999999",
			},
			password: "password123",
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByEmail(gomock.Any()).Return(&models.User{}, nil)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
			assert.Equal(t, fmt.Errorf("SERVICE: User with email exists"), err)
		},
	},
	{
		testName: "invalid email",
		inputData: struct {
			user     *models.User
			password string
		}{
			user: &models.User{
				Email:       "not-valid",
				Name:        "Test",
				Surname:     "Test",
				Address:     "Test",
				PhoneNumber: "+79999999999",
			},
			password: "password123",
		},
		prepare: func(fields *userServiceFields) {},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
	{
		testName: "invalid name",
		inputData: struct {
			user     *models.User
			password string
		}{
			user: &models.User{
				Email:       "test@gmail.com",
				Name:        "",
				Surname:     "",
				Address:     "Test",
				PhoneNumber: "+79999999999",
			},
			password: "password123",
		},
		prepare: func(fields *userServiceFields) {},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
	{
		testName: "invalid address",
		inputData: struct {
			user     *models.User
			password string
		}{
			user: &models.User{
				Address:     "",
				Email:       "test@gmail.com",
				Name:        "Test",
				Surname:     "Test",
				PhoneNumber: "+79999999999",
			},
			password: "password123",
		},
		prepare: func(fields *userServiceFields) {},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
	{
		testName: "invalid phone number",
		inputData: struct {
			user     *models.User
			password string
		}{
			user: &models.User{
				PhoneNumber: "123",
				Email:       "test@gmail.com",
				Name:        "Test",
				Surname:     "Test",
				Address:     "Test",
			},
			password: "password123",
		},
		prepare: func(fields *userServiceFields) {},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
	{
		testName: "invalid password",
		inputData: struct {
			user     *models.User
			password string
		}{
			user: &models.User{
				Email:       "test@gmail.com",
				Name:        "Test",
				Surname:     "Test",
				Address:     "Test",
				PhoneNumber: "+79999999999",
			},
			password: "admin",
		},
		prepare: func(fields *userServiceFields) {},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
}

func TestUserServiceRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initUserServiceFields(ctrl)
	service := initUserService(fields)

	for _, test := range testUserRegisterSuccess {
		t.Run(test.testName, func(t *testing.T) {
			test.prepare(fields)
			user, err := service.Register(test.inputData.user, test.inputData.password)
			test.checkOutput(t, user, err)
		})
	}

	for _, test := range testUserRegisterFail {
		t.Run(test.testName, func(t *testing.T) {
			test.prepare(fields)
			user, err := service.Register(test.inputData.user, test.inputData.password)
			test.checkOutput(t, user, err)
		})
	}
}

var testUserUpdatePersonalInformationSuccess = []struct {
	testName  string
	inputData struct {
		id          uuid.UUID
		name        string
		surname     string
		email       string
		address     string
		phoneNumber string
		password    string
	}
	prepare     func(fields *userServiceFields)
	checkOutput func(t *testing.T, user *models.User, err error)
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
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "+79999999999",
			password:    "password123",
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(&models.User{}, nil)
			fields.hash.EXPECT().GetHash(gomock.Any()).Return("password123", nil)
			fields.userRepoMock.EXPECT().Update(gomock.Any()).Return(&models.User{
				Name:        "Test",
				Surname:     "Test",
				Email:       "test@gmail.com",
				Address:     "Test",
				PhoneNumber: "+79999999999",
				Password:    "password123",
			}, nil)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, "Test", user.Name)
			assert.Equal(t, "Test", user.Surname)
			assert.Equal(t, "test@gmail.com", user.Email)
			assert.Equal(t, "Test", user.Address)
			assert.Equal(t, "+79999999999", user.PhoneNumber)
		},
	},
}

var testUserUpdatePersonalInformationFail = []struct {
	testName  string
	inputData struct {
		id          uuid.UUID
		name        string
		surname     string
		email       string
		address     string
		phoneNumber string
		password    string
	}
	prepare     func(fields *userServiceFields)
	checkOutput func(t *testing.T, user *models.User, err error)
}{
	{
		testName: "user not found",
		inputData: struct {
			id          uuid.UUID
			name        string
			surname     string
			email       string
			address     string
			phoneNumber string
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "+79999999999",
			password:    "password123",
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(nil, repository_errors.DoesNotExist)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
			assert.Nil(t, user)
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
			password    string
		}{
			id:          uuid.New(),
			name:        "", //invalid name
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "+79999999999",
			password:    "password123",
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(&models.User{}, nil)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
			assert.Nil(t, user)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
	{
		testName: "invalid surname",
		inputData: struct {
			id          uuid.UUID
			name        string
			surname     string
			email       string
			address     string
			phoneNumber string
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "", //invalid surname
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "+79999999999",
			password:    "password123",
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(&models.User{}, nil)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
			assert.Nil(t, user)
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
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "not-valid", //invalid email
			address:     "Test",
			phoneNumber: "+79999999999",
			password:    "password123",
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(&models.User{}, nil)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
			assert.Nil(t, user)
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
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "", //invalid address
			phoneNumber: "+79999999999",
			password:    "password123",
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(&models.User{}, nil)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
			assert.Nil(t, user)
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
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "123", //invalid phone number
			password:    "password123",
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(&models.User{}, nil)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
			assert.Nil(t, user)
			assert.Equal(t, fmt.Errorf("SERVICE: Invalid input"), err)
		},
	},
	{
		testName: "update failed",
		inputData: struct {
			id          uuid.UUID
			name        string
			surname     string
			email       string
			address     string
			phoneNumber string
			password    string
		}{
			id:          uuid.New(),
			name:        "Test",
			surname:     "Test",
			email:       "test@gmail.com",
			address:     "Test",
			phoneNumber: "+79999999999",
			password:    "password123",
		},
		prepare: func(fields *userServiceFields) {
			fields.userRepoMock.EXPECT().GetUserByID(gomock.Any()).Return(&models.User{}, nil)
			fields.hash.EXPECT().GetHash(gomock.Any()).Return("password123", nil)
			fields.userRepoMock.EXPECT().Update(gomock.Any()).Return(nil, repository_errors.UpdateError)
		},
		checkOutput: func(t *testing.T, user *models.User, err error) {
			assert.Error(t, err)
			assert.Nil(t, user)
			assert.Equal(t, repository_errors.UpdateError, err)
		},
	},
}

func TestUserServiceUpdatePersonalInformation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fields := initUserServiceFields(ctrl)
	service := initUserService(fields)

	for _, test := range testUserUpdatePersonalInformationSuccess {
		t.Run(test.testName, func(t *testing.T) {
			test.prepare(fields)
			user, err := service.Update(test.inputData.id, test.inputData.name, test.inputData.surname, test.inputData.email, test.inputData.address, test.inputData.phoneNumber, test.inputData.password)
			test.checkOutput(t, user, err)
		})
	}

	for _, test := range testUserUpdatePersonalInformationFail {
		t.Run(test.testName, func(t *testing.T) {
			test.prepare(fields)
			user, err := service.Update(test.inputData.id, test.inputData.name, test.inputData.surname, test.inputData.email, test.inputData.address, test.inputData.phoneNumber, test.inputData.password)
			test.checkOutput(t, user, err)
		})
	}
}
