// Package postgres provides repository implementations for data persistence
// using PostgreSQL database. It includes repositories for managing workers,
// users, tasks, orders, and categories.
package postgres

import (
	"database/sql"
	"teamdev/config"
	"teamdev/internal/repository/repository_errors"
	"teamdev/internal/repository/repository_interfaces"

	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
)

// PostgresConnection encapsulates a PostgreSQL database connection
// and associated configuration. It provides the foundation for
// creating and accessing PostgreSQL repositories.
type PostgresConnection struct {
	DB     *sql.DB       // Raw database connection to PostgreSQL
	Config config.Config // Application configuration parameters
}

// NewPostgresConnection creates a new PostgreSQL connection using the provided configuration.
// It establishes a connection to the database and validates that the connection works.
//
// Parameters:
//   - Postgres: Database connection parameters (host, port, credentials, etc.)
//   - logger: Logger for recording connection events
//
// Returns:
//   - *PostgresConnection: Initialized connection object if successful
//   - error: repository_errors.ConnectionError if connection fails
func NewPostgresConnection(Postgres config.DbConnectionFlags, logger *log.Logger) (*PostgresConnection, error) {
	fields := new(PostgresConnection)
	var err error

	fields.Config.DBFlags = Postgres

	fields.DB, err = fields.Config.DBFlags.InitPostgresDB(logger)
	if err != nil {
		logger.Error("POSTGRES! Error parse config for postgreSQL")
		return nil, repository_errors.ConnectionError
	}

	logger.Info("POSTGRES! Successfully create postgres repository fields")

	return fields, nil
}

// CreateUserRepository constructs a new UserRepository with the connection.
// This factory method provides a properly initialized repository implementation
// that satisfies the IUserRepository interface.
//
// Parameters:
//   - fields: Initialized PostgreSQL connection
//
// Returns:
//   - repository_interfaces.IUserRepository: Ready-to-use user repository
func CreateUserRepository(fields *PostgresConnection) repository_interfaces.IUserRepository {
	dbx := sqlx.NewDb(fields.DB, "pgx")

	return NewUserRepository(dbx)
}

// CreateWorkerRepository constructs a new WorkerRepository with the connection.
// This factory method provides a properly initialized repository implementation
// that satisfies the IWorkerRepository interface.
//
// Parameters:
//   - fields: Initialized PostgreSQL connection
//
// Returns:
//   - repository_interfaces.IWorkerRepository: Ready-to-use worker repository
func CreateWorkerRepository(fields *PostgresConnection) repository_interfaces.IWorkerRepository {
	dbx := sqlx.NewDb(fields.DB, "pgx")

	return NewWorkerRepository(dbx)
}

// CreateOrderRepository constructs a new OrderRepository with the connection.
// This factory method provides a properly initialized repository implementation
// that satisfies the IOrderRepository interface.
//
// Parameters:
//   - fields: Initialized PostgreSQL connection
//
// Returns:
//   - repository_interfaces.IOrderRepository: Ready-to-use order repository
func CreateOrderRepository(fields *PostgresConnection) repository_interfaces.IOrderRepository {
	dbx := sqlx.NewDb(fields.DB, "pgx")

	return NewOrderRepository(dbx)
}

// CreateTaskRepository constructs a new TaskRepository with the connection.
// This factory method provides a properly initialized repository implementation
// that satisfies the ITaskRepository interface.
//
// Parameters:
//   - fields: Initialized PostgreSQL connection
//
// Returns:
//   - repository_interfaces.ITaskRepository: Ready-to-use task repository
func CreateTaskRepository(fields *PostgresConnection) repository_interfaces.ITaskRepository {
	dbx := sqlx.NewDb(fields.DB, "pgx")

	return NewTaskRepository(dbx)
}

// CreateCategoryRepository constructs a new CategoryRepository with the connection.
// This factory method provides a properly initialized repository implementation
// that satisfies the ICategoryRepository interface.
//
// Parameters:
//   - fields: Initialized PostgreSQL connection
//
// Returns:
//   - repository_interfaces.ICategoryRepository: Ready-to-use category repository
func CreateCategoryRepository(fields *PostgresConnection) repository_interfaces.ICategoryRepository {
	dbx := sqlx.NewDb(fields.DB, "pgx")

	return NewCategoryRepository(dbx)
}
