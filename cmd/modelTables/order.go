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

// Orders renders a slice of Order entities in a formatted table on the console.
// It displays order information including creation date, status, address, and user rating
// in an aligned tabular format for better readability.
//
// The function automatically adjusts column widths to accommodate data while
// truncating excessively long strings to maintain display consistency. It uses
// the tabwriter package to ensure proper alignment of columns.
//
// Parameters:
//   - orders: A slice of models.Order entities to display in the table
//
// Returns:
//   - error: Any error that occurs during formatting or output operations
func Orders(orders []models.Order) error {
	var err error

	// Calculate maximum widths for variable-length fields
	maxAddressLen, maxStatusLen := 0, 0
	for _, order := range orders {
		if len(order.Address) > maxAddressLen {
			maxAddressLen = len(order.Address)
		}
		statusLen := len(models.OrderStatuses[order.Status])
		if statusLen > maxStatusLen {
			maxStatusLen = statusLen
		}
	}

	// Initialize tabwriter for formatted columnar output
	t := new(tabwriter.Writer)

	t.Init(os.Stdout, 2, 4, 5, ' ', 0)

	// Write the table header
	_, err = fmt.Fprintf(t, "\n %s\t%s\t%s\t%s\t%s",
		"№", "Дата создания", "Статус", "Адрес", "Оценка")
	if err != nil {
		fmt.Println(err)
	}

	// Write each order as a table row
	for i, order := range orders {
		_, err = fmt.Fprintf(t, "\n %d\t%s\t%s\t%s\t%d",
			i+1, order.CreationDate.Format("2006-01-02"), cmdUtils.TruncateString(models.OrderStatuses[order.Status], 20), cmdUtils.TruncateString(order.Address, 20), order.Rate)
		if err != nil {
			return err
		}
	}

	// Flush buffered output to standard output
	err = t.Flush()
	if err != nil {
		return err
	}

	return nil
}
