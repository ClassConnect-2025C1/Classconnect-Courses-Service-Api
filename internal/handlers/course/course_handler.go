package course

import (
	"templateGo/internal/handlers/ai"
	"templateGo/internal/handlers/notification"
	"templateGo/internal/repositories"
)

// courseHandlerImpl implements CourseHandler interface
type courseHandlerImpl struct {
	repo         repositories.CourseRepository
	notification *notification.NotificationClient
	aiAnalyzer   ai.FeedbackAnalyzer
}

// NewCourseHandler creates a new CourseHandler
func NewCourseHandler(repo repositories.CourseRepository, noti *notification.NotificationClient, analyzer ai.FeedbackAnalyzer) CourseHandler {
	return &courseHandlerImpl{
		repo:         repo,
		notification: noti,
		aiAnalyzer:   analyzer,
	}
}
