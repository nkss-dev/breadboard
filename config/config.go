package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseUrl string
}

func GetConfig() Config {
	DatabaseUrl := os.Getenv("DATABASE_URL")
	if DatabaseUrl == "" {
		log.Fatalln("DATABASE_URL is not set")
	}

	return Config{DatabaseUrl: DatabaseUrl}
}
