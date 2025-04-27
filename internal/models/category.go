// Package models provides data structures representing the core domain entities
// of the PikaClean application, including workers, users, tasks, orders, and categories.
package models

// Category represents a classification grouping for cleaning tasks.
// Each category has a unique identifier and a descriptive name.
// Categories are used to organize and filter the cleaning services offered.
type Category struct {
	ID   int    // Unique identifier for the category
	Name string // Descriptive name of the category
}
