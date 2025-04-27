// Package interfaces provides validation functions for the business logic layer
// of the PikaClean application. These validators ensure that user inputs and
// data modifications meet the required constraints before persistence operations.
package interfaces

import (
	"github.com/google/uuid"
	"net/mail"
	"regexp"
	"teamdev/internal/models"
	"time"
)

// validName checks if a name (first name or surname) is valid.
// A valid name has at least one character.
//
// Parameters:
//   - name: The name string to validate
//
// Returns:
//   - bool: True if the name is valid, false otherwise
func validName(name string) bool {
	return len(name) > 0
}

// validPrice checks if a price value is valid.
// A valid price must be greater than zero.
//
// Parameters:
//   - price: The price value to validate
//
// Returns:
//   - bool: True if the price is valid, false otherwise
func validPrice(price float64) bool {
	return price > 0
}

// validCategory checks if a category ID is valid.
// A valid category ID must be between 1 and 8 inclusive.
//
// Parameters:
//   - category: The category ID to validate
//
// Returns:
//   - bool: True if the category ID is valid, false otherwise
func validCategory(category int) bool {
	return category > 0 && category < 9
}

// validEmail checks if an email address is valid.
// Uses the mail.ParseAddress function to verify email format.
//
// Parameters:
//   - email: The email address to validate
//
// Returns:
//   - bool: True if the email is valid, false otherwise
func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// validAddress checks if a physical address is valid.
// A valid address has at least one character.
//
// Parameters:
//   - address: The address string to validate
//
// Returns:
//   - bool: True if the address is valid, false otherwise
func validAddress(address string) bool {
	return len(address) > 0
}

// validPhoneNumber checks if a phone number is valid.
// A valid phone number must be in the format: +{country code}{10 digits}
// Example: +12345678901 or +11234567890
//
// Parameters:
//   - phoneNumber: The phone number to validate
//
// Returns:
//   - bool: True if the phone number is valid, false otherwise
func validPhoneNumber(phoneNumber string) bool {
	re := regexp.MustCompile(`^\+\d{1,3}\d{10}$`)
	return re.MatchString(phoneNumber)
}

// validPassword checks if a password is valid.
// A valid password must be at least 8 characters long and contain
// at least one letter and one number.
//
// Parameters:
//   - password: The password string to validate
//
// Returns:
//   - bool: True if the password is valid, false otherwise
func validPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	reLetter := regexp.MustCompile(`[a-zA-Z]`)
	reNumber := regexp.MustCompile(`[0-9]`)

	return reLetter.MatchString(password) && reNumber.MatchString(password)
}

// validRole checks if a role ID is valid.
// A valid role ID must be either 1 (worker) or 2 (admin).
//
// Parameters:
//   - role: The role ID to validate
//
// Returns:
//   - bool: True if the role ID is valid, false otherwise
func validRole(role int) bool {
	return role > 0 && role < 3
}

// validDeadline checks if an order deadline is valid.
// A valid deadline must be in the future.
//
// Parameters:
//   - deadline: The deadline time to validate
//
// Returns:
//   - bool: True if the deadline is valid, false otherwise
func validDeadline(deadline time.Time) bool {
	return deadline.After(time.Now())
}

// validTasksNumber checks if an order contains at least one task.
// A valid order must have at least one associated task.
//
// Parameters:
//   - tasks: Slice of ordered tasks to validate
//
// Returns:
//   - bool: True if there's at least one task, false otherwise
func validTasksNumber(tasks []models.OrderedTask) bool {
	return len(tasks) > 0
}

// validStatus checks if an order status code is valid.
// A valid status must be one of the defined status constants in models.
//
// Parameters:
//   - status: The status code to validate
//
// Returns:
//   - bool: True if the status code is valid, false otherwise
func validStatus(status int) bool {
	return status == models.NewOrderStatus || status == models.InProgressOrderStatus || status == models.CompletedOrderStatus || status == models.CancelledOrderStatus
}

// validRate checks if a user rating is valid.
// A valid rating must be between 0 and 5 inclusive.
//
// Parameters:
//   - rate: The rating value to validate
//
// Returns:
//   - bool: True if the rating is valid, false otherwise
func validRate(rate int) bool {
	return rate >= 0 && rate <= 5
}

// taskIsAttachedToOrder checks if a task is attached to an order.
// This validates whether a specific task is associated with an order
// when managing order-task relationships.
//
// Parameters:
//   - taskID: UUID of the task to check
//   - tasks: Slice of tasks to search through
//
// Returns:
//   - bool: True if the task is attached to the order, false otherwise
func taskIsAttachedToOrder(taskID uuid.UUID, tasks []models.Task) bool {
	for _, task := range tasks {
		if task.ID == taskID {
			return true
		}
	}
	return false
}
