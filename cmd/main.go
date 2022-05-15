package main

import (
	"NKSS-backend/config"
	"NKSS-backend/pkg"
	"fmt"
)

func main() {
	config := config.GetConfig()
	app := &pkg.App{}
	app.Initialise(config)

	fmt.Println("API Online")
	app.Run()
}
