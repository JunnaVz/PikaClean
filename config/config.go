// Package config provides configuration structures and utilities for the PikaClean application.
// It handles loading and managing configuration parameters from environment variables
// and other sources, allowing the application to be configured for different environments.
package config

import "os"

// Config represents the main application configuration.
// It contains all settings needed to run the PikaClean application,
// including database connection parameters, server settings, and logging configuration.
type Config struct {
	DBFlags  DbConnectionFlags `mapstructure:"postgres"` // Database connection parameters
	Address  string            `mapstructure:"address"`  // Server bind address
	Port     string            `mapstructure:"port"`     // Server listen port
	LogLevel string            `mapstructure:"loglevel"` // Logging verbosity level (debug, info)
	LogFile  string            `mapstructure:"logfile"`  // Path to log file
	Mode     string            `mapstructure:"mode"`     // Application mode (development, production)
	DBType   string            `mapstructure:"dbtype"`   // Database type (postgres, etc.)
}

// ParseConfig loads configuration values from environment variables into the Config struct.
// This method allows for flexible configuration through environment variables, making the
// application suitable for containerized deployments and different environments.
//
// Returns:
//   - error: Currently always returns nil, but could be extended to validate configuration
func (c *Config) ParseConfig() error {
	// Set database connection parameters from environment variables
	c.DBFlags.Host = os.Getenv("POSTGRES_HOST")
	c.DBFlags.User = os.Getenv("POSTGRES_USER")
	c.DBFlags.Password = os.Getenv("POSTGRES_PASSWORD")
	c.DBFlags.Port = os.Getenv("POSTGRES_PORT")
	c.DBFlags.DBName = os.Getenv("POSTGRES_DATABASE")
	c.Address = os.Getenv("ADDRESS")
	c.Port = os.Getenv("PORT")
	c.LogLevel = os.Getenv("LOGLEVEL")
	c.LogFile = os.Getenv("LOGFILE")
	c.Mode = os.Getenv("MODE")
	c.DBType = os.Getenv("DBTYPE")

	return nil
}
