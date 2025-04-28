package server

import (
	"github.com/gin-contrib/sessions/cookie"
	"html/template"
	"teamdev/internal/registry"
	"teamdev/middleware"
	"teamdev/utils"

	"github.com/gin-contrib/sessions"
	//"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Services struct {
	Services *registry.Services
}

func RunServer(app *registry.App) error {
	s := Services{
		app.Services,
	}

	router := s.setupRouter(app)

	gin.SetMode(gin.DebugMode)

	//port := app.Config.Port
	//address := app.Config.Address
	//port := os.Getenv("PORT")
	//address := os.Getenv("ADDRESS")
	err := router.Run("127.0.0.1" + ":" + "8080")
	return err
}

func (s *Services) setupRouter(app *registry.App) *gin.Engine {
	authMiddleware := middleware.NewMiddleware(*app)

	router := gin.Default()

	router.SetFuncMap(template.FuncMap{
		"formatDate":    utils.FormatDate,
		"displayStatus": utils.DisplayStatus,
	})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	//store := sessions.NewCookieStore([]byte("secret"))
	//router.Use(sessions.Sessions("mysession", store))

	router.LoadHTMLGlob("templates/**/*")

	router.GET("/", s.index)
	router.GET("/prices", s.priceList)

	authGroup := router.Group("/auth")
	{
		authGroup.GET("/signup", s.signupGet)
		authGroup.POST("/signup", s.signupPost)

		authGroup.GET("/signin", s.signinGet)
		authGroup.POST("/signin", s.signinPost)

		authGroup.GET("/logout", s.logout)
	}

	workerAuthGroup := router.Group("/worker-auth")
	{
		workerAuthGroup.GET("/signin", s.workerSigninGet)
		workerAuthGroup.POST("/signin", s.workerSigninPost)
	}

	usersGroup := router.Group("/users")
	usersGroup.Use(authMiddleware.AuthMiddleware())
	{
		usersGroup.GET("/profile", s.profile)
		usersGroup.GET("/change-password", s.changePasswordGet)
		usersGroup.POST("/change-password", s.changePasswordPost)

		usersGroup.GET("/edit-profile", s.editProfileGet)
		usersGroup.POST("/edit-profile", s.editProfilePost)
	}

	userOrderGroup := usersGroup.Group("/orders")
	userOrderGroup.Use(authMiddleware.AuthMiddleware())
	{
		userOrderGroup.GET("/create", s.createOrderGet)
		userOrderGroup.POST("/create", s.createOrderPost)

		userOrderGroup.GET("/in-progress", s.inProgressOrders)
		userOrderGroup.GET("/completed", s.completedOrders)
		userOrderGroup.GET("/:id", s.orderGet)
		userOrderGroup.POST("/:id/rate", s.rateOrderApiPost)
		userOrderGroup.POST("/:id/cancel", s.cancelOrderApiPost)
	}

	workerGroup := router.Group("/worker")
	workerGroup.Use(authMiddleware.WorkerMiddleware())
	{
		workerGroup.GET("/", s.dashboard)
		workerGroup.GET("/profile", s.workerProfile)
		workerGroup.GET("/directory", s.workersDirectory)
		workerGroup.GET("/create", s.createWorkerGet)
		workerGroup.POST("/create", s.createWorkerPost)
		workerGroup.GET("/:id", s.workerDetails)
		workerGroup.GET("/orders/history", s.ordersHistory)
		workerGroup.GET("/orders/:id", s.orderDetails)
		workerGroup.POST("/orders/:id/status", s.changeStatusOrderApiPost)
		workerGroup.POST("/orders/:id/worker", s.changeWorkerApiPost)
		workerGroup.GET("/:id/edit", s.editWorkerGet)
		workerGroup.POST("/:id/edit", s.editWorkerPost)
		workerGroup.GET("/change-password", s.changeWorkerPasswordGet)
		workerGroup.POST("/change-password", s.changeWorkerPasswordPost)
	}

	servicesGroup := router.Group("/services")
	servicesGroup.Use(authMiddleware.WorkerMiddleware())
	{
		servicesGroup.GET("/", s.services)
		servicesGroup.GET("/create", s.createServiceGet)
		servicesGroup.POST("/create", s.createServicePost)
		servicesGroup.GET("/:id", s.editServiceGet)
		servicesGroup.POST("/:id", s.editServicePost)
	}

	categoriesGroup := router.Group("/categories")
	categoriesGroup.Use(authMiddleware.WorkerMiddleware())
	{
		categoriesGroup.GET("/create", s.createCategoryGet)
		categoriesGroup.POST("/create", s.createCategoryPost)
		categoriesGroup.GET("/:id", s.editCategoryGet)
		categoriesGroup.POST("/:id", s.editCategoryPost)
	}

	return router
}
