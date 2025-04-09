package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID           uuid.UUID
	WorkerID     uuid.UUID
	UserID       uuid.UUID
	Status       int
	Address      string
	CreationDate time.Time
	Deadline     time.Time
	Rate         int
}

const NoStatus = 0
const NewOrderStatus = 1
const InProgressOrderStatus = 2
const CompletedOrderStatus = 3
const CancelledOrderStatus = 4

var OrderStatuses = map[int]string{
	NoStatus:              "Не определен",
	NewOrderStatus:        "Новый",
	InProgressOrderStatus: "В процессе",
	CompletedOrderStatus:  "Завершен",
	CancelledOrderStatus:  "Отменен",
}
