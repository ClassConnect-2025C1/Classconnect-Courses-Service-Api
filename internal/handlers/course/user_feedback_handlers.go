package course

import (
	"net/http"
	"templateGo/internal/model"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
)

// CreateUserFeedback adds feedback for a user in a course
// @Summary Add feedback for a user in a course
// @Description Submit feedback for a specific user in a course context
// @Tags feedback
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Param user_id path string true "User ID"
// @Param feedback body model.FeedbackRequest true "Feedback information"
// @Success 201 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/user/{user_id}/feedback [post]
func (h *courseHandlerImpl) CreateUserFeedback(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	studentID, ok := h.getUserID(c)
	if !ok {
		return
	}

	// Check if the course exists
	course, ok := h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	// Check if student is enrolled in the course
	isEnrolled, err := h.repo.IsUserEnrolled(courseID, studentID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error checking enrollment")
		return
	}

	if !isEnrolled {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Request", "User is not enrolled in this course")
		return
	}

	// Parse request body
	var request struct {
		Comment string `json:"comment" binding:"required"`
		Rating  uint   `json:"rating" binding:"required,min=1,max=5"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}

	feedback := &model.UserFeedback{
		CourseID:    courseID,
		StudentID:   studentID,
		CourseTitle: course.Title,
		Comment:     request.Comment,
		Rating:      request.Rating,
	}

	if err := h.repo.CreateUserFeedback(feedback); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error creating feedback")
		return
	}

	h.notification.SendNotification(studentID, course.Title, "feedback")
	c.JSON(http.StatusCreated, gin.H{"message": "Feedback created successfully"})
}

// GetUserFeedbacks retrieves all feedback for a specific user
// @Summary Get all feedback for a user
// @Description Retrieve all feedback submitted for a specific user
// @Tags feedback
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /user/{user_id}/feedbacks [get]
func (h *courseHandlerImpl) GetUserFeedbacks(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	feedbacks, err := h.repo.GetUserFeedbacks(userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving user feedbacks")
		return
	}

	// Format the response
	response := make([]gin.H, 0, len(feedbacks))
	for _, feedback := range feedbacks {
		response = append(response, gin.H{
			"id":           feedback.ID,
			"course_id":    feedback.CourseID,
			"student_id":   feedback.StudentID,
			"course_title": feedback.CourseTitle,
			"comment":      feedback.Comment,
			"rating":       feedback.Rating,
			"created_at":   feedback.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
