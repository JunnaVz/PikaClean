package taskViews

import (
	"fmt"
	utils "lab3/cmd/cmdUtils"
	"lab3/cmd/views/stringConst"
	"lab3/internal/models"
	"lab3/internal/registry"
)

func Update(services registry.Services, task models.Task) (*models.Task, error) {
	var name = utils.EndlessReadRow(stringConst.NameRequest)
	var price = utils.EndlessReadFloat64(stringConst.PriceRequest)
	var category = utils.EndlessReadInt(stringConst.CategoryRequest)

	updatedTask, err := services.TaskService.Update(task.ID, category, name, price)

	fmt.Println("Услуга успешно обновлена")
	return updatedTask, err
}
