package server

import (
	"log"
	"net/http"
	"teamdev/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Services) services(c *gin.Context) {
	worker := s.authenticatedWorker(c)

	categories, err := s.Services.CategoryService.GetAll()
	if err != nil {
		log.Printf("Error getting categories: %v", err)
		categories = []models.Category{}
	}

	prices := make(map[models.Category][]models.Task)

	for _, category := range categories {
		tasks, err := s.Services.CategoryService.GetTasksInCategory(category.ID)
		if err != nil {
			log.Printf("Error getting tasks in category %s: %v", category.Name, err)
			continue
		}
		prices[category] = tasks
	}
	c.HTML(http.StatusOK, "servicesList", gin.H{
		"title":  "Доступные услуги",
		"worker": worker,
		"prices": prices,
	})
}

func (s *Services) createServiceGet(c *gin.Context) {
	worker := s.authenticatedWorker(c)
	categories, err := s.Services.CategoryService.GetAll()
	if err != nil {
		log.Printf("Error getting categories: %v", err)
		categories = []models.Category{}
	}

	c.HTML(http.StatusOK, "createService", gin.H{
		"title":      "Создать услугу",
		"worker":     worker,
		"categories": categories,
	})
}

type ServiceFormData struct {
	Name           string  `form:"name"`
	PricePerSingle float64 `form:"pricePerSingle"`
	Category       int     `form:"category"`
}

func (s *Services) createServicePost(c *gin.Context) {
	worker := s.authenticatedWorker(c)

	var data ServiceFormData
	if err := c.Bind(&data); err != nil {
		c.HTML(http.StatusBadRequest, "createService", gin.H{
			"worker":   worker,
			"title":    "Создать услугу",
			"error":    err.Error(),
			"formData": data,
		})
		return
	}

	_, err := s.Services.TaskService.Create(data.Name, data.PricePerSingle, data.Category)
	if err != nil {
		c.HTML(http.StatusBadRequest, "createService", gin.H{
			"worker":   worker,
			"title":    "Создать услугу",
			"error":    err.Error(),
			"formData": data,
		})
		return
	}

	c.Redirect(http.StatusFound, "/services")
}

func (s *Services) editServiceGet(c *gin.Context) {
	worker := s.authenticatedWorker(c)
	serviceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "createService", gin.H{
			"worker": worker,
			"title":  "Создать услугу",
			"error":  "Неверный идентификатор услуги",
		})
		return
	}

	service, err := s.Services.TaskService.GetTaskByID(serviceID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "createService", gin.H{
			"worker": worker,
			"title":  "Создать услугу",
			"error":  "Услуга не найдена",
		})
		return
	}

	formData := ServiceFormData{
		Name:           service.Name,
		PricePerSingle: service.PricePerSingle,
		Category:       service.Category,
	}

	c.HTML(http.StatusOK, "createService", gin.H{
		"title":    "Создать услугу",
		"worker":   worker,
		"formData": formData,
	})
}

func (s *Services) editServicePost(c *gin.Context) {
	worker := s.authenticatedWorker(c)
	serviceID, err := uuid.Parse(c.Param("id"))

	var data ServiceFormData
	if err := c.Bind(&data); err != nil {
		c.HTML(http.StatusBadRequest, "createService", gin.H{
			"worker":   worker,
			"title":    "Изменить услугу",
			"error":    err.Error(),
			"formData": data,
		})
		return
	}

	_, err = s.Services.TaskService.Update(serviceID, data.Category, data.Name, data.PricePerSingle)
	if err != nil {
		c.HTML(http.StatusBadRequest, "createService", gin.H{
			"worker":   worker,
			"title":    "Изменить услугу",
			"error":    err.Error(),
			"formData": data,
		})
		return
	}

	c.Redirect(http.StatusFound, "/services")
}
