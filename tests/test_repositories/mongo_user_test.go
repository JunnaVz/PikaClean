package test_repositories

import (
	"fmt"
	"lab3/internal/models"
	"lab3/internal/repository/mongodb"
	"testing"

	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func TestMongoUserRepositoryCreate(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}
	userRepository := mongodb.CreateUserRepository(&fields)

	for _, test := range testUserRepositoryCreateSuccess {
		t.Run(test.TestName, func(t *testing.T) {
			createdUser, err := userRepository.Create(test.InputData.user)
			test.CheckOutput(t, test.InputData, createdUser, err)
		})
	}
}

func TestMongoUserRepositoryGetByID(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}
	userRepository := mongodb.CreateUserRepository(&fields)

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

func TestMongoUserRepositoryGetByEmail(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}
	userRepository := mongodb.CreateUserRepository(&fields)

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

func TestMongoUserRepositoryUpdate(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}
	userRepository := mongodb.CreateUserRepository(&fields)

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

func TestMongoUserRepositoryDelete(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}
	userRepository := mongodb.CreateUserRepository(&fields)

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

func TestMongoUserRepositoryGetAll(t *testing.T) {
	dbContainer, db := SetupTestDatabaseMongo()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	fields := mongodb.MongoConnection{DB: db}
	userRepository := mongodb.CreateUserRepository(&fields)

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
