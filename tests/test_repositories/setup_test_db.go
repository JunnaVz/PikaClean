package test_repositories

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"path/filepath"
	"strings"
)

const (
	USER     = "postgres"
	PASSWORD = "0252"
	DBNAME   = "ppo"
)

func SetupTestDatabase() (testcontainers.Container, *sql.DB) {
	ctx := context.Background()

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

	dbContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		fmt.Printf("Failed to start container: %v\n", err)
		return nil, nil
	}

	host, err := dbContainer.Host(ctx)
	if err != nil {
		fmt.Printf("Failed to get container host: %v\n", err)
		return dbContainer, nil
	}

	port, err := dbContainer.MappedPort(ctx, "5432")
	if err != nil {
		fmt.Printf("Failed to get container port: %v\n", err)
		return dbContainer, nil
	}

	dsnPGConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port.Int(), USER, PASSWORD, DBNAME)

	db, err := sql.Open("pgx", dsnPGConn)
	if err != nil {
		fmt.Printf("Failed to open database connection: %v\n", err)
		return dbContainer, nil
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("Failed to ping database: %v\n", err)
		return dbContainer, nil
	}

	err = initializeDatabase(db)
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		return dbContainer, nil
	}

	return dbContainer, db
}

func initializeDatabase(db *sql.DB) error {
	// для прохода тестов, в базе не должны быть данными о workers, users, orders, order_contains_tasks
	_, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;`)
	if err != nil {
		return fmt.Errorf("failed to create uuid extension: %v", err)
	}

	// Read schema file
	schemaPath := filepath.Join("..", "..", "db", "sql", "init.sql")
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %v", err)
	}

	// Split and execute schema statements
	statements := strings.Split(string(schema), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" ||
			strings.HasPrefix(stmt, "--") ||
			strings.Contains(stmt, "COPY") ||
			strings.Contains(stmt, "CSV HEADER") {
			continue
		}

		_, err = db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %v\nStatement: %s", err, stmt)
		}
	}

	err = insertTestData(db)
	if err != nil {
		return fmt.Errorf("failed to insert test data: %v", err)
	}

	return nil
}

func insertTestData(db *sql.DB) error {
	// Insert minimal required test data
	tasks := []struct {
		id             string
		name           string
		pricePerSingle float64
		category       int
	}{
		{"daa09f13-0ba4-4511-a105-0e612ca11603", "Task 1", 300, 1},
		{"3068fe74-e9fc-40ac-9674-e0bef4f83083", "Task 2", 500, 1},
	}

	for _, task := range tasks {
		_, err := db.Exec(`
			INSERT INTO tasks (id, name, price_per_single, category)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO NOTHING
		`, task.id, task.name, task.pricePerSingle, task.category)
		if err != nil {
			return fmt.Errorf("failed to insert task data: %v", err)
		}
	}

	return nil
}
