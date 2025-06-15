package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"templateGo/internal/logger"
	"templateGo/internal/metrics"
	"templateGo/internal/repositories"
	controller "templateGo/internal/services"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Global clients
var datadogLogger *logger.DatadogLogger
var datadogMetrics *metrics.DatadogMetricsClient

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("No .env file found or couldn't load it, using system environment variables (expected in deployment environments)")
	} else {
		log.Println("Loaded environment variables from .env file")
	}

	apiKey := os.Getenv("DATADOG_API_KEY")

	if apiKey == "" {
		log.Println("Warning: DATADOG_API_KEY not set, logging to Datadog is disabled")
	} else {
		// Setting the site explicitly by environment variable
		os.Setenv("DATADOG_SITE", "us5.datadoghq.com")

		datadogLogger = logger.NewDatadogLogger(apiKey)
		datadogMetrics = metrics.NewDatadogMetricsClient(apiKey)

		// Log application startup
		err := datadogLogger.Info("Application starting up", map[string]any{
			"version":     "1.0.0",
			"environment": os.Getenv("ENVIRONMENT"),
		}, []string{"startup", "init"})

		if err != nil {
			log.Printf("Error logging to Datadog: %v", err)
		}
	}

	// Connect to database
	dbManager := repositories.NewDatabaseManager()
	if err := dbManager.ConnectDB(); err != nil {
		logError("Failed to connect to database", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbManager.CloseDB()

	// Setup routes with logger and metrics
	mux := controller.SetupRoutes(datadogLogger, datadogMetrics)

	// Configurar CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"}, // Permitir tu frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true, // Si tu frontend necesita enviar cookies o auth
	}).Handler(httpLoggerMiddleware(mux))

	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logInfo("Server starting", map[string]interface{}{
		"port": port,
	})

	fmt.Printf("Server listening on port %s\n", port)
	// Usamos corsHandler en vez de mux directamente
	if err := http.ListenAndServe(":"+port, corsHandler); err != nil {
		logError("Server failed to start", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("Error starting server: %v", err)
	}
}

// Helper function to log errors
func logError(message string, attributes map[string]interface{}) {
	log.Printf("ERROR: %s %v", message, attributes)
	if datadogLogger != nil {
		if err := datadogLogger.Error(message, attributes, nil); err != nil {
			log.Printf("Failed to send error log to Datadog: %v", err)
		}
	}
}

// Helper function to log info
func logInfo(message string, attributes map[string]interface{}) {
	log.Printf("INFO: %s %v", message, attributes)
	if datadogLogger != nil {
		if err := datadogLogger.Info(message, attributes, nil); err != nil {
			log.Printf("Failed to send info log to Datadog: %v", err)
		}
	}
}

// HTTP middleware to log requests
func httpLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request
		logInfo("HTTP Request", map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
		})

		// Wrap the response writer to capture status code
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default to 200 OK
		}

		// Call the next handler
		next.ServeHTTP(lrw, r)

		// Log response
		duration := time.Since(start).Milliseconds()
		attributes := map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status_code": lrw.statusCode,
			"duration_ms": duration,
		}

		if lrw.statusCode >= 400 {
			logError("HTTP Error Response", attributes)
		} else {
			logInfo("HTTP Response", attributes)
		}
	})
}

// Custom response writer to capture status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
