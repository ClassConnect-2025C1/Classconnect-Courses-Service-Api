package main

// Example of how to use the updated service controller with proper lifecycle management

/*
import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"templateGo/internal/services"
	"templateGo/internal/logger"
	"templateGo/internal/metrics"
)

func main() {
	// Initialize logger and metrics
	ddLogger := logger.NewDatadogLogger()
	ddMetrics := metrics.NewDatadogMetricsClient()

	// Setup routes and services - now returns a ServiceManager
	serviceManager := services.SetupRoutes(ddLogger, ddMetrics)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: serviceManager, // ServiceManager implements http.Handler
	}

	// Start server in a goroutine
	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Shutdown all managed services (including statistics service)
	if err := serviceManager.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down services: %v", err)
	}

	log.Println("Server gracefully stopped")
}
*/

// Key changes:
// 1. SetupRoutes now returns a *ServiceManager instead of http.Handler
// 2. ServiceManager implements http.Handler so it can be used as server handler
// 3. ServiceManager handles starting/stopping the statistics service automatically
// 4. Proper graceful shutdown with context timeout
// 5. All background workers are properly cleaned up on shutdown
