package main

import (
	"NKSS-backend/config"
	"NKSS-backend/pkg"
	"fmt"
)

func main() {
	config := config.GetConfig()
	app := pkg.NewApp(&config)

	fmt.Println("API Online")
	app.Run()
}
