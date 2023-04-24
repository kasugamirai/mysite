// main.go

package main

import (
	"fmt"
	"log"

	"xy.com/mysite/config"
	"xy.com/mysite/database"
	"xy.com/mysite/handlers"
	"xy.com/mysite/routes"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize the database connection
	err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Set up the Gin router
	router := routes.SetupRouter()

	// Start the HandleMessages goroutine for chat functionality
	go handlers.HandleMessages() // new add

	// Start the server
	port := config.Instance.Server.Port
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on port %s\n", port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
