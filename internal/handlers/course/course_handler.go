package course

import (
	"templateGo/internal/handlers/ai"
	"templateGo/internal/handlers/notification"
	"templateGo/internal/metrics"
	"templateGo/internal/repositories"
)

// courseHandlerImpl implements CourseHandler interface
type courseHandlerImpl struct {
	repo          repositories.CourseRepository
	notification  *notification.NotificationClient
	aiAnalyzer    ai.FeedbackAnalyzer
	metricsClient *metrics.DatadogMetricsClient
}

// NewCourseHandler creates a new CourseHandler
func NewCourseHandler(
	repo repositories.CourseRepository,
	notification *notification.NotificationClient,
	aiAnalyzer ai.FeedbackAnalyzer,
	metricsClient *metrics.DatadogMetricsClient,
) CourseHandler {
	return &courseHandlerImpl{
		repo:          repo,
		notification:  notification,
		aiAnalyzer:    aiAnalyzer,
		metricsClient: metricsClient,
	}
}
