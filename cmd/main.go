package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"templateGo/config"
	sql "templateGo/internal/repositories"
	controller "templateGo/internal/services"

	"github.com/rs/cors" // Importa la librería CORS
)

func main() {
	// Load environment variables from .env file
	config.LoadEnv()

	// Connect to database
	if err := sql.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer sql.CloseDB()

	// Setup routes
	mux := controller.SetupRoutes()

	// Configurar CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"}, // Permitir tu frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true, // Si tu frontend necesita enviar cookies o auth
	}).Handler(mux)

	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server listening on port %s\n", port)
	// Usamos corsHandler en vez de mux directamente
	if err := http.ListenAndServe(":"+port, corsHandler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
