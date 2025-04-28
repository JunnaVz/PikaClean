package server

import (
	"net/http"
	"strconv"
	"teamdev/internal/models"

	"github.com/gin-contrib/sessions"
	//"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Services) authenticatedWorker(c *gin.Context) *models.Worker {
	session := sessions.Default(c)
	sessionID := session.Get("workerID")

	if sessionID != nil {
		strWorkerID, ok := sessionID.(string)
		if ok {
			userId, err := uuid.Parse(strWorkerID)
			if err == nil {
				worker, err := s.Services.WorkerService.GetWorkerByID(userId)
				if err == nil {
					return worker
				}
			}
		}
	}
	return nil
}

func (s *Services) workerSigninGet(c *gin.Context) {
	c.HTML(200, "signin", gin.H{
		"title": "Вход для исполнителя",
	})
}

func (s *Services) workerSigninPost(c *gin.Context) {
	var data loginFormData
	if err := c.Bind(&data); err != nil {
		c.HTML(http.StatusBadRequest, "signin", gin.H{
			"title": "Вход для исполнителя",
			"error": err.Error(),
		})
		return
	}

	// try to login
	worker, err := s.Services.WorkerService.Login(data.Email, data.Password)
	if err != nil {
		c.HTML(http.StatusBadRequest, "signin", gin.H{
			"title":    "Вход для исполнителя",
			"error":    "Неверный пароль или исполнитель с таким email не существует",
			"formData": data,
		})
		return
	}

	// Set the session.
	session := sessions.Default(c)
	session.Set("workerID", worker.ID.String())
	ok := session.Save()
	if ok != nil {
		c.HTML(http.StatusBadRequest, "signin", gin.H{
			"title":    "Вход для исполнителя",
			"error":    "Не удалось сохранить сессию",
			"formData": data,
		})
		return
	}

	c.Redirect(http.StatusFound, "/worker/profile")
}

func (s *Services) workerProfile(c *gin.Context) {
	worker := s.authenticatedWorker(c)

	if worker.Role == models.ManagerRole {
		c.HTML(200, "worker-profile", gin.H{
			"title":  "Профиль менеджера",
			"worker": worker,
		})
		return
	}

	avgRate, _ := s.Services.WorkerService.GetAverageOrderRate(worker)

	c.HTML(200, "worker-profile", gin.H{
		"title":   "Профиль исполнителя",
		"worker":  worker,
		"avgRate": avgRate,
	})
}

func (s *Services) adminDashboard(worker *models.Worker) gin.H {
	var result = gin.H{
		"title":  "Панель администратора",
		"worker": worker,
	}

	ordersWithoutWorker, _ := s.Services.OrderService.Filter(map[string]string{"worker_id": "null", "status": "1,2"})
	ordersInProgress, _ := s.Services.OrderService.Filter(map[string]string{"worker_id": "not null", "status": "1,2"})

	ordersWithoutWorkerData := make([]orderData, len(ordersWithoutWorker))
	for i, o := range ordersWithoutWorker {
		user, _ := s.Services.UserService.GetUserByID(o.UserID)
		ordersWithoutWorkerData[i] = orderData{
			ID:           o.ID,
			User:         user,
			Status:       models.OrderStatuses[o.Status],
			Address:      o.Address,
			CreationDate: o.CreationDate.Format("2006-01-02 15:04:05"),
			Deadline:     o.Deadline.Format("2006-01-02"),
			Rate:         o.Rate,
		}
	}

	ordersInProgressData := make([]orderData, len(ordersInProgress))
	for i, o := range ordersInProgress {
		user, _ := s.Services.UserService.GetUserByID(o.UserID)
		ordersInProgressData[i] = orderData{
			ID:           o.ID,
			User:         user,
			Status:       models.OrderStatuses[o.Status],
			Address:      o.Address,
			CreationDate: o.CreationDate.Format("2006-01-02 15:04:05"),
			Deadline:     o.Deadline.Format("2006-01-02"),
			Rate:         o.Rate,
		}
	}

	result["ordersWithoutWorker"] = ordersWithoutWorkerData
	result["ordersInProgress"] = ordersInProgressData

	return result
}

func (s *Services) masterDashboard(worker *models.Worker) gin.H {
	var result = gin.H{
		"title":  "Панель мастера",
		"worker": worker,
	}

	params := map[string]string{
		"status":    "1,2",
		"worker_id": worker.ID.String(),
	}
	inProgressOrders, _ := s.Services.OrderService.Filter(params)

	inProgressOrdersData := make([]orderData, len(inProgressOrders))
	for i, o := range inProgressOrders {
		user, _ := s.Services.UserService.GetUserByID(o.UserID)
		inProgressOrdersData[i] = orderData{
			ID:           o.ID,
			User:         user,
			Status:       models.OrderStatuses[o.Status],
			Address:      o.Address,
			CreationDate: o.CreationDate.Format("2006-01-02 15:04:05"),
			Deadline:     o.Deadline.Format("2006-01-02"),
			Rate:         o.Rate,
		}
	}

	result["ordersInProgress"] = inProgressOrdersData
	return result
}

func (s *Services) dashboard(c *gin.Context) {
	worker := s.authenticatedWorker(c)

	if worker.Role == models.ManagerRole {
		c.HTML(200, "adminDashboard", s.adminDashboard(worker))
		return
	}

	c.HTML(200, "masterDashboard", s.masterDashboard(worker))
}

type workerData struct {
	ID          uuid.UUID
	Name        string
	Surname     string
	Address     string
	PhoneNumber string
	Email       string
	AverageRate float64
}

func (s *Services) workersDirectory(c *gin.Context) {
	worker := s.authenticatedWorker(c)

	managers, _ := s.Services.WorkerService.GetWorkersByRole(models.ManagerRole)
	workers, _ := s.Services.WorkerService.GetWorkersByRole(models.MasterRole)

	workersData := make([]workerData, len(workers))
	for i, w := range workers {
		avgRate, _ := s.Services.WorkerService.GetAverageOrderRate(&w)
		workersData[i] = workerData{
			ID:          w.ID,
			Name:        w.Name,
			Surname:     w.Surname,
			Address:     w.Address,
			PhoneNumber: w.PhoneNumber,
			Email:       w.Email,
			AverageRate: avgRate,
		}
	}

	if worker.Role == models.ManagerRole {
		c.HTML(200, "workersDirectory", gin.H{
			"title":    "Список исполнителей",
			"worker":   worker,
			"managers": managers,
			"workers":  workersData,
		})
		return
	}

	c.HTML(403, "workersDirectory", gin.H{"title": "Список исполнителей", "error": "Доступ запрещен!"})
}

func (s *Services) createWorkerGet(c *gin.Context) {
	worker := s.authenticatedWorker(c)

	if worker.Role == models.ManagerRole {
		c.HTML(200, "createWorker", gin.H{
			"title":  "Добавление исполнителя",
			"worker": worker,
		})
		return
	}

	c.HTML(403, "createWorker", gin.H{"title": "Добавление исполнителя", "error": "Доступ запрещен!"})
}

type createWorkerFormData struct {
	Name        string `form:"name" binding:"required"`
	Surname     string `form:"surname" binding:"required"`
	Address     string `form:"address" binding:"required"`
	PhoneNumber string `form:"phone_number" binding:"required"`
	Email       string `form:"email" binding:"required"`
	Password    string `form:"password" binding:"required"`
	Role        string `form:"role" binding:"required"`
}

func (s *Services) createWorkerPost(c *gin.Context) {
	worker := s.authenticatedWorker(c)

	if worker.Role != models.ManagerRole {
		c.HTML(403, "createWorker", gin.H{"title": "Добавление исполнителя", "error": "Доступ запрещен!"})
		return
	}

	var data createWorkerFormData
	if err := c.Bind(&data); err != nil {
		c.HTML(http.StatusBadRequest, "createWorker", gin.H{
			"title": "Добавление исполнителя",
			"error": err.Error(),
		})
		return
	}

	role, _ := strconv.Atoi(data.Role)

	newWorker := models.Worker{
		Name:        data.Name,
		Surname:     data.Surname,
		Address:     data.Address,
		PhoneNumber: data.PhoneNumber,
		Email:       data.Email,
		Role:        role,
		Password:    data.Password,
	}

	_, err := s.Services.WorkerService.Create(&newWorker, newWorker.Password)
	if err != nil {
		c.HTML(http.StatusBadRequest, "createWorker", gin.H{
			"title": "Добавление исполнителя",
			"error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/worker/directory")
}

type orderData struct {
	ID           uuid.UUID
	User         *models.User
	Status       string
	Address      string
	CreationDate string
	Deadline     string
	Rate         int
}

func (s *Services) workerDetails(c *gin.Context) {
	worker := s.authenticatedWorker(c)

	if worker.Role != models.ManagerRole {
		c.HTML(403, "workerDetails", gin.H{"title": "Информация об исполнителе", "error": "Доступ запрещен!"})
		return
	}

	workerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "workerDetails", gin.H{
			"title": "Информация об исполнителе",
			"error": "Неверный идентификатор исполнителя",
		})
		return
	}

	workerDetails, err := s.Services.WorkerService.GetWorkerByID(workerID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "workerDetails", gin.H{
			"title": "Информация об исполнителе",
			"error": "Исполнитель не найден",
		})
		return
	}

	params := map[string]string{
		"worker_id": workerID.String(),
		"status":    "1,2",
	}
	inProgressOrders, _ := s.Services.OrderService.Filter(params)

	inProgressOrdersData := make([]orderData, len(inProgressOrders))
	for i, o := range inProgressOrders {
		user, _ := s.Services.UserService.GetUserByID(o.UserID)
		inProgressOrdersData[i] = orderData{
			ID:           o.ID,
			User:         user,
			Status:       models.OrderStatuses[o.Status],
			Address:      o.Address,
			CreationDate: o.CreationDate.Format("2006-01-02 15:04:05"),
			Deadline:     o.Deadline.Format("2006-01-02"),
			Rate:         o.Rate,
		}
	}

	params["status"] = "3"
	completedOrders, _ := s.Services.OrderService.Filter(params)

	avgRate, _ := s.Services.WorkerService.GetAverageOrderRate(workerDetails)

	c.HTML(200, "workerDetails", gin.H{
		"worker":           worker,
		"title":            "Информация об исполнителе",
		"workerDetails":    workerDetails,
		"inProgressOrders": inProgressOrdersData,
		"completedOrders":  completedOrders,
		"avgRate":          avgRate,
	})
}

func (s *Services) ordersHistory(c *gin.Context) {
	worker := s.authenticatedWorker(c)

	params := map[string]string{
		"status": "3,4",
	}

	if worker.Role == models.MasterRole {
		params["worker_id"] = worker.ID.String()
	}

	orders, _ := s.Services.OrderService.Filter(params)

	ordersData := make([]orderData, len(orders))
	for i, o := range orders {
		user, _ := s.Services.UserService.GetUserByID(o.UserID)
		ordersData[i] = orderData{
			ID:           o.ID,
			User:         user,
			Status:       models.OrderStatuses[o.Status],
			Address:      o.Address,
			CreationDate: o.CreationDate.Format("2006-01-02 15:04:05"),
			Deadline:     o.Deadline.Format("2006-01-02"),
			Rate:         o.Rate,
		}
	}

	c.HTML(200, "ordersHistory", gin.H{
		"title":  "История заказов",
		"worker": worker,
		"orders": ordersData,
	})
}

func (s *Services) orderDetails(c *gin.Context) {
	worker := s.authenticatedWorker(c)

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "changeStatus", gin.H{
			"title":  "Информация о заказе",
			"error":  "Неверный идентификатор заказа",
			"worker": worker,
		})
		return
	}

	order, err := s.Services.OrderService.GetOrderByID(orderID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "changeStatus", gin.H{
			"title":  "Информация о заказе",
			"error":  "Заказ не найден",
			"worker": worker,
		})
		return
	}

	if worker.Role != models.ManagerRole && order.WorkerID != worker.ID {
		c.HTML(403, "changeStatus", gin.H{"title": "Информация о заказе", "error": "Доступ запрещен!", "worker": worker})
		return
	}

	user, _ := s.Services.UserService.GetUserByID(order.UserID)
	tasks, _ := s.Services.OrderService.GetTasksInOrder(orderID)
	orderedTasks := make([]models.OrderedTask, 0)
	for _, task := range tasks {
		quantity, _ := s.Services.OrderService.GetTaskQuantity(order.ID, task.ID)
		orderedTasks = append(orderedTasks, models.OrderedTask{
			Task:     &task,
			Quantity: quantity,
		})

	}

	if worker.Role == models.MasterRole {
		c.HTML(200, "changeStatus", gin.H{
			"title":  "Информация о заказе",
			"worker": worker,
			"order":  order,
			"user":   user,
			"tasks":  orderedTasks,
		})
		return
	} else if worker.Role == models.ManagerRole {
		workers, _ := s.Services.WorkerService.GetWorkersByRole(models.MasterRole)
		c.HTML(200, "changeStatus", gin.H{
			"title":         "Информация о заказе",
			"worker":        worker,
			"order":         order,
			"user":          user,
			"workersSelect": workers,
			"tasks":         orderedTasks,
		})
		return
	}

	c.HTML(200, "orderDetails", gin.H{
		"title":  "Информация о заказе",
		"worker": worker,
		"order":  order,
		"user":   user,
	})
}

type editWorkerData struct {
	Name        string `form:"name"`
	Surname     string `form:"surname"`
	Email       string `form:"email"`
	Address     string `form:"address"`
	PhoneNumber string `form:"phone_number"`
	Role        int    `form:"role"`
}

func (s *Services) editWorkerGet(c *gin.Context) {
	authWorker := s.authenticatedWorker(c)
	workerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.HTML(400, "editWorker", gin.H{
			"title":  "Редактировать профиль",
			"worker": authWorker,
			"error":  "Неверный идентификатор исполнителя",
		})
		return
	}

	editedWorker, err := s.Services.WorkerService.GetWorkerByID(workerID)
	if err != nil {
		c.HTML(400, "editWorker", gin.H{
			"title":  "Редактировать профиль",
			"worker": authWorker,
			"error":  "Исполнитель не найден",
		})
		return
	}

	c.HTML(200, "editWorker", gin.H{
		"title":  "Редактировать профиль",
		"worker": authWorker,
		"formData": editWorkerData{
			Name:        editedWorker.Name,
			Surname:     editedWorker.Surname,
			Email:       editedWorker.Email,
			Address:     editedWorker.Address,
			PhoneNumber: editedWorker.PhoneNumber,
			Role:        editedWorker.Role,
		},
	})
}

func (s *Services) editWorkerPost(c *gin.Context) {
	authWorker := s.authenticatedWorker(c)
	workerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.HTML(400, "editWorker", gin.H{
			"title":  "Редактировать профиль",
			"worker": authWorker,
			"error":  "Неверный идентификатор исполнителя",
		})
		return
	}

	editedWorker, err := s.Services.WorkerService.GetWorkerByID(workerID)
	if err != nil {
		c.HTML(400, "editWorker", gin.H{
			"title":  "Редактировать профиль",
			"worker": authWorker,
			"error":  "Исполнитель не найден",
		})
		return
	}

	var data editWorkerData
	err = c.Bind(&data)
	if err != nil {
		c.HTML(400, "editWorker", gin.H{
			"title":    "Редактировать профиль",
			"worker":   authWorker,
			"error":    err.Error(),
			"formData": data,
		})
		return
	}

	_, updateErr := s.Services.WorkerService.Update(
		editedWorker.ID,
		data.Name,
		data.Surname,
		data.Email,
		data.Address,
		data.PhoneNumber,
		data.Role,
		editedWorker.Password,
	)

	if updateErr != nil {
		c.HTML(400, "editWorker", gin.H{
			"title":  "Редактировать профиль",
			"worker": authWorker,
			"formData": editWorkerData{
				Name:        data.Name,
				Surname:     data.Surname,
				Email:       data.Email,
				Address:     data.Address,
				PhoneNumber: data.PhoneNumber,
				Role:        data.Role,
			},
			"error": updateErr.Error(),
		})
		return
	}

	c.Redirect(302, "/worker/"+workerID.String())
}

func (s *Services) changeWorkerPasswordGet(c *gin.Context) {
	worker := s.authenticatedWorker(c)

	c.HTML(200, "changePassword", gin.H{
		"title":  "Изменить пароль",
		"worker": worker,
	})
}

func (s *Services) changeWorkerPasswordPost(c *gin.Context) {
	authWorker := s.authenticatedWorker(c)

	var data changePasswordData
	if err := c.Bind(&data); err != nil {
		c.HTML(400, "changePassword", gin.H{
			"title":  "Изменить пароль",
			"worker": authWorker,
			"error":  err.Error(),
		})
		return
	}

	if data.NewPassword != data.NewPassword2 {
		c.HTML(400, "changePassword", gin.H{
			"title": "Изменить пароль",
			"auth":  authWorker,
			"error": "Новые пароли не совпадают",
		})
		return
	}

	worker, err := s.Services.WorkerService.Login(authWorker.Email, data.OldPassword)

	if err != nil {
		c.HTML(400, "changePassword", gin.H{
			"title": "Изменить пароль",
			"auth":  authWorker,
			"error": "Старый пароль неверен",
		})
		return
	}

	updatedWorker, updateErr := s.Services.WorkerService.Update(
		worker.ID,
		worker.Name,
		worker.Surname,
		worker.Email,
		worker.Address,
		worker.PhoneNumber,
		worker.Role,
		data.NewPassword,
	)

	if updateErr != nil {
		c.HTML(400, "changePassword", gin.H{
			"title": "Изменить пароль",
			"auth":  updatedWorker,
			"error": updateErr.Error(),
		})
		return
	}

	c.Redirect(302, "/worker/profile")
}
