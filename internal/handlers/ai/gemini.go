package ai

import "templateGo/internal/model"

// FeedbackAnalyzer represents a service that can analyze course feedback
type FeedbackAnalyzer interface {
	// AnalyzeFeedback analyzes a collection of course feedback and returns insights
	AnalyzeFeedback(courseTitle string, feedbacks []model.CourseFeedback) (string, error)
}
