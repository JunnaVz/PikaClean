// Package orderViews provides view functions for managing cleaning service orders
// in the PikaClean application. It contains functions that encapsulate business
// logic for common order-related operations and UI workflows.
package orderViews

import (
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

// CancelOrder updates an existing order to the cancelled status.
// It uses the application's order service to change the status while
// maintaining other order properties unchanged.
//
// Parameters:
//   - services: Service container providing access to business logic services
//   - order: The order to be cancelled
//
// Returns:
//   - error: Any error that occurred during the cancellation process,
//     or nil if the operation was successful
func CancelOrder(services registry.Services, order *models.Order) error {
	_, err := services.OrderService.Update(order.ID, models.CancelledOrderStatus, order.Rate, order.WorkerID)
	if err != nil {
		return err
	}

	return nil
}
