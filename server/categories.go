package server

import (
	"net/http"
	"strconv"
	"teamdev/internal/models"

	"github.com/gin-gonic/gin"
)

func (s *Services) createCategoryGet(c *gin.Context) {
	worker := s.authenticatedWorker(c)

	c.HTML(http.StatusOK, "createCategory", gin.H{
		"title":  "Создать категорию",
		"worker": worker,
	})
}

type CategoryFormData struct {
	Name string `form:"name"`
}

func (s *Services) createCategoryPost(c *gin.Context) {
	worker := s.authenticatedWorker(c)

	var data CategoryFormData
	if err := c.Bind(&data); err != nil {
		c.HTML(http.StatusBadRequest, "createCategory", gin.H{
			"worker":   worker,
			"title":    "Создать категорию",
			"error":    err.Error(),
			"formData": data,
		})
		return
	}

	_, err := s.Services.CategoryService.Create(data.Name)
	if err != nil {
		c.HTML(http.StatusBadRequest, "createCategory", gin.H{
			"worker":   worker,
			"title":    "Создать категорию",
			"error":    err.Error(),
			"formData": data,
		})
		return
	}

	c.Redirect(http.StatusFound, "/services")
}

func (s *Services) editCategoryGet(c *gin.Context) {
	worker := s.authenticatedWorker(c)
	serviceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "createCategory", gin.H{
			"worker": worker,
			"title":  "Создать категорию",
			"error":  "Неверный идентификатор категории",
		})
		return
	}

	service, err := s.Services.CategoryService.GetByID(serviceID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "createCategory", gin.H{
			"worker": worker,
			"title":  "Создать категорию",
			"error":  "Категория не найдена",
		})
		return
	}

	formData := CategoryFormData{
		Name: service.Name,
	}

	c.HTML(http.StatusOK, "createCategory", gin.H{
		"title":    "Создать категорию",
		"worker":   worker,
		"formData": formData,
	})
}

func (s *Services) editCategoryPost(c *gin.Context) {
	worker := s.authenticatedWorker(c)
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "createCategory", gin.H{
			"worker": worker,
			"title":  "Создать категорию",
			"error":  "Неверный идентификатор категории",
		})
		return
	}

	var data CategoryFormData
	if err := c.Bind(&data); err != nil {
		c.HTML(http.StatusBadRequest, "createCategory", gin.H{
			"worker":   worker,
			"title":    "Изменить категорию",
			"error":    err.Error(),
			"formData": data,
		})
		return
	}

	_, err = s.Services.CategoryService.Update(&models.Category{
		ID:   categoryID,
		Name: data.Name,
	})
	if err != nil {
		c.HTML(http.StatusBadRequest, "createCategory", gin.H{
			"worker":   worker,
			"title":    "Изменить категорию",
			"error":    err.Error(),
			"formData": data,
		})
		return
	}

	c.Redirect(http.StatusFound, "/services")
}
