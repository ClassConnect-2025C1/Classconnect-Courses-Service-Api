package ai

import "templateGo/internal/model"

// FeedbackAnalyzer represents a service that can analyze course feedback
type FeedbackAnalyzer interface {
	// GenerateCourseFeedbackAnalysis analyzes a collection of course feedback and returns insights
	GenerateCourseFeedbackAnalysis(courseTitle string, feedbacks []model.CourseFeedback) (string, error)

	// GenerateGradeAndFeedback generates a grade and feedback for a submission
	GenerateGradeAndFeedback(assignmentDescription string, submissionFiles []model.SubmissionFile) (int, string, error)

	// GenerateUserFeedbackAnalysis generates AI analysis for user feedback
	GenerateUserFeedbackAnalysis(feedbacks []model.UserFeedback) (string, error)

	// GenerateCourseSuggestionsBasedOnStats generates course suggestions based on statistics
	GenerateCourseSuggestionsBasedOnStats(lastGradeTendency string, lastSubmissionRateTendency string, averageGrade float64) (string, error)
}
