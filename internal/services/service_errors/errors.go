// Package service_errors defines standard error types for the PikaClean application's
// service layer. This package provides a centralized set of error definitions to ensure
// consistent error handling and reporting throughout the business logic components.
package service_errors

import "errors"

var (
	// InvalidName indicates that a name field doesn't meet validation requirements
	// (e.g., empty string, too short, contains invalid characters).
	InvalidName = errors.New("invalid name")

	// InvalidPrice indicates that a price value is not acceptable
	// (e.g., negative, zero when it should be positive, or too high).
	InvalidPrice = errors.New("invalid price")

	// InvalidCategory indicates that a category reference is not valid
	// (e.g., category doesn't exist or is not applicable in the context).
	InvalidCategory = errors.New("invalid category")

	// InvalidEmail indicates that an email address doesn't meet format requirements
	// or violates domain-specific rules (e.g., not a valid RFC 5322 format).
	InvalidEmail = errors.New("invalid email")

	// InvalidAddress indicates that a physical address doesn't meet validation requirements
	// (e.g., empty, missing required components, or invalid format).
	InvalidAddress = errors.New("invalid address")

	// InvalidPhoneNumber indicates that a phone number doesn't meet format requirements
	// (e.g., too short, contains non-numeric characters when it shouldn't).
	InvalidPhoneNumber = errors.New("invalid phone number")

	// InvalidPassword indicates that a password doesn't meet security requirements
	// (e.g., too short, missing required character types, or too common).
	InvalidPassword = errors.New("invalid password")

	// InvalidRole indicates that a role identifier is not recognized
	// or is not applicable in the given context.
	InvalidRole = errors.New("invalid role")

	// InvalidAddressOrder indicates that an order's service address is invalid
	// (e.g., empty, unreachable, or outside of service area).
	InvalidAddressOrder = errors.New("invalid address of the order")

	// InvalidDeadlineOrder indicates that an order's deadline is not acceptable
	// (e.g., in the past, too soon to be fulfilled, or too far in the future).
	InvalidDeadlineOrder = errors.New("invalid deadline of the order")

	// EmptyTasksOrder indicates an attempt to create or process an order with no tasks.
	EmptyTasksOrder = errors.New("order has no tasks")

	// NotUnique indicates a violation of a uniqueness constraint
	// (e.g., duplicate email address, phone number, or task name).
	NotUnique = errors.New("such row already exists")

	// MismatchedPassword indicates that provided passwords do not match
	// (e.g., during password confirmation in registration or change password flows).
	MismatchedPassword = errors.New("passwords do not match")

	// InvalidReference indicates a reference to a non-existent entity
	// (e.g., user ID, worker ID, task ID that doesn't exist in the database).
	InvalidReference = errors.New("invalid reference")

	// OrderIsNotCompleted indicates an operation that requires a completed order
	// is being performed on an order that hasn't reached completed status.
	OrderIsNotCompleted = errors.New("order is not completed")

	// OrderIsAlreadyCompleted indicates an attempt to modify an order that has
	// already been marked as completed.
	OrderIsAlreadyCompleted = errors.New("order is already completed")

	// RatingOutOfRange indicates that a rating value is outside the acceptable range
	// (e.g., not between 0-5 stars).
	RatingOutOfRange = errors.New("rating is out of range")

	// InvalidOrderStatus indicates that an order status value is not recognized
	// or represents an invalid state transition.
	InvalidOrderStatus = errors.New("invalid order status")

	// TaskIsNotAttachedToOrder indicates an operation on a task within an order
	// where the task is not actually part of the order.
	TaskIsNotAttachedToOrder = errors.New("task is not attached to the order")

	// TaskIsAlreadyAttachedToOrder indicates an attempt to add a task to an order
	// when that task is already part of the order.
	TaskIsAlreadyAttachedToOrder = errors.New("task is already attached to the order")

	// NegativeQuantity indicates an attempt to set a negative quantity value
	// for a task in an order or other quantity field.
	NegativeQuantity = errors.New("quantity is negative")
)
