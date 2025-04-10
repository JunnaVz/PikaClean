package service_errors

import "errors"

var (
	InvalidName                  = errors.New("invalid name")
	InvalidPrice                 = errors.New("invalid price")
	InvalidCategory              = errors.New("invalid category")
	InvalidEmail                 = errors.New("invalid email")
	InvalidAddress               = errors.New("invalid address")
	InvalidPhoneNumber           = errors.New("invalid phone number")
	InvalidPassword              = errors.New("invalid password")
	InvalidRole                  = errors.New("invalid role")
	InvalidAddressOrder          = errors.New("invalid address of the order")
	InvalidDeadlineOrder         = errors.New("invalid deadline of the order")
	EmptyTasksOrder              = errors.New("order has no tasks")
	NotUnique                    = errors.New("such row already exists")
	MismatchedPassword           = errors.New("passwords do not match")
	InvalidReference             = errors.New("invalid reference")
	OrderIsNotCompleted          = errors.New("order is not completed")
	OrderIsAlreadyCompleted      = errors.New("order is already completed")
	RatingOutOfRange             = errors.New("rating is out of range")
	InvalidOrderStatus           = errors.New("invalid order status")
	TaskIsNotAttachedToOrder     = errors.New("task is not attached to the order")
	TaskIsAlreadyAttachedToOrder = errors.New("task is already attached to the order")
	NegativeQuantity             = errors.New("quantity is negative")
)
