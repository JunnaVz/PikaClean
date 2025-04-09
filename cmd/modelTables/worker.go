package modelTables

import (
	"fmt"
	"lab3/internal/models"
	"lab3/internal/registry"
	"os"
	"text/tabwriter"
)

func Workers(services registry.Services, workers []models.Worker) error {
	var err error

	t := new(tabwriter.Writer)
	t.Init(os.Stdout, 1, 4, 2, ' ', 0)

	_, err = fmt.Fprintf(t, "\n %s\t%s\t%s\t%s\t%s\t%s\n",
		"№", "Имя", "Роль", "Телефон", "Email", "Ср. оценка")
	if err != nil {
		fmt.Println(err)
	}

	for i, worker := range workers {
		workersRate, _ := services.WorkerService.GetAverageOrderRate(&worker)

		fmt.Fprintf(t, " %d\t%s\t%s\t%s\t%s\t%f\n",
			i+1, worker.FullName(), worker.DisplayRole(), worker.PhoneNumber, worker.Email, workersRate)
	}

	err = t.Flush()
	if err != nil {
		return err
	}

	return nil
}
