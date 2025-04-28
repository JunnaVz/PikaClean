// Package workerViews provides user interface functions for the PikaClean application
// focused on worker-related operations including profile viewing, updating, and management.
// This file contains functionality for listing and managing workers.
package workerViews

import (
	"fmt"
	"teamdev/cmd/modelTables"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// getAllWorkers displays a list of all workers in the system and allows a manager
// to select a worker profile to update. It retrieves the worker list from the
// service layer and displays it in a formatted table. Then it provides an interactive
// menu for the manager to select workers for profile modification.
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - manager: The authenticated manager worker object who is accessing the list
//
// Returns:
//   - error: Any error that occurred during retrieval or display of worker list
func getAllWorkers(services registry.Services, manager *models.Worker) error {
	workers, err := services.WorkerService.GetAllWorkers()

	if err != nil {
		return err
	}

	err = modelTables.Workers(services, workers)
	if err != nil {
		return err
	}

	var action int
	for {
		fmt.Print("Введите номер работника, чтобы изменить его профиль или 0, чтобы выйти\n")

		_, err = fmt.Scanf("%d", &action)
		if err != nil {
			fmt.Println(err)
		}

		if action == 0 {
			return nil
		}

		if action < 1 || action > len(workers) {
			fmt.Println("Неверный номер")
		}

		err = Update(services, workers[action-1].ID, manager)

		if err != nil {
			return err
		}
	}
}
