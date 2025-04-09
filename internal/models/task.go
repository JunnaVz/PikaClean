package models

import "github.com/google/uuid"

type Task struct {
	ID             uuid.UUID
	Name           string
	PricePerSingle float64
	Category       int
}

var TaskCategories = [8]string{
	"Генеральная уборка",
	"Послестроительная уборка",
	"Мытье окон",
	"Ежедневная уборка офисов",
	"Поддерживающая уборка",
	"Химчистка ковров и мебели",
	"Уход за твердыми полами",
	"Глубинная Эко Чистка",
}

func GetCategoryName(category int) string {
	switch category {
	case 1:
		return TaskCategories[0]
	case 2:
		return TaskCategories[1]
	case 3:
		return TaskCategories[2]
	case 4:
		return TaskCategories[3]
	case 5:
		return TaskCategories[4]
	case 6:
		return TaskCategories[5]
	case 7:
		return TaskCategories[6]
	case 8:
		return TaskCategories[7]
	default:
		return "Неизвестная категория"
	}
}
