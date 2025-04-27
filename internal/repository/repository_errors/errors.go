// Package repository_errors defines standardized errors for database operations.
// These errors are returned by repository implementations to provide consistent
// error handling across the application when interacting with data storage.
package repository_errors

import "errors"

var (
	// InsertError is returned when a database insert operation fails.
	// This may occur due to constraint violations, connection issues, or other database errors.
	InsertError = errors.New("DB ERROR: Insert operation was not successful")

	// DeleteError is returned when a database delete operation fails.
	// This may occur due to referential integrity constraints or other database errors.
	DeleteError = errors.New("DB ERROR: Delete operation was not successful")

	// SelectError is returned when a database select operation fails.
	// This typically indicates an error in the query or connection problems.
	SelectError = errors.New("DB ERROR: Select operation was not successful")

	// UpdateError is returned when a database update operation fails.
	// This may occur due to constraint violations, concurrent modifications, or other database errors.
	UpdateError = errors.New("DB ERROR: Update operation was not successful")

	// DoesNotExist is returned when a requested entity cannot be found in the database.
	// This indicates that the specified row or record does not exist in the data store.
	DoesNotExist = errors.New("GET operation has failed. Such row does not exist")

	// TransactionBeginError is returned when starting a database transaction fails.
	// This typically indicates database connection or permission issues.
	TransactionBeginError = errors.New("DB ERROR: Transaction begin error")

	// TransactionRollbackError is returned when rolling back a database transaction fails.
	// This is a serious error as it may leave the database in an inconsistent state.
	TransactionRollbackError = errors.New("DB ERROR: Transaction rollback error")

	// TransactionCommitError is returned when committing a database transaction fails.
	// This may occur due to constraint violations discovered during commit or connection issues.
	TransactionCommitError = errors.New("DB ERROR: Transaction commit error")

	// ConnectionError is returned when establishing a database connection fails.
	// This may be due to incorrect credentials, network issues, or database server problems.
	ConnectionError = errors.New("DB ERROR: Connection error")
)
