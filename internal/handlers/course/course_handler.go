package course

import (
	"templateGo/internal/handlers/ai"
	"templateGo/internal/handlers/notification"
	"templateGo/internal/metrics"
	"templateGo/internal/queue"
	"templateGo/internal/repositories"
)

// courseHandlerImpl implements CourseHandler interface
type courseHandlerImpl struct {
	repo              repositories.CourseRepository
	notification      *notification.NotificationClient
	aiAnalyzer        ai.FeedbackAnalyzer
	metricsClient     *metrics.DatadogMetricsClient
	statisticsService *queue.StatisticsService
}

// NewCourseHandler creates a new CourseHandler
func NewCourseHandler(
	repo repositories.CourseRepository,
	notification *notification.NotificationClient,
	aiAnalyzer ai.FeedbackAnalyzer,
	metricsClient *metrics.DatadogMetricsClient,
	statisticsService *queue.StatisticsService,
) CourseHandler {
	return &courseHandlerImpl{
		repo:              repo,
		notification:      notification,
		aiAnalyzer:        aiAnalyzer,
		metricsClient:     metricsClient,
		statisticsService: statisticsService,
	}
}
