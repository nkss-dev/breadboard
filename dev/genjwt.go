package main

import (
	"fmt"
	"log"
	"os"

	m "breadboard/middleware"
)

func main() {
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
