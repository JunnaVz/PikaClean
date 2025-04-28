package server

import (
	"strconv"
	"teamdev/internal/models"
	"teamdev/utils"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Services) profile(c *gin.Context) {
	c.HTML(200, "profile", gin.H{
		"title": "Ваш профиль",
		"auth":  s.authenticatedUser(c),
	})
}

type changePasswordData struct {
	OldPassword  string `form:"oldPassword"`
	NewPassword  string `form:"newPassword"`
	NewPassword2 string `form:"newPassword2"`
}

func (s *Services) changePasswordGet(c *gin.Context) {
	c.HTML(200, "changePassword", gin.H{
		"title": "Изменить пароль",
		"auth":  s.authenticatedUser(c),
	})
}

func (s *Services) changePasswordPost(c *gin.Context) {
	var data changePasswordData
	if err := c.Bind(&data); err != nil {
		c.HTML(400, "changePassword", gin.H{
			"title": "Изменить пароль",
			"auth":  s.authenticatedUser(c),
			"error": err.Error(),
		})
		return
	}

	if data.NewPassword != data.NewPassword2 {
		c.HTML(400, "changePassword", gin.H{
			"title": "Изменить пароль",
			"auth":  s.authenticatedUser(c),
			"error": "Новые пароли не совпадают",
		})
		return
	}

	authUser := s.authenticatedUser(c)

	user, err := s.Services.UserService.Login(authUser.Email, data.OldPassword)

	if err != nil {
		c.HTML(400, "changePassword", gin.H{
			"title": "Изменить пароль",
			"auth":  s.authenticatedUser(c),
			"error": "Старый пароль неверен",
		})
		return
	}

	updatedUser, updateErr := s.Services.UserService.Update(
		user.ID,
		user.Name,
		user.Surname,
		user.Email,
		user.Address,
		user.PhoneNumber,
		data.NewPassword,
	)

	if updateErr != nil {
		c.HTML(400, "changePassword", gin.H{
			"title": "Изменить пароль",
			"auth":  updatedUser,
			"error": updateErr.Error(),
		})
		return
	}

	c.Redirect(302, "/users/profile")
}

type editProfileData struct {
	Name        string `form:"InputName"`
	Surname     string `form:"InputSurname"`
	Email       string `form:"InputEmail"`
	Address     string `form:"InputAddress"`
	PhoneNumber string `form:"InputPhone"`
}

func (s *Services) editProfileGet(c *gin.Context) {
	authUser := s.authenticatedUser(c)
	c.HTML(200, "editProfile", gin.H{
		"title": "Редактировать профиль",
		"auth":  authUser,
		"formData": editProfileData{
			Name:        authUser.Name,
			Surname:     authUser.Surname,
			Email:       authUser.Email,
			Address:     authUser.Address,
			PhoneNumber: authUser.PhoneNumber,
		},
	})
}

func (s *Services) editProfilePost(c *gin.Context) {
	var data editProfileData
	if err := c.Bind(&data); err != nil {
		c.HTML(400, "editProfile", gin.H{
			"title":    "Редактировать профиль",
			"auth":     s.authenticatedUser(c),
			"error":    err.Error(),
			"formData": data,
		})
		return
	}

	authUser := s.authenticatedUser(c)

	updatedUser, updateErr := s.Services.UserService.Update(
		authUser.ID,
		data.Name,
		data.Surname,
		data.Email,
		data.Address,
		data.PhoneNumber,
		authUser.Password,
	)

	if updateErr != nil {
		c.HTML(400, "editProfile", gin.H{
			"title":    "Редактировать профиль",
			"auth":     updatedUser,
			"error":    updateErr.Error(),
			"formData": data,
		})
		return
	}

	c.Redirect(302, "/users/profile")
}

func (s *Services) createOrderGet(c *gin.Context) {
	prices := make(map[models.Category][]models.Task)
	categories, err := s.Services.CategoryService.GetAll()
	if err != nil {
		log.Printf("Error getting categories: %v", err)
		categories = []models.Category{}
	}

	for i, category := range categories {
		tasks, err := s.Services.CategoryService.GetTasksInCategory(category.ID)
		if err != nil {
			log.Printf("Error getting tasks in category %d: %v", i, err)
			continue
		}
		prices[category] = tasks
	}
	c.HTML(200, "createOrder", gin.H{
		"title":    "Создать заказ",
		"auth":     s.authenticatedUser(c),
		"prices":   prices,
		"category": categories,
	})
}

type createOrderData struct {
	SameAddress bool   `form:"sameAddress"`
	Deadline    string `form:"deadlineInput"`
	Address     string `form:"addressInput"`
	Tasks       map[string]string
	Confirmed   bool `form:"confirmed"`
}

func (s *Services) createOrderPost(c *gin.Context) {
	data := createOrderData{
		SameAddress: utils.ParseHtmlToggle(c.DefaultPostForm("sameAddress", "off")),
		Deadline:    c.PostForm("deadlineInput"),
		Address:     c.PostForm("addressInput"),
		Tasks:       c.PostFormMap("tasks"),
		Confirmed:   c.DefaultPostForm("confirmed", "false") == "true",
	}

	if data.SameAddress {
		data.Address = s.authenticatedUser(c).Address
	}

	authUser := s.authenticatedUser(c)

	var orderedTasks []models.OrderedTask
	var totalPrice float64 = 0
	for taskID, taskAmount := range data.Tasks {
		parsedID, _ := uuid.Parse(taskID)
		taskObj, err := s.Services.TaskService.GetTaskByID(parsedID)

		quantity, err := strconv.Atoi(taskAmount)
		if err == nil && quantity > 0 {
			totalPrice += taskObj.PricePerSingle * float64(quantity)
			orderedTasks = append(orderedTasks, models.OrderedTask{
				Task:     taskObj,
				Quantity: quantity,
			})
		} else {
			log.Printf("Error getting task by ID %s: %v", taskID, err)
		}
	}

	if !data.Confirmed {
		sum := float64(0)
		for _, task := range orderedTasks {
			sum += task.Task.PricePerSingle * float64(task.Quantity)
		}
		c.HTML(200, "confirmOrder", gin.H{
			"title":      "Подтвердить заказ",
			"auth":       authUser,
			"address":    data.Address,
			"deadline":   data.Deadline,
			"tasks":      orderedTasks,
			"sum":        sum,
			"totalPrice": totalPrice,
		})
		return
	}

	_, err := s.Services.OrderService.CreateOrder(
		authUser.ID,
		data.Address,
		utils.ConvertStringToTime(data.Deadline),
		orderedTasks,
	)

	if err != nil {
		c.HTML(400, "createOrder", gin.H{
			"title": "Создать заказ",
			"auth":  authUser,
			"error": err.Error(),
		})
		return
	}

	c.Redirect(302, "/users/profile")
}

type OrderItem struct {
	ID           uuid.UUID
	TotalPrice   float64
	Worker       *models.Worker
	User         *models.User
	Status       string
	Address      string
	CreationDate time.Time
	Deadline     time.Time
	Rate         int
}

func (s *Services) getOrdersList(params map[string]string) ([]OrderItem, error) {
	orders, err := s.Services.OrderService.Filter(params)
	if err != nil {
		return nil, err
	}

	ordersList := make([]OrderItem, 0)
	for _, order := range orders {
		worker, _ := s.Services.WorkerService.GetWorkerByID(order.WorkerID)
		user, _ := s.Services.UserService.GetUserByID(order.UserID)
		totalPrice, _ := s.Services.OrderService.GetTotalPrice(order.ID)
		ordersList = append(ordersList, OrderItem{
			ID:           order.ID,
			TotalPrice:   totalPrice,
			Worker:       worker,
			User:         user,
			Status:       models.OrderStatuses[order.Status],
			Address:      order.Address,
			CreationDate: order.CreationDate,
			Deadline:     order.Deadline,
			Rate:         order.Rate,
		})
	}

	return ordersList, nil
}

func (s *Services) inProgressOrders(c *gin.Context) {
	authUser := s.authenticatedUser(c)

	params := map[string]string{
		"status":  "1,2",
		"user_id": authUser.ID.String(),
	}

	orders, err := s.getOrdersList(params)

	if err != nil {
		c.HTML(500, "error", gin.H{
			"title": "Ошибка",
			"auth":  authUser,
			"error": err.Error(),
		})
		return
	}

	c.HTML(200, "ordersList", gin.H{
		"title":  "Активные заказы",
		"auth":   authUser,
		"orders": orders,
	})
}

func (s *Services) completedOrders(c *gin.Context) {
	authUser := s.authenticatedUser(c)

	params := map[string]string{
		"status":  "3,4",
		"user_id": authUser.ID.String(),
	}

	orders, err := s.getOrdersList(params)

	if err != nil {
		c.HTML(500, "error", gin.H{
			"title": "Ошибка",
			"auth":  authUser,
			"error": err.Error(),
		})
		return
	}

	c.HTML(200, "ordersList", gin.H{
		"title":  "Завершенные заказы",
		"auth":   authUser,
		"orders": orders,
	})
}

func (s *Services) orderGet(c *gin.Context) {
	authUser := s.authenticatedUser(c)
	orderID, _ := uuid.Parse(c.Param("id"))

	order, err := s.Services.OrderService.GetOrderByID(orderID)
	if err != nil || order.UserID != authUser.ID {
		c.HTML(500, "orderDetails", gin.H{
			"title": "Ошибка",
			"auth":  authUser,
			"error": "Такого заказа не существует или у вас нет прав на его просмотр",
		})
		return
	}

	worker, _ := s.Services.WorkerService.GetWorkerByID(order.WorkerID)
	totalPrice, _ := s.Services.OrderService.GetTotalPrice(order.ID)
	tasks, _ := s.Services.OrderService.GetTasksInOrder(order.ID)

	orderedTasks := make([]models.OrderedTask, 0)
	for _, task := range tasks {
		quantity, _ := s.Services.OrderService.GetTaskQuantity(order.ID, task.ID)
		orderedTasks = append(orderedTasks, models.OrderedTask{
			Task:     &task,
			Quantity: quantity,
		})

	}

	c.HTML(200, "orderDetails", gin.H{
		"title":      "Заказ",
		"auth":       authUser,
		"order":      order,
		"worker":     worker,
		"tasks":      orderedTasks,
		"totalPrice": totalPrice,
	})
}
