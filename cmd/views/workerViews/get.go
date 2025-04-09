package workerViews

import (
	"fmt"
	"lab3/internal/models"
	"lab3/internal/registry"
)

func Get(service registry.Services, worker *models.Worker) error {
	workerFromDB, err := service.WorkerService.GetWorkerByID(worker.ID)
	if err != nil {
		return err
	}

	fmt.Print("\nWorker info:\n")
	fmt.Printf("Роль: %s\nEmail: %s\nИмя: %s\nФамилия: %s\nТелефон: %s\nАдрес: %s\n", models.WorkerRole[workerFromDB.Role], workerFromDB.Email, workerFromDB.Name, workerFromDB.Surname, workerFromDB.PhoneNumber, workerFromDB.Address)
	fmt.Print("----------------\n")
	return nil
}
