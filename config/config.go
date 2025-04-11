package config

import "os"

type Config struct {
	DBFlags  DbConnectionFlags `mapstructure:"postgres"`
	Address  string            `mapstructure:"address"`
	Port     string            `mapstructure:"port"`
	LogLevel string            `mapstructure:"loglevel"`
	LogFile  string            `mapstructure:"logfile"`
	Mode     string            `mapstructure:"mode"`
	DBType   string            `mapstructure:"dbtype"`
}

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
