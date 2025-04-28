// Package workerViews provides user interface functions for the PikaClean application
// focused on worker-related operations like assignment to orders, profile management,
// and administrative worker functions. It handles worker interactions through the command line
// interface for operations performed by cleaning staff and administrators.
package workerViews

import (
	"fmt"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// Get displays detailed worker profile information on the command line.
// It retrieves the worker data from the database and formats it for display,
// showing details such as role, contact information, and personal details.
//
// Parameters:
//   - service: Service container providing access to business logic services,
//     particularly the WorkerService for worker data retrieval
//   - worker: Worker entity whose detailed information should be displayed
//
// Returns:
//   - error: Any error that occurred during worker information retrieval,
//     such as database errors or if the worker doesn't exist
func Get(service registry.Services, worker *models.Worker) error {
	workerFromDB, err := service.WorkerService.GetWorkerByID(worker.ID)
	if err != nil {
		return err
	}

	fmt.Print("\nWorker info:\n")
	fmt.Printf("Роль: %s\nEmail: %s\nИмя: %s\nФамилия: %s\nТелефон: %s\nАдрес: %s\n", models.WorkerRole[workerFromDB.Role], workerFromDB.Email, workerFromDB.Name, workerFromDB.Surname, workerFromDB.PhoneNumber, workerFromDB.Address)
	fmt.Print("----------------\n")
	return nil
}
