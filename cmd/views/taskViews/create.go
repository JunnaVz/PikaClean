package taskViews

import (
	utils "teamdev/cmd/cmdUtils"
	"teamdev/cmd/views/stringConst"
	"teamdev/internal/registry"
)

func Create(services registry.Services) error {
	var name = utils.EndlessReadWord(stringConst.NameRequest)
	var price = utils.EndlessReadFloat64(stringConst.PriceRequest)
	var category = utils.EndlessReadInt(stringConst.CategoryRequest)

	_, err := services.TaskService.Create(name, price, category)
	if err != nil {
		println(err.Error())
	}

	println("Услуга успешно создана")
	return nil
}
