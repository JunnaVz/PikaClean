// Package modelTables provides functionality for displaying domain models
// in tabular format for command-line interfaces in the PikaClean application.
// It offers formatted tabular output for various entities such as orders,
// users, workers, and tasks, making data easily readable in terminal displays.
package modelTables

import (
	"fmt"
	"os"
	"teamdev/cmd/cmdUtils"
	"teamdev/internal/models"
	"text/tabwriter"
)

// Tasks renders a slice of Task entities in a formatted table on the console.
// It displays task information including name, price per unit, and category
// in an aligned tabular format for better readability.
//
// The function automatically adjusts column widths to accommodate data while
// truncating excessively long strings to maintain display consistency. It uses
// the tabwriter package to ensure proper alignment of columns.
//
// Parameters:
//   - tasks: A slice of models.Task entities to display in the table
//
// Returns:
//   - error: Any error that occurs during formatting or output operations
func Tasks(tasks []models.Task) error {
	var err error

	// Calculate maximum widths for variable-length fields
	maxNameLen, maxPriceLen, maxCategoryLen := 0, 0, 0
	for _, task := range tasks {
		if len(task.Name) > maxNameLen {
			maxNameLen = len(task.Name)
		}
		priceLen := len(fmt.Sprintf("%.2f", task.PricePerSingle))
		if priceLen > maxPriceLen {
			maxPriceLen = priceLen
		}
		categoryLen := len(models.GetCategoryName(task.Category))
		if categoryLen > maxCategoryLen {
			maxCategoryLen = categoryLen
		}
	}

	// Initialize tabwriter for formatted columnar output
	t := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	t.Init(os.Stdout, 2, 4, 5, ' ', 0)

	// Write the table header
	_, err = fmt.Fprintf(t, "\n %s\t%s\t%s\t%s",
		"№", "Название", "Стоимость", "Категория")
	if err != nil {
		fmt.Println(err)
	}

	// Write each task as a table row
	for i, task := range tasks {
		_, err = fmt.Fprintf(t, "\n %d\t%s\t%.2f\t%s\t",
			i+1, cmdUtils.TruncateString(task.Name, 27), task.PricePerSingle, cmdUtils.TruncateString(models.GetCategoryName(task.Category), 27))
		if err != nil {
			return err
		}
	}

	// Flush buffered output to standard output
	err = t.Flush()
	if err != nil {
		return err
	}

	// Add a visual separator after the table
	fmt.Printf("\n-----------\n\n")
	return nil
}
