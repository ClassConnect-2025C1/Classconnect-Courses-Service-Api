package course

import (
	"fmt"
	"net/http"
	"os"
	"templateGo/internal/model"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
)

// CreateCourseFeedback handles feedback submission for a course
// @Summary Submit feedback for a course
// @Description Submit feedback and rating for a specific course
// @Tags feedback
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Param feedback body model.FeedbackRequest true "Feedback information"
// @Success 201 {object} model.SuccessResponse{data=model.FeedbackResponse}
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/feedback [post]
func (h *courseHandlerImpl) CreateCourseFeedback(c *gin.Context) {
	courseID, ok := h.getCourseID(c)

	if !ok {
		return
	}

	// Check if course exists
	_, ok = h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	var req model.CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}

	// Check if user is enrolled in the course
	isEnrolled, err := h.repo.IsUserEnrolled(courseID, userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error checking enrollment status")
		return
	}

	if !isEnrolled {
		utils.NewErrorResponse(c, http.StatusForbidden, "Forbidden", "Only enrolled users can provide feedback")
		return
	}

	feedback := &model.CourseFeedback{
		CourseID: courseID,
		Rating:   req.Rating,
		Comment:  req.Comment,
		Summary:  req.Summary,
	}

	if err := h.repo.CreateFeedback(feedback); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error creating feedback")
		return
	}

	// Track feedback creation metric
	if h.metricsClient != nil {
		// Add relevant tags for better filtering and visualization
		tags := []string{
			fmt.Sprintf("course_id:%d", courseID),
			fmt.Sprintf("rating:%d", req.Rating),
			fmt.Sprintf("user_id:%s", userID),
			fmt.Sprintf("environment:%s", os.Getenv("ENVIRONMENT")),
		}

		fmt.Printf("Sending metric: classconnect.course.feedback.created with tags: %v\n", tags)
		if err := h.metricsClient.IncrementCounter("classconnect.course.feedback.created", tags); err != nil {
			fmt.Printf("Error sending course feedback creation metric: %v\n", err)
		} else {
			fmt.Println("Successfully sent feedback creation metric to Datadog")
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Feedback submitted successfully"})

}

// GetCourseFeedbacks returns all feedback for a course
// @Summary Get all feedback for a course
// @Description Retrieve all feedback submitted for a specific course
// @Tags feedback
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/feedbacks [get]
func (h *courseHandlerImpl) GetCourseFeedbacks(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	// Check if course exists
	_, ok = h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	feedbackList, err := h.repo.GetFeedbacksForCourse(courseID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving feedback")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": feedbackList})
}

// GetAICourseFeedbackAnalysis returns AI-generated analysis of course feedback
// @Summary Get AI-generated analysis of course feedbacks
// @Description Get AI-powered analysis and insights from course feedback
// @Tags feedback
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/ai-feedback-analysis [get]
func (h *courseHandlerImpl) GetAICourseFeedbackAnalysis(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	// Check if course exists
	course, ok := h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	// Get all feedback for this course
	feedbacks, err := h.repo.GetFeedbacksForCourse(courseID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving feedback")
		return
	}

	// If there's no feedback, return an appropriate message
	if len(feedbacks) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("No feedback available for analysis for course '%s'", course.Title),
		})
		return
	}

	// Use the AI analyzer to analyze the feedback
	analysis, err := h.aiAnalyzer.GenerateCourseFeedbackAnalysis(course.Title, feedbacks)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "AI Analysis Error", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":           analysis,
		"feedback_count": len(feedbacks),
	})
}

// GetAIUserFeedbackAnalysis returns AI-generated analysis of user feedback
// @Summary Get AI-generated analysis of user feedbacks
// @Description Get AI-powered analysis and insights from user feedback
// @Tags feedback
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /user/{user_id}/ai-feedback-analysis [get]
func (h *courseHandlerImpl) GetAIUserFeedbackAnalysis(c *gin.Context) {
	// Get user ID from url
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	// Get all feedback for the user
	feedbacks, err := h.repo.GetUserFeedbacks(userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving user feedback")
		return
	}

	// If there's no feedback, return an appropriate message
	if len(feedbacks) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No feedback available for analysis",
		})
		return
	}

	// Use the AI analyzer to analyze the feedback
	analysis, err := h.aiAnalyzer.GenerateUserFeedbackAnalysis(feedbacks)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "AI Analysis Error", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":           analysis,
		"feedback_count": len(feedbacks),
	})
}
