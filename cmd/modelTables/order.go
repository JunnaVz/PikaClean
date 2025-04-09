package modelTables

import (
	"fmt"
	"lab3/cmd/cmdUtils"
	"lab3/internal/models"
	"os"
	"text/tabwriter"
)

func Orders(orders []models.Order) error {
	var err error

	maxAddressLen, maxStatusLen := 0, 0
	for _, order := range orders {
		if len(order.Address) > maxAddressLen {
			maxAddressLen = len(order.Address)
		}
		statusLen := len(models.OrderStatuses[order.Status])
		if statusLen > maxStatusLen {
			maxStatusLen = statusLen
		}
	}

	t := new(tabwriter.Writer)

	t.Init(os.Stdout, 2, 4, 5, ' ', 0)

	_, err = fmt.Fprintf(t, "\n %s\t%s\t%s\t%s\t%s",
		"№", "Дата создания", "Статус", "Адрес", "Оценка")
	if err != nil {
		fmt.Println(err)
	}

	for i, order := range orders {
		_, err = fmt.Fprintf(t, "\n %d\t%s\t%s\t%s\t%d",
			i+1, order.CreationDate.Format("2006-01-02"), cmdUtils.TruncateString(models.OrderStatuses[order.Status], 20), cmdUtils.TruncateString(order.Address, 20), order.Rate)
		if err != nil {
			return err
		}
	}

	err = t.Flush()
	if err != nil {
		return err
	}

	return nil
}
