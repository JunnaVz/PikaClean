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

var testUserRepositoryCreateSuccess = []struct {
	TestName  string
	InputData struct {
		user *models.User
	}
	CheckOutput func(t *testing.T, inputData struct{ user *models.User }, createdUser *models.User, err error)
}{
	{
		TestName: "create success test",
		InputData: struct {
			user *models.User
		}{
			&models.User{
				ID:          uuid.New(),
				Name:        "First Name",
				Surname:     "Last Name",
				Address:     "Address",
				PhoneNumber: "+79999999999",
				Email:       "test@email.com",
				Password:    "hashed_password",
			},
		},
		CheckOutput: func(t *testing.T, inputData struct{ user *models.User }, createdUser *models.User, err error) {
			require.NoError(t, err)
			require.Equal(t, inputData.user.Name, createdUser.Name)
			require.Equal(t, inputData.user.Surname, createdUser.Surname)
			require.Equal(t, inputData.user.Address, createdUser.Address)
			require.Equal(t, inputData.user.PhoneNumber, createdUser.PhoneNumber)
			require.Equal(t, inputData.user.Email, createdUser.Email)
			require.Equal(t, inputData.user.Password, createdUser.Password)
		},
	},
}

func TestUserRepositoryCreate(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}
	userRepository := postgres.CreateUserRepository(&fields)

	for _, test := range testUserRepositoryCreateSuccess {
		t.Run(test.TestName, func(t *testing.T) {
			createdUser, err := userRepository.Create(test.InputData.user)
			test.CheckOutput(t, test.InputData, createdUser, err)
		})
	}
}

var testUserRepositoryGetByIDSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdUser *models.User, receivedUser *models.User, err error)
}{
	{
		TestName: "get by id success test",
		CheckOutput: func(t *testing.T, createdUser *models.User, receivedUser *models.User, err error) {
			require.NoError(t, err)
			require.Equal(t, createdUser.ID, receivedUser.ID)
		},
	},
}

func TestUserRepositoryGetByID(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}
	userRepository := postgres.CreateUserRepository(&fields)

	for _, test := range testUserRepositoryGetByIDSuccess {
		t.Run(test.TestName, func(t *testing.T) {
			createdUser, err := userRepository.Create(&models.User{
				Name:        "First Name",
				Surname:     "Last Name",
				Address:     "Address",
				PhoneNumber: "+79999999999",
				Email:       "email@test.ru",
				Password:    "hashed_password",
			})
			require.NoError(t, err)

			receivedUser, err := userRepository.GetUserByID(createdUser.ID)
			test.CheckOutput(t, createdUser, receivedUser, err)

		})
	}
}

var testUserRepositoryGetByEmailSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdUser *models.User, receivedUser *models.User, err error)
}{
	{
		TestName: "get by email success test",
		CheckOutput: func(t *testing.T, createdUser *models.User, receivedUser *models.User, err error) {
			require.NoError(t, err)
			require.Equal(t, createdUser.Email, receivedUser.Email)
		},
	},
}

func TestUserRepositoryGetByEmail(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}
	userRepository := postgres.CreateUserRepository(&fields)

	for _, test := range testUserRepositoryGetByEmailSuccess {
		t.Run(test.TestName, func(t *testing.T) {
			createdUser, err := userRepository.Create(&models.User{
				Name:        "First Name",
				Surname:     "Last Name",
				Address:     "Address",
				PhoneNumber: "+79999999999",
				Email:       "test@email.com",
				Password:    "hashed_password",
			})
			require.NoError(t, err)

			receivedUser, err := userRepository.GetUserByEmail(createdUser.Email)
			test.CheckOutput(t, createdUser, receivedUser, err)

		})
	}
}

var testUserRepositoryUpdateSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdUser *models.User, updatedUser *models.User, err error)
}{
	{
		TestName: "update success test",
		CheckOutput: func(t *testing.T, createdUser *models.User, updatedUser *models.User, err error) {
			require.NoError(t, err)
			require.Equal(t, createdUser.ID, updatedUser.ID)
			require.NotEqual(t, createdUser.Name, updatedUser.Name)
			require.NotEqual(t, createdUser.Surname, updatedUser.Surname)
			require.NotEqual(t, createdUser.Address, updatedUser.Address)
			require.NotEqual(t, createdUser.PhoneNumber, updatedUser.PhoneNumber)
			require.NotEqual(t, createdUser.Email, updatedUser.Email)
			require.NotEqual(t, createdUser.Password, updatedUser.Password)
		},
	},
}

func TestUserRepositoryUpdate(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}
	userRepository := postgres.CreateUserRepository(&fields)

	for _, test := range testUserRepositoryUpdateSuccess {
		t.Run(test.TestName, func(t *testing.T) {
			createdUser, err := userRepository.Create(&models.User{
				Name:        "First Name",
				Surname:     "Last Name",
				Address:     "Address",
				PhoneNumber: "+79999999999",
				Email:       "test@email.com",
				Password:    "hashed_password",
			})
			require.NoError(t, err)

			updatedUser, err := userRepository.Update(&models.User{
				ID:          createdUser.ID,
				Name:        "New First Name",
				Surname:     "New Last Name",
				Address:     "New Address",
				PhoneNumber: "+79999999998",
				Email:       "new@email.com",
				Password:    "new_hashed_password",
			})

			test.CheckOutput(t, createdUser, updatedUser, err)
		})
	}
}

var testUserRepositoryDeleteSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdUser *models.User, err error)
}{
	{
		TestName: "delete success test",
		CheckOutput: func(t *testing.T, createdUser *models.User, err error) {
			require.NoError(t, err)
		},
	},
}

var testUserRepositoryDeleteFailure = []struct {
	TestName    string
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "delete non-existent user test",
		CheckOutput: func(t *testing.T, err error) {
			require.Nil(t, err)
		},
	},
}

func TestUserRepositoryDelete(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}
	userRepository := postgres.CreateUserRepository(&fields)

	for _, test := range testUserRepositoryDeleteSuccess {
		t.Run(test.TestName, func(t *testing.T) {
			createdUser, err := userRepository.Create(&models.User{
				Name:        "First Name",
				Surname:     "Last Name",
				Address:     "Address",
				PhoneNumber: "+79999999999",
				Email:       "test@email.com",
				Password:    "hashed_password",
			})
			require.NoError(t, err)

			err = userRepository.Delete(createdUser.ID)
			test.CheckOutput(t, createdUser, err)

			_, err = userRepository.GetUserByID(createdUser.ID)
			require.Error(t, err)
		})
	}

	for _, test := range testUserRepositoryDeleteFailure {
		t.Run(test.TestName, func(t *testing.T) {
			err := userRepository.Delete(uuid.New())
			test.CheckOutput(t, err)
		})
	}
}

var testUserRepositoryGetAllSuccess = []struct {
	TestName    string
	CheckOutput func(t *testing.T, createdUsers []models.User, receivedUsers []models.User, err error)
}{
	{
		TestName: "get all success test",
		CheckOutput: func(t *testing.T, createdUsers []models.User, receivedUsers []models.User, err error) {
			require.NoError(t, err)
			require.Equal(t, len(createdUsers), len(receivedUsers))
		},
	},
}

func TestUserRepositoryGetAll(t *testing.T) {
	dbContainer, db := SetupTestDatabase()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := postgres.PostgresConnection{DB: db}
	userRepository := postgres.CreateUserRepository(&fields)

	for _, test := range testUserRepositoryGetAllSuccess {
		t.Run(test.TestName, func(t *testing.T) {
			createdUsers := make([]models.User, 0)
			for i := 0; i < 10; i++ {
				user, err := userRepository.Create(&models.User{
					Name:        fmt.Sprintf("First Name %d", i+1),
					Surname:     fmt.Sprintf("Last Name %d", i+1),
					Address:     fmt.Sprintf("Address   %d", i+1),
					PhoneNumber: fmt.Sprintf("+7999999999%d", i),
					Email:       fmt.Sprintf("test%d@email.com", i),
					Password:    "hashed_password",
				})
				require.NoError(t, err)
				createdUsers = append(createdUsers, *user)
			}
			receivedUsers, err := userRepository.GetAllUsers()
			test.CheckOutput(t, createdUsers, receivedUsers, err)
		})
	}
}
