package course

import (
	"net/http"
	"templateGo/internal/model"
	"templateGo/internal/utils"

	"slices"

	"github.com/gin-gonic/gin"
)

// CreateUserFeedback adds feedback for a user in a course
func (h *courseHandlerImpl) CreateUserFeedback(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	studentID, ok := h.getUserID(c)
	if !ok {
		return
	}

	// Get instructor ID from token
	instructorID, ok := h.getUserIDFromToken(c)
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

	// Verify the instructor is the course creator or a teaching assistant
	if instructorID != course.CreatedBy && !contains(course.TeachingAssistants, instructorID) {
		utils.NewErrorResponse(c, http.StatusForbidden, "Forbidden", "Only instructors or TAs can provide feedback")
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
		CourseID:     courseID,
		StudentID:    studentID,
		InstructorID: instructorID,
		Comment:      request.Comment,
		Rating:       request.Rating,
	}

	if err := h.repo.CreateUserFeedback(feedback); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error creating feedback")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Feedback created successfully"})
}

// GetUserFeedbacks retrieves all feedback for a specific user
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
			"id":            feedback.ID,
			"course_id":     feedback.CourseID,
			"student_id":    feedback.StudentID,
			"instructor_id": feedback.InstructorID,
			"comment":       feedback.Comment,
			"rating":        feedback.Rating,
			"created_at":    feedback.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// Helper function to check if a string is in a slice
func contains(slice []string, str string) bool {
	return slices.Contains(slice, str)
}
