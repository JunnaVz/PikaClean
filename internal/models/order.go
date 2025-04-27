// Package models provides data structures representing the core domain entities
// of the PikaClean application, including workers, users, tasks, and orders.
package models

import (
	"github.com/google/uuid"
	"time"
)

// Order represents a cleaning service request from a user.
// It contains information about who placed the order, who is assigned to fulfill it,
// when it should be completed, and its current status in the workflow.
type Order struct {
	ID           uuid.UUID // Unique identifier for the order
	WorkerID     uuid.UUID // ID of the worker assigned to fulfill the order
	UserID       uuid.UUID // ID of the user who placed the order
	Status       int       // Current status of the order (see status constants)
	Address      string    // Location where cleaning services should be performed
	CreationDate time.Time // When the order was created in the system
	Deadline     time.Time // When the order should be completed by
	Rate         int       // Customer satisfaction rating (0-5)
}

// NoStatus indicates an order with an undefined status.
const NoStatus = 0

// NewOrderStatus indicates a newly created order that hasn't been assigned yet.
const NewOrderStatus = 1

// InProgressOrderStatus indicates an order that has been assigned and is being worked on.
const InProgressOrderStatus = 2

// CompletedOrderStatus indicates an order that has been successfully fulfilled.
const CompletedOrderStatus = 3

// CancelledOrderStatus indicates an order that was cancelled before completion.
const CancelledOrderStatus = 4

// OrderStatuses maps numeric status codes to human-readable status descriptions.
// Used for displaying order status in the user interface.
var OrderStatuses = map[int]string{
	NoStatus:              "Не определен", // Undefined
	NewOrderStatus:        "Новый",        // New
	InProgressOrderStatus: "В процессе",   // In progress
	CompletedOrderStatus:  "Завершен",     // Completed
	CancelledOrderStatus:  "Отменен",      // Cancelled
}
