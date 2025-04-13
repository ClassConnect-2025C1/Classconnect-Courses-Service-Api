package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"templateGo/config"
	"templateGo/controller"
	"templateGo/dbConfig/sql"
)

func main() {
	// Load environment variables from .env file
	config.LoadEnv()

	// Connect to database
	if err := sql.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer sql.CloseDB()

	mux := controller.SetupRoutes()

	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
