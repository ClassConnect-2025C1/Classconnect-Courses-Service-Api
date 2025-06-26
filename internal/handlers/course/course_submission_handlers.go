package course

import (
	"log"
	"net/http"
	"templateGo/internal/model"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
)

// PutSubmissionOfCurrentUser creates or updates a submission
func (h *courseHandlerImpl) PutSubmissionOfCurrentUser(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}
	assignmentID, ok := h.getAssignmentID(c)
	if !ok {
		return
	}
	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}
	userEmail, ok := h.getUserEmailFromToken(c)
	if !ok {
		return
	}
	// Check if course exists
	_, ok = h.getCourseByID(c, courseID)
	if !ok {
		return
	}
	var req model.CreateSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}
	submission := &model.Submission{
		CourseID:     courseID,
		AssignmentID: assignmentID,
		UserID:       userID,
		Content:      req.Content,
		Files:        req.Files,
	}

	if err := h.repo.PutSubmission(submission); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error creating submission")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Submission created/updated successfully"})

	h.CalculateAndStoreCourseStatistics(courseID, userID, userEmail)
	h.CalculateAndStoreUserCourseStatistics(courseID, userID, userEmail)
}

// DeleteSubmissionOfCurrentUser removes a user's submission
func (h *courseHandlerImpl) DeleteSubmissionOfCurrentUser(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}
	assignmentID, ok := h.getAssignmentID(c)
	if !ok {
		return
	}
	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}
	userEmail, ok := h.getUserEmailFromToken(c)
	if !ok {
		return
	}
	// Check if course exists
	_, ok = h.getCourseByID(c, courseID)
	if !ok {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Course not found")
		return
	}
	// Check if assignment exists
	assignment, ok := h.getAssignmentByID(c, assignmentID)
	if !ok {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Assignment not found")
		return
	}
	// Check if assignment belongs to the course
	if assignment.CourseID != courseID {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "Assignment does not belong to this course")
		return
	}
	// First get the submission of the user to check if it exists
	submission, err := h.repo.GetSubmissionByUserID(courseID, assignmentID, userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Submission not found")
		return
	}
	if err := h.repo.DeleteSubmission(submission.ID); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error deleting submission")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Submission deleted successfully"})

	h.CalculateAndStoreCourseStatistics(courseID, userID, userEmail)
	h.CalculateAndStoreUserCourseStatistics(courseID, userID, userEmail)
}

// GetSubmissionOfCurrentUser returns the current user's submission
func (h *courseHandlerImpl) GetSubmissionOfCurrentUser(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}
	assignmentID, ok := h.getAssignmentID(c)
	if !ok {
		return
	}
	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}
	submission, err := h.repo.GetSubmissionByUserID(courseID, assignmentID, userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Submission not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": submission})
}

// GetSubmissionByUserID returns a specific user's submission
func (h *courseHandlerImpl) GetSubmissionByUserID(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}
	assignmentID, ok := h.getAssignmentID(c)
	if !ok {
		return
	}
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}
	submission, err := h.repo.GetSubmissionByUserID(courseID, assignmentID, userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Submission not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": submission})

}

// GetSubmissions returns all submissions for an assignment
func (h *courseHandlerImpl) GetSubmissions(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}
	if !h.isCourseCreatorOrAssistant(c, courseID) {
		utils.NewErrorResponse(c, http.StatusForbidden, "Forbidden", "You do not have permission to access this resource")
		return
	}
	assignmentID, ok := h.getAssignmentID(c)
	if !ok {
		return
	}
	submissions, err := h.repo.GetSubmissions(courseID, assignmentID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving submissions")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": submissions})
}

// GradeSubmission allows grading a submission with feedback
func (h *courseHandlerImpl) GradeSubmission(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}
	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}
	userEmail, ok := h.getUserEmailFromToken(c)
	if !ok {
		return
	}
	submissionID, ok := h.getSubmissionID(c)
	if !ok {
		return
	}
	submission, err := h.repo.GetSubmission(submissionID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Submission not found")
		return
	}
	var req model.GradeSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}
	submission.Grade = req.Grade
	submission.Feedback = req.Feedback

	if err := h.repo.PutSubmission(submission); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error grading submission")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Submission graded successfully"})

	h.CalculateAndStoreCourseStatistics(courseID, userID, userEmail)
	h.CalculateAndStoreUserCourseStatistics(courseID, userID, userEmail)
}

// GetAIGeneratedGrade retrieves AI-generated grade for a submission
func (h *courseHandlerImpl) GetAIGeneratedGradeAndFeedback(c *gin.Context) {
	submissionID, ok := h.getSubmissionID(c)
	if !ok {
		return
	}
	submission, err := h.repo.GetSubmission(submissionID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Submission not found")
		return
	}
	assignmentID, ok := h.getAssignmentID(c)
	if !ok {
		return
	}
	// Check if assignment exists
	assignment, ok := h.getAssignmentByID(c, assignmentID)
	if !ok {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Assignment not found")
		return
	}
	// Simulate AI-generated grade retrieval
	aiGrade, aiFeedback, err := h.aiAnalyzer.GenerateGradeAndFeedback(assignment.Description, submission.Files)
	if err != nil {
		// print it to the console for debugging
		log.Printf("Error generating AI grade/feedback: %v", err)

		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error generating AI grade/feedback")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"grade":    aiGrade,
		"feedback": aiFeedback,
	}})
}
