package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	m "breadboard/middleware"

	"github.com/joho/godotenv"
)

func main() {
	// Load the environment variables from files (if any)
	envs, _ := filepath.Glob("*.env")
	if len(envs) > 0 {
		if err := godotenv.Load(envs...); err != nil {
			log.Fatalln("Error while loading .env --> ", err)
		}
	}
	HMACSecret := []byte(os.Getenv("HMAC_SECRET"))
	if HMACSecret == nil {
		log.Fatalln("HMAC_SECRET is not set")
	}

	arglen := len(os.Args[1:])

	if arglen == 2 {
		fmt.Println("Generated jwt:", m.CreateJWT(os.Args[1], os.Args[2], HMACSecret))
		return
	} else {
		fmt.Fprintf(os.Stderr, "Error! Insufficient arguments provided!\n")
	}
}
