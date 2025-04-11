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

type Services struct {
	UserService     service_interfaces.IUserService
	WorkerService   service_interfaces.IWorkerService
	TaskService     service_interfaces.ITaskService
	OrderService    service_interfaces.IOrderService
	CategoryService service_interfaces.ICategoryService
}

type Repositories struct {
	UserRepository     repository_interfaces.IUserRepository
	WorkerRepository   repository_interfaces.IWorkerRepository
	TaskRepository     repository_interfaces.ITaskRepository
	OrderRepository    repository_interfaces.IOrderRepository
	CategoryRepository repository_interfaces.ICategoryRepository
}

type App struct {
	Config       config.Config
	Repositories *Repositories
	Services     *Services
	Logger       *log.Logger
}

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

func (a *App) Run() error {
	err := a.Init()

	if err != nil {
		a.Logger.Error("Error init app", "err", err)
		return err
	}

	return nil
}
