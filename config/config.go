package config

import (
	"os"
)

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Database string
	Host     string
	Password string
	Port     string
	User     string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Database: os.Getenv("pgdatabase"),
			Host:     os.Getenv("pghost"),
			Password: os.Getenv("pgpassword"),
			Port:     os.Getenv("pgport"),
			User:     os.Getenv("pguser"),
		},
	}
}
