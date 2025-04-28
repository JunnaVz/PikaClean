package utils

import (
	"teamdev/internal/models"
)

func DisplayStatus(statusNum int) string {
	return models.OrderStatuses[statusNum]
}
