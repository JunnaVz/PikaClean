// Package models provides data structures representing the core domain entities
// of the PikaClean application, including workers, users, tasks, and orders.
package models

import "github.com/google/uuid"

// Task represents a cleaning service that can be ordered by clients.
// Each task has a unique identifier, name, price, and belongs to a specific category.
type Task struct {
	ID             uuid.UUID // Unique identifier for the task
	Name           string    // Name of the cleaning task
	PricePerSingle float64   // Price per unit of the task
	Category       int       // Category identifier (1-8 matching TaskCategories)
}

// TaskCategories defines the available cleaning service categories offered by PikaClean.
// The index (plus 1) corresponds to the category identifier used in the Task struct.
var TaskCategories = [8]string{
	"Генеральная уборка",        // General cleaning
	"Послестроительная уборка",  // Post-construction cleaning
	"Мытье окон",                // Window washing
	"Ежедневная уборка офисов",  // Daily office cleaning
	"Поддерживающая уборка",     // Maintenance cleaning
	"Химчистка ковров и мебели", // Carpet and furniture dry cleaning
	"Уход за твердыми полами",   // Hard floor care
	"Глубинная Эко Чистка",      // Deep eco cleaning
}

// GetCategoryName returns the human-readable name of a task category based on its ID.
// If the category ID is not valid (1-8), it returns "Неизвестная категория" (Unknown category).
func GetCategoryName(category int) string {
	switch category {
	case 1:
		return TaskCategories[0]
	case 2:
		return TaskCategories[1]
	case 3:
		return TaskCategories[2]
	case 4:
		return TaskCategories[3]
	case 5:
		return TaskCategories[4]
	case 6:
		return TaskCategories[5]
	case 7:
		return TaskCategories[6]
	case 8:
		return TaskCategories[7]
	default:
		return "Неизвестная категория"
	}
}
