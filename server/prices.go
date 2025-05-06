package server

import (
	"teamdev/internal/models"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

func (s *Services) priceList(c *gin.Context) {
	authUser := s.authenticatedUser(c)
	worker := s.authenticatedWorker(c)

	categories, err := s.Services.CategoryService.GetAll()
	if err != nil {
		log.Printf("Error getting categories: %v", err)
		categories = []models.Category{}
	}

	prices := make(map[string][]models.Task)
	for _, category := range categories {
		tasks, err := s.Services.TaskService.GetTasksInCategory(category.ID)
		if err != nil {
			log.Printf("Error getting tasks in category %s: %v", category.Name, err)
			continue
		}
		prices[category.Name] = tasks
	}

	c.HTML(200, "prices", gin.H{
		"auth":   authUser,
		"worker": worker,
		"title":  "Услуги",
		"prices": prices,
	})
}
