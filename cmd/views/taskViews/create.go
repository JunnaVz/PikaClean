// Package taskViews provides view functions for managing cleaning tasks
// in the PikaClean application. It contains functions that encapsulate user interaction
// for creating, updating, and displaying cleaning service tasks.
package taskViews

import (
	utils "teamdev/cmd/cmdUtils"
	"teamdev/cmd/views/stringConst"
	"teamdev/internal/registry"
)

// Create handles the creation of a new cleaning task by collecting required information
// from the user and persisting it through the task service. It prompts the user for
// the task name, price, and category, then attempts to create the task in the database.
//
// Parameters:
//   - services: Service container providing access to business logic services
//
// Returns:
//   - error: Any error that occurred during task creation,
//     or nil if the operation was successful
func Create(services registry.Services) error {
	var name = utils.EndlessReadWord(stringConst.NameRequest)
	var price = utils.EndlessReadFloat64(stringConst.PriceRequest)
	var category = utils.EndlessReadInt(stringConst.CategoryRequest)

	_, err := services.TaskService.Create(name, price, category)
	if err != nil {
		println(err.Error())
	}

	println("Услуга успешно создана")
	return nil
}
