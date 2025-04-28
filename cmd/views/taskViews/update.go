// Package taskViews provides view functions for managing cleaning tasks
// in the PikaClean application. It contains functions that encapsulate user interaction
// for creating, updating, and displaying cleaning service tasks.
package taskViews

import (
	"fmt"
	utils "teamdev/cmd/cmdUtils"
	"teamdev/cmd/views/stringConst"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// Update handles the process of modifying an existing cleaning task. It prompts the user
// for updated information including name, price, and category, then persists the changes
// through the task service. After successful update, a confirmation message is displayed.
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - task: The existing task to be updated with its current values
//
// Returns:
//   - *models.Task: The updated task with new values after modification
//   - error: Any error that occurred during task update,
//     or nil if the operation was successful
func Update(services registry.Services, task models.Task) (*models.Task, error) {
	var name = utils.EndlessReadRow(stringConst.NameRequest)
	var price = utils.EndlessReadFloat64(stringConst.PriceRequest)
	var category = utils.EndlessReadInt(stringConst.CategoryRequest)

	updatedTask, err := services.TaskService.Update(task.ID, category, name, price)

	fmt.Println("Услуга успешно обновлена")
	return updatedTask, err
}
