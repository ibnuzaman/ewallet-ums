package main

import (
	"log"

	"github.com/ibnuzaman/ewallet-ums/cmd"
	"github.com/ibnuzaman/ewallet-ums/database"
	"github.com/ibnuzaman/ewallet-ums/helpers"
)

func main() {
	// Setup configuration first
	if err := helpers.SetupConfig(); err != nil {
		log.Printf("Warning: Error loading configuration: %v", err)
	}

	// Setup logger after config
	helpers.SetupLogger()

	// Initialize database connection
	db, err := database.InitPostgres()
	if err != nil {
		helpers.Logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := database.ClosePostgres(); err != nil {
			helpers.Logger.Errorf("Error closing database: %v", err)
		}
	}()

	helpers.Logger.Infof("Database initialized: %v", db.Stats())

	// Start HTTP server
	cmd.ServerHTTP()
}
