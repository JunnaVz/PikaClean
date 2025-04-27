// Package modelTables provides functionality for displaying domain models
// in tabular format for command-line interfaces in the PikaClean application.
// It offers formatted tabular output for various entities such as orders,
// users, workers, and tasks, making data easily readable in terminal displays.
package modelTables

import (
	"fmt"
	"os"
	"teamdev/internal/models"
	"teamdev/internal/registry"
	"text/tabwriter"
)

// Workers renders a slice of Worker entities in a formatted table on the console.
// It displays worker information including name, role, phone number, email address,
// and average rating in an aligned tabular format for better readability.
//
// The function uses the tabwriter package to ensure proper alignment of columns,
// and leverages the worker service to obtain additional data such as average ratings.
//
// Parameters:
//   - services: Registry Services container providing access to business logic services
//   - workers: A slice of models.Worker entities to display in the table
//
// Returns:
//   - error: Any error that occurs during formatting or output operations
func Workers(services registry.Services, workers []models.Worker) error {
	var err error

	// Initialize tabwriter for formatted columnar output
	t := new(tabwriter.Writer)
	t.Init(os.Stdout, 1, 4, 2, ' ', 0)

	// Write the table header
	_, err = fmt.Fprintf(t, "\n %s\t%s\t%s\t%s\t%s\t%s\n",
		"№", "Имя", "Роль", "Телефон", "Email", "Ср. оценка")
	if err != nil {
		fmt.Println(err)
	}

	// Write each worker as a table row
	for i, worker := range workers {
		// Obtain the worker's average rating from completed orders
		workersRate, _ := services.WorkerService.GetAverageOrderRate(&worker)

		fmt.Fprintf(t, " %d\t%s\t%s\t%s\t%s\t%f\n",
			i+1, worker.FullName(), worker.DisplayRole(), worker.PhoneNumber, worker.Email, workersRate)
	}

	// Flush buffered output to standard output
	err = t.Flush()
	if err != nil {
		return err
	}

	return nil
}
