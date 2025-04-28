// Package config provides configuration structures and utilities for the PikaClean application.
// It handles loading and managing configuration parameters from environment variables
// and other sources, allowing the application to be configured for different environments.
// This file specifically contains database connection utilities.
package config

import (
	"database/sql"
	"fmt"
	_ "time"

	"github.com/charmbracelet/log"
	_ "github.com/jackc/pgx/v4/stdlib" // Import pgx driver
	_ "golang.org/x/net/context"
)

// DbConnectionFlags contains the necessary parameters for connecting to a PostgreSQL database.
// These fields are populated from environment variables or configuration files.
type DbConnectionFlags struct {
	Host     string `mapstructure:"host"`     // Database server hostname or IP address
	User     string `mapstructure:"user"`     // Username for database authentication
	Password string `mapstructure:"password"` // Password for database authentication
	Port     string `mapstructure:"port"`     // Port number the database server is listening on
	DBName   string `mapstructure:"dbname"`   // Name of the database to connect to
}

// InitPostgresDB establishes a connection to a PostgreSQL database using the
// parameters provided in DbConnectionFlags. It also configures the connection
// pool with appropriate settings for the application's needs.
//
// Parameters:
//   - logger: Logger for recording connection events and errors
//
// Returns:
//   - *sql.DB: Initialized database connection pool
//   - error: Connection error if the database cannot be reached or configured
func (p *DbConnectionFlags) InitPostgresDB(logger *log.Logger) (*sql.DB, error) {
	logger.Debug("POSTGRES! Start init postgreSQL", "user", p.User, "DBName", p.DBName,
		"host", p.Host, "port", p.Port)

	dsnPGConn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		p.User, p.DBName, p.Password,
		p.Host, p.Port)

	db, err := sql.Open("pgx", dsnPGConn)
	if err != nil {
		logger.Fatal("POSTGRES! Error in method open")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal("POSTGRES! Error in method ping")
		return nil, err
	}

	db.SetMaxOpenConns(10)

	logger.Info("POSTGRES! Successfully init postgreSQL")
	return db, nil
}
