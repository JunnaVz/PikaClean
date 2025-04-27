//go install golang.org/x/tools/cmd/godoc@latest
//godoc -http=:6060
//http://localhost:6060/pkg/teamdev/
//go doc teamdev/internal/models
//go doc teamdev/internal/registry
//go doc teamdev/internal/repository
//go doc teamdev/internal/services

// Package main provides the entry point for the PikaClean application.
//
// It initializes the application configuration, sets up service dependencies,
// ensures a default admin user exists, and runs the application in the
// configured mode (currently supporting command-line interface).
package main

import (
	"fmt"
	"os"
	"teamdev/cmd"
	"teamdev/internal/models"
	"teamdev/internal/registry"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

// main is the entry point for the TeamDev application.
// It initializes the application, sets up dependencies,
// creates a default admin user if one doesn't exist,
// and runs either in command-line mode or exits with an error
// if an invalid mode is specified.
func main() {
	app := registry.App{}

	err := godotenv.Load()
	if err != nil {
		log.Warn("No .env file found")
	}

	err = app.Config.ParseConfig()
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

// initAdmin ensures there is at least one admin user in the system.
// If no admin users exist, it creates a default admin user using environment variables.
//
// Parameters:
//   - services: Application services registry containing worker service
//
// Returns:
//   - error: If there was an error retrieving workers or creating the admin user
//
// Environment Variables Used:
//   - ADMIN_EMAIL: Email address for the default admin
//   - ADMIN_NAME: First name of the default admin
//   - ADMIN_SURNAME: Last name of the default admin
//   - ADMIN_PHONE: Phone number for the default admin (environment variable appears misnamed)
//   - ADMIN_ADDRESS: Physical address of the default admin
//   - ADMIN_PASSWORD: Password for the default admin user
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
			PhoneNumber: os.Getenv("ADMIN_PHONE"),
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
