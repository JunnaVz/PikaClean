package taskViews

import (
	utils "lab3/cmd/cmdUtils"
	"lab3/cmd/views/stringConst"
	"lab3/internal/registry"
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
