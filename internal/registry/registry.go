// Package registry provides application initialization and dependency management
// for the PikaClean application. It creates and wires together all repositories
// and services needed by the application.
package registry

import (
	"os"
	"teamdev/config"
	"teamdev/internal/repository/postgres"
	"teamdev/internal/repository/repository_interfaces"
	services "teamdev/internal/services"
	"teamdev/internal/services/service_interfaces"
	"teamdev/password_hash"

	"github.com/charmbracelet/log"
)

// Services encapsulates all business logic services used by the application.
// It provides a central point of access to functionality like user management,
// worker management, task management, order processing, and category management.
type Services struct {
	UserService     service_interfaces.IUserService     // Handles user-related business logic
	WorkerService   service_interfaces.IWorkerService   // Handles worker-related business logic
	TaskService     service_interfaces.ITaskService     // Handles cleaning task-related business logic
	OrderService    service_interfaces.IOrderService    // Handles order processing business logic
	CategoryService service_interfaces.ICategoryService // Handles category management business logic
}

// Repositories encapsulates all data access objects used by the application.
// It provides structured access to the underlying data storage systems.
type Repositories struct {
	UserRepository     repository_interfaces.IUserRepository     // Handles user data persistence
	WorkerRepository   repository_interfaces.IWorkerRepository   // Handles worker data persistence
	TaskRepository     repository_interfaces.ITaskRepository     // Handles cleaning task data persistence
	OrderRepository    repository_interfaces.IOrderRepository    // Handles order data persistence
	CategoryRepository repository_interfaces.ICategoryRepository // Handles category data persistence
}

// App is the main application container that holds configuration,
// repositories, services, and logging facilities.
type App struct {
	Config       config.Config // Application configuration parameters
	Repositories *Repositories // Data access layer
	Services     *Services     // Business logic layer
	Logger       *log.Logger   // Application logging facility
}

// postgresRepositoriesInitialization creates and initializes all PostgreSQL-based repositories.
// It sets up database connections and creates repository instances for each domain entity.
func (a *App) postgresRepositoriesInitialization(fields *postgres.PostgresConnection) *Repositories {
	r := &Repositories{
		UserRepository:     postgres.CreateUserRepository(fields),
		WorkerRepository:   postgres.CreateWorkerRepository(fields),
		TaskRepository:     postgres.CreateTaskRepository(fields),
		OrderRepository:    postgres.CreateOrderRepository(fields),
		CategoryRepository: postgres.CreateCategoryRepository(fields),
	}
	a.Logger.Info("Success initialization of repositories")
	return r
}

// servicesInitialization creates and initializes all business logic services.
// It connects services with their required repositories and utilities.
func (a *App) servicesInitialization(r *Repositories) *Services {
	passwordHash := password_hash.NewPasswordHash()

	s := &Services{
		UserService:     services.NewUserService(r.UserRepository, passwordHash, a.Logger),
		WorkerService:   services.NewWorkerService(r.WorkerRepository, passwordHash, a.Logger),
		OrderService:    services.NewOrderService(r.OrderRepository, r.WorkerRepository, r.TaskRepository, r.UserRepository, a.Logger),
		TaskService:     services.NewTaskService(r.TaskRepository, a.Logger),
		CategoryService: services.NewCategoryService(r.CategoryRepository, r.TaskRepository, a.Logger),
	}
	a.Logger.Info("Success initialization of services")

	return s
}

// initLogger configures the application's logging system based on configuration settings.
// It creates log files, sets appropriate log levels, and configures log formatting.
func (a *App) initLogger() {
	f, err := os.OpenFile(a.Config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	Logger := log.New(f)

	log.SetFormatter(log.LogfmtFormatter)
	Logger.SetReportTimestamp(true)
	Logger.SetReportCaller(true)

	if a.Config.LogLevel == "debug" {
		Logger.SetLevel(log.DebugLevel)
	} else if a.Config.LogLevel == "info" {
		Logger.SetLevel(log.InfoLevel)
	} else {
		log.Fatal("Error log level")
	}

	Logger.Info("Success initialization of new Logger!")

	a.Logger = Logger
}

// Init initializes the entire application stack including logger, repositories and services.
// It prepares the application for running by establishing all necessary connections and dependencies.
// Returns an error if initialization fails.
func (a *App) Init() error {
	a.initLogger()

	if a.Config.DBType == "postgres" {
		fields, err := postgres.NewPostgresConnection(a.Config.DBFlags, a.Logger)
		if err != nil {
			a.Logger.Fatal("Error create postgres repository fields", "err", err)
			return err
		}

		a.Repositories = a.postgresRepositoriesInitialization(fields)
		a.Services = a.servicesInitialization(a.Repositories)
	}

	return nil
}

// Run executes the application initialization sequence and prepares it for operation.
// Returns an error if the application fails to initialize properly.
func (a *App) Run() error {
	err := a.Init()

	if err != nil {
		a.Logger.Error("Error init app", "err", err)
		return err
	}

	return nil
}
