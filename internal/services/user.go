// Package interfaces provides service layer implementations for the PikaClean application,
// handling the business logic between API controllers and data repositories.
// This package contains concrete implementations of the service interfaces defined
// in the service_interfaces package, with implementations for tasks, orders, users,
// workers, and categories.
package interfaces

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"teamdev/internal/models"
	"teamdev/internal/repository/repository_errors"
	"teamdev/internal/repository/repository_interfaces"
	"teamdev/internal/services/service_interfaces"
	"teamdev/password_hash"
)

// UserService implements the IUserService interface and provides
// business logic for user account management in the application.
// It handles user registration, authentication, data retrieval and updates.
type UserService struct {
	UserRepository repository_interfaces.IUserRepository // Repository for persistent user operations
	hash           password_hash.PasswordHash            // Utility for password hashing and verification
	logger         *log.Logger                           // Logger for recording service activity
}

// NewUserService creates a new UserService instance with the provided dependencies.
//
// Parameters:
//   - UserRepository: Repository for user data access operations
//   - hash: Password hashing utility for secure password storage
//   - logger: Logger for recording service activity and errors
//
// Returns:
//   - service_interfaces.IUserService: A fully initialized user service
func NewUserService(UserRepository repository_interfaces.IUserRepository, hash password_hash.PasswordHash, logger *log.Logger) service_interfaces.IUserService {
	return &UserService{
		UserRepository: UserRepository,
		hash:           hash,
		logger:         logger,
	}
}

// GetUserByID retrieves a specific user by their unique identifier.
//
// Parameters:
//   - id: UUID of the user to retrieve
//
// Returns:
//   - *models.User: Retrieved user entity
//   - error: Any retrieval errors
func (u UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	user, err := u.UserRepository.GetUserByID(id)

	if err != nil {
		u.logger.Error("SERVICE-REPOSITORY: GetUserByID method failed", "id", id, "error", err)
		return nil, err
	}

	u.logger.Info("SERVICE: Successfully got user with GetUserByID", "id", id)
	return user, nil
}

// GetUserByEmail retrieves a user by their email address.
//
// Parameters:
//   - email: Email address of the user to retrieve
//
// Returns:
//   - *models.User: Retrieved user entity
//   - error: Any retrieval errors
func (u UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := u.UserRepository.GetUserByEmail(email)

	if err != nil {
		u.logger.Error("SERVICE-REPOSITORY: GetUserByEmail method failed", "email", email, "error", err)
		return nil, err
	}

	u.logger.Info("SERVICE: Successfully got user with GetUserByEmail", "email", email)
	return user, nil
}

// checkIfUserWithEmailExists verifies if a user with the given email already exists.
//
// Parameters:
//   - email: Email address to check
//
// Returns:
//   - *models.User: User entity if found, nil if not found
//   - error: Any errors encountered during the check
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

// Register creates a new user account with validated information and a secure password.
//
// Parameters:
//   - user: User entity with personal information
//   - password: Plain text password to be hashed and stored
//
// Returns:
//   - *models.User: Created user with assigned ID if successful
//   - error: Validation or persistence errors if they occur
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

// Login authenticates a user with email and password.
//
// Parameters:
//   - email: User's email address
//   - password: Plain text password for verification
//
// Returns:
//   - *models.User: Authenticated user entity if successful
//   - error: Authentication errors if credentials are invalid
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

// Update modifies an existing user's information with validated data.
//
// Parameters:
//   - id: UUID of the user to update
//   - name: New first name
//   - surname: New last name
//   - email: New email address
//   - address: New physical address
//   - phoneNumber: New contact phone number
//   - password: New password (will be hashed if changed)
//
// Returns:
//   - *models.User: Updated user after changes
//   - error: Validation or persistence errors if they occur
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
