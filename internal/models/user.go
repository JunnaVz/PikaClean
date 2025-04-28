// Package models provides data structures representing the core domain entities
// of the PikaClean application, including workers, users, tasks, and orders.
// This package defines the structure and behavior of these entities
// without any database implementation details.
package models

import "github.com/google/uuid"

// User represents a client of the cleaning service.
// Users can place orders for cleaning services and have
// their contact information stored in the system.
type User struct {
	ID          uuid.UUID // Unique identifier for the user
	Name        string    // First name of the user
	Surname     string    // Last name of the user
	Address     string    // Physical address where cleaning services might be performed
	PhoneNumber string    // Contact phone number for notifications and communication
	Email       string    // Email address used for account access and communication
	Password    string    // Hashed password for authentication
}
