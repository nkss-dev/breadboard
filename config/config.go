package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseUrl string
	HMACSecret  []byte
}

func GetConfig() Config {
	DatabaseUrl := os.Getenv("DATABASE_URL")
	HMACSecret := []byte(os.Getenv("HMAC_SECRET"))
	if DatabaseUrl == "" {
		log.Fatalln("DATABASE_URL is not set")
	}
	if HMACSecret == nil {
		log.Fatalln("HMAC_SECRET is not set")
	}

	return Config{DatabaseUrl: DatabaseUrl, HMACSecret: HMACSecret}
}
