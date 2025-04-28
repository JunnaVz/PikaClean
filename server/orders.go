package server

import (
	"encoding/json"
	"strconv"
	"teamdev/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Rating struct {
	Rating string `json:"rating"`
}

func (s *Services) rateOrderApiPost(c *gin.Context) {
	strOrderID := c.Param("id")

	orderID, err := uuid.Parse(strOrderID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	rating := Rating{}
	err = json.NewDecoder(c.Request.Body).Decode(&rating)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid rating",
		})
		return
	}

	ratingInt, err := strconv.Atoi(rating.Rating)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid rating",
		})
		return

	}

	order, err := s.Services.OrderService.GetOrderByID(orderID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Order not found",
		})
		return
	}

	if order.UserID != s.authenticatedUser(c).ID {
		c.JSON(400, gin.H{
			"error": "You are not the owner of this order",
		})
		return
	}

	if order.Status != 3 {
		c.JSON(400, gin.H{
			"error": "Order is not completed",
		})
		return
	}

	_, err = s.Services.OrderService.Update(order.ID, order.Status, ratingInt, order.WorkerID)

	c.JSON(200, gin.H{
		"message": "Rating updated",
	})
}

func (s *Services) cancelOrderApiPost(c *gin.Context) {
	strOrderID := c.Param("id")

	orderID, err := uuid.Parse(strOrderID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	order, err := s.Services.OrderService.GetOrderByID(orderID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Order not found",
		})
		return
	}

	if order.UserID != s.authenticatedUser(c).ID {
		c.JSON(400, gin.H{
			"error": "You are not the owner of this order",
		})
		return
	}

	if order.Status != 1 {
		c.JSON(400, gin.H{
			"error": "Order is not new",
		})
		return
	}

	_, err = s.Services.OrderService.Update(order.ID, models.CancelledOrderStatus, models.NoStatus, order.WorkerID)

	c.JSON(200, gin.H{
		"message": "Order cancelled",
	})
}

type statusData struct {
	Status string `json:"status"`
}

func (s *Services) changeStatusOrderApiPost(c *gin.Context) {
	strOrderID := c.Param("id")

	orderID, err := uuid.Parse(strOrderID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	status := statusData{}
	err = json.NewDecoder(c.Request.Body).Decode(&status)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid status",
		})
		return
	}

	statusInt, err := strconv.Atoi(status.Status)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid status",
		})
		return

	}

	order, err := s.Services.OrderService.GetOrderByID(orderID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Order not found",
		})
		return
	}

	if order.WorkerID != s.authenticatedWorker(c).ID && s.authenticatedWorker(c).Role != models.ManagerRole {
		c.JSON(400, gin.H{
			"error": "You are not the owner of this order",
		})
		return
	}

	if order.Status > statusInt {
		c.JSON(400, gin.H{
			"error": "Invalid status",
		})
		return
	}

	_, err = s.Services.OrderService.Update(order.ID, statusInt, order.Rate, order.WorkerID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Status updated",
	})
}

type changeWorkerData struct {
	WorkerID string `json:"workerId"`
}

func (s *Services) changeWorkerApiPost(c *gin.Context) {
	strOrderID := c.Param("id")

	orderID, err := uuid.Parse(strOrderID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	order, err := s.Services.OrderService.GetOrderByID(orderID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Order not found",
		})
		return
	}

	data := changeWorkerData{}
	err = json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid worker ID",
		})
		return
	}

	if data.WorkerID == "" {
		c.JSON(400, gin.H{
			"error": "Invalid worker ID",
		})
		return
	}

	workerID, err := uuid.Parse(data.WorkerID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid worker ID",
		})
		return
	}

	worker, err := s.Services.WorkerService.GetWorkerByID(workerID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Worker not found",
		})
		return
	}

	if worker.Role != models.MasterRole {
		c.JSON(400, gin.H{
			"error": "Worker is not a master",
		})
		return
	}

	_, err = s.Services.OrderService.Update(order.ID, order.Status, order.Rate, workerID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Worker assigned",
	})
}
