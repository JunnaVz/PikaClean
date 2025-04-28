// Package models provides data structures representing the core domain entities
// of the PikaClean application, including workers, users, tasks, and orders.
package models

// OrderedTask represents a task that has been ordered by a client with a specific quantity.
// It combines a Task with a quantity value, representing how many units of the task
// have been requested as part of an order.
type OrderedTask struct {
	Task     *Task // Reference to the cleaning task being ordered
	Quantity int   // Number of units of the task ordered by the client
}
