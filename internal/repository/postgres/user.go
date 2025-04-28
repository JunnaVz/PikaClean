// Package postgres provides repository implementations for data persistence
// using PostgreSQL database. It includes repositories for managing workers,
// users, tasks, orders, and categories.
package postgres

import (
	"database/sql"
	"errors"
	"teamdev/internal/models"
	"teamdev/internal/repository/repository_errors"
	"teamdev/internal/repository/repository_interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// UserDB represents a user entity as stored in the PostgreSQL database.
// It maps directly to the columns in the users table.
type UserDB struct {
	ID          uuid.UUID `db:"id"`           // Unique identifier for the user
	Name        string    `db:"name"`         // First name of the user
	Surname     string    `db:"surname"`      // Last name of the user
	Address     string    `db:"address"`      // Physical address of the user
	PhoneNumber string    `db:"phone_number"` // Contact phone number
	Email       string    `db:"email"`        // Email address, used as username for login
	Password    string    `db:"password"`     // Hashed password for authentication
}

// UserRepository implements the IUserRepository interface for PostgreSQL.
// It provides methods for creating, updating, and retrieving user records.
type UserRepository struct {
	db *sqlx.DB // Database connection
}

// NewUserRepository creates a new UserRepository instance with the provided
// database connection.
//
// Parameters:
//   - db: An initialized sqlx.DB connection to PostgreSQL
//
// Returns:
//   - repository_interfaces.IUserRepository: Repository implementation
func NewUserRepository(db *sqlx.DB) repository_interfaces.IUserRepository {
	return &UserRepository{db: db}
}

// copyUserResultToModel converts a UserDB database entity to a models.User domain entity.
//
// Parameters:
//   - userDB: Database entity to convert
//
// Returns:
//   - *models.User: Corresponding domain entity
func copyUserResultToModel(userDB *UserDB) *models.User {
	return &models.User{
		ID:          userDB.ID,
		Name:        userDB.Name,
		Surname:     userDB.Surname,
		Address:     userDB.Address,
		PhoneNumber: userDB.PhoneNumber,
		Email:       userDB.Email,
		Password:    userDB.Password,
	}
}

// Create inserts a new user record into the database.
//
// Parameters:
//   - user: User entity to be created
//
// Returns:
//   - *models.User: Created user with assigned ID
//   - error: repository_errors.InsertError if the operation fails
func (u UserRepository) Create(user *models.User) (*models.User, error) {
	query := `INSERT INTO users(name, surname, address, phone_number, email, password) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

	var userID uuid.UUID
	err := u.db.QueryRow(query, user.Name, user.Surname, user.Address, user.PhoneNumber, user.Email, user.Password).Scan(&userID)

	if err != nil {
		return nil, repository_errors.InsertError
	}

	return &models.User{
		ID:          userID,
		Name:        user.Name,
		Surname:     user.Surname,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Password:    user.Password,
	}, nil
}

// Delete removes a user record and all associated orders from the database.
// This operation is performed as a transaction to maintain data integrity.
//
// Parameters:
//   - id: UUID of the user to delete
//
// Returns:
//   - error: repository_errors.TransactionBeginError, repository_errors.TransactionRollbackError,
//     repository_errors.DeleteError, or repository_errors.TransactionCommitError if the operation fails
func (u UserRepository) Delete(id uuid.UUID) error {
	// Start a new transaction
	tx, err := u.db.Begin()
	if err != nil {
		return repository_errors.TransactionBeginError
	}

	// Delete the records in the orders table that reference the user
	_, err = tx.Exec(`DELETE FROM orders WHERE user_id = $1;`, id)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return repository_errors.TransactionRollbackError
		}
		return repository_errors.DeleteError
	}

	// Delete the user
	_, err = tx.Exec(`DELETE FROM users WHERE id = $1;`, id)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return repository_errors.TransactionRollbackError
		}
		return repository_errors.DeleteError
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return repository_errors.TransactionCommitError
	}

	return nil
}

// Update modifies an existing user record in the database.
//
// Parameters:
//   - user: User entity with updated values
//
// Returns:
//   - *models.User: Updated user after the operation
//   - error: repository_errors.UpdateError if the operation fails
func (u UserRepository) Update(user *models.User) (*models.User, error) {
	query := `UPDATE users SET name = $1, surname = $2, email = $3, phone_number = $4, address = $5, password = $6 WHERE users.id = $7 RETURNING id, name, surname, address, phone_number, email, password;`

	var updatedUser models.User
	err := u.db.QueryRow(query, user.Name, user.Surname, user.Email, user.PhoneNumber, user.Address, user.Password, user.ID).Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Surname, &updatedUser.Address, &updatedUser.PhoneNumber, &updatedUser.Email, &updatedUser.Password)
	if err != nil {
		return nil, repository_errors.UpdateError
	}
	return &updatedUser, nil
}

// GetUserByID retrieves a user by their unique identifier.
//
// Parameters:
//   - id: UUID of the user to retrieve
//
// Returns:
//   - *models.User: Retrieved user entity
//   - error: repository_errors.DoesNotExist if no user found,
//     repository_errors.SelectError for other failures
func (u UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	query := `SELECT * FROM users WHERE id = $1;`
	userDB := &UserDB{}
	err := u.db.Get(userDB, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository_errors.DoesNotExist
	} else if err != nil {
		return nil, repository_errors.SelectError
	}

	userModels := copyUserResultToModel(userDB)

	return userModels, nil
}

// GetUserByEmail retrieves a user by their email address.
//
// Parameters:
//   - email: Email address to search for
//
// Returns:
//   - *models.User: Retrieved user entity
//   - error: repository_errors.DoesNotExist if no user found,
//     repository_errors.SelectError for other failures
func (u UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email = $1;`
	userDB := &UserDB{}
	err := u.db.Get(userDB, query, email)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository_errors.DoesNotExist
	} else if err != nil {
		return nil, repository_errors.SelectError
	}

	userModels := copyUserResultToModel(userDB)

	return userModels, nil
}

// GetAllUsers retrieves all users from the database.
// For security reasons, this method does not return user passwords.
//
// Returns:
//   - []models.User: Slice of all user entities
//   - error: repository_errors.SelectError if the operation fails
func (u UserRepository) GetAllUsers() ([]models.User, error) {
	query := `SELECT name, surname, address, phone_number, email FROM users;`
	var userDB []UserDB

	err := u.db.Select(&userDB, query)

	if err != nil {
		return nil, repository_errors.SelectError
	}

	var userModels []models.User
	for i := range userDB {
		user := copyUserResultToModel(&userDB[i])
		userModels = append(userModels, *user)
	}

	return userModels, nil
}
