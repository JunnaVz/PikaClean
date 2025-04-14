package test_repositories

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	USER     = "admin"
	PASSWORD = "admin"
	DBNAME   = "ppo"
)

func SetupTestDatabase() (testcontainers.Container, *sql.DB) {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       DBNAME,
			"POSTGRES_PASSWORD": PASSWORD,
			"POSTGRES_USER":     USER,
		},
	}

	dbContainer, _ := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})

	host, _ := dbContainer.Host(context.Background())
	port, _ := dbContainer.MappedPort(context.Background(), "5432")

	dsnPGConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port.Int(), USER, PASSWORD, DBNAME)
	db, err := sql.Open("pgx", dsnPGConn)
	if err != nil {
		fmt.Println("Failed to open database connection: ", err)
		return dbContainer, nil
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to ping database: ", err)
		return dbContainer, nil
	}

	text, err := os.ReadFile("../../db/sql/init.sql")
	if err != nil {
		return dbContainer, nil
	}

	_, err = db.Exec(string(text))
	if err != nil {
		fmt.Println(err)
		return dbContainer, nil
	}

	return dbContainer, db
}