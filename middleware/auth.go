package middleware

import (
	"net/http"
	"teamdev/internal/registry"

	"github.com/google/uuid"

	//"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	Services *registry.Services
}

func NewMiddleware(registry registry.App) *Middleware {
	return &Middleware{Services: registry.Services}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		session := sessions.Default(c)
		sessionID := session.Get("userID")
		if sessionID == nil {
			c.Redirect(http.StatusMovedPermanently, "/auth/signin")
			c.Abort()
			return
		}

		strUserId, ok := sessionID.(string)
		if !ok {
			c.Redirect(http.StatusMovedPermanently, "/auth/signin")
			c.Abort()
			return
		}
		// Check if the user exists
		userId, err := uuid.Parse(strUserId)
		user, err := m.Services.UserService.GetUserByID(userId)
		if err != nil || user.ID == uuid.Nil {
			c.Redirect(http.StatusMovedPermanently, "/auth/signin")
			c.Abort()
			return
		}

		c.Set("userID", user.ID)
		c.Next()
	}
}

func (m *Middleware) WorkerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		session := sessions.Default(c)
		sessionID := session.Get("workerID")
		if sessionID == nil {
			c.Redirect(http.StatusMovedPermanently, "/worker-auth/signin")
			c.Abort()
			return
		}

		strWorkerId, ok := sessionID.(string)
		if !ok {
			c.Redirect(http.StatusMovedPermanently, "/worker-auth/signin")
			c.Abort()
			return
		}
		// Check if the user exists
		workerID, err := uuid.Parse(strWorkerId)
		worker, err := m.Services.WorkerService.GetWorkerByID(workerID)
		if err != nil || worker.ID == uuid.Nil {
			c.Redirect(http.StatusMovedPermanently, "/worker-auth/signin")
			c.Abort()
			return
		}

		c.Set("workerID", worker.ID)
		c.Next()
	}
}
