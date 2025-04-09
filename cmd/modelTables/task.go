package modelTables

import (
	"fmt"
	"lab3/cmd/cmdUtils"
	"lab3/internal/models"
	"os"
	"text/tabwriter"
)

func Tasks(tasks []models.Task) error {
	var err error

	maxNameLen, maxPriceLen, maxCategoryLen := 0, 0, 0
	for _, task := range tasks {
		if len(task.Name) > maxNameLen {
			maxNameLen = len(task.Name)
		}
		priceLen := len(fmt.Sprintf("%.2f", task.PricePerSingle))
		if priceLen > maxPriceLen {
			maxPriceLen = priceLen
		}
		categoryLen := len(models.GetCategoryName(task.Category))
		if categoryLen > maxCategoryLen {
			maxCategoryLen = categoryLen
		}
	}

	t := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	t.Init(os.Stdout, 2, 4, 5, ' ', 0)

	_, err = fmt.Fprintf(t, "\n %s\t%s\t%s\t%s",
		"№", "Название", "Стоимость", "Категория")
	if err != nil {
		fmt.Println(err)
	}

	for i, task := range tasks {
		_, err = fmt.Fprintf(t, "\n %d\t%s\t%.2f\t%s\t",
			i+1, cmdUtils.TruncateString(task.Name, 27), task.PricePerSingle, cmdUtils.TruncateString(models.GetCategoryName(task.Category), 27))
		if err != nil {
			return err
		}
	}

	err = t.Flush()
	if err != nil {
		return err
	}

	fmt.Printf("\n-----------\n\n")
	return nil
}
