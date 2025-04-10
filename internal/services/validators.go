package interfaces

import (
	"github.com/google/uuid"
	"lab3/internal/models"
	"net/mail"
	"regexp"
	"time"
)

func validName(name string) bool {
	return len(name) > 0
}

func validPrice(price float64) bool {
	return price > 0
}

func validCategory(category int) bool {
	return category > 0 && category < 9
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func validAddress(address string) bool {
	return len(address) > 0
}

func validPhoneNumber(phoneNumber string) bool {
	re := regexp.MustCompile(`^\+\d{1,3}\d{10}$`)
	return re.MatchString(phoneNumber)
}

func validPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	reLetter := regexp.MustCompile(`[a-zA-Z]`)
	reNumber := regexp.MustCompile(`[0-9]`)

	return reLetter.MatchString(password) && reNumber.MatchString(password)
}

func validRole(role int) bool {
	return role > 0 && role < 3
}

func validDeadline(deadline time.Time) bool {
	return deadline.After(time.Now())
}

func validTasksNumber(tasks []models.OrderedTask) bool {
	return len(tasks) > 0
}

func validStatus(status int) bool {
	return status == models.NewOrderStatus || status == models.InProgressOrderStatus || status == models.CompletedOrderStatus || status == models.CancelledOrderStatus
}

func validRate(rate int) bool {
	return rate >= 0 && rate <= 5
}

func taskIsAttachedToOrder(taskID uuid.UUID, tasks []models.Task) bool {
	for _, task := range tasks {
		if task.ID == taskID {
			return true
		}
	}
	return false
}
