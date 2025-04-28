package server

import "github.com/gin-gonic/gin"

func (s *Services) index(c *gin.Context) {
	c.HTML(200, "index", gin.H{
		"title":  "Домашняя страница",
		"auth":   s.authenticatedUser(c),
		"worker": s.authenticatedWorker(c),
	})
}
