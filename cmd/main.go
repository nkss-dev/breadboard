package main

import (
	"log"
	"os"
	"path/filepath"

	"nkssbackend"

	"github.com/joho/godotenv"
)

func init() {
	// Load the environment variables from files (if any)
	envs, _ := filepath.Glob(".env*")
	if len(envs) > 0 {
		if err := godotenv.Load(envs...); err != nil {
			log.Fatalln("Error while loading .env --> ", err)
		}
	}

	// Verify that all required environment variables are set
	envVars := []string{"DATABASE_URL", "HMAC_SECRET", "PORT"}
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			log.Fatalln(envVar + " not set")
		}
	}
}

func main() {
	server := nkssbackend.NewServer()
	log.Println("Listening at http://localhost:" + os.Getenv("PORT"))
	server.Run()
}
