package main

import (
	"log"

	"github.com/ibnuzaman/ewallet-ums/cmd"
	"github.com/ibnuzaman/ewallet-ums/helpers"
)

func main() {
	// Setup configuration first
	if err := helpers.SetupConfig(); err != nil {
		log.Printf("Warning: Error loading configuration: %v", err)
	}

	// Setup logger after config
	helpers.SetupLogger()

	// Start HTTP server
	cmd.ServerHttp()
}
