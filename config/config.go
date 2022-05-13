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
			Database: os.Getenv("database"),
			Host:     os.Getenv("host"),
			Password: os.Getenv("password"),
			Port:     os.Getenv("port"),
			User:     os.Getenv("user"),
		},
	}
}
