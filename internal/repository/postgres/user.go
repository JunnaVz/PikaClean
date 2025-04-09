package postgres

import (
	"database/sql"
	"errors"
	"lab3/internal/models"
	"lab3/internal/repository/repository_errors"
	"lab3/internal/repository/repository_interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserDB struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Surname     string    `db:"surname"`
	Address     string    `db:"address"`
	PhoneNumber string    `db:"phone_number"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository_interfaces.IUserRepository {
	return &UserRepository{db: db}
}

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

	//// Check if the order was actually deleted
	//rowsAffected, err := result.RowsAffected()
	//if err != nil {
	//	err := tx.Rollback()
	//	if err != nil {
	//		return repository_errors.TransactionRollbackError
	//	}
	//	return repository_errors.DeleteError
	//}
	//
	//if rowsAffected == 0 {
	//	return errors.New("no order found to delete")
	//}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return repository_errors.TransactionCommitError
	}

	return nil
}

func (u UserRepository) Update(user *models.User) (*models.User, error) {
	query := `UPDATE users SET name = $1, surname = $2, email = $3, phone_number = $4, address = $5, password = $6 WHERE users.id = $7 RETURNING id, name, surname, address, phone_number, email, password;`

	var updatedUser models.User
	err := u.db.QueryRow(query, user.Name, user.Surname, user.Email, user.PhoneNumber, user.Address, user.Password, user.ID).Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Surname, &updatedUser.Address, &updatedUser.PhoneNumber, &updatedUser.Email, &updatedUser.Password)
	if err != nil {
		return nil, repository_errors.UpdateError
	}
	return &updatedUser, nil
}

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
