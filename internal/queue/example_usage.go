package queue

// Example of how to initialize and use the statistics service
// This should be integrated into your main application initialization

/*
import (
	"templateGo/internal/queue"
	"templateGo/internal/repositories"
	"templateGo/internal/handlers/ai"
	"templateGo/internal/handlers/course"
)

func main() {
	// Initialize your existing dependencies
	repo := repositories.NewCourseRepository()
	aiAnalyzer := ai.NewGeminiAnalyzer() // or whatever your AI analyzer implementation is
	notification := notification.NewNotificationClient()
	metricsClient := metrics.NewDatadogMetricsClient()

	// Initialize the statistics service
	statisticsService := queue.NewStatisticsService(repo, aiAnalyzer)

	// Start the statistics service (this starts the background workers)
	statisticsService.Start()

	// Initialize the course handler with the statistics service
	courseHandler := course.NewCourseHandler(
		repo,
		notification,
		aiAnalyzer,
		metricsClient,
		statisticsService,
	)

	// Set up your routes with the courseHandler
	// ...

	// Make sure to stop the statistics service when shutting down
	defer statisticsService.Stop()

	// Start your server
	// ...
}
*/

// Usage example:
// When a submission is created/updated/graded, the statistics calculations
// will automatically be enqueued and processed in the background by worker goroutines.
// This prevents blocking the HTTP response while statistics are being calculated.

// The queue provides:
// - Asynchronous processing of statistics calculations
// - Automatic retry mechanism (3 retries by default)
// - Multiple worker goroutines for parallel processing
// - Graceful shutdown handling
// - Configurable buffer size and worker count
