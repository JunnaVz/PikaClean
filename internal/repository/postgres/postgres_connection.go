package postgres

import (
	"database/sql"
	"lab3/config"
	"lab3/internal/repository/repository_errors"
	"lab3/internal/repository/repository_interfaces"

	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
)

type PostgresConnection struct {
	DB     *sql.DB
	Config config.Config
}

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

func CreateUserRepository(fields *PostgresConnection) repository_interfaces.IUserRepository {
	dbx := sqlx.NewDb(fields.DB, "pgx")

	return NewUserRepository(dbx)
}

func CreateWorkerRepository(fields *PostgresConnection) repository_interfaces.IWorkerRepository {
	dbx := sqlx.NewDb(fields.DB, "pgx")

	return NewWorkerRepository(dbx)
}

func CreateOrderRepository(fields *PostgresConnection) repository_interfaces.IOrderRepository {
	dbx := sqlx.NewDb(fields.DB, "pgx")

	return NewOrderRepository(dbx)
}

func CreateTaskRepository(fields *PostgresConnection) repository_interfaces.ITaskRepository {
	dbx := sqlx.NewDb(fields.DB, "pgx")

	return NewTaskRepository(dbx)
}

func CreateCategoryRepository(fields *PostgresConnection) repository_interfaces.ICategoryRepository {
	dbx := sqlx.NewDb(fields.DB, "pgx")

	return NewCategoryRepository(dbx)
}
