package interfaces

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"lab3/internal/models"
	"lab3/internal/repository/repository_errors"
	"lab3/internal/repository/repository_interfaces"
	"lab3/internal/services/service_interfaces"
	"lab3/password_hash"
)

type UserService struct {
	UserRepository repository_interfaces.IUserRepository
	hash           password_hash.PasswordHash
	logger         *log.Logger
}

func NewUserService(UserRepository repository_interfaces.IUserRepository, hash password_hash.PasswordHash, logger *log.Logger) service_interfaces.IUserService {
	return &UserService{
		UserRepository: UserRepository,
		hash:           hash,
		logger:         logger,
	}
}

func (u UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	user, err := u.UserRepository.GetUserByID(id)

	if err != nil {
		u.logger.Error("SERVICE-REPOSITORY: GetUserByID method failed", "id", id, "error", err)
		return nil, err
	}

	u.logger.Info("SERVICE: Successfully got user with GetUserByID", "id", id)
	return user, nil
}

func (u UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := u.UserRepository.GetUserByEmail(email)

	if err != nil {
		u.logger.Error("SERVICE-REPOSITORY: GetUserByEmail method failed", "email", email, "error", err)
		return nil, err
	}

	u.logger.Info("SERVICE: Successfully got user with GetUserByEmail", "email", email)
	return user, nil
}

func (u UserService) checkIfUserWithEmailExists(email string) (*models.User, error) {
	u.logger.Info("SERVICE: Checking if user with email exists", "email", email)
	tempUser, err := u.UserRepository.GetUserByEmail(email)

	if err != nil && errors.Is(err, repository_errors.DoesNotExist) {
		u.logger.Info("SERVICE: User with email does not exist", "email", email)
		return nil, nil
	} else if err != nil {
		u.logger.Error("SERVICE: GetUserByEmail method failed", "email", email, "error", err)
		return nil, err
	} else {
		u.logger.Info("SERVICE: User with email exists", "email", email)
		return tempUser, nil
	}
}

func (u UserService) Register(user *models.User, password string) (*models.User, error) {
	u.logger.Infof("SERVICE: validate user with email %s", user.Email)
	if !validName(user.Name) {
		u.logger.Error("SERVICE: Invalid name")
		return nil, fmt.Errorf("SERVICE: Invalid name")
	}

	if !validName(user.Surname) {
		u.logger.Error("SERVICE: Invalid surname")
		return nil, fmt.Errorf("SERVICE: Invalid surname")
	}

	if !validEmail(user.Email) {
		u.logger.Error("SERVICE: Invalid email")
		return nil, fmt.Errorf("SERVICE: Invalid email")
	}

	if !validAddress(user.Address) {
		u.logger.Error("SERVICE: Invalid address")
		return nil, fmt.Errorf("SERVICE: Invalid address")
	}

	if !validPhoneNumber(user.PhoneNumber) {
		u.logger.Error("SERVICE: Invalid phone number")
		return nil, fmt.Errorf("SERVICE: Invalid phone number")
	}

	if !validPassword(password) {
		u.logger.Error("SERVICE: Invalid password")
		return nil, fmt.Errorf("SERVICE: Invalid password")
	}

	u.logger.Infof("SERVICE: Checking if user with email %s exists", user.Email)
	tempUser, err := u.checkIfUserWithEmailExists(user.Email)
	if err != nil {
		u.logger.Error("SERVICE: Error occurred during checking if user with email exists")
		return nil, err
	} else if tempUser != nil {
		u.logger.Info("SERVICE: User with email exists", "email", user.Email)
		return nil, fmt.Errorf("SERVICE: User with email exists")
	}

	u.logger.Infof("SERVICE: Creating new user: %s %s", user.Name, user.Surname)
	hashedPassword, err := u.hash.GetHash(password)
	if err != nil {
		u.logger.Error("SERVICE: Error occurred during password hashing")
		return nil, err
	} else {
		user.Password = hashedPassword
	}

	createdUser, err := u.UserRepository.Create(user)
	if err != nil {
		u.logger.Error("SERVICE: Create method failed", "error", err)
		return nil, err
	}

	u.logger.Info("SERVICE: Successfully created new user with ", "id", createdUser.ID)

	return createdUser, nil
}

func (u UserService) Login(email, password string) (*models.User, error) {
	u.logger.Infof("SERVICE: Checking if user with email %s exists", email)
	tempUser, err := u.checkIfUserWithEmailExists(email)
	if err != nil {
		u.logger.Error("SERVICE: Error occurred during checking if user with email exists")
		return nil, err
	} else if tempUser == nil {
		u.logger.Info("SERVICE: User with email does not exist", "email", email)
		return nil, fmt.Errorf("SERVICE: User with email does not exist")
	}

	u.logger.Infof("SERVICE: Checking if password is correct for user with email %s", email)
	isPasswordCorrect := u.hash.CompareHashAndPassword(tempUser.Password, password)
	if !isPasswordCorrect {
		u.logger.Info("SERVICE: Password is incorrect for user with email", "email", email)
		return nil, fmt.Errorf("SERVICE: Password is incorrect for user with email")
	}

	u.logger.Info("SERVICE: Successfully logged in user with email", "email", email)
	return tempUser, nil
}

func (u UserService) Update(id uuid.UUID, name string, surname string, email string, address string, phoneNumber string, password string) (*models.User, error) {
	user, err := u.UserRepository.GetUserByID(id)
	if err != nil {
		u.logger.Error("SERVICE: GetUserByID method failed", "id", id, "error", err)
		return nil, err
	}

	if !validName(name) || !validName(surname) || !validEmail(email) || !validAddress(address) || !validPhoneNumber(phoneNumber) || !validPassword(password) {
		u.logger.Error("SERVICE: Invalid input")
		return nil, fmt.Errorf("SERVICE: Invalid input")
	}

	user.Name = name
	user.Surname = surname
	user.Email = email
	user.Address = address
	user.PhoneNumber = phoneNumber

	if user.Password != password {
		hashedPassword, hashErr := u.hash.GetHash(password)
		if hashErr != nil {
			u.logger.Error("SERVICE: Error occurred during password hashing")
			return nil, hashErr
		} else {
			user.Password = hashedPassword
		}
	}

	user, err = u.UserRepository.Update(user)
	if err != nil {
		u.logger.Error("SERVICE: Update method failed", "error", err)
		return nil, err
	}

	u.logger.Info("SERVICE: Successfully updated user personal information", "user", user)
	return user, nil
}
