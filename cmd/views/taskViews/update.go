package taskViews

import (
	"fmt"
	utils "teamdev/cmd/cmdUtils"
	"teamdev/cmd/views/stringConst"
	"teamdev/internal/models"
	"teamdev/internal/registry"
)

func Update(services registry.Services, task models.Task) (*models.Task, error) {
	var name = utils.EndlessReadRow(stringConst.NameRequest)
	var price = utils.EndlessReadFloat64(stringConst.PriceRequest)
	var category = utils.EndlessReadInt(stringConst.CategoryRequest)

	updatedTask, err := services.TaskService.Update(task.ID, category, name, price)

	fmt.Println("Услуга успешно обновлена")
	return updatedTask, err
}
