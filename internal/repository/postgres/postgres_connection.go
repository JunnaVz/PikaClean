package postgres

import (
	"PikaClean/config"
	"PikaClean/internal/repository/repository_errors"
	"PikaClean/internal/repository/repository_interfaces"
	"database/sql"

	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
)

// PostgresConnection - структура, представляющая соединение с PostgreSQL
type PostgresConnection struct {
	DB     *sql.DB       // Соединение с базой данных
	Config config.Config // Конфигурация для базы данных
}

// NewPostgresConnection - конструктор для создания нового соединения с PostgreSQL
// Используется флаг конфигурации Postgres для инициализации соединения
func NewPostgresConnection(Postgres config.DbConnectionFlags, logger *log.Logger) (*PostgresConnection, error) {
	fields := new(PostgresConnection)
	var err error

	// Устанавливаем флаги конфигурации для соединения с базой данных
	fields.Config.DBFlags = Postgres

	// Инициализируем подключение к базе данных PostgreSQL
	fields.DB, err = fields.Config.DBFlags.InitPostgresDB(logger)
	if err != nil {
		// Логируем ошибку, если соединение не удалось
		logger.Error("POSTGRES! Error parse config for postgreSQL")
		return nil, repository_errors.ConnectionError // Возвращаем ошибку подключения
	}

	// Логируем успешное создание репозитория PostgreSQL
	logger.Info("POSTGRES! Successfully create postgres repository fields")

	// Возвращаем объект соединения с PostgreSQL
	return fields, nil
}

// CreateUserRepository - функция для создания репозитория пользователей для PostgreSQL
func CreateUserRepository(fields *PostgresConnection) repository_interfaces.IUserRepository {
	// Создаем объект sqlx.Db для работы с PostgreSQL
	dbx := sqlx.NewDb(fields.DB, "pgx")

	// Возвращаем новый репозиторий пользователей
	return NewUserRepository(dbx)
}

// CreateWorkerRepository - функция для создания репозитория рабочих для PostgreSQL
func CreateWorkerRepository(fields *PostgresConnection) repository_interfaces.IWorkerRepository {
	// Создаем объект sqlx.Db для работы с PostgreSQL
	dbx := sqlx.NewDb(fields.DB, "pgx")

	// Возвращаем новый репозиторий рабочих
	return NewWorkerRepository(dbx)
}

// CreateOrderRepository - функция для создания репозитория заказов для PostgreSQL
func CreateOrderRepository(fields *PostgresConnection) repository_interfaces.IOrderRepository {
	// Создаем объект sqlx.Db для работы с PostgreSQL
	dbx := sqlx.NewDb(fields.DB, "pgx")

	// Возвращаем новый репозиторий заказов
	return NewOrderRepository(dbx)
}

// CreateTaskRepository - функция для создания репозитория задач для PostgreSQL
func CreateTaskRepository(fields *PostgresConnection) repository_interfaces.ITaskRepository {
	// Создаем объект sqlx.Db для работы с PostgreSQL
	dbx := sqlx.NewDb(fields.DB, "pgx")

	// Возвращаем новый репозиторий задач
	return NewTaskRepository(dbx)
}

// CreateCategoryRepository - функция для создания репозитория категорий для PostgreSQL
func CreateCategoryRepository(fields *PostgresConnection) repository_interfaces.ICategoryRepository {
	// Создаем объект sqlx.Db для работы с PostgreSQL
	dbx := sqlx.NewDb(fields.DB, "pgx")

	// Возвращаем новый репозиторий категорий
	return NewCategoryRepository(dbx)
}
