// Package models provides data structures representing the core domain entities
// of the PikaClean application, including workers, users, tasks, and orders.
// This package defines the structure and behavior of these entities
// without any database implementation details.
package models

import "github.com/google/uuid"

// Worker represents a staff member of the cleaning service.
// Workers can be either managers who oversee operations or
// cleaning masters who perform the actual cleaning tasks.
type Worker struct {
	ID          uuid.UUID // Unique identifier for the worker
	Name        string    // First name of the worker
	Surname     string    // Last name of the worker
	Address     string    // Physical address of the worker
	PhoneNumber string    // Contact phone number
	Email       string    // Email address used for communication and login
	Role        int       // Role identifier (ManagerRole or MasterRole)
	Password    string    // Hashed password for authentication
}

// ManagerRole is a constant indicating that a worker has manager privileges.
// Managers can oversee operations and have administrative access.
const ManagerRole = 1

// MasterRole is a constant indicating that a worker performs cleaning services.
// Masters are assigned to and complete cleaning orders.
const MasterRole = 2

// WorkerRole maps role identifiers to human-readable role names.
// This mapping is used for display purposes in the user interface.
var WorkerRole = map[int]string{
	ManagerRole: "Менеджер",
	MasterRole:  "Мастер",
}

// DisplayRole returns the human-readable name of the worker's role.
// Returns either "Менеджер" (Manager) or "Мастер" (Master) based on the worker's Role field.
func (w Worker) DisplayRole() string {
	return WorkerRole[w.Role]
}

// FullName returns the complete name of the worker by combining
// the first name and surname with a space in between.
// Used for displaying the worker's name in the user interface.
func (w Worker) FullName() string {
	return w.Name + " " + w.Surname
}
