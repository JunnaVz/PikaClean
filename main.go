package main

import (
	"fmt"
	"os"
	"teamdev/cmd"
	"teamdev/internal/models"
	"teamdev/internal/registry"

	"github.com/charmbracelet/log"
)

func main() {
	app := registry.App{}

	err := app.Config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run()

	if err != nil {
		fmt.Println("Error")
		log.Fatal(err)
	}

	err = initAdmin(app.Services)
	if err != nil {
		log.Fatal(err)
		return
	}

	if app.Config.Mode == "cmd" {
		cmdErr := cmd.RunMenu(app.Services)
		if cmdErr != nil {
			log.Fatal(cmdErr)
			return
		}
	} else {
		log.Error("Wrong app mode", "mode", app.Config.Mode)
	}
}

func initAdmin(services *registry.Services) error {
	admins, err := services.WorkerService.GetWorkersByRole(models.ManagerRole)
	if err != nil {
		return err
	}

	if len(admins) == 0 {
		defaultAdmin := &models.Worker{
			Email:       os.Getenv("ADMIN_EMAIL"),
			Name:        os.Getenv("ADMIN_NAME"),
			Surname:     os.Getenv("ADMIN_SURNAME"),
			Role:        models.ManagerRole,
			PhoneNumber: os.Getenv("ADMIN_ROLE"),
			Address:     os.Getenv("ADMIN_ADDRESS"),
		}
		_, err = services.WorkerService.Create(defaultAdmin, os.Getenv("ADMIN_PASSWORD"))
		if err != nil {
			return err
		}

		log.Info("Default admin created")
	}

	return nil
}
